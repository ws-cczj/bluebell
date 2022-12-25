package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户还未登录")

// getCurrentUser 获取当前的userID
func getCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
