package redis

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// redis Key 使用命名空间的方式，方便查询和分割
const (
	KeyPrefix = "bluebell:"
	// User Key
	KeyUserToken     = "username:"       // string;用户token记录 have expire_time
	KeyUserCommunity = "user_community:" // set;记录每个用户管理的社区;参数是uid
	KeyUserPostNums  = "user_post:"      // set;记录每个用户管理的帖子;参数是uid
	// Post Key
	KeyPostTimeZSet    = "post:time"   // zset;帖子及发表时间
	KeyPostScoreZSet   = "post:score"  // zset;帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" // zset;记录用户投票的参数类型;参数是 pid
	// Community Key
	KeyCommunitySetPF = "community:"     // set;记录每个社区中的帖子;PF表示前缀;参数是 cid
	KeyCommunityNums  = "community_nums" // string;记录社区总数 have expire_time
	// Comment Key
	KeyCommentFavorite  = "comment:"        // hset;记录评论的点赞;key cid; key: uid,value: to_uid(点赞者->被点赞者)
	KeyCommentTimeZSet  = "comment:time:"   // zset;根据时间记录帖子的评论id;参数是pid
	KeyCommentScoreZSet = "comment:score:"  // zset;根据分数进行排序评论id;参数是pid
	KeyCommentFather    = "comment:father:" // list;记录一个父评论存储的所有子评论; key fid -> childId...
)

// addKeyPrefix 添加对象前缀
func addKeyPrefix(keys ...string) string {
	var build strings.Builder
	build.WriteString(KeyPrefix)
	for _, key := range keys {
		build.WriteString(key)
	}
	return build.String()
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
	return redis.Z{float64(score.(int)), member}
}

// redisZS 对redisZStore进行封装
func redisZS(compare string) redis.ZStore {
	return redis.ZStore{Aggregate: compare}
}

// redisZRBy 对redisZRByScore进行封装
func redisZRBy(min, max int64) redis.ZRangeBy {
	return redis.ZRangeBy{
		Min: stvI64toa(min),
		Max: stvI64toa(max),
	}
}
