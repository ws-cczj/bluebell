package api

import (
	"bluebell/pkg/e"
	"bluebell/service"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// UserRegisterHandler 响应用户注册
func UserRegisterHandler(c *gin.Context) {
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

// UserLoginHandler 响应用户登录
func UserLoginHandler(c *gin.Context) {
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

// UserCommunityHandler 获取该用户管理的社区列表
func UserCommunityHandler(c *gin.Context) {
	uid, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, e.TokenInvalidAuth)
		return
	}
	// 根据用户id去查询社区
	data, err := service.UserCommunityList(uid)
	if err != nil {
		zap.L().Error("service CommunityList method err", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// UserPostsHandler 获取该用户发布的帖子
func UserPostsHandler(c *gin.Context) {
	uid, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("getCurrentUserId method err", zap.Error(err))
		ResponseError(c, e.TokenInvalidAuth)
		return
	}
	data, err := service.UserPostList(uid)
	if err != nil {
		zap.L().Error("service UserPostList method err", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
