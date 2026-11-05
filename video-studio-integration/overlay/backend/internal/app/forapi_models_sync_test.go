package app

import (
	"testing"

	"yingce/backend/internal/model"
)

func TestForAPIModelProtocolMapping(t *testing.T) {
	cases := []struct {
		in         ForAPIModel
		capability string
		protocol   model.ChannelInterfaceType
	}{
		{ForAPIModel{ID: "claude-sonnet-5-5", Capability: "text"}, "text", model.ChannelInterfaceChatCompletion},
		{ForAPIModel{ID: "gpt-image-2", Capability: "image"}, "image", model.ChannelInterfaceOpenAIImage},
		{ForAPIModel{ID: "doubao-seedance-2-0", Capability: "video"}, "video", model.ChannelInterfaceVolcengineArkVideo},
		{ForAPIModel{ID: "grok-imagine-video", Capability: "video"}, "video", model.ChannelInterfaceXAIVideo},
		{ForAPIModel{ID: "unknown", Capability: ""}, "text", model.ChannelInterfaceChatCompletion},
	}
	for _, c := range cases {
		capability, protocol := forAPIModelProtocol(c.in)
		if capability != c.capability || protocol != c.protocol {
			t.Errorf("%s: got (%s, %s), want (%s, %s)", c.in.ID, capability, protocol, c.capability, c.protocol)
		}
	}
}

func TestForAPIChannelModelRequestIsFreeAndEnabled(t *testing.T) {
	for _, m := range []ForAPIModel{
		{ID: "glm-5.3", Capability: "text"},
		{ID: "gpt-image-2", Capability: "image"},
		{ID: "grok-imagine-video", Capability: "video"},
	} {
		req := forAPIChannelModelRequest(m)
		if req.ModelKey != m.ID || req.ProviderModelKey != m.ID {
			t.Errorf("%s: model key = %q/%q", m.ID, req.ModelKey, req.ProviderModelKey)
		}
		if req.Enabled == nil || !*req.Enabled || !req.PriceConfigured {
			t.Errorf("%s: must be enabled with price configured", m.ID)
		}
		if req.UnitPriceMicrocredits != 0 || req.InputTokenPriceMicrocredits != 0 || req.OutputTokenPriceMicrocredits != 0 {
			t.Errorf("%s: studio credits must be 0, billing happens in forapi", m.ID)
		}
		if req.CapabilityConfig == nil {
			t.Fatalf("%s: capability config missing", m.ID)
		}
		switch m.Capability {
		case "text":
			if req.CapabilityConfig.Text == nil || req.BillingMode != "token" {
				t.Errorf("%s: text config/billing wrong", m.ID)
			}
		case "image":
			if req.CapabilityConfig.Image == nil {
				t.Errorf("%s: image config missing", m.ID)
			}
		case "video":
			if req.CapabilityConfig.Video == nil {
				t.Errorf("%s: video config missing", m.ID)
			}
		}
	}
}
