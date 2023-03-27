package routes

import (
	"bluebell/api"
	"bluebell/logger"
	"bluebell/middleware"
	silr "bluebell/serializer"
	"bluebell/settings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	if settings.Conf.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := silr.InitTrans("zh"); err != nil {
		zap.L().Error("init translation fail!", zap.Error(err))
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true),
		middleware.RateLimitMiddleware(settings.Conf.GenInterval, settings.Conf.MaxCaps))

	//pprof.Register(r)
	// api/v1
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", api.UserRegisterHandler)
		v1.POST("/login", api.UserLoginHandler)
		// community 社区
		v1.GET("/community", api.CommunityHandler)
		v1.GET("/community/:cid", api.CommunityDetailHandler)
		v1.GET("/community/:cid/posts", api.CommunityPostHandler)
		// post 帖子
		v1.GET("/posts", api.PostListHandler)
		v1.GET("/post/:pid", api.PostDetailHandler)
		// comment 评论
		v1.GET("/comment/:pid", api.CommentListHandler)
	}
	// jwt auth 用户认证
	v1.Use(middleware.JWTAuthMiddleware())
	{
		// user 用户
		v1.GET("/user/communitys", api.UserCommunityHandler)
		v1.GET("/user/posts", api.UserPostsHandler)
		v1.POST("user/follow", api.UserFollowHandler)
		v1.GET("user/:uid/to_follow", api.UserToFollowListHandler)
		v1.GET("user/:uid/follow", api.UserFollowListHandler)
		// community 社区
		v1.POST("/community", api.CommunityCreateHandler)
		// post 帖子
		v1.POST("/post", api.PostPublishHandler)
		v1.PUT("/post/:pid", api.PostPutHandler)
		v1.DELETE("/post/:pid", api.PostDeleteHandler)
		v1.POST("/votes", api.PostVotesHandler)
		// comment 评论
		v1.POST("/comment", api.CommentPublishHandler)
		v1.DELETE("/comment", api.CommentDeleteHandler)
		v1.POST("/comment/favorite", api.CommentFavoriteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		silr.ResponseNotFound(c)
	})
	return r
}
