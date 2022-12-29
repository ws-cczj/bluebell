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
	CodeRepeatLogin
	CodePostVoteExpired
	// -- Token错误
	TokenFailGenerate
	TokenFailVerify
	TokenNullNeedLogin
	TokenInvalidAuth
)
