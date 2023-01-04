package redis

import (
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

const (
	KeyCPInterExpired = 60 * time.Second
	AggregateSUM      = "SUM"
	AggregateMAX      = "MAX"
	AggregateMIN      = "MIN"
	ZCountMIN         = "1"
	ZCountMAX         = "1"
)

// CreatePost 创建帖子的时间和初始分数
func CreatePost(pid, community int64) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(addKeyPrefix(KeyPostTimeZSet), redisZ(time.Now().Unix(), pid))
	pipeline.ZAdd(addKeyPrefix(KeyPostScoreZSet), redisZ(time.Now().Unix(), pid))
	pipeline.SAdd(addKeyPrefix(KeyCommunitySetPF+stvI64toa(community)), pid)
	_, err = pipeline.Exec()
	return
}

// DeletePost 删除帖子信息
func DeletePost(pid, community int64) (err error) {
	pipeline := rdb.Pipeline()
	pipeline.ZRem(addKeyPrefix(KeyPostTimeZSet), pid)
	pipeline.ZRem(addKeyPrefix(KeyPostScoreZSet), pid)
	pipeline.SRem(addKeyPrefix(KeyCommunitySetPF+stvI64toa(community)), pid)
	_, err = pipeline.Exec()
	return
}

// GetPostVote 获取帖子的票数
func GetPostVote(pid int64) int64 {
	return rdb.ZCount(addKeyPrefix(KeyPostVotedZSetPF+stvI64toa(pid)), "1", "1").Val()
}

// GetPostIds 根据顺序查询帖子列表
func GetPostIds(page, size int64, key string) (ids []string, err error) {
	start := (page - 1) * size
	end := start + size - 1
	return rdb.ZRevRange(addKeyPrefix(key), start, end).Result()
}

// GetPostVotes 获取帖子的票数
func GetPostVotes(ids []string) (tickets []int64, err error) {
	pipe := rdb.TxPipeline()
	tickets = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := addKeyPrefix(KeyPostVotedZSetPF + id)
		pipe.ZCount(key, ZCountMAX, ZCountMIN)
	}
	cmders, err := pipe.Exec()
	if err != nil {
		return
	}
	for _, cmder := range cmders {
		ticket := cmder.(*redis.IntCmd).Val()
		tickets = append(tickets, ticket)
	}
	return
}

// GetCommunityPostIds 获取社区的帖子ids
func GetCommunityPostIds(page, size, cid int64, orderkey string) (ids []string, err error) {
	key := addKeyPrefix(orderkey) + stvI64toa(cid)
	ckey := addKeyPrefix(KeyCommunitySetPF + stvI64toa(cid))
	// -- 设置key缓存，减少 zinterStore的消耗
	if rdb.Exists(key).Val() < 1 {
		pipe := rdb.TxPipeline()
		pipe.ZInterStore(key, redisZS(AggregateMAX), ckey, addKeyPrefix(orderkey))
		pipe.Expire(key, KeyCPInterExpired)
		_, err = pipe.Exec()
		if err != nil {
			zap.L().Error("pipe ZInterStore or Expire exec is failed",
				zap.Error(err))
			return
		}
	}
	return GetPostIds(page, size, orderkey+stvI64toa(cid))
}
