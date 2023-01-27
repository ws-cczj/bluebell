package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

const secret = "cczjblog.top"

// Md5Password 对password进行md5加密处理
func Md5Password(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
