package serializer

import (
	"bluebell/pkg/e"
)

// Response 基础序列化器
type Response struct {
	Status e.ResCode   `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Msg    interface{} `json:"msg"`
}
