package oauthprovider

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// ---- 面向前端（JWT 已登录） ----

func authorizeRequestFrom(c *gin.Context) AuthorizeRequest {
	pick := func(k string) string {
		if v := c.PostForm(k); v != "" {
			return v
		}
		return c.Query(k)
	}
	return AuthorizeRequest{
		ClientID:            pick("client_id"),
		RedirectURI:         pick("redirect_uri"),
		ResponseType:        pick("response_type"),
		State:               pick("state"),
		CodeChallenge:       pick("code_challenge"),
		CodeChallengeMethod: pick("code_challenge_method"),
	}
}

// Check GET /api/v1/oauth/authorize —— 授权页加载时调用，判断要不要显示确认页
func (h *Handler) Check(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	check, err := h.svc.CheckAuthorize(c.Request.Context(), subject.UserID, authorizeRequestFrom(c))
	if err != nil {
		writeAuthorizeError(c, err)
		return
	}
	response.Success(c, check)
}

type approveBody struct {
	ClientID            string `json:"client_id"`
	RedirectURI         string `json:"redirect_uri"`
	ResponseType        string `json:"response_type"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
}

// Approve POST /api/v1/oauth/authorize/approve —— 返回带 code 的回调地址，由前端跳转
func (h *Handler) Approve(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	var body approveBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	target, err := h.svc.Approve(c.Request.Context(), subject.UserID, AuthorizeRequest(body))
	if err != nil {
		writeAuthorizeError(c, err)
		return
	}
	response.Success(c, gin.H{"redirect_to": target})
}

// StudioKeyIDs GET /api/v1/oauth/studio-keys —— 「API 密钥」页用来标出工作台专用密钥
func (h *Handler) StudioKeyIDs(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	ids, err := h.svc.StudioKeyIDs(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]int64, 0, len(ids))
	for id := range ids {
		out = append(out, id)
	}
	response.Success(c, gin.H{"ids": out})
}

// Public GET /api/v1/oauth/public —— 前端顶部菜单用：开关、名称、地址
func (h *Handler) Public(c *gin.Context) {
	response.Success(c, h.svc.PublicConfig(c.Request.Context()))
}

// ---- 管理后台 ----

type adminConfigView struct {
	Config
	HasClientSecret bool `json:"has_client_secret"`
}

func (h *Handler) AdminGet(c *gin.Context) {
	cfg, has, err := h.svc.AdminConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, adminConfigView{Config: cfg, HasClientSecret: has})
}

func (h *Handler) AdminSave(c *gin.Context) {
	var body Config
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	if err := h.svc.SaveConfig(c.Request.Context(), body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	h.AdminGet(c)
}

// ---- 面向客户端后端（OAuth 协议格式，不走统一的 {code,data} 包装） ----

func clientCredentials(c *gin.Context) (string, string) {
	if id, secret, ok := c.Request.BasicAuth(); ok {
		return id, secret
	}
	return c.PostForm("client_id"), c.PostForm("client_secret")
}

func bearer(c *gin.Context) string {
	v := c.GetHeader("Authorization")
	if len(v) > 7 && strings.EqualFold(v[:7], "Bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return ""
}

// Token POST /oauth/token
func (h *Handler) Token(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	id, secret := clientCredentials(c)
	resp, err := h.svc.Exchange(c.Request.Context(), TokenRequest{
		GrantType:    c.PostForm("grant_type"),
		Code:         c.PostForm("code"),
		RedirectURI:  c.PostForm("redirect_uri"),
		CodeVerifier: c.PostForm("code_verifier"),
		ClientID:     id,
		ClientSecret: secret,
	})
	if err != nil {
		writeOAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UserInfo GET /oauth/userinfo
func (h *Handler) UserInfo(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	info, err := h.svc.UserInfo(c.Request.Context(), bearer(c))
	if err != nil {
		c.Header("WWW-Authenticate", `Bearer error="invalid_token"`)
		writeOAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// StudioKey GET /oauth/studio-key —— 需要客户端密钥（Basic）+ 访问令牌（X-Access-Token）
func (h *Handler) StudioKey(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	id, secret, _ := c.Request.BasicAuth()
	key, err := h.svc.StudioKey(c.Request.Context(), id, secret, c.GetHeader("X-Access-Token"))
	if err != nil {
		writeOAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"api_key": key})
}

// StudioModels GET /oauth/studio-models —— 需要客户端密钥（Basic）
func (h *Handler) StudioModels(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	id, secret, _ := c.Request.BasicAuth()
	models, err := h.svc.StudioModels(c.Request.Context(), id, secret)
	if err != nil {
		writeOAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func writeOAuthError(c *gin.Context, err error) {
	code, status := "server_error", http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrInvalidClient):
		code, status = "invalid_client", http.StatusUnauthorized
	case errors.Is(err, ErrInvalidGrant):
		code, status = "invalid_grant", http.StatusBadRequest
	case errors.Is(err, ErrInvalidRequest):
		code, status = "invalid_request", http.StatusBadRequest
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrUserInactive):
		code, status = "invalid_token", http.StatusUnauthorized
	case errors.Is(err, ErrKeyUnavailable):
		code, status = "key_disabled", http.StatusForbidden
	case errors.Is(err, ErrDisabled):
		code, status = "unauthorized_client", http.StatusForbidden
	}
	c.JSON(status, gin.H{"error": code})
}

// authorizeErrors 给授权页的错误配上用户看得懂的说明和稳定的 reason
var authorizeErrors = []struct {
	err     error
	reason  string
	message string
}{
	{ErrDisabled, "OAUTH_DISABLED", "视频工作台接入尚未开启"},
	{ErrInvalidClient, "OAUTH_INVALID_CLIENT", "未知的应用"},
	{ErrInvalidRedirect, "OAUTH_INVALID_REDIRECT", "回调地址不匹配"},
	{ErrInvalidRequest, "OAUTH_INVALID_REQUEST", "授权参数不完整"},
	{ErrUserInactive, "OAUTH_USER_INACTIVE", "账号已被停用"},
	{ErrKeyUnavailable, "OAUTH_KEY_DISABLED", "你的「工作台」密钥已停用，请在「API 密钥」页重新启用后再试"},
}

func writeAuthorizeError(c *gin.Context, err error) {
	for _, e := range authorizeErrors {
		if errors.Is(err, e.err) {
			response.ErrorWithDetails(c, http.StatusBadRequest, e.message, e.reason, nil)
			return
		}
	}
	response.ErrorFrom(c, err)
}
