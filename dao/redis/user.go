package redis

import (
	"time"
)

const RTExpiredDuration = 3 * time.Hour

// SetSingleUserToken 设置redis单用户token
func SetSingleUserToken(username string, token string) (err error) {
	err = rdb.Set(username, token, RTExpiredDuration).Err()
	return
}

// GetSingleUserToken 获取redis单用户token
func GetSingleUserToken(username string) (string, error) {
	return rdb.Get(username).Result()
}
