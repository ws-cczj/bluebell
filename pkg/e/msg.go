package e

type ResCode int64

var CodeMsgMap = map[ResCode]string{
	CodeSUCCESS:            "success",
	CodeServerBusy:         "系统繁忙，请稍后再试",
	CodeInvalidParams:      "请求参数错误",
	CodeFailEncryption:     "用户密码加密失败",
	CodeExistUser:          "用户已经存在!",
	CodeNotExistUser:       "用户不存在!",
	CodeNotComparePassword: "用户密码不匹配",

	ErrorExecDatabase:  "数据库执行操作失败!",
	ErrorQueryDatabase: "数据库查询失败!",
}

// Msg 获取状态码对应信息
func (c ResCode) Msg() string {
	msg, ok := CodeMsgMap[c]
	if ok {
		return msg
	}
	return CodeMsgMap[CodeServerBusy]
}
