// ForAPI 模型目录同步（forapi 集成层 overlay，不属于上游）。
//
// 工作台的「ForAPI」桥接渠道不再由管理员手工加模型：后台定时从 forapi 的
// /oauth/studio-models 拉取「工作台」密钥分组可用的模型，自动建成渠道模型。
// forapi 那边删掉或移出分组的模型，这里会被停用（不删除，保留历史任务引用）。
// 计费由 forapi 按用户自己的密钥结算，所以这里的积分价格一律为 0。
package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"yingce/backend/internal/auth"
	"yingce/backend/internal/model"
	"yingce/backend/internal/outbound"
)

const forAPIModelSyncInterval = 10 * time.Minute

// ForAPIModel 与 forapi /oauth/studio-models 的返回项一致。
type ForAPIModel struct {
	ID         string `json:"id"`
	Platform   string `json:"platform"`
	Capability string `json:"capability"`
}

// ForAPIModelSyncResult 汇总一次同步做了什么，便于日志和测试核对。
type ForAPIModelSyncResult struct {
	Created  []string
	Enabled  []string
	Disabled []string
}

// RunForAPIModelSync 启动后立即同步一次，之后每 10 分钟同步一次，直到 ctx 结束。
func (s *Service) RunForAPIModelSync(ctx context.Context) {
	ticker := time.NewTicker(forAPIModelSyncInterval)
	defer ticker.Stop()
	for {
		if result, err := s.SyncForAPIModels(ctx); err != nil {
			slog.Warn("ForAPI 模型同步失败", "error", err)
		} else if len(result.Created)+len(result.Enabled)+len(result.Disabled) > 0 {
			slog.Info("ForAPI 模型已同步", "created", result.Created, "enabled", result.Enabled, "disabled", result.Disabled)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// SyncForAPIModels 拉取 forapi 模型目录并对齐桥接渠道的模型。
func (s *Service) SyncForAPIModels(ctx context.Context) (*ForAPIModelSyncResult, error) {
	tokenURL, clientID, clientSecret, err := s.authDomain().ForAPIClient()
	if err != nil {
		return &ForAPIModelSyncResult{}, nil // 桥接未启用时静默跳过
	}
	models, err := fetchForAPIModels(ctx, tokenURL, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return s.applyForAPIModels(models)
}

func fetchForAPIModels(ctx context.Context, tokenURL, clientID, clientSecret string) ([]ForAPIModel, error) {
	target, err := auth.ForAPIModelsURL(tokenURL)
	if err != nil {
		return nil, err
	}
	if _, err := outbound.ValidateOutboundURL(target); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Accept", "application/json")
	resp, err := outbound.OutboundHTTPClient(20 * time.Second).Do(req)
	if err != nil {
		return nil, fmt.Errorf("拉取 ForAPI 模型目录失败：%w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("拉取 ForAPI 模型目录失败：HTTP %d", resp.StatusCode)
	}
	var payload struct {
		Models []ForAPIModel `json:"models"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, errors.New("ForAPI 模型目录格式无效")
	}
	return payload.Models, nil
}

// forAPIModelProtocol 把 forapi 的能力类别映射成工作台的请求协议。
// 文本走 OpenAI 兼容 chat/completions；图片走 OpenAI images；
// 视频按 forapi 实际支持的两条入口：Seedance 走方舟任务接口，其余走 xAI 视频接口。
func forAPIModelProtocol(m ForAPIModel) (string, model.ChannelInterfaceType) {
	switch strings.TrimSpace(m.Capability) {
	case "image":
		return "image", model.ChannelInterfaceOpenAIImage
	case "video":
		if strings.Contains(strings.ToLower(m.ID), "seedance") {
			return "video", model.ChannelInterfaceVolcengineArkVideo
		}
		return "video", model.ChannelInterfaceXAIVideo
	default:
		return "text", model.ChannelInterfaceChatCompletion
	}
}

func (s *Service) forAPIBridgeChannel() (*model.ModelChannel, error) {
	channels, err := s.repo.SystemChannels(true)
	if err != nil {
		return nil, err
	}
	for index := range channels {
		if err := s.decryptSystemChannelSecrets(&channels[index]); err != nil {
			continue
		}
		if isForAPIBridgeChannel(&channels[index]) {
			return &channels[index], nil
		}
	}
	return nil, errors.New("没有找到 ForAPI 桥接渠道")
}

// forAPISyncActor 用一个活跃管理员的身份写入，保证经过与后台手工保存相同的校验和审计。
func (s *Service) forAPISyncActor() (*model.User, error) {
	users, _, err := s.repo.AdminUsers("", model.UserRoleAdmin, model.UserStatusActive, "", 1, 0)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("没有可用的管理员账号执行模型同步")
	}
	return &users[0], nil
}

func forAPIChannelModelRequest(m ForAPIModel) ChannelModelRequest {
	capability, protocol := forAPIModelProtocol(m)
	enabled := true
	currency := model.PriceCurrencyCredits
	defaults := DefaultModelCapabilityConfigForModel(string(protocol), m.ID)
	config := &ModelCapabilityConfig{Version: 1}
	billing := "fixed_request"
	switch capability {
	case "text":
		config.Text = defaults.Text
		billing = "token"
	case "image":
		config.Image = defaults.Image
	case "video":
		config.Video = defaults.Video
	}
	return ChannelModelRequest{
		ModelKey: m.ID, ProviderModelKey: m.ID, DisplayName: m.ID,
		Capability: capability, Protocol: string(protocol),
		BillingMode: billing, PriceConfigured: true, PriceCurrency: &currency,
		Enabled: &enabled, CapabilityConfig: config,
	}
}

func (s *Service) applyForAPIModels(models []ForAPIModel) (*ForAPIModelSyncResult, error) {
	result := &ForAPIModelSyncResult{}
	channel, err := s.forAPIBridgeChannel()
	if err != nil {
		return nil, err
	}
	actor, err := s.forAPISyncActor()
	if err != nil {
		return nil, err
	}
	wanted := make(map[string]bool, len(models))
	var firstErr error
	for _, m := range models {
		id := strings.TrimPrefix(strings.TrimSpace(m.ID), "models/")
		if id == "" || wanted[id] {
			continue
		}
		m.ID = id
		wanted[id] = true
		existing, lookupErr := s.repo.ChannelModelByKeyIncludingDisabled(channel.ID, id)
		if lookupErr == nil && existing != nil {
			if !existing.Enabled {
				if _, err := s.SetAdminChannelModelEnabled(actor, channel.ID, existing.ID, true); err != nil {
					firstErr = errors.Join(firstErr, fmt.Errorf("%s：%w", id, err))
					continue
				}
				result.Enabled = append(result.Enabled, id)
			}
			continue
		}
		if _, err := s.SaveAdminChannelModel(actor, channel.ID, "", forAPIChannelModelRequest(m)); err != nil {
			// 单个模型失败不影响其他模型；错误汇总后返回给日志。
			firstErr = errors.Join(firstErr, fmt.Errorf("%s：%w", id, err))
			continue
		}
		result.Created = append(result.Created, id)
	}
	current, err := s.repo.ChannelModels(channel.ID, true)
	if err != nil {
		return result, errors.Join(firstErr, err)
	}
	for _, item := range current {
		if item.Enabled && !wanted[item.ModelKey] {
			if _, err := s.SetAdminChannelModelEnabled(actor, channel.ID, item.ID, false); err != nil {
				firstErr = errors.Join(firstErr, fmt.Errorf("%s：%w", item.ModelKey, err))
				continue
			}
			result.Disabled = append(result.Disabled, item.ModelKey)
		}
	}
	return result, firstErr
}
