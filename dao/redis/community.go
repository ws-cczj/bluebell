package redis

import (
	"time"

	"go.uber.org/zap"
)

const CommunityNumsCache = 5 * time.Minute

// SetCommunityNums 设置社区总数缓存
func SetCommunityNums(cNums int64) error {
	return rdb.Set(addKeyPrefix(KeyCommunityNums), cNums, CommunityNumsCache).Err()
}

// GetCommunitys 获取总社区数
func GetCommunitys() (string, error) {
	return rdb.Get(addKeyPrefix(KeyCommunityNums)).Result()
}

// GetCommunityPosts 获取该社区下的帖子数
func GetCommunityPosts(cid int64) (pidNums int64, err error) {
	return rdb.SCard(addKeyPrefix(KeyCommunitySetPF, stvI64toa(cid))).Result()
}

// GetCommunityPostIds 获取社区的帖子ids
func GetCommunityPostIds(page, size, cid int64, orderkey string) (pids []string, err error) {
	key := addKeyPrefix(orderkey, stvI64toa(cid))
	ckey := addKeyPrefix(KeyCommunitySetPF, stvI64toa(cid))
	// -- 设置key缓存，减少 ZinterStore的消耗, 也避免了资源的浪费
	if rdb.Exists(key).Val() < 1 {
		pipe := rdb.TxPipeline()
		pipe.ZInterStore(key, redisZS(AggregateMAX), ckey, addKeyPrefix(orderkey))
		pipe.Expire(key, KeyZInterExpired)
		_, err = pipe.Exec()
		if err != nil {
			zap.L().Error("pipe ZInterStore or Expire exec is failed",
				zap.Error(err))
			return
		}
	}
	return GetPostIds(page, size, orderkey+stvI64toa(cid))
}
