package api

import (
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"

	"github.com/ws-cczj/gee"

	"go.uber.org/zap"
)

// PostVotesHandler 帖子投票
func PostVotesHandler(c *gee.Context) {
	v := new(service.PostVote)
	if err := c.ShouldBind(v); err != nil {
		zap.L().Error("postVote ShouldBind method failed", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	userID, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("getCurrentUser method Error", zap.Error(err))
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	res, err := v.Build(userID)
	if err != nil {
		zap.L().Error("voteBuild method Error", zap.Error(err))
		silr.ResponseErrorWithRes(c, res)
		return
	}
	silr.ResponseSuccess(c, nil)
}
