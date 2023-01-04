package middleware

import (
	"bluebell/api"
	"bluebell/dao/redis"
	"bluebell/pkg/e"
	"bluebell/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		rtoken := c.Request.Header.Get("Grant_type")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			api.ResponseError(c, e.TokenNullNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			api.ResponseError(c, e.TokenInvalidAuth)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.VerifyToken(parts[1])
		if err != nil {
			if newAtoken, newRtoken, err := jwt.RefreshToken(parts[1], rtoken); err == nil {
				if err = redis.SetSingleUserToken(mc.Username, newAtoken); err == nil {
					c.Writer.Header().Set("newAtoken", newAtoken)
					c.Writer.Header().Set("newRtoken", newRtoken)
					c.Set(api.ContextUserIDKey, mc.UserID)
					c.Next()
				}
			}
			api.ResponseError(c, e.TokenFailVerify)
			c.Abort()
			return
		}
		// 通过获取redis中的token来校验是否单用户登录
		token, err := redis.GetSingleUserToken(mc.Username)
		if err != nil {
			api.ResponseError(c, e.CodeServerBusy)
			c.Abort()
			return
		}
		if token != parts[1] {
			api.ResponseError(c, e.CodeRepeatLogin)
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(api.ContextUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
	}
}
