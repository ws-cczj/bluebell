package api

import (
	"bluebell/models"
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommentPublishHandler 评论发布
func CommentPublishHandler(c *gin.Context) {
	// 1. 绑定参数
	comment := new(models.CommentDetail)
	if err := c.ShouldBind(comment); err != nil {
		zap.L().Error("comment publish params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			silr.ResponseError(c, e.CodeInvalidParams)
			return
		}
		silr.ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 查找当前请求用户的uid和uname
	uid, err := getCurrentUserId(c)
	if err != nil {
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	uname, err := getCurrentUsername(c)
	if err != nil {
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	comment.Comment.AuthorId = uid
	comment.Comment.AuthorName = uname
	// 3. 发布评论
	if err = service.PublishComment(comment); err != nil {
		zap.L().Error("service publish method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
	}
	silr.ResponseSuccess(c, nil)
}

// CommentFavoriteHandler 评论点赞或取消点赞
func CommentFavoriteHandler(c *gin.Context) {
	favorite := new(models.Favorite)
	if err := c.ShouldBind(favorite); err != nil {
		zap.L().Error("Comment Favorite params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			silr.ResponseError(c, e.CodeInvalidParams)
			return
		}
		silr.ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	userID, err := getCurrentUserId(c)
	if err != nil {
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	if err = service.FavoriteBuild(favorite, userID); err != nil {
		zap.L().Error("service favoriteBuild method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// CommentDeleteHandler 删除评论
func CommentDeleteHandler(c *gin.Context) {
	commentD := new(models.CommentDelete)
	if err := c.ShouldBind(commentD); err != nil {
		zap.L().Error("CommentDeleteHandler ShouldBind method err", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			silr.ResponseError(c, e.CodeInvalidParams)
			return
		}
		silr.ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	if err := service.DeleteComment(commentD); err != nil {
		zap.L().Error("service DeleteComment method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// CommentListHandler 获取所有评论信息
func CommentListHandler(c *gin.Context) {
	pid, err := getParamId(c, "pid")
	if err != nil {
		silr.ResponseError(c, e.CodeInvalidParams)
		return
	}
	order := c.Query("order")
	data, err := service.GetCommentList(pid, order)
	if err != nil {
		zap.L().Error("service favoriteBuild method err", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}
