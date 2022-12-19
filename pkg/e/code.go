package e

const (
	// -- 状态码
	CodeSUCCESS ResCode = 1000 + iota
	CodeInvalidParams
	CodeServerBusy
	CodeExistUser
	CodeNotExistUser
	CodeFailEncryption
	CodeNotComparePassword
	// -- 数据库错误
	ErrorQueryDatabase
	ErrorExecDatabase
)
