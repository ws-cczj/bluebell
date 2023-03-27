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
	KeyUserCommunity = "user_community:" // list;记录每个用户管理的社区;key 是uid,value: cid
	KeyUserPost      = "user_post:"      // list;记录每个用户管理的帖子;key 是uid,value: pid
	// 可以根据集合的并集来看出互相关注的关系 follow 粉丝表, to_follow 关注表
	KeyUserFollow   = "user_follow:"    // list;记录每个用户被谁关注了 key to_uid, value: uid
	KeyUserToFollow = "user_to_follow:" // list;记录每个用户关注了谁 key: uid, value: to_uid
	// Post Key
	KeyPostTimeZSet    = "post:time"   // zset;帖子及发表时间
	KeyPostScoreZSet   = "post:score"  // zset;帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" // zset;记录用户投票的参数类型;参数是 pid
	// Community Key
	KeyCommunitySetPF = "community:"     // set;记录每个社区中的帖子;PF表示前缀;key: cid value: pid
	KeyCommunityNums  = "community_nums" // string;记录社区总数 have expire_time
	// Comment Key
	KeyCommentFavorite  = "comment:"        // hset;记录评论的点赞;key commentId; key: uid,value: to_uid(点赞者->被点赞者)
	KeyCommentTimeZSet  = "comment:time:"   // zset;根据时间记录帖子的评论id;key是pid value: fid和time
	KeyCommentScoreZSet = "comment:score:"  // zset;根据分数进行排序评论id;key是pid, value: fid和score
	KeyCommentFather    = "comment:father:" // list;记录一个父评论存储的所有子评论; key fid -> childId...
)

// addKeyPrefix 添加对象前缀
func addKeyPrefix(keys ...string) string {
	var build strings.Builder
	build.Grow(len(keys) + 1)
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

func stvItoa(id int) string {
	return strconv.Itoa(id)
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
