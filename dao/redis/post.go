package redis

import (
	"time"

	"github.com/go-redis/redis"
)

const (
	KeyZInterExpired = time.Minute // ZinterStore 联合查询的生成临时缓存的过期时间
	AggregateMAX     = "MAX"
	ZRangeMinINF     = "-INF"
	ZRangeMaxINF     = "INF"
	ZCountMIN        = "1"
	ZCountMAX        = "1"
)

// PostCreate 创建帖子的时间和初始分数
func PostCreate(uid, pid, cid int64) (err error) {
	pipe := rdb.TxPipeline()
	pipe.ZAdd(addKeyPrefix(KeyPostTimeZSet), redisZ(time.Now().Unix(), pid))
	pipe.ZAdd(addKeyPrefix(KeyPostScoreZSet), redisZ(time.Now().Unix(), pid))
	pipe.SAdd(addKeyPrefix(KeyCommunitySetPF, stvI64toa(cid)), pid)
	pipe.LPush(addKeyPrefix(KeyUserPost, stvI64toa(uid)), pid)
	_, err = pipe.Exec()
	return
}

// PostDelete 删除帖子的所有信息 状态为 4
func PostDelete(uid, pid, cid int64) (err error) {
	pipe := rdb.Pipeline()
	pipe.ZRem(addKeyPrefix(KeyPostTimeZSet), pid)
	pipe.ZRem(addKeyPrefix(KeyPostScoreZSet), pid)
	// 删除帖子投票记录
	pipe.ZRemRangeByScore(addKeyPrefix(KeyPostVotedZSetPF, stvI64toa(pid)), ZRangeMinINF, ZRangeMaxINF)
	// 删除社区集合中的帖子
	pipe.SRem(addKeyPrefix(KeyCommunitySetPF, stvI64toa(cid)), pid)
	// 删除用户列表中的帖子
	pipe.LRem(addKeyPrefix(KeyUserPost, stvI64toa(uid)), 0, pid)
	_, err = pipe.Exec()
	return
}

// PostExpire 处理过期帖子信息 状态为 3
func PostExpire(pids []string) (err error) {
	keyT := addKeyPrefix(KeyPostTimeZSet)
	keyS := addKeyPrefix(KeyPostScoreZSet)
	pipe := rdb.Pipeline()
	for _, pid := range pids {
		pipe.ZRem(keyT, pid)
		pipe.ZRem(keyS, pid)
		pipe.ZRemRangeByScore(addKeyPrefix(KeyPostVotedZSetPF, pid), ZRangeMinINF, ZRangeMaxINF)
	}
	_, err = pipe.Exec()
	return
}

// PostCommentDelete 删除帖子的所有评论信息
func PostCommentDelete(pid int64) (err error) {
	// hash 删除过于麻烦，暂时不删除
	keyT := addKeyPrefix(KeyCommentTimeZSet, stvI64toa(pid))
	keyS := addKeyPrefix(KeyCommentScoreZSet, stvI64toa(pid))
	fids := rdb.ZRevRange(keyT, 0, -1).Val()
	pipe := rdb.Pipeline()
	// 删除所有父评论中的子评论
	for _, fid := range fids {
		pipe.LTrim(addKeyPrefix(KeyCommentFather, fid), 1, 0)
	}
	// 删除帖子的评论排序数据
	pipe.ZRemRangeByScore(keyT, ZRangeMinINF, ZRangeMaxINF)
	pipe.ZRemRangeByScore(keyS, ZRangeMinINF, ZRangeMaxINF)
	_, err = pipe.Exec()
	return
}

// GetPostVote 获取帖子的票数
func GetPostVote(pid int64) int64 {
	return rdb.ZCount(addKeyPrefix(KeyPostVotedZSetPF, stvI64toa(pid)), ZCountMAX, ZCountMIN).Val()
}

// GetPostIds 根据顺序查询帖子列表
func GetPostIds(page, size int64, key string) (ids []string, err error) {
	start := (page - 1) * size
	end := start + size - 1
	return rdb.ZRevRange(addKeyPrefix(key), start, end).Result()
}

// GetPostVotes 获取帖子的票数
func GetPostVotes(pids []string) (tickets []uint32, err error) {
	pipe := rdb.TxPipeline()
	tickets = make([]uint32, 0, len(pids))
	for _, pid := range pids {
		key := addKeyPrefix(KeyPostVotedZSetPF, pid)
		pipe.ZCount(key, ZCountMAX, ZCountMIN)
	}
	cmders, err := pipe.Exec()
	if err != nil {
		return
	}
	for _, cmder := range cmders {
		ticket := cmder.(*redis.IntCmd).Val()
		tickets = append(tickets, uint32(ticket))
	}
	return
}

// GetPostExpired 获取已经过期的帖子集合
func GetPostExpired(min, max int64) (pids []string, err error) {
	return rdb.ZRangeByScore(addKeyPrefix(KeyPostTimeZSet),
		redisZRBy(min-OneWeekPostTime, max-OneWeekPostTime)).Result()
}
