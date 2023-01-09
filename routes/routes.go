package routes

import (
	"bluebell/api"
	"bluebell/logger"
	"bluebell/middleware"
	"bluebell/settings"

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
	r.Use(logger.GinLogger(), logger.GinRecovery(true),
		middleware.RateLimitMiddleware(cfg.RateLimit.GenInterval, cfg.RateLimit.MaxCaps))

	//pprof.Register(r)
	// api/v1
	v1 := r.Group("/api/v1")
	v1.POST("/register", api.UserRegister)
	v1.POST("/login", api.UserLogin)
	// community 社区
	v1.GET("/community", api.CommunityHandler)
	v1.GET("/community/:id", api.CommunityDetailHandler)
	v1.GET("/community/:id/posts", api.CommunityPostHandler)
	// post 帖子
	v1.GET("/posts", api.PostListHandler)
	v1.GET("/post/:id", api.PostDetailHandler)
	// jwt auth 用户认证
	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.POST("/post", api.PostPublishHandler)
		v1.PUT("/post/:id", api.PostPutHandler)
		v1.DELETE("/post/:id", api.PostDeleteHandler)
		v1.POST("/votes", api.PostVotesHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		api.ResponseNotFound(c)
	})
	return r
}
