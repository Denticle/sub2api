package app

import (
	"encoding/json"
	"strings"
	"testing"

	"yingce/backend/internal/auth"
	"yingce/backend/internal/model"
	"yingce/backend/internal/outbound/outboundtest"
)

func storeBridgeKey(t *testing.T, svc *Service, userID, key string) {
	t.Helper()
	encrypted, err := svc.encryptSettingSecret(key)
	if err != nil {
		t.Fatal(err)
	}
	value, _ := json.Marshal(map[string]string{"apiKey": encrypted})
	if err := svc.repo.SaveSystemSetting(&model.SystemSetting{Key: auth.ForAPIBridgeKeySettingPrefix + userID, ValueJSON: string(value)}); err != nil {
		t.Fatal(err)
	}
}

func newBridgeChannelTestService(t *testing.T, apiKey string) (*Service, model.ChannelModel) {
	t.Helper()
	outboundtest.PublicDNS(t, "forapi.cloud")
	svc, db := newChannelModelTestService(t)
	svc.dataDir = t.TempDir()
	if err := db.AutoMigrate(&model.SystemSetting{}); err != nil {
		t.Fatal(err)
	}
	channel := model.ModelChannel{
		ID: "forapi", Scope: model.ChannelScopeSystem, Enabled: true, Name: "ForAPI",
		BaseURL: "https://forapi.cloud/v1", APIKey: apiKey, APIFormat: "openai", ModelsJSON: `["gpt-image-2"]`,
	}
	if err := svc.encryptSystemChannelSecrets(&channel); err != nil {
		t.Fatal(err)
	}
	item := model.ChannelModel{
		ID: "model-1", ChannelID: channel.ID, ModelKey: "gpt-image-2", ProviderModelKey: "gpt-image-2",
		Capability: "image", Protocol: model.ChannelInterfaceChatCompletion, Enabled: true,
	}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}
	return svc, item
}

func TestBridgeChannelUsesEachUsersOwnKey(t *testing.T) {
	svc, item := newBridgeChannelTestService(t, ForAPIUserKeySentinel)
	storeBridgeKey(t, svc, "alice", "sk-alice")
	storeBridgeKey(t, svc, "bob", "sk-bob")
	for user, want := range map[string]string{"alice": "sk-alice", "bob": "sk-bob"} {
		config, err := svc.resolveProviderConfig(user, providerConfig{ChannelID: "forapi", Model: item.ModelKey})
		if err != nil {
			t.Fatal(err)
		}
		if config.APIKey != want {
			t.Fatalf("%s resolved key %q, want %q", user, config.APIKey, want)
		}
	}
}

func TestBridgeChannelFailsWithoutUserKey(t *testing.T) {
	svc, item := newBridgeChannelTestService(t, ForAPIUserKeySentinel)
	storeBridgeKey(t, svc, "alice", "sk-alice")
	for _, user := range []string{"carol", ""} {
		_, err := svc.resolveProviderConfig(user, providerConfig{ChannelID: "forapi", Model: item.ModelKey})
		if err == nil || !strings.Contains(err.Error(), "ForAPI") {
			t.Fatalf("user %q: err = %v, want missing-key error", user, err)
		}
	}
}

func TestOrdinaryChannelKeepsItsOwnKey(t *testing.T) {
	svc, item := newBridgeChannelTestService(t, "sk-channel")
	storeBridgeKey(t, svc, "alice", "sk-alice")
	config, err := svc.resolveProviderConfig("alice", providerConfig{ChannelID: "forapi", Model: item.ModelKey})
	if err != nil {
		t.Fatal(err)
	}
	if config.APIKey != "sk-channel" {
		t.Fatalf("resolved key %q, want channel key", config.APIKey)
	}
}
