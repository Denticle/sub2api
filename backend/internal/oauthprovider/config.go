// Package oauthprovider 让 forapi 充当一个最小的 OAuth 2.0 授权服务器。
//
// 只服务一个预先配置的客户端（视频工作台），只支持授权码 + PKCE S256，
// 不做动态客户端注册、不签发 refresh token。访问令牌只在客户端登录回调的
// 那一瞬间用来拉用户信息和取专用密钥，之后客户端用自己的会话。
package oauthprovider

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

const settingKey = "oauth_provider_studio"

// Config 是管理员在后台配置的唯一客户端。
type Config struct {
	Enabled      bool   `json:"enabled"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"` // 存储时加密
	RedirectURI  string `json:"redirect_uri"`  // 精确匹配
	// AppName 显示在授权确认页
	AppName string `json:"app_name"`
	// AppURL 是工作台首页，前端顶部菜单用它拼登录入口
	AppURL string `json:"app_url"`
	// KeyGroupID 已弃用：专用密钥改为自动模式，按用户权限选组。
	// 字段保留以兼容已保存的配置，不再参与建密钥与启用校验。
	KeyGroupID *int64 `json:"key_group_id"`
	// KeyQuota 专用密钥的额度上限（0 = 不限）
	KeyQuota float64 `json:"key_quota"`
}

// PublicConfig 是给前端看的部分，不含任何密钥。
type PublicConfig struct {
	Enabled bool   `json:"enabled"`
	AppName string `json:"app_name"`
	AppURL  string `json:"app_url"`
}

var (
	ErrDisabled        = errors.New("oauth provider is disabled")
	ErrInvalidClient   = errors.New("invalid client")
	ErrInvalidRedirect = errors.New("redirect_uri does not match")
)

type configStore struct {
	settings  service.SettingRepository
	encryptor service.SecretEncryptor
}

func (s *configStore) load(ctx context.Context) (Config, error) {
	raw, err := s.settings.GetValue(ctx, settingKey)
	if err != nil || strings.TrimSpace(raw) == "" {
		// 没配置过：视为关闭
		return Config{}, nil
	}
	var cfg Config
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return Config{}, err
	}
	if cfg.ClientSecret != "" && s.encryptor != nil {
		plain, err := s.encryptor.Decrypt(cfg.ClientSecret)
		if err != nil {
			return Config{}, err
		}
		cfg.ClientSecret = plain
	}
	return cfg, nil
}

func (s *configStore) save(ctx context.Context, next Config) error {
	if s.encryptor != nil && next.ClientSecret != "" {
		enc, err := s.encryptor.Encrypt(next.ClientSecret)
		if err != nil {
			return err
		}
		next.ClientSecret = enc
	}
	raw, err := json.Marshal(next)
	if err != nil {
		return err
	}
	return s.settings.Set(ctx, settingKey, string(raw))
}

func validateConfig(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}
	if strings.TrimSpace(cfg.ClientID) == "" || strings.TrimSpace(cfg.ClientSecret) == "" {
		return errors.New("启用前请填写 Client ID 和 Client Secret")
	}
	if len(cfg.ClientSecret) < 32 {
		return errors.New("Client Secret 至少 32 个字符")
	}
	if cfg.KeyQuota < 0 {
		return errors.New("专用密钥额度不能小于 0")
	}
	for name, raw := range map[string]string{"回调地址": cfg.RedirectURI, "工作台地址": cfg.AppURL} {
		u, err := url.Parse(strings.TrimSpace(raw))
		if err != nil || u.Host == "" || (u.Scheme != "https" && !(u.Scheme == "http" && isLoopback(u.Hostname()))) {
			return errors.New(name + "必须是 HTTPS 地址")
		}
	}
	return nil
}

func isLoopback(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func (c Config) public() PublicConfig {
	return PublicConfig{Enabled: c.Enabled, AppName: c.AppName, AppURL: strings.TrimRight(c.AppURL, "/")}
}
