package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/open4go/log"
	"net/http"
)

// AccessMiddleware 验证cookie并且将解析出来的账号
// 通过账号获取角色
// 通过角色判断其是否具有该api的访问权限
// 用户登陆完成后会将权限配置信息写入 redis 数据库完成
// 通过hget api/path/ role boolean
func AccessMiddleware(key []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("jwt")
		if cookie == "" {
			log.Log(c.Request.Context()).
				WithField("cookie", "empty").
				Error("cookie name as jwt no found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

		if err != nil {
			log.Log(c.Request.Context()).WithField("call", "ParseWithClaims").
				Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*jwt.StandardClaims)
		loginInfo, err := LoadLoginInfo(claims.Issuer)
		if err != nil {
			log.Log(c.Request.Context()).WithField("call", "LoadLoginInfo").
				Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 写入解析客户的jwt token后得到的数据
		c.Request.Header.Set("MerchantId", loginInfo.Namespace)
		c.Request.Header.Set("AccountId", loginInfo.AccountId)
		c.Request.Header.Set("UserId", loginInfo.UserId)
		c.Request.Header.Set("UserName", loginInfo.UserName)
		c.Request.Header.Set("Avatar", loginInfo.Avatar)
		c.Request.Header.Set("LoginType", loginInfo.LoginType)
		c.Request.Header.Set("LoginLevel", loginInfo.LoginLevel)

		c.Next()
	}
}
