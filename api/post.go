package api

import (
	"bluebell/pkg/e"
	"bluebell/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostPublishHandler 帖子发布
func PostPublishHandler(c *gin.Context) {
	// 1. 获取创建帖子的数据
	p := new(service.PublishService)
	var err error
	if err = c.ShouldBind(p); err != nil {
		zap.L().Error("post publish params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, e.CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 将数据插入数据库中
	p.AuthorId, err = getCurrentUser(c)
	if err != nil {
		ResponseError(c, e.TokenInvalidAuth)
		return
	}
	if res, err := p.PublishPost(); err != nil {
		ResponseErrorWithRes(c, res)
		zap.L().Error("publish post is failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, nil)
}

// PostDetailHandler 根据帖子ID获取帖子详情
func PostDetailHandler(c *gin.Context) {
	pid, err := getPostId(c)
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	p := &service.PostService{}
	if err = p.PostDetailById(pid); err != nil {
		zap.L().Error("PostDetailById select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, p)
}

// PostPutHandler 修改帖子
func PostPutHandler(c *gin.Context) {
	// 1. 获取修改帖子的数据
	p := new(service.PublishService)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("post put params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, e.CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 获取帖子ID
	pid, err := getPostId(c)
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	if res, err := p.PostPut(pid); err != nil {
		ResponseErrorWithRes(c, res)
		zap.L().Error("PostPut is failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, nil)
}

// PostDeleteHandler 删除帖子
func PostDeleteHandler(c *gin.Context) {
	pid, err := getPostId(c)
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	if err = service.DeletePost(pid); err != nil {
		zap.L().Error("delete post failed", zap.Int64("pid", pid), zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// PostListHandler 顺序获取所有帖子
func PostListHandler(c *gin.Context) {
	page, size, order := getPostListInfo(c)
	p := &service.PostService{}
	data, err := p.PostListInOrder(page, size, order)
	if err != nil {
		zap.L().Error("PostListInOrder select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityPostHandler 获取社区的帖子
func CommunityPostHandler(c *gin.Context) {
	page, size, order := getPostListInfo(c)
	cid, err := getPostId(c)

	p := &service.PostService{}
	data, err := p.CommunityPostListInOrder(page, size, cid, order)
	if err != nil {
		zap.L().Error("CommunityPostListInOrder select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
