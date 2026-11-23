package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 与 open4go/middle 的 AbortSessionExpired / AbortPermissionDenied 保持同一响应格式：
//   - 401 code=session_expired：没有可信身份（JWT 未通过），需要重新登录；
//   - 403 code=permission_denied：已登录但没有该接口权限（含账号暂无角色，避免无角色账号登录后被反复踢回登录页）；
//   - 503 code=permission_unavailable：读角色/权限缓存出错，可重试。
const (
	CodeSessionExpired        = "session_expired"
	CodePermissionDenied      = "permission_denied"
	CodePermissionUnavailable = "permission_unavailable"
)

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
