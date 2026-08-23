package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open4go/log"
)

func handleSingleResource(c *gin.Context) bool {
	id := c.Param("_id")
	return id != ""
}

func isWebSocketRequest(c *gin.Context) bool {
	if c.IsWebsocket() {
		return true
	}
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}
	return isWSEndpoint(path)
}

func lookupPath(c *gin.Context) (fullPath string, isSingle bool) {
	fullPath = c.FullPath()
	if fullPath == "" {
		fullPath = c.Request.URL.Path
	}
	isSingle = handleSingleResource(c)
	if isSingle {
		fullPath = strings.TrimSuffix(fullPath, "/:_id")
	}
	return fullPath, isSingle
}

// AccessMiddleware 在 route 部分应用该中间件。
// WebSocket 握手只需校验 GET 权限（List 或 Read），不要求 Create/Update/Delete。
func AccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("UserID")
		roleManager := RoleManager{
			"system:auth",
			nil,
		}

		fullPath, isSingleResource := lookupPath(c)
		isWS := isWebSocketRequest(c)
		method := c.Request.Method
		if isWS {
			method = http.MethodGet
		}

		baseLog := log.Log(c.Request.Context()).
			WithField("request_path", c.FullPath()).
			WithField("lookup_path", fullPath).
			WithField("url_path", c.Request.URL.Path).
			WithField("request_method", c.Request.Method).
			WithField("check_method", method).
			WithField("is_ws", isWS).
			WithField("is_single_resource", isSingleResource).
			WithField("userID", userID).
			WithField("userIP", c.ClientIP())

		if userID == "" {
			baseLog.Error("Access denied: missing UserID header, cannot resolve account roles")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if isWS {
			baseLog.Debug("websocket request: permission check uses GET only (List or Read)")
		}

		ctx := withWSFlag(c.Request.Context(), isWS)
		verify, err := roleManager.Verify(ctx, fullPath, userID, method, isSingleResource)
		if err != nil {
			baseLog.WithField("error", err).
				Error("Verification error: failed to check permission for this endpoint")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !verify {
			baseLog.Error("Access denied: user does not have permission to access this endpoint")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
