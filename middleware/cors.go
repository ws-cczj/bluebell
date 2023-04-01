package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func isAllowOrigin(origin string) bool {
	return true
}

func Cors() func(c *gin.Context) {
	return func(c *gin.Context) {
		method := c.Request.Method
		header := c.Writer.Header()
		origin := c.Request.Header.Get("Origin")
		// 如果origin为 "" 说明不是跨域，null 不等于 ""
		if origin != "" {
			// 解决cdn缓存问题
			header.Add("Vary", "Origin")
			if !isAllowOrigin(origin) {
				// 该请求不允许
				c.Abort()
				return
			}
			// 说明这个请求是一个预请求
			if method == http.MethodOptions {
				reqMethod := header.Get("Access-Control-Request-Method")
				if reqMethod != "PUT" && reqMethod != "DELETE" {
					// 如果方法不对，不提供Access Origin字段。
					c.Abort()
					return
				}
				// 必须，接受指定域的请求，可以使用*不加以限制，但不安全
				//header.Set("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Origin", origin)
				// 预请求最大缓存时间
				c.Header("Access-Control-Max-Age", "86400")
				// 必须，设置服务器支持的所有跨域请求的方法
				c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
				// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
				c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
				// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
				c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
				// 可选，是否允许后续请求携带认证信息Cookie，该值只能是true，不需要则不设置
				c.Header("Access-Control-Allow-Credentials", "true")
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			header.Set("Access-Control-Allow-Origin", origin)
			// 必须，设置服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
			// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
			// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
			c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
			// 可选，是否允许后续请求携带认证信息Cookie，该值只能是true，不需要则不设置
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		c.Next()
	}
}
