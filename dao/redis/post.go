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
func PostCreate(uid, pid, cid string) (err error) {
	pipe := rdb.TxPipeline()
	pipe.ZAdd(addKeyPrefix(KeyPostTimeZSet), redisZ(time.Now().Unix(), pid))
	pipe.ZAdd(addKeyPrefix(KeyPostScoreZSet), redisZ(time.Now().Unix(), pid))
	pipe.SAdd(addKeyPrefix(KeyCommunitySetPF, cid), pid)
	pipe.LPush(addKeyPrefix(KeyUserPost, uid), pid)
	_, err = pipe.Exec()
	return
}

// PostDelete 对于还未过期帖子的删除 状态为 4
func PostDelete(uid, pid, cid string) (err error) {
	pipe := rdb.Pipeline()
	// 删除社区集合中的帖子
	pipe.SRem(addKeyPrefix(KeyCommunitySetPF, cid), pid)
	// 删除帖子投票记录以及帖子信息
	pipe.ZRem(addKeyPrefix(KeyPostTimeZSet), pid)
	pipe.ZRem(addKeyPrefix(KeyPostScoreZSet), pid)
	pipe.ZRemRangeByScore(addKeyPrefix(KeyPostVotedZSetPF, pid), ZRangeMinINF, ZRangeMaxINF)
	// 删除用户列表中的帖子
	pipe.LRem(addKeyPrefix(KeyUserPost, uid), 0, pid)
	_, err = pipe.Exec()
	return
}

// PostExpiredDelete 对于已经过期帖子的删除
func PostExpiredDelete(uid, pid, cid string) (err error) {
	pipe := rdb.TxPipeline()
	pipe.SRem(addKeyPrefix(KeyCommunitySetPF, cid), pid)
	pipe.ZRem(addKeyPrefix(KeyPostTimeZSet), pid)
	pipe.ZRem(addKeyPrefix(KeyPostScoreZSet), pid)
	pipe.LRem(addKeyPrefix(KeyUserPost, uid), 0, pid)
	_, err = pipe.Exec()
	return
}

// PostExpire 处理过期帖子信息 状态为 3,只清除对应的票数
func PostExpire(pids []string) (err error) {
	pipe := rdb.Pipeline()
	for _, pid := range pids {
		pipe.ZRemRangeByScore(addKeyPrefix(KeyPostVotedZSetPF, pid), ZRangeMinINF, ZRangeMaxINF)
	}
	_, err = pipe.Exec()
	return
}

// PostCommentDelete 删除帖子的所有评论信息
func PostCommentDelete(pid string) (err error) {
	// hash 删除过于麻烦，暂时不删除
	keyT := addKeyPrefix(KeyCommentTimeZSet, pid)
	keyS := addKeyPrefix(KeyCommentScoreZSet, pid)
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
func GetPostVote(pid string) int64 {
	return rdb.ZCount(addKeyPrefix(KeyPostVotedZSetPF, pid), ZCountMAX, ZCountMIN).Val()
}

// GetPostIds 根据顺序查询帖子列表
func GetPostIds(page, size int64, key string) (ids []string, err error) {
	start := (page - 1) * size
	end := start + size - 1
	return rdb.ZRevRange(addKeyPrefix(key), start, end).Result()
}

// GetPostVotes 获取帖子的票数
func GetPostVotes(pids []string) (tickets []int, err error) {
	pipe := rdb.TxPipeline()
	tickets = make([]int, 0, len(pids))
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
		tickets = append(tickets, int(ticket))
	}
	return
}

// GetPostExpired 获取已经过期的帖子集合
func GetPostExpired(min, max int64) (pids []string, err error) {
	return rdb.ZRangeByScore(addKeyPrefix(KeyPostTimeZSet),
		redisZRBy(min-OneWeekPostTime, max-OneWeekPostTime)).Result()
}
