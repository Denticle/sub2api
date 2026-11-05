// ForAPI 登录桥接（forapi 集成层 overlay，不属于上游）。
//
// 用户通过「Linux.do 登录」入口（实际指向 forapi 的 OAuth 服务）登录时，
// 顺带用同一个 access token 向 forapi 换取该用户的「工作台」专用 API 密钥，
// 加密存进 system_settings。生成任务只用这把密钥调用 forapi，费用记在用户自己名下。
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"yingce/backend/internal/kernel"
	"yingce/backend/internal/model"
	"yingce/backend/internal/outbound"
)

// ForAPIBridgeKeySettingPrefix 与 app 包读取时使用的前缀保持一致。
const ForAPIBridgeKeySettingPrefix = "forapi_bridge_key:"

// forAPIStudioKeyURL 从 token 地址推出同站的 studio-key 地址：
// https://forapi.cloud/oauth/token → https://forapi.cloud/oauth/studio-key
func forAPIStudioKeyURL(tokenURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(tokenURL))
	if err != nil || parsed.Host == "" {
		return "", errors.New("ForAPI Token 地址无效")
	}
	parsed.Path = strings.TrimSuffix(parsed.Path, "/token") + "/studio-key"
	parsed.RawQuery = ""
	return parsed.String(), nil
}

func fetchForAPIStudioKey(setting linuxDOSettingValue, accessToken string) (string, error) {
	target, err := forAPIStudioKeyURL(setting.TokenURL)
	if err != nil {
		return "", err
	}
	if _, err := outbound.ValidateOutboundURL(target); err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(setting.ClientID, setting.ClientSecret)
	req.Header.Set("X-Access-Token", accessToken)
	req.Header.Set("Accept", "application/json")
	resp, err := outbound.OutboundHTTPClient(20 * time.Second).Do(req)
	if err != nil {
		return "", fmt.Errorf("获取 ForAPI 密钥失败：%w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return "", err
	}
	var payload struct {
		APIKey string `json:"api_key"`
		Error  string `json:"error"`
	}
	_ = json.Unmarshal(body, &payload)
	if resp.StatusCode == http.StatusForbidden && payload.Error == "key_disabled" {
		return "", kernel.Forbidden("你在 ForAPI 的「工作台」密钥已停用，请到 ForAPI 的 API 密钥页启用后再登录")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || payload.APIKey == "" {
		return "", fmt.Errorf("获取 ForAPI 密钥失败：HTTP %d", resp.StatusCode)
	}
	return payload.APIKey, nil
}

// ForAPIClient 返回同步模型目录所需的 forapi 地址与客户端凭据（复用 Linux.do 登录配置）。
func (s *Service) ForAPIClient() (tokenURL, clientID, clientSecret string, err error) {
	_, setting, err := s.readLinuxDOSetting()
	if err != nil {
		return "", "", "", err
	}
	if !setting.Enabled || setting.ClientID == "" || setting.ClientSecret == "" {
		return "", "", "", errors.New("ForAPI 登录尚未启用")
	}
	return setting.TokenURL, setting.ClientID, setting.ClientSecret, nil
}

// ForAPIModelsURL 从 token 地址推出同站的 studio-models 地址。
func ForAPIModelsURL(tokenURL string) (string, error) {
	keyURL, err := forAPIStudioKeyURL(tokenURL)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(keyURL, "/studio-key") + "/studio-models", nil
}

// storeForAPIBridgeKey 拿不到密钥就让登录失败：没有密钥的会话无法生成，
// 也不能回退到渠道上的公共密钥。
func (s *Service) storeForAPIBridgeKey(setting linuxDOSettingValue, accessToken string, userID string) error {
	key, err := fetchForAPIStudioKey(setting, accessToken)
	if err != nil {
		return err
	}
	encrypted, err := s.host.EncryptSecret(key)
	if err != nil {
		return err
	}
	value, err := json.Marshal(map[string]string{"apiKey": encrypted})
	if err != nil {
		return err
	}
	now := time.Now()
	return s.repo.SaveSystemSetting(&model.SystemSetting{
		Key:       ForAPIBridgeKeySettingPrefix + userID,
		ValueJSON: string(value),
		UpdatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	})
}
