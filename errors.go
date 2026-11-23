package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 与 open4go/middle 的 AbortSessionExpired / AbortPermissionDenied 保持同一响应格式：
//   - 401 code=session_expired：没有可信身份或角色缓存已失效（与 JWT 同为 24h），需要重新登录；
//   - 403 code=permission_denied：已登录但没有该接口权限；
//   - 503 code=permission_unavailable：读角色/权限缓存出错，可重试。
const (
	CodeSessionExpired        = "session_expired"
	CodePermissionDenied      = "permission_denied"
	CodePermissionUnavailable = "permission_unavailable"
)

// ErrNoRolesCached 账号在缓存中没有角色（未登录或角色缓存已随登录过期）。
var ErrNoRolesCached = errors.New("auth: account has no roles in cache")

func abortSessionExpired(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    CodeSessionExpired,
		"message": "登录已过期，请重新登录",
	})
}

func abortPermissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"code":    CodePermissionDenied,
		"message": "没有访问该功能的权限",
	})
}

func abortPermissionUnavailable(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"code":    CodePermissionUnavailable,
		"message": "权限服务暂不可用，请稍后重试",
	})
}
