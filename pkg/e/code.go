package e

const (
	// -- 状态码
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistIdentifier    = 202
	InvalidParams         = 400
	ERROR                 = 500

	// -- 成员错误
	ErrorExistUser          = 10002
	ErrorNotExistUser       = 10003
	ErrorFailEncryption     = 10006
	ErrorNotComparePassword = 10007

	// -- 数据库错误
	ErrorInitDatabase  = 400000
	ErrorQueryDatabase = 400001
	ErrorExecDatabase  = 400002
	ErrorNullDatabase  = 400003
)
