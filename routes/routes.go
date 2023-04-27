package routes

import (
	"bluebell/api"
	"bluebell/middleware"
	silr "bluebell/serializer"
	"bluebell/settings"

	"github.com/ws-cczj/gee"

	"go.uber.org/zap"
)

func Setup() *gee.Engine {

	r := gee.Default(gee.WithExitOp(true),
		gee.WithReleaseMode(settings.Conf.Mode == "release"),
		gee.WithMiddlewares(gee.Cors(), gee.Logger(), gee.Recover(),
			middleware.RateLimit(settings.Conf.GenInterval, settings.Conf.MaxCaps)))

	if err := silr.InitTrans("zh", r); err != nil {
		zap.L().Error("init translation fail!", zap.Error(err))
	}

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
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuthMiddleware())
	{
		// user 用户
		admin.GET("/user/communitys", api.UserCommunityHandler)
		admin.GET("/user/posts", api.UserPostsHandler)
		admin.POST("/user/follow", api.UserFollowHandler)
		admin.GET("/user/:uid/to_follow", api.UserToFollowListHandler)
		admin.GET("/user/:uid/follow", api.UserFollowListHandler)
		// community 社区
		admin.POST("/community", api.CommunityCreateHandler)
		// post 帖子
		admin.POST("/post", api.PostPublishHandler)
		admin.PUT("/post/:pid", api.PostPutHandler)
		admin.DELETE("/post/:pid", api.PostDeleteHandler)
		admin.POST("/votes", api.PostVotesHandler)
		// comment 评论
		admin.POST("/comment", api.CommentPublishHandler)
		admin.DELETE("/comment", api.CommentDeleteHandler)
		admin.POST("/comment/favorite", api.CommentFavoriteHandler)
	}

	return r
}
