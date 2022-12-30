package routes

import (
	"bluebell/api"
	"bluebell/logger"
	"bluebell/middleware"
	"bluebell/settings"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *settings.AppConfig) *gin.Engine {
	if cfg.Mode == gin.ReleaseMode {
		gin.SetMode(cfg.Mode)
	}
	if err := api.InitTrans("zh"); err != nil {
		zap.L().Error("init translation fail!", zap.Error(err))
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	v1.POST("/register", api.UserRegister)
	v1.POST("/login", api.UserLogin)

	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/community", api.CommunityHandler)
		v1.GET("/community/:id", api.CommunityDetailHandler)
		v1.GET("/community/:id/posts", api.CommunityPostHandler)

		v1.POST("/post", api.PostPublishHandler)
		v1.GET("/post/:id", api.PostDetailHandler)
		v1.GET("/posts", api.PostListHandler)
		v1.GET("/postsOrder", api.PostListOrderHandler)

		v1.POST("/votes", api.PostVotesHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"Msg": "请求的路径不存在",
		})
	})
	return r
}
