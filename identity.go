package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open4go/model"
	"github.com/spf13/viper"
)

// VerifiedUserIDKey 与 open4go/middle.VerifiedUserIDKey 一致：middle.JWTMiddleware
// 校验 token 成功后把 token 内的 UserID 写入 gin.Context 与 request context。
const VerifiedUserIDKey = "OPEN4GO_VERIFIED_USER_ID"

// 身份来源配置 auth.identitySource：
//
//	verified（默认）  只信任 middle.JWTMiddleware 校验后写入的身份（需要新版 middle）
//	context           读取 request context 中的 model.OperatorKey（兼容旧版 middle，需确保前面挂了 JWTMiddleware）
//	header            读取 UserID 请求头（旧行为，仅限网关已剥离并重写该头的部署，不推荐）
const (
	cfgIdentitySource = "auth.identitySource"

	IdentityVerified = "verified"
	IdentityContext  = "context"
	IdentityHeader   = "header"
)

func identitySource() string {
	switch v := strings.ToLower(strings.TrimSpace(viper.GetString(cfgIdentitySource))); v {
	case IdentityContext, IdentityHeader:
		return v
	default:
		return IdentityVerified
	}
}

func stringFromCtx(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	if c.Request != nil {
		if s, ok := c.Request.Context().Value(key).(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

// resolveUserID 返回用于权限校验的调用者身份。
func resolveUserID(c *gin.Context) (string, string) {
	src := identitySource()
	switch src {
	case IdentityHeader:
		return strings.TrimSpace(c.GetHeader("UserID")), src
	case IdentityContext:
		return stringFromCtx(c, model.OperatorKey), src
	default:
		return stringFromCtx(c, VerifiedUserIDKey), src
	}
}
