package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 替换gin中默认的logger
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		if c.Writer.Status() != http.StatusOK {
			// 记录异常信息
			zap.L().Error(query,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		}
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("httpRequest", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				zap.L().Error("[Recovery from panic]", zap.Any("error", err))
				request := strings.Split(string(httpRequest), "\r\n")
				if stack {
					split := strings.Split(string(debug.Stack()), "\n\t")
					for _, str := range request {
						zap.L().Error("[Recovery from request panic]", zap.String("request", str))
					}
					for _, str := range split {
						zap.L().Error("[Recovery from Stack panic]", zap.String("stack", str))
					}
				} else {
					for _, str := range request {
						zap.L().Error("[Recovery from request panic]", zap.String("request", str))
					}
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
