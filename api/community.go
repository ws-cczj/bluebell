package api

import (
	"bluebell/pkg/e"
	"bluebell/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 返回社区的列表信息
func CommunityHandler(c *gin.Context) {
	data, err := service.CommunityList()
	if err != nil {
		zap.L().Error("CommunityList select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 通过ID获取到详细的社区情况
func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("ParseInt data is invalid", zap.Error(err))
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	data, err := service.CommunityDetailById(id)
	if err != nil {
		zap.L().Error("CommunityDetailById select data is failed", zap.Error(err))
		ResponseError(c, e.CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
