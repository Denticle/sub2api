package middleware

import (
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// JWTAuthMiddleware JWT 认证中间件类型
type JWTAuthMiddleware gin.HandlerFunc

// OptionalJWTAuthMiddleware 可选 JWT 认证中间件类型：匿名放行，带 token 严格校验
type OptionalJWTAuthMiddleware gin.HandlerFunc

// AdminAuthMiddleware 管理员认证中间件类型
type AdminAuthMiddleware gin.HandlerFunc

// APIKeyAuthMiddleware API Key 认证中间件类型
type APIKeyAuthMiddleware gin.HandlerFunc

// ProviderSet 中间件层的依赖注入
var ProviderSet = wire.NewSet(
	NewJWTAuthMiddleware,
	NewOptionalJWTAuthMiddleware,
	NewAdminAuthMiddleware,
	ProvideAPIKeyAuthMiddleware,
	NewAuditLogMiddleware,
	NewStepUpAuthMiddleware,
)

// ProvideAPIKeyAuthMiddleware 是 wire 用的构造函数。
//
// 依赖 *service.GroupRouter 是为了让 wire 在服务启动时构建选组服务：
// ProvideGroupRouter 会把它挂到 apiKeyService 上，鉴权中间件（含 Google 变体）
// 通过 apiKeyService.GroupRouter() 使用。没有消费者时 wire 会裁掉该 provider。
func ProvideAPIKeyAuthMiddleware(apiKeyService *service.APIKeyService, subscriptionService *service.SubscriptionService, _ *service.GroupRouter, cfg *config.Config) APIKeyAuthMiddleware {
	return NewAPIKeyAuthMiddleware(apiKeyService, subscriptionService, cfg)
}
