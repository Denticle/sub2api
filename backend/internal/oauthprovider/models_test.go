package oauthprovider

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type fakeGroups []service.PlazaGroup

func (f fakeGroups) ListGroups(context.Context) ([]service.PlazaGroup, error) { return f, nil }

func pricing(mode service.BillingMode) *service.ChannelModelPricing {
	return &service.ChannelModelPricing{BillingMode: mode}
}

func TestStudioModelsUnionsAllActiveGroupsWithCapabilities(t *testing.T) {
	f := newFixture(t)
	f.svc.groups = fakeGroups{
		{ID: 7, Models: []service.PlazaModel{
			{Name: "gpt-image-2", Platform: "openai", Pricing: pricing(service.BillingModeImage)},
			{Name: "claude-sonnet-5-5", Platform: "anthropic", Pricing: pricing(service.BillingModeToken)},
			{Name: "seedance-2-0", Platform: "openai"}, // 未定价：按名字判断
		}},
		{ID: 8, IsExclusive: true, Models: []service.PlazaModel{
			{Name: "grok-imagine-video", Platform: "grok", Pricing: pricing(service.BillingModeVideo)},
			{Name: "claude-sonnet-5-5", Platform: "anthropic"}, // 跨分组同名只保留一个
		}},
	}
	got, err := f.svc.StudioModels(context.Background(), testClientID, testClientSecret)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"claude-sonnet-5-5": "text", "gpt-image-2": "image", "grok-imagine-video": "video", "seedance-2-0": "video"}
	if len(got) != len(want) {
		t.Fatalf("got %d models: %+v", len(got), got)
	}
	for i, m := range got {
		if want[m.ID] != m.Capability {
			t.Errorf("%s capability = %q, want %q", m.ID, m.Capability, want[m.ID])
		}
		if i > 0 && got[i-1].ID >= m.ID {
			t.Errorf("models must be sorted by id: %q before %q", got[i-1].ID, m.ID)
		}
	}
}

func TestStudioModelsRequiresClientCredentials(t *testing.T) {
	f := newFixture(t)
	f.svc.groups = fakeGroups{}
	if _, err := f.svc.StudioModels(context.Background(), testClientID, "wrong-secret"); !errors.Is(err, ErrInvalidClient) {
		t.Fatalf("err = %v, want ErrInvalidClient", err)
	}
}

// 旧配置里遗留的 KeyGroupID 不再收窄目录。
func TestStudioModelsIgnoresLegacyKeyGroup(t *testing.T) {
	f := newFixture(t)
	cfg, err := f.svc.config.load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	legacy := int64(7)
	cfg.KeyGroupID = &legacy
	if err := f.svc.config.save(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	f.svc.groups = fakeGroups{
		{ID: 7, Models: []service.PlazaModel{{Name: "gpt-5", Platform: "openai"}}},
		{ID: 9, Models: []service.PlazaModel{{Name: "claude-sonnet-5-5", Platform: "anthropic"}}},
	}
	got, err := f.svc.StudioModels(context.Background(), testClientID, testClientSecret)
	if err != nil || len(got) != 2 {
		t.Fatalf("got %v, %v; want models from both groups", got, err)
	}
}
