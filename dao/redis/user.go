package redis

import (
	"time"
)

const RTExpiredDuration = 610 * time.Second
const TestRTExpiredDuration = 24 * 3600 * time.Second

// SetSingleUserToken 设置redis单用户token
func SetSingleUserToken(username string, token string) (err error) {
	err = rdb.Set(addKeyPrefix(KeyUserToken, username), token, TestRTExpiredDuration).Err()
	return
}

// GetSingleUserToken 获取redis单用户token
func GetSingleUserToken(username string) (string, error) {
	return rdb.Get(addKeyPrefix(KeyUserToken, username)).Result()
}

// SetUserCommunity 设置user对应的社区集合
func SetUserCommunity(uid, cid int64) error {
	return rdb.SAdd(addKeyPrefix(KeyUserCommunity, stvI64toa(uid)), cid).Err()
}

// GetUserCommunitys 获取该用户管理的社区数
func GetUserCommunitys(uid int64) (pidNums int64, err error) {
	return rdb.SCard(addKeyPrefix(KeyUserCommunity, stvI64toa(uid))).Result()
}
