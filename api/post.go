package api

import (
	"bluebell/models"
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"

	"github.com/ws-cczj/gee"

	"go.uber.org/zap"
)

// PostPublishHandler 帖子发布
func PostPublishHandler(c *gee.Context) {
	// 1. 获取创建帖子的数据
	p := new(service.Publish)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("post publish params is not illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	// 2. 获取用户Id
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
	// 3. 将数据插入数据库中
	if err = p.Publish(uid, uname); err != nil {
		silr.ResponseError(c, e.CodeServerBusy)
		zap.L().Error("publish post is failed", zap.Error(err))
		return
	}
	silr.ResponseSuccess(c, nil)
}

// PostDetailHandler 根据帖子ID获取帖子详情
func PostDetailHandler(c *gee.Context) {
	pid := c.Param("pid")
	p := new(service.PostAll)
	if err := p.DetailById(pid); err != nil {
		zap.L().Error("PostDetailById select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, p)
}

// PostPutHandler 修改帖子
func PostPutHandler(c *gee.Context) {
	// 1. 获取修改帖子的数据
	p := new(models.PostPut)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("post put params is not illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	// 2. 获取帖子ID
	pid := c.Param("pid")
	if res, err := service.NewPostInstance().Put(pid, p); err != nil {
		silr.ResponseErrorWithRes(c, res)
		zap.L().Error("PostPut is failed", zap.Error(err))
		return
	}
	silr.ResponseSuccess(c, nil)
}

// PostDeleteHandler 删除帖子
func PostDeleteHandler(c *gee.Context) {
	p := new(models.PostDelete)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("post delete params is not illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
		return
	}
	uid, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("user token Verify fail", zap.Error(err))
		silr.ResponseError(c, e.TokenInvalidAuth)
		return
	}
	if err = service.NewPostInstance().Delete(uid, p); err != nil {
		zap.L().Error("delete post failed", zap.String("pid", p.PostId), zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// PostListHandler 顺序获取所有帖子
func PostListHandler(c *gee.Context) {
	page, size, order := getPostListInfo(c)
	data, err := service.NewPostInstance().ListInOrder(page, size, order)
	if err != nil {
		zap.L().Error("PostListInOrder select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}
