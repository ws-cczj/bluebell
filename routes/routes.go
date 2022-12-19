package routes

import (
	"bluebell/api"
	"bluebell/logger"
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
	r.GET("/index", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "routes build ok")
	})
	user := r.Group("/user")
	{
		user.POST("/register", api.UserRegister)
		user.POST("/login", api.UserLogin)
	}
	return r
}
