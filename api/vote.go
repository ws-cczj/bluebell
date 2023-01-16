package api

import (
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVotesHandler 帖子投票
func PostVotesHandler(c *gin.Context) {
	v := new(service.PostVoteService)
	if err := c.ShouldBind(v); err != nil {
		zap.L().Error("postVote ShouldBind method failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if ok {
		}
		silr.ResponseErrorWithMsg(c, e.CodeInvalidParams, errs.Translate(trans))
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
