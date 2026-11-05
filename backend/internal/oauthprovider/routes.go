package oauthprovider

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/middleware"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(NewService, NewHandler)

// RegisterRoutes 挂载全部 OAuth 相关路由。
//
//   - /oauth/token、/oauth/userinfo、/oauth/studio-key、/oauth/studio-models：给客户端后端调用，OAuth 协议格式
//   - /api/v1/oauth/*：给本站前端（授权确认页、顶部菜单、API 密钥页）
//   - /api/v1/admin/oauth-provider：管理后台配置
func RegisterRoutes(
	r *gin.Engine,
	v1 *gin.RouterGroup,
	h *Handler,
	jwtAuth servermiddleware.JWTAuthMiddleware,
	adminAuth servermiddleware.AdminAuthMiddleware,
	stepUpAuth servermiddleware.StepUpAuthMiddleware,
	auditLog servermiddleware.AuditLogMiddleware,
	redisClient *redis.Client,
) {
	limiter := middleware.NewRateLimiter(redisClient)
	failClose := middleware.RateLimitOptions{FailureMode: middleware.RateLimitFailClose}

	oauth := r.Group("/oauth")
	oauth.POST("/token", limiter.LimitWithOptions("oauth-token", 30, time.Minute, failClose), h.Token)
	oauth.GET("/userinfo", limiter.LimitWithOptions("oauth-userinfo", 60, time.Minute, failClose), h.UserInfo)
	oauth.GET("/studio-key", limiter.LimitWithOptions("oauth-studio-key", 30, time.Minute, failClose), h.StudioKey)
	oauth.GET("/studio-models", limiter.LimitWithOptions("oauth-studio-models", 30, time.Minute, failClose), h.StudioModels)

	v1.GET("/oauth/public", h.Public)

	user := v1.Group("/oauth")
	user.Use(gin.HandlerFunc(jwtAuth))
	user.Use(gin.HandlerFunc(auditLog))
	user.GET("/authorize", h.Check)
	user.POST("/authorize/approve", limiter.LimitWithOptions("oauth-approve", 30, time.Minute, failClose), h.Approve)
	user.GET("/studio-keys", h.StudioKeyIDs)

	admin := v1.Group("/admin/oauth-provider")
	admin.Use(gin.HandlerFunc(adminAuth))
	admin.Use(gin.HandlerFunc(auditLog))
	admin.GET("", h.AdminGet)
	// 改客户端密钥等同于改登录凭据，要求二次验证
	admin.PUT("", gin.HandlerFunc(stepUpAuth), h.AdminSave)
}
