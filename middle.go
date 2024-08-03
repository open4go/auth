package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/open4go/log"
	"net/http"
	"strings"
)

func handleSingleResource(c *gin.Context) bool {
	// 检查路由参数 id 是否存在
	id := c.Param("_id")
	log.Log(c.Request.Context()).WithField("_id", id).Debug("Checking for single resource with parameter _id")
	return id != ""
}

// AccessMiddleware 在route部分应用该中间件
func AccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("UserID")
		roleManager := RoleManager{
			"system:auth",
			nil,
		}

		fullPath := c.FullPath()
		isSingleResource := handleSingleResource(c)
		if isSingleResource {
			fullPath = strings.TrimSuffix(fullPath, "/:_id")
		}

		verify, err := roleManager.Verify(c.Request.Context(), fullPath, userID, c.Request.Method, isSingleResource)
		if err != nil {
			log.Log(c.Request.Context()).
				WithField("request_path", c.FullPath()).
				WithField("request_method", c.Request.Method).
				WithField("userID", userID).
				WithField("userIP", c.ClientIP()).
				WithField("error", err).
				Error("Verification error: user does not have permission to access this endpoint")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !verify {
			log.Log(c.Request.Context()).
				WithField("request_path", c.FullPath()).
				WithField("request_method", c.Request.Method).
				WithField("userID", userID).
				WithField("userIP", c.ClientIP()).
				Error("Access denied: user does not have permission to access this endpoint")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
