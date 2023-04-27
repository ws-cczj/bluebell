package api

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"

	"github.com/ws-cczj/gee"

	"go.uber.org/zap"
)

// UserRegisterHandler 响应用户注册
func UserRegisterHandler(c *gee.Context) {
	u := new(models.UserRegister)
	if err := c.ShouldBind(u); err != nil {
		zap.L().Error("register params is not illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	if err := service.NewUserInstance().Register(u); err != nil {
		zap.L().Error("service UserRegister method err", zap.Error(err))
		if err == mysql.ErrorUserExist {
			silr.ResponseError(c, e.CodeExistUser)
			return
		}
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// UserLoginHandler 响应用户登录
func UserLoginHandler(c *gee.Context) {
	u := new(models.UserLogin)
	if err := c.ShouldBind(u); err != nil {
		zap.L().Error("register params is not illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	atoken, rtoken, err := service.NewUserInstance().Login(u)
	if err != nil {
		zap.L().Error("user login method err", zap.Error(err))
		if err == mysql.ErrNoRows {
			silr.ResponseError(c, e.CodeExistUser)
			return
		} else if err == mysql.ErrorNotComparePwd {
			silr.ResponseError(c, e.CodeNotComparePassword)
			return
		}
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, silr.ResponseUserLogin{UserId: u.UserId, Atoken: atoken, Rtoken: rtoken})
}

// UserFollowHandler 用户关注
func UserFollowHandler(c *gee.Context) {
	u := new(models.UserFollow)
	if err := c.ShouldBind(u); err != nil {
		zap.L().Error("userFollowHandler method param is illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	uid, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("UserFollowHandler getCurrentUserId method err", zap.Error(err))
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	if err = service.NewUserInstance().FollowBuild(uid, u); err != nil {
		zap.L().Error("service userFollowBuild method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// UserToFollowListHandler 用户关注列表
func UserToFollowListHandler(c *gee.Context) {
	uid := c.Param("uid")
	data, err := service.NewUserInstance().ToFollowList(uid)
	if err != nil {
		zap.L().Error("service userToFollowList method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// UserFollowListHandler 用户粉丝列表
func UserFollowListHandler(c *gee.Context) {
	uid := c.Param("uid")
	data, err := service.NewUserInstance().FollowList(uid)
	if err != nil {
		zap.L().Error("service userToFollowList method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// UserCommunityHandler 用户管理的社区
func UserCommunityHandler(c *gee.Context) {
	uid, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("getCurrentUserId method err", zap.Error(err))
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	// 根据用户id去查询社区
	data, err := service.NewUserInstance().CommunityList(uid)
	if err != nil {
		zap.L().Error("service CommunityList method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// UserPostsHandler 用户发布的帖子
func UserPostsHandler(c *gee.Context) {
	uid, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("getCurrentUserId method err", zap.Error(err))
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	data, err := service.NewUserInstance().PostList(uid)
	if err != nil {
		zap.L().Error("service UserPostList method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}
