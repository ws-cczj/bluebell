package api

import (
	"bluebell/pkg/e"
	silr "bluebell/serializer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseError 错误响应体
func ResponseError(c *gin.Context, code e.ResCode) {
	c.JSON(http.StatusServiceUnavailable, silr.Response{
		Status: code,
		Data:   nil,
		Msg:    code.Msg(),
	})
}

// ResponseErrorWithCode 错误响应体自指定code
func ResponseErrorWithCode(c *gin.Context, htc int, code e.ResCode) {
	c.JSON(htc, silr.Response{
		Status: code,
		Data:   nil,
		Msg:    code.Msg(),
	})
}

// ResponseErrorWithRes 带有响应体的响应体
func ResponseErrorWithRes(c *gin.Context, res silr.Response) {
	c.JSON(http.StatusServiceUnavailable, res)
}

// ResponseErrorWithMsg 带有消息的错误响应体
func ResponseErrorWithMsg(c *gin.Context, code e.ResCode, Msg interface{}) {
	c.JSON(int(code), silr.Response{
		Status: code,
		Data:   nil,
		Msg:    Msg,
	})
}

// ResponseSuccess 响应成功
func ResponseSuccess(c *gin.Context, data interface{}) {
	code := e.CodeSUCCESS
	c.JSON(http.StatusOK, silr.Response{
		Status: code,
		Data:   data,
		Msg:    code.Msg(),
	})
}
