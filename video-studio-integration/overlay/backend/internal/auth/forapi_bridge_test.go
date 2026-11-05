package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"yingce/backend/internal/model"
	"yingce/backend/internal/outbound"
)

// fakeForAPI 模拟 forapi 的 /oauth/token、/oauth/userinfo、/oauth/studio-key。
func fakeForAPI(t *testing.T, studioStatus int, studioBody string, seen *http.Header) {
	t.Helper()
	t.Setenv("CANVAS_ALLOWED_PRIVATE_UPSTREAM_HOSTS", "127.0.0.1")
	transport := outbound.OutboundHTTPClient(time.Second).Transport.(*http.Transport)
	transport.CloseIdleConnections()
	previousDial := transport.DialContext
	transport.DialContext = func(context.Context, string, string) (net.Conn, error) {
		client, server := net.Pipe()
		go func() {
			defer server.Close()
			req, err := http.ReadRequest(bufio.NewReader(server))
			if err != nil {
				t.Error(err)
				return
			}
			defer req.Body.Close()
			_, _ = io.Copy(io.Discard, req.Body)
			status, body := 200, `{"sub":"42","username":"member"}`
			switch req.URL.Path {
			case "/oauth/token":
				body = `{"access_token":"test-token"}`
			case "/oauth/studio-key":
				*seen = req.Header.Clone()
				status, body = studioStatus, studioBody
			}
			_, _ = fmt.Fprintf(server, "HTTP/1.1 %d %s\r\nContent-Type: application/json\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", status, http.StatusText(status), len(body), body)
		}()
		return client, nil
	}
	t.Cleanup(func() { transport.CloseIdleConnections(); transport.DialContext = previousDial })
}

func setupForAPILogin(t *testing.T) *Service {
	t.Helper()
	svc, db := newRegistrationTestService(t)
	if err := svc.repo.CreateOAuthUser(&model.User{ID: "member", Username: "member", Role: model.UserRoleUser, Status: model.UserStatusActive}, &model.UserIdentity{ID: "identity", UserID: "member", Provider: "linuxdo", Subject: "42"}); err != nil {
		t.Fatal(err)
	}
	setting, err := json.Marshal(linuxDOSettingValue{Enabled: true, ClientID: "client", ClientSecret: "secret", ClientAuthMethod: "client_secret_basic", TokenURL: "http://127.0.0.1/oauth/token", UserInfoURL: "http://127.0.0.1/oauth/userinfo", SubjectField: "sub", UsernameField: "username"})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.SystemSetting{Key: linuxDOSettingKey, ValueJSON: string(setting)}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.OAuthState{ID: "state", Provider: "linuxdo", StateHash: HashToken("state-value"), ExpiresAt: time.Now().Add(time.Hour)}).Error; err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestForAPILoginStoresUserStudioKey(t *testing.T) {
	var seen http.Header
	fakeForAPI(t, 200, `{"api_key":"sk-member"}`, &seen)
	svc := setupForAPILogin(t)
	result, err := svc.CompleteLinuxDOLogin("state-value", "code", "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Session == nil {
		t.Fatal("login did not create a session")
	}
	if seen.Get("X-Access-Token") != "test-token" || !strings.HasPrefix(seen.Get("Authorization"), "Basic ") {
		t.Fatalf("studio-key request missing credentials: %v", seen)
	}
	setting, err := svc.repo.SystemSetting(ForAPIBridgeKeySettingPrefix + "member")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(setting.ValueJSON, "sk-member") {
		t.Fatalf("stored key = %s", setting.ValueJSON)
	}
}

func TestForAPILoginFailsWithoutStudioKey(t *testing.T) {
	for name, tc := range map[string]struct {
		status int
		body   string
		want   string
	}{
		"key disabled":  {403, `{"error":"key_disabled"}`, "已停用"},
		"server error":  {500, `{"error":"server_error"}`, "HTTP 500"},
		"empty api_key": {200, `{"api_key":""}`, "获取 ForAPI 密钥失败"},
	} {
		t.Run(name, func(t *testing.T) {
			var seen http.Header
			fakeForAPI(t, tc.status, tc.body, &seen)
			svc := setupForAPILogin(t)
			_, err := svc.CompleteLinuxDOLogin("state-value", "code", "")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want containing %q", err, tc.want)
			}
			if setting, _ := svc.repo.SystemSettingOptional(ForAPIBridgeKeySettingPrefix + "member"); setting != nil {
				t.Fatal("failed login must not store a key")
			}
		})
	}
}

func TestPasswordLoginOnlyForAdmins(t *testing.T) {
	svc, db := newRegistrationTestService(t)
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}
	for _, user := range []model.User{
		{ID: "admin", Username: "admin", PasswordHash: hash, Role: model.UserRoleAdmin, Status: model.UserStatusActive},
		{ID: "member", Username: "member", PasswordHash: hash, Role: model.UserRoleUser, Status: model.UserStatusActive},
	} {
		if err := db.Create(&user).Error; err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.Login(LoginRequest{Username: "admin", Password: "password123"}); err != nil {
		t.Fatalf("admin password login: %v", err)
	}
	if _, err := svc.Login(LoginRequest{Username: "member", Password: "password123"}); err == nil || !strings.Contains(err.Error(), "ForAPI") {
		t.Fatalf("member password login err = %v", err)
	}
}

func TestStudioKeyURLFollowsTokenURL(t *testing.T) {
	got, err := forAPIStudioKeyURL("https://forapi.cloud/oauth/token?x=1")
	if err != nil || got != "https://forapi.cloud/oauth/studio-key" {
		t.Fatalf("got %q, %v", got, err)
	}
}
