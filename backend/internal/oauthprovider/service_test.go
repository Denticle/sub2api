package oauthprovider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// ---- 测试替身 ----

type memSettings struct {
	mu   sync.Mutex
	data map[string]string
}

func (m *memSettings) Get(ctx context.Context, key string) (*service.Setting, error) { return nil, nil }
func (m *memSettings) GetValue(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[key], nil
}
func (m *memSettings) Set(ctx context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = map[string]string{}
	}
	m.data[key] = value
	return nil
}
func (m *memSettings) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	return nil, nil
}
func (m *memSettings) SetMultiple(ctx context.Context, settings map[string]string) error { return nil }
func (m *memSettings) GetAll(ctx context.Context) (map[string]string, error)             { return nil, nil }
func (m *memSettings) Delete(ctx context.Context, key string) error                      { return nil }

// prefixEncryptor 让测试能看出「存进去的是密文」
type prefixEncryptor struct{}

func (prefixEncryptor) Encrypt(p string) (string, error) { return "enc:" + p, nil }
func (prefixEncryptor) Decrypt(c string) (string, error) {
	if !strings.HasPrefix(c, "enc:") {
		return "", errors.New("not encrypted")
	}
	return strings.TrimPrefix(c, "enc:"), nil
}

type memRepo struct {
	mu       sync.Mutex
	nextID   int64
	codes    map[string]*authorizationCode
	tokens   map[string]*accessToken
	tokenOf  map[int64][]string
	consents map[string]bool
	keys     map[string]int64
}

func newMemRepo() *memRepo {
	return &memRepo{codes: map[string]*authorizationCode{}, tokens: map[string]*accessToken{}, tokenOf: map[int64][]string{}, consents: map[string]bool{}, keys: map[string]int64{}}
}

func ck(userID int64, clientID string) string { return clientID + "/" + string(rune(userID)) }

func (r *memRepo) createCode(ctx context.Context, h string, c authorizationCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	c.ID = r.nextID
	r.codes[h] = &c
	return nil
}
func (r *memRepo) consumeCode(ctx context.Context, h string, now time.Time) (*authorizationCode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.codes[h]
	if !ok {
		return nil, ErrNotFound
	}
	snapshot := *c
	if c.ConsumedAt == nil {
		t := now
		c.ConsumedAt = &t
	}
	return &snapshot, nil
}
func (r *memRepo) createToken(ctx context.Context, h string, t accessToken, codeID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[h] = &t
	r.tokenOf[codeID] = append(r.tokenOf[codeID], h)
	return nil
}
func (r *memRepo) revokeTokensByCode(ctx context.Context, codeID int64, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, h := range r.tokenOf[codeID] {
		t := now
		r.tokens[h].RevokedAt = &t
	}
	return nil
}
func (r *memRepo) getToken(ctx context.Context, h string) (*accessToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tokens[h]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *t
	return &cp, nil
}
func (r *memRepo) hasConsent(ctx context.Context, u int64, c string) (bool, error) {
	return r.consents[ck(u, c)], nil
}
func (r *memRepo) saveConsent(ctx context.Context, u int64, c string) error {
	r.consents[ck(u, c)] = true
	return nil
}
func (r *memRepo) getClientKeyID(ctx context.Context, u int64, c string) (int64, error) {
	id, ok := r.keys[ck(u, c)]
	if !ok {
		return 0, ErrNotFound
	}
	return id, nil
}
func (r *memRepo) setClientKeyID(ctx context.Context, u int64, c string, id int64) error {
	r.keys[ck(u, c)] = id
	return nil
}
func (r *memRepo) deleteExpired(ctx context.Context, before time.Time) error { return nil }

type memUsers map[int64]*service.User

func (m memUsers) GetByID(ctx context.Context, id int64) (*service.User, error) {
	u, ok := m[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

type memKeys struct {
	nextID  int64
	byID    map[int64]*service.APIKey
	created int
}

func (m *memKeys) Create(ctx context.Context, userID int64, req service.CreateAPIKeyRequest) (*service.APIKey, error) {
	m.nextID++
	m.created++
	k := &service.APIKey{ID: m.nextID, UserID: userID, Key: "sk-studio-" + string(rune('a'+m.nextID)), Name: req.Name, GroupID: req.GroupID, GroupMode: req.GroupMode, Quota: req.Quota, Status: service.StatusActive}
	m.byID[k.ID] = k
	return k, nil
}
func (m *memKeys) GetByID(ctx context.Context, id int64) (*service.APIKey, error) {
	k, ok := m.byID[id]
	if !ok {
		return nil, service.ErrAPIKeyNotFound
	}
	return k, nil
}

// ---- 夹具 ----

const (
	testClientID     = "video-studio"
	testClientSecret = "0123456789abcdef0123456789abcdef"
	testRedirect     = "https://video.forapi.cloud/api/auth/linuxdo/callback"
	testVerifier     = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
)

type fixture struct {
	svc   *Service
	repo  *memRepo
	keys  *memKeys
	users memUsers
	clock time.Time
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{
		repo:  newMemRepo(),
		keys:  &memKeys{byID: map[int64]*service.APIKey{}},
		users: memUsers{1: {ID: 1, Username: "alice", Email: "alice@example.com", Status: service.StatusActive}, 2: {ID: 2, Username: "bob", Email: "bob@example.com", Status: service.StatusActive}},
		clock: time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC),
	}
	f.svc = &Service{
		config: configStore{settings: &memSettings{}, encryptor: prefixEncryptor{}},
		repo:   f.repo,
		users:  f.users,
		keys:   f.keys,
		now:    func() time.Time { return f.clock },
	}
	if err := f.svc.SaveConfig(context.Background(), Config{
		Enabled: true, ClientID: testClientID, ClientSecret: testClientSecret,
		RedirectURI: testRedirect, AppName: "视频工作台", AppURL: "https://video.forapi.cloud",
	}); err != nil {
		t.Fatalf("save config: %v", err)
	}
	return f
}

func challengeOf(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (f *fixture) authorizeReq() AuthorizeRequest {
	return AuthorizeRequest{ClientID: testClientID, RedirectURI: testRedirect, ResponseType: "code", State: "st-123", CodeChallenge: challengeOf(testVerifier), CodeChallengeMethod: "S256"}
}

func (f *fixture) approve(t *testing.T, userID int64) string {
	t.Helper()
	target, err := f.svc.Approve(context.Background(), userID, f.authorizeReq())
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	u, _ := url.Parse(target)
	if u.Query().Get("state") != "st-123" {
		t.Fatalf("state not echoed: %s", target)
	}
	return u.Query().Get("code")
}

func (f *fixture) exchange(code, verifier string) (*TokenResponse, error) {
	return f.svc.Exchange(context.Background(), TokenRequest{GrantType: "authorization_code", Code: code, RedirectURI: testRedirect, CodeVerifier: verifier, ClientID: testClientID, ClientSecret: testClientSecret})
}

// ---- 用例 ----

func TestFullFlow(t *testing.T) {
	f := newFixture(t)
	code := f.approve(t, 1)
	tok, err := f.exchange(code, testVerifier)
	if err != nil {
		t.Fatalf("exchange: %v", err)
	}
	if tok.TokenType != "Bearer" || tok.ExpiresIn != 600 || tok.AccessToken == "" {
		t.Fatalf("bad token response: %+v", tok)
	}
	info, err := f.svc.UserInfo(context.Background(), tok.AccessToken)
	if err != nil {
		t.Fatalf("userinfo: %v", err)
	}
	if info.Sub != "1" || info.Username != "alice" || info.Email != "alice@example.com" {
		t.Fatalf("bad userinfo: %+v", info)
	}
	key, err := f.svc.StudioKey(context.Background(), testClientID, testClientSecret, tok.AccessToken)
	if err != nil || key == "" {
		t.Fatalf("studio key: %q %v", key, err)
	}
}

func TestClientSecretStoredEncrypted(t *testing.T) {
	f := newFixture(t)
	raw, _ := f.svc.config.settings.GetValue(context.Background(), settingKey)
	var stored Config
	_ = json.Unmarshal([]byte(raw), &stored)
	if stored.ClientSecret != "enc:"+testClientSecret {
		t.Fatalf("secret not encrypted at rest: %q", stored.ClientSecret)
	}
	cfg, has, _ := f.svc.AdminConfig(context.Background())
	if !has || cfg.ClientSecret != "" {
		t.Fatalf("admin config must not leak secret: %+v has=%v", cfg, has)
	}
}

func TestPKCEMismatchRejected(t *testing.T) {
	f := newFixture(t)
	code := f.approve(t, 1)
	if _, err := f.exchange(code, "wrong-verifier-wrong-verifier-wrong-verifier"); !errors.Is(err, ErrInvalidGrant) {
		t.Fatalf("want invalid_grant, got %v", err)
	}
}

func TestCodeReplayRevokesIssuedToken(t *testing.T) {
	f := newFixture(t)
	code := f.approve(t, 1)
	tok, err := f.exchange(code, testVerifier)
	if err != nil {
		t.Fatalf("first exchange: %v", err)
	}
	if _, err := f.exchange(code, testVerifier); !errors.Is(err, ErrInvalidGrant) {
		t.Fatalf("replay should fail, got %v", err)
	}
	if _, err := f.svc.UserInfo(context.Background(), tok.AccessToken); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("token issued from replayed code must be revoked, got %v", err)
	}
}

func TestExpiredCodeRejected(t *testing.T) {
	f := newFixture(t)
	code := f.approve(t, 1)
	f.clock = f.clock.Add(61 * time.Second)
	if _, err := f.exchange(code, testVerifier); !errors.Is(err, ErrInvalidGrant) {
		t.Fatalf("want invalid_grant, got %v", err)
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	f := newFixture(t)
	tok, _ := f.exchange(f.approve(t, 1), testVerifier)
	f.clock = f.clock.Add(11 * time.Minute)
	if _, err := f.svc.UserInfo(context.Background(), tok.AccessToken); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("want invalid_token, got %v", err)
	}
}

func TestRedirectMustMatchExactly(t *testing.T) {
	f := newFixture(t)
	for _, bad := range []string{testRedirect + "/", "https://evil.example/cb", testRedirect + "?x=1"} {
		req := f.authorizeReq()
		req.RedirectURI = bad
		if _, err := f.svc.Approve(context.Background(), 1, req); !errors.Is(err, ErrInvalidRedirect) {
			t.Fatalf("%s: want invalid redirect, got %v", bad, err)
		}
	}
	// 交换时的 redirect_uri 也必须一致
	code := f.approve(t, 1)
	_, err := f.svc.Exchange(context.Background(), TokenRequest{GrantType: "authorization_code", Code: code, RedirectURI: "https://evil.example/cb", CodeVerifier: testVerifier, ClientID: testClientID, ClientSecret: testClientSecret})
	if !errors.Is(err, ErrInvalidGrant) {
		t.Fatalf("want invalid_grant, got %v", err)
	}
}

func TestPlainPKCEAndMissingStateRejected(t *testing.T) {
	f := newFixture(t)
	req := f.authorizeReq()
	req.CodeChallengeMethod = "plain"
	if _, err := f.svc.Approve(context.Background(), 1, req); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("plain pkce: want invalid_request, got %v", err)
	}
	req = f.authorizeReq()
	req.State = ""
	if _, err := f.svc.Approve(context.Background(), 1, req); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("missing state: want invalid_request, got %v", err)
	}
}

func TestWrongClientSecretRejected(t *testing.T) {
	f := newFixture(t)
	code := f.approve(t, 1)
	_, err := f.svc.Exchange(context.Background(), TokenRequest{GrantType: "authorization_code", Code: code, RedirectURI: testRedirect, CodeVerifier: testVerifier, ClientID: testClientID, ClientSecret: strings.Repeat("x", 32)})
	if !errors.Is(err, ErrInvalidClient) {
		t.Fatalf("want invalid_client, got %v", err)
	}
	tok, _ := f.exchange(code, testVerifier)
	_ = tok
	if _, err := f.svc.StudioKey(context.Background(), testClientID, "nope", "whatever"); !errors.Is(err, ErrInvalidClient) {
		t.Fatalf("studio key without client secret: want invalid_client, got %v", err)
	}
}

func TestDisabledRejectsEverything(t *testing.T) {
	f := newFixture(t)
	code := f.approve(t, 1)
	_ = f.svc.SaveConfig(context.Background(), Config{Enabled: false, ClientID: testClientID, RedirectURI: testRedirect, AppURL: "https://video.forapi.cloud"})
	if _, err := f.svc.Approve(context.Background(), 1, f.authorizeReq()); !errors.Is(err, ErrDisabled) {
		t.Fatalf("approve: want disabled, got %v", err)
	}
	if _, err := f.exchange(code, testVerifier); !errors.Is(err, ErrDisabled) {
		t.Fatalf("exchange: want disabled, got %v", err)
	}
}

func TestInactiveUserCannotAuthorize(t *testing.T) {
	f := newFixture(t)
	f.users[1].Status = "disabled"
	if _, err := f.svc.Approve(context.Background(), 1, f.authorizeReq()); !errors.Is(err, ErrUserInactive) {
		t.Fatalf("want inactive, got %v", err)
	}
}

func TestStudioKeyReusedNotRecreated(t *testing.T) {
	f := newFixture(t)
	tok1, _ := f.exchange(f.approve(t, 1), testVerifier)
	k1, _ := f.svc.StudioKey(context.Background(), testClientID, testClientSecret, tok1.AccessToken)
	tok2, _ := f.exchange(f.approve(t, 1), testVerifier)
	k2, _ := f.svc.StudioKey(context.Background(), testClientID, testClientSecret, tok2.AccessToken)
	if k1 != k2 || f.keys.created != 1 {
		t.Fatalf("key should be reused: %q %q created=%d", k1, k2, f.keys.created)
	}
	if k := f.keys.byID[1]; k.Name != studioKeyName || k.GroupMode != service.APIKeyGroupModeAuto || k.GroupID != nil {
		t.Fatalf("studio key must be an auto-mode key without a fixed group: %+v", k)
	}
}

func TestStudioKeyDisabledIsNotSilentlyReplaced(t *testing.T) {
	f := newFixture(t)
	f.approve(t, 1)
	f.keys.byID[1].Status = "inactive"
	if _, err := f.svc.Approve(context.Background(), 1, f.authorizeReq()); !errors.Is(err, ErrKeyUnavailable) {
		t.Fatalf("want key unavailable, got %v", err)
	}
	if f.keys.created != 1 {
		t.Fatalf("must not create a replacement key, created=%d", f.keys.created)
	}
}

func TestStudioKeyRecreatedAfterDeletion(t *testing.T) {
	f := newFixture(t)
	f.approve(t, 1)
	delete(f.keys.byID, 1)
	f.approve(t, 1)
	if f.keys.created != 2 {
		t.Fatalf("deleted key should be recreated, created=%d", f.keys.created)
	}
}

func TestUsersDoNotShareKeys(t *testing.T) {
	f := newFixture(t)
	tokA, _ := f.exchange(f.approve(t, 1), testVerifier)
	tokB, _ := f.exchange(f.approve(t, 2), testVerifier)
	ka, _ := f.svc.StudioKey(context.Background(), testClientID, testClientSecret, tokA.AccessToken)
	kb, _ := f.svc.StudioKey(context.Background(), testClientID, testClientSecret, tokB.AccessToken)
	if ka == kb {
		t.Fatalf("two users got the same key %q", ka)
	}
	ia, _ := f.svc.UserInfo(context.Background(), tokA.AccessToken)
	ib, _ := f.svc.UserInfo(context.Background(), tokB.AccessToken)
	if ia.Sub == ib.Sub {
		t.Fatalf("subjects collide")
	}
}

func TestConsentRemembered(t *testing.T) {
	f := newFixture(t)
	check, _ := f.svc.CheckAuthorize(context.Background(), 1, f.authorizeReq())
	if check.Consented {
		t.Fatal("first visit should need consent")
	}
	f.approve(t, 1)
	check, _ = f.svc.CheckAuthorize(context.Background(), 1, f.authorizeReq())
	if !check.Consented {
		t.Fatal("consent should be remembered")
	}
}

func TestConfigValidation(t *testing.T) {
	f := newFixture(t)
	base := Config{Enabled: true, ClientID: "c", ClientSecret: testClientSecret, RedirectURI: testRedirect, AppURL: "https://video.forapi.cloud"}
	cases := map[string]func(c *Config){
		"short secret":   func(c *Config) { c.ClientSecret = "short" },
		"http redirect":  func(c *Config) { c.RedirectURI = "http://video.forapi.cloud/cb" },
		"negative quota": func(c *Config) { c.KeyQuota = -1 },
	}
	for name, mut := range cases {
		c := base
		mut(&c)
		if err := f.svc.SaveConfig(context.Background(), c); err == nil {
			t.Fatalf("%s: expected validation error", name)
		}
	}
	// 专用密钥改为自动模式后，不再需要选分组
	if err := f.svc.SaveConfig(context.Background(), base); err != nil {
		t.Fatalf("enabling without a key group should be allowed: %v", err)
	}
	// 留空 secret 沿用旧值
	c := base
	c.ClientSecret = ""
	if err := f.svc.SaveConfig(context.Background(), c); err != nil {
		t.Fatalf("empty secret should keep existing: %v", err)
	}
}
