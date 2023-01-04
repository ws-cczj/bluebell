package middleware

import (
	"net/http"
	"time"

	"github.com/juju/ratelimit"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware 令牌桶限流中间件
func RateLimitMiddleware(fillInterval, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(time.Duration(fillInterval)*time.Second, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		// TakeAvailable 每次返回取出的令牌数，如果取出来的
		// 令牌数为0，说明没有令牌可以通过，就会被限流，直接返回，否则通行！
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
