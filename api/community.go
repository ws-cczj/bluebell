package api

import (
	"bluebell/models"
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"
	"strconv"

	"github.com/ws-cczj/gee"

	"go.uber.org/zap"
)

// CommunityCreateHandler 创建社区
func CommunityCreateHandler(c *gee.Context) {
	community := models.NewCommunityDetail()
	if err := c.ShouldBind(community); err != nil {
		zap.L().Error("Community Create params is not illegal", zap.Error(err))
		silr.ResponseValidatorError(c, err)
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
	community.AuthorId = uid
	community.AuthorName = uname
	if err = service.NewCommunityInstance().Create(community); err != nil {
		zap.L().Error("CommunityCreate method err",
			zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// CommunityHandler 返回社区的列表信息
func CommunityHandler(c *gee.Context) {
	data, err := service.NewCommunityInstance().List()
	if err != nil {
		zap.L().Error("CommunityList select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// CommunityDetailHandler 通过ID获取到详细的社区情况
func CommunityDetailHandler(c *gee.Context) {
	cidStr := c.Param("cid")
	cid, err := strconv.ParseInt(cidStr, 10, 32)
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		silr.ResponseError(c, e.CodeInvalidParams)
		return
	}
	data, err := service.NewCommunityInstance().DetailById(int(cid))
	if err != nil {
		zap.L().Error("CommunityDetailById select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// CommunityPostHandler 获取该社区的所有帖子
func CommunityPostHandler(c *gee.Context) {
	page, size, order := getPostListInfo(c)
	cid := c.Param("cid")
	data, err := service.NewCommunityInstance().PostListInOrder(page, size, cid, order)
	if err != nil {
		zap.L().Error("CommunityPostListInOrder select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}
