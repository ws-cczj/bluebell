package api

import (
	"bluebell/models"
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"bluebell/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityCreateHandler 创建社区
func CommunityCreateHandler(c *gin.Context) {
	community := new(models.CommunityDetail)
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
	community.Community.AuthorId = uid
	community.Community.AuthorName = uname
	if err = service.CommunityCreate(community); err != nil {
		zap.L().Error("CommunityCreate method err",
			zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, nil)
}

// CommunityHandler 返回社区的列表信息
func CommunityHandler(c *gin.Context) {
	data, err := service.CommunityList()
	if err != nil {
		zap.L().Error("CommunityList select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// CommunityDetailHandler 通过ID获取到详细的社区情况
func CommunityDetailHandler(c *gin.Context) {
	cid, err := getParamId(c, "cid")
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		silr.ResponseError(c, e.CodeInvalidParams)
		return
	}
	data, err := service.CommunityDetailById(cid)
	if err != nil {
		zap.L().Error("CommunityDetailById select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}

// CommunityPostHandler 获取该社区的所有帖子
func CommunityPostHandler(c *gin.Context) {
	page, size, order := getPostListInfo(c)
	cid, err := getParamId(c, "cid")
	if err != nil {
		zap.L().Error("Param is illgal", zap.Error(err))
		silr.ResponseError(c, e.CodeInvalidParams)
		return
	}
	data, err := service.CommunityPostListInOrder(page, size, cid, order)
	if err != nil {
		zap.L().Error("CommunityPostListInOrder select data is failed", zap.Error(err))
		silr.ResponseError(c, e.CodeServerBusy)
		return
	}
	silr.ResponseSuccess(c, data)
}
