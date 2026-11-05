// ForAPI 密钥桥接（forapi 集成层 overlay，不属于上游）。
//
// 渠道密钥填写为 ForAPIUserKeySentinel 的系统渠道，不使用渠道上保存的密钥，
// 而是换成「发起这次请求的用户」登录时从 forapi 取得的专用密钥。
// 找不到用户密钥时直接失败，绝不回退到其他密钥，避免费用记到别人名下。
package app

import (
	"encoding/json"
	"errors"
	"strings"

	"yingce/backend/internal/auth"
	"yingce/backend/internal/kernel"
	"yingce/backend/internal/model"
)

// ForAPIUserKeySentinel 写进渠道「API Key」字段的占位值；它不是有效密钥，
// 万一被原样发往上游也只会被拒绝。
const ForAPIUserKeySentinel = "forapi-user-key"

func isForAPIBridgeChannel(channel *model.ModelChannel) bool {
	return channel != nil && strings.TrimSpace(channel.APIKey) == ForAPIUserKeySentinel
}

func errForAPIKeyMissing() error {
	return kernel.Forbidden("没有找到你的 ForAPI 密钥，请退出后从 ForAPI 顶部菜单「视频工作台」重新进入")
}

// ForAPIUserKey 读取用户登录时保存的 forapi 专用密钥。
func (s *Service) ForAPIUserKey(userID string) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", errForAPIKeyMissing()
	}
	setting, err := s.repo.SystemSettingOptional(auth.ForAPIBridgeKeySettingPrefix + userID)
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", errForAPIKeyMissing()
	}
	var value struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.Unmarshal([]byte(setting.ValueJSON), &value); err != nil {
		return "", errors.New("ForAPI 密钥记录损坏")
	}
	key, err := s.decryptSettingSecret(value.APIKey)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(key) == "" {
		return "", errForAPIKeyMissing()
	}
	return key, nil
}

// ChannelAPIKeyForUser 返回这次请求应使用的渠道密钥：桥接渠道用用户自己的，其余照旧。
func (s *Service) ChannelAPIKeyForUser(channel *model.ModelChannel, userID string) (string, error) {
	if !isForAPIBridgeChannel(channel) {
		return channel.APIKey, nil
	}
	return s.ForAPIUserKey(userID)
}
