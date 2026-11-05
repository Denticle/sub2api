package oauthprovider

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

const (
	codeTTL  = 60 * time.Second
	tokenTTL = 10 * time.Minute
	// studioKeyName 是托管密钥在「API 密钥」页显示的名字
	studioKeyName = "工作台"
)

var (
	ErrInvalidRequest = errors.New("invalid_request")
	ErrInvalidGrant   = errors.New("invalid_grant")
	ErrInvalidToken   = errors.New("invalid_token")
	ErrUserInactive   = errors.New("user is not active")
	ErrKeyUnavailable = errors.New("studio key is disabled")
)

// userReader / keyManager 只声明本模块用到的方法，便于测试替换。
type userReader interface {
	GetByID(ctx context.Context, id int64) (*service.User, error)
}

type keyManager interface {
	Create(ctx context.Context, userID int64, req service.CreateAPIKeyRequest) (*service.APIKey, error)
	GetByID(ctx context.Context, id int64) (*service.APIKey, error)
}

type Service struct {
	config configStore
	repo   repository
	users  userReader
	keys   keyManager
	groups groupLister
	now    func() time.Time
}

func NewService(db *sql.DB, settings service.SettingRepository, encryptor service.SecretEncryptor, users *service.UserService, keys *service.APIKeyService, plaza *service.ModelPlazaService) *Service {
	return &Service{
		config: configStore{settings: settings, encryptor: encryptor},
		repo:   newSQLRepository(db),
		users:  users,
		keys:   keys,
		groups: plaza,
		now:    time.Now,
	}
}

// ---- 配置 ----

func (s *Service) PublicConfig(ctx context.Context) PublicConfig {
	cfg, err := s.config.load(ctx)
	if err != nil {
		return PublicConfig{}
	}
	return cfg.public()
}

// AdminConfig 返回配置，Client Secret 只告诉有没有。
func (s *Service) AdminConfig(ctx context.Context) (Config, bool, error) {
	cfg, err := s.config.load(ctx)
	if err != nil {
		return Config{}, false, err
	}
	has := cfg.ClientSecret != ""
	cfg.ClientSecret = ""
	return cfg, has, nil
}

// SaveConfig 保存配置；ClientSecret 留空表示沿用旧值。
func (s *Service) SaveConfig(ctx context.Context, next Config) error {
	current, err := s.config.load(ctx)
	if err != nil {
		return err
	}
	next.ClientID = strings.TrimSpace(next.ClientID)
	next.ClientSecret = strings.TrimSpace(next.ClientSecret)
	next.RedirectURI = strings.TrimSpace(next.RedirectURI)
	next.AppURL = strings.TrimSpace(next.AppURL)
	next.AppName = strings.TrimSpace(next.AppName)
	if next.ClientSecret == "" {
		next.ClientSecret = current.ClientSecret
	}
	if err := validateConfig(next); err != nil {
		return err
	}
	return s.config.save(ctx, next)
}

// ---- 授权 ----

type AuthorizeRequest struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type AuthorizeCheck struct {
	AppName string `json:"app_name"`
	// Consented 为 true 时前端不需要显示确认页，直接调 Approve
	Consented bool `json:"consented"`
}

func (s *Service) validateAuthorize(ctx context.Context, req AuthorizeRequest) (Config, error) {
	cfg, err := s.config.load(ctx)
	if err != nil {
		return Config{}, err
	}
	if !cfg.Enabled {
		return Config{}, ErrDisabled
	}
	if subtle.ConstantTimeCompare([]byte(req.ClientID), []byte(cfg.ClientID)) != 1 {
		return Config{}, ErrInvalidClient
	}
	if req.RedirectURI != cfg.RedirectURI {
		return Config{}, ErrInvalidRedirect
	}
	if req.ResponseType != "code" || req.CodeChallengeMethod != "S256" || !validChallenge(req.CodeChallenge) || req.State == "" || len(req.State) > 512 {
		return Config{}, ErrInvalidRequest
	}
	return cfg, nil
}

// CheckAuthorize 只校验参数，不签发任何东西；前端用它决定是否显示确认页。
func (s *Service) CheckAuthorize(ctx context.Context, userID int64, req AuthorizeRequest) (*AuthorizeCheck, error) {
	cfg, err := s.validateAuthorize(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := s.requireActiveUser(ctx, userID); err != nil {
		return nil, err
	}
	consented, err := s.repo.hasConsent(ctx, userID, cfg.ClientID)
	if err != nil {
		return nil, err
	}
	return &AuthorizeCheck{AppName: cfg.AppName, Consented: consented}, nil
}

// Approve 记录授权、准备好专用密钥，签发授权码，返回带 code 和 state 的回调地址。
func (s *Service) Approve(ctx context.Context, userID int64, req AuthorizeRequest) (string, error) {
	cfg, err := s.validateAuthorize(ctx, req)
	if err != nil {
		return "", err
	}
	if err := s.requireActiveUser(ctx, userID); err != nil {
		return "", err
	}
	// 授权时就把密钥建好：建不出来（例如分组不可用）应该在这里报错，
	// 而不是等用户进了工作台第一次生成才失败。
	if _, err := s.ensureStudioKey(ctx, cfg, userID); err != nil {
		return "", err
	}
	if err := s.repo.saveConsent(ctx, userID, cfg.ClientID); err != nil {
		return "", err
	}
	code, err := randomToken()
	if err != nil {
		return "", err
	}
	if err := s.repo.createCode(ctx, hashToken(code), authorizationCode{
		ClientID:      cfg.ClientID,
		UserID:        userID,
		RedirectURI:   cfg.RedirectURI,
		CodeChallenge: req.CodeChallenge,
		ExpiresAt:     s.now().Add(codeTTL),
	}); err != nil {
		return "", err
	}
	target, err := url.Parse(cfg.RedirectURI)
	if err != nil {
		return "", err
	}
	q := target.Query()
	q.Set("code", code)
	q.Set("state", req.State)
	target.RawQuery = q.Encode()
	return target.String(), nil
}

// ---- 令牌 ----

type TokenRequest struct {
	GrantType    string
	Code         string
	RedirectURI  string
	CodeVerifier string
	ClientID     string
	ClientSecret string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (s *Service) Exchange(ctx context.Context, req TokenRequest) (*TokenResponse, error) {
	cfg, err := s.authenticateClient(ctx, req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}
	if req.GrantType != "authorization_code" || req.Code == "" || req.CodeVerifier == "" {
		return nil, ErrInvalidRequest
	}
	now := s.now()
	code, err := s.repo.consumeCode(ctx, hashToken(req.Code), now)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrInvalidGrant
	}
	if err != nil {
		return nil, err
	}
	if code.ConsumedAt != nil {
		// 授权码被重放：作废由它签发的令牌
		_ = s.repo.revokeTokensByCode(ctx, code.ID, now)
		return nil, ErrInvalidGrant
	}
	if code.ClientID != cfg.ClientID || code.RedirectURI != req.RedirectURI || now.After(code.ExpiresAt) {
		return nil, ErrInvalidGrant
	}
	if subtle.ConstantTimeCompare([]byte(pkceChallenge(req.CodeVerifier)), []byte(code.CodeChallenge)) != 1 {
		return nil, ErrInvalidGrant
	}
	if err := s.requireActiveUser(ctx, code.UserID); err != nil {
		return nil, ErrInvalidGrant
	}
	token, err := randomToken()
	if err != nil {
		return nil, err
	}
	if err := s.repo.createToken(ctx, hashToken(token), accessToken{
		ClientID:  cfg.ClientID,
		UserID:    code.UserID,
		ExpiresAt: now.Add(tokenTTL),
	}, code.ID); err != nil {
		return nil, err
	}
	return &TokenResponse{AccessToken: token, TokenType: "Bearer", ExpiresIn: int(tokenTTL.Seconds())}, nil
}

// ---- 用户信息 ----

type UserInfo struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Picture  string `json:"picture,omitempty"`
}

func (s *Service) UserInfo(ctx context.Context, bearer string) (*UserInfo, error) {
	token, err := s.validToken(ctx, bearer)
	if err != nil {
		return nil, err
	}
	user, err := s.users.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if user.Status != service.StatusActive {
		return nil, ErrUserInactive
	}
	name := strings.TrimSpace(user.Username)
	if name == "" {
		name = strings.Split(user.Email, "@")[0]
	}
	return &UserInfo{
		Sub:      strconv.FormatInt(user.ID, 10),
		Username: name,
		Name:     name,
		Email:    user.Email,
		Picture:  publicAvatar(user.AvatarURL),
	}, nil
}

// StudioKey 给客户端后端取该用户的专用密钥。需要客户端密钥 + 访问令牌，
// 浏览器拿不到客户端密钥，所以这个接口不会被页面直接调用。
func (s *Service) StudioKey(ctx context.Context, clientID, clientSecret, bearer string) (string, error) {
	cfg, err := s.authenticateClient(ctx, clientID, clientSecret)
	if err != nil {
		return "", err
	}
	token, err := s.validToken(ctx, bearer)
	if err != nil {
		return "", err
	}
	if token.ClientID != cfg.ClientID {
		return "", ErrInvalidToken
	}
	key, err := s.ensureStudioKey(ctx, cfg, token.UserID)
	if err != nil {
		return "", err
	}
	return key.Key, nil
}

// ensureStudioKey 复用已有的专用密钥；没有或已被删除时新建。
// 用户停用了它就原样报错——停用是用户的明确意图，不能悄悄换一把新的。
func (s *Service) ensureStudioKey(ctx context.Context, cfg Config, userID int64) (*service.APIKey, error) {
	id, err := s.repo.getClientKeyID(ctx, userID, cfg.ClientID)
	if err == nil {
		key, getErr := s.keys.GetByID(ctx, id)
		if getErr == nil && key.UserID == userID {
			if !key.IsActive() {
				return nil, ErrKeyUnavailable
			}
			return key, nil
		}
		if getErr != nil && !errors.Is(getErr, service.ErrAPIKeyNotFound) {
			return nil, getErr
		}
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	// 自动模式：请求时按用户当前权限选组（订阅优先、倍率低优先），管理员不用维护分组。
	key, err := s.keys.Create(ctx, userID, service.CreateAPIKeyRequest{
		Name:      studioKeyName,
		GroupMode: service.APIKeyGroupModeAuto,
		Quota:     cfg.KeyQuota,
	})
	if err != nil {
		return nil, fmt.Errorf("create studio key: %w", err)
	}
	if err := s.repo.setClientKeyID(ctx, userID, cfg.ClientID, key.ID); err != nil {
		return nil, err
	}
	return key, nil
}

// IsStudioKey 告诉「API 密钥」页哪把是工作台专用的。
func (s *Service) StudioKeyIDs(ctx context.Context, userID int64) (map[int64]bool, error) {
	cfg, err := s.config.load(ctx)
	if err != nil || cfg.ClientID == "" {
		return map[int64]bool{}, err
	}
	id, err := s.repo.getClientKeyID(ctx, userID, cfg.ClientID)
	if errors.Is(err, ErrNotFound) {
		return map[int64]bool{}, nil
	}
	if err != nil {
		return nil, err
	}
	return map[int64]bool{id: true}, nil
}

// Cleanup 清理过期的授权码和令牌；由后台定时调用。
func (s *Service) Cleanup(ctx context.Context) error {
	return s.repo.deleteExpired(ctx, s.now().Add(-time.Hour))
}

// ---- 内部 ----

func (s *Service) authenticateClient(ctx context.Context, clientID, clientSecret string) (Config, error) {
	cfg, err := s.config.load(ctx)
	if err != nil {
		return Config{}, err
	}
	if !cfg.Enabled {
		return Config{}, ErrDisabled
	}
	idOK := subtle.ConstantTimeCompare([]byte(clientID), []byte(cfg.ClientID)) == 1
	secretOK := subtle.ConstantTimeCompare([]byte(clientSecret), []byte(cfg.ClientSecret)) == 1
	if !idOK || !secretOK || cfg.ClientSecret == "" {
		return Config{}, ErrInvalidClient
	}
	return cfg, nil
}

func (s *Service) validToken(ctx context.Context, bearer string) (*accessToken, error) {
	bearer = strings.TrimSpace(bearer)
	if bearer == "" {
		return nil, ErrInvalidToken
	}
	token, err := s.repo.getToken(ctx, hashToken(bearer))
	if err != nil {
		return nil, ErrInvalidToken
	}
	if token.RevokedAt != nil || s.now().After(token.ExpiresAt) {
		return nil, ErrInvalidToken
	}
	return token, nil
}

func (s *Service) requireActiveUser(ctx context.Context, userID int64) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Status != service.StatusActive {
		return ErrUserInactive
	}
	return nil
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// validChallenge：S256 的 challenge 是 32 字节摘要的 base64url，固定 43 个字符
func validChallenge(v string) bool {
	if len(v) != 43 {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(v)
	return err == nil
}

// publicAvatar 只透出外部可访问的头像地址；内嵌的 data: 头像不外传
func publicAvatar(v string) string {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "https://") {
		return v
	}
	return ""
}
