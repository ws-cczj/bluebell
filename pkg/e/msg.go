package e

var MsgFlags = map[int]string{
	SUCCESS:               "ok",
	UpdatePasswordSuccess: "修改密码成功",
	NotExistIdentifier:    "该第三方账号未绑定",
	ERROR:                 "系统繁忙，请稍后再试",
	InvalidParams:         "请求参数错误",

	ErrorFailEncryption:     "用户密码加密失败",
	ErrorExistUser:          "用户已经存在!",
	ErrorNotExistUser:       "用户不存在!",
	ErrorNotComparePassword: "用户密码不匹配",

	ErrorExecDatabase:  "数据库执行操作失败!",
	ErrorInitDatabase:  "数据库初始化失败!",
	ErrorQueryDatabase: "数据库查询失败!",
	ErrorNullDatabase:  "数据查询结果为空",
}

// -- GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
