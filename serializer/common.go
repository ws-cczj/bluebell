package serializer

import (
	"bluebell/pkg/e"
	"net/http"

	"github.com/ws-cczj/gee"

	"github.com/go-playground/validator/v10"
)

// Response 基础序列化器
type Response struct {
	Status e.ResCode   `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Msg    interface{} `json:"msg"`
}

// ResponseError 错误响应体
func ResponseError(c *gee.Context, code e.ResCode) {
	c.JSON(http.StatusServiceUnavailable, Response{
		Status: code,
		Msg:    code.Msg(),
	})
}

// ResponseErrorWithRes 带有响应体的响应体
func ResponseErrorWithRes(c *gee.Context, res Response) {
	c.JSON(http.StatusServiceUnavailable, res)
}

// ResponseSuccess 响应成功
func ResponseSuccess(c *gee.Context, data interface{}) {
	code := e.CodeSUCCESS
	switch data.(type) {
	case ResponseUserLogin:
		c.JSON(http.StatusOK, data)
	case ResponseUserFollow:
		c.JSON(http.StatusOK, data)
	default:
		c.JSON(http.StatusOK, Response{
			Status: code,
			Data:   data,
			Msg:    code.Msg(),
		})
	}
}

// ResponseValidatorError 处理翻译器错误请求
func ResponseValidatorError(c *gee.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ResponseError(c, e.CodeInvalidParams)
		return
	}
	c.JSON(http.StatusBadRequest, Response{
		Status: http.StatusBadRequest,
		Data:   nil,
		Msg:    removeTopStruct(errs.Translate(trans)),
	})
}
