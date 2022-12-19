package serializer

// -- Response 基础序列化器
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Msg    interface{} `json:"msg,omitempty"`
}
