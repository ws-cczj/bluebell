package redis

import (
	"strconv"

	"github.com/go-redis/redis"
)

// redis Key 使用命名空间的方式，方便查询和分割
const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"   // zset;帖子及发表时间
	KeyPostScoreZSet   = "post:score"  // zset;帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" // zset;记录用户投票的参数类型;参数是 post id
	KeyCommunitySetPF  = "community:"  // set;记录每个社区中的帖子;参数是 community
)

// addKeyPrefix 添加对象前缀
func addKeyPrefix(key string) string {
	return KeyPrefix + key
}

// stvI64toa 将int64转为ascii码
func stvI64toa(id int64) string {
	return strconv.FormatInt(id, 10)
}

// redisZ 对redisZ进行封装
func redisZ(score, member any) redis.Z {
	switch score.(type) {
	case float64:
		return redis.Z{
			Score:  score.(float64),
			Member: member,
		}
	case int64:
		return redis.Z{
			Score:  float64(score.(int64)),
			Member: member,
		}
	}
	return redis.Z{}
}

// redisZS 对redisZStore进行封装
func redisZS(compare string) redis.ZStore {
	return redis.ZStore{
		Aggregate: compare,
	}
}
