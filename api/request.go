package api

import (
	"errors"
	"strconv"

	"github.com/ws-cczj/gee"
)

const ContextUserIDKey = "user_id"
const ContextUsernameKey = "username"

var ErrorUserNotLogin = errors.New("用户还未登录")

// getCurrentUserId 获取当前的userId
func getCurrentUserId(c *gee.Context) (userID string, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(string)
	if !ok {
		err = ErrorUserNotLogin
	}
	return
}

// getCurrentUserId 获取当前用户的username
func getCurrentUsername(c *gee.Context) (username string, err error) {
	uname, ok := c.Get(ContextUsernameKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	username, ok = uname.(string)
	if !ok {
		err = ErrorUserNotLogin
	}
	return
}

// getPostListInfo 获取帖子的参数信息
func getPostListInfo(c *gee.Context) (page, size int64, order string) {
	var err error
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	order = c.Query("order")
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return
}
