package api

import (
	"bluebell/pkg/e"
	"bluebell/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostPublishHandler 帖子发布
func PostPublishHandler(c *gin.Context) {
	// 1. 获取创建帖子的数据
	p := new(service.PublishService)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("register params is not illegal", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, e.CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, http.StatusBadRequest, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 将数据插入数据库中
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, e.TokenInvalidAuth)
		return
	}
	p.AuthorId = userID
	res, err := p.PublishPost()
	if err != nil {
		ResponseErrorWithRes(c, res)
		zap.L().Error("publish post is failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, nil)
}

// PostDetailHandler 根据帖子ID获取帖子详情
func PostDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	p := &service.PostService{}
	err = p.PostDetailById(pid)
	if err != nil {
		zap.L().Error("PostDetailById select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, p)
}

// PostListHandler 获取所有的帖子
func PostListHandler(c *gin.Context) {
	page, size := getPostListInfo(c)
	p := &service.PostService{}
	data, err := p.PostList(page, size)
	if err != nil {
		zap.L().Error("PostList select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
