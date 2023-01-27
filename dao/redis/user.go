package redis

import (
	"time"
)

const (
	RTExpiredDuration     = 610 * time.Second
	TestRTExpiredDuration = 24 * 3600 * time.Second
)

// ------Token
// SetSingleUserToken 设置redis单用户token
func SetSingleUserToken(username string, token string) (err error) {
	err = rdb.Set(addKeyPrefix(KeyUserToken, username), token, TestRTExpiredDuration).Err()
	return
}

// GetSingleUserToken 获取redis单用户token
func GetSingleUserToken(username string) (string, error) {
	return rdb.Get(addKeyPrefix(KeyUserToken, username)).Result()
}

// ------User Community
// SetUserCommunity 设置user对应的社区集合
func SetUserCommunity(cNums int64, uid string, cid int) (err error) {
	pipe := rdb.Pipeline()
	pipe.Set(addKeyPrefix(KeyCommunityNums), cNums, CommunityNumsCache)
	pipe.LPush(addKeyPrefix(KeyUserCommunity, uid), cid)
	_, err = pipe.Exec()
	return
}

// GetUserCommunitys 获取该用户管理的社区数
func GetUserCommunitys(uid string) int64 {
	return rdb.LLen(addKeyPrefix(KeyUserCommunity, uid)).Val()
}

// -----------User Post
// UserDeletePost 用户删除帖子
func UserDeletePost(uid, pid string) error {
	return rdb.LRem(addKeyPrefix(KeyUserPost, uid), 0, pid).Err()
}

// GetUserPostNums 获取该用户管理的帖子的数目
func GetUserPostNums(uid string) int64 {
	return rdb.LLen(addKeyPrefix(KeyUserPost, uid)).Val()
}

// -------------------------------------User Following
// SetUserToFollow 建立用户关注表 to_follow 去跟随 也就是关注. follow 粉丝表
func SetUserToFollow(uid, to_uid string) (err error) {
	pipe := rdb.TxPipeline()
	pipe.LPush(addKeyPrefix(KeyUserToFollow, uid), to_uid)
	pipe.LPush(addKeyPrefix(KeyUserFollow, to_uid), uid)
	_, err = pipe.Exec()
	return
}

// CancelUserToFollow 用户取消关注
func CancelUserToFollow(uid, to_uid string) (err error) {
	pipe := rdb.TxPipeline()
	pipe.LRem(addKeyPrefix(KeyUserToFollow, uid), 0, to_uid)
	pipe.LRem(addKeyPrefix(KeyUserFollow, to_uid), 0, uid)
	_, err = pipe.Exec()
	return
}

// GetUserToFollows 获取用户的关注表 range 根据时间倒序
func GetUserToFollows(uid string) ([]string, error) {
	return rdb.LRange(addKeyPrefix(KeyUserToFollow, uid), 0, -1).Result()
}

// GetUserFollows 获取用户的粉丝表 range 根据时间倒序
func GetUserFollows(uid string) ([]string, error) {
	return rdb.LRange(addKeyPrefix(KeyUserFollow, uid), 0, -1).Result()
}
