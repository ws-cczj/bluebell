package api

import (
	"bluebell/pkg/e"
	"bluebell/service"

	"github.com/gin-gonic/gin"
)

// UserRegister 响应用户注册
func UserRegister(c *gin.Context) {
	var svc service.RegisterService
	code := e.SUCCESS
	if err := c.ShouldBind(&svc); err != nil {
		code = e.ERROR
		c.JSON(code, ErrorResponse(err))
		return
	}
	res, err := svc.Register()
	if err != nil {
		code = res.Status
		c.JSON(code, res)
		return
	}
	c.JSON(code, res)
}
