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
	var err error
	// 1. 获取创建帖子的数据
	p := new(service.Publish)
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
	// 2. 获取用户Id
	uid, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, e.TokenInvalidAuth)
		return
	}
	uname, err := getCurrentUsername(c)
	if err != nil {
		ResponseError(c, e.TokenInvalidAuth)
		return
	}
	// 3. 将数据插入数据库中
	if err = p.PublishPost(uid, uname); err != nil {
		ResponseError(c, e.CodeServerBusy)
		zap.L().Error("publish post is failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, nil)
}

// PostDetailHandler 根据帖子ID获取帖子详情
func PostDetailHandler(c *gin.Context) {
	pid, err := getParamId(c, "pid")
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	p := new(service.PostService)
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
	p := new(service.Publish)
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
	pid, err := getParamId(c, "pid")
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
	pid, err := getParamId(c, "pid")
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
	data, err := service.PostListInOrder(page, size, order)
	if err != nil {
		zap.L().Error("PostListInOrder select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
