package api

import (
	"bluebell/pkg/e"
	"bluebell/service"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// UserRegister 响应用户注册
func UserRegister(c *gin.Context) {
	var svc service.RegisterService
	if err := c.ShouldBind(&svc); err != nil {
		zap.L().Error("register params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, e.CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	res, err := svc.Register()
	if err != nil {
		ResponseErrorWithRes(c, res)
		return
	}
	ResponseSuccess(c, nil)
}

// UserLogin 响应用户登录
func UserLogin(c *gin.Context) {
	var svc service.LoginService
	if err := c.ShouldBind(&svc); err != nil {
		zap.L().Error("register params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, e.CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	res, err := svc.Login()
	if err != nil {
		ResponseErrorWithRes(c, res)
		return
	}
	ResponseSuccess(c, res.Data)
}
