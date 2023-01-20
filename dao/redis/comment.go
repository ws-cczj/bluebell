package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// ---------------------------Comment Info
// CreatePostComment 创建帖子评论的相关信息
func CreatePostComment(pid, commentId int64) (err error) {
	pidStr := stvI64toa(pid)
	pipe := rdb.TxPipeline()
	pipe.ZAdd(addKeyPrefix(KeyCommentTimeZSet, pidStr), redisZ(time.Now().Unix(), commentId))
	pipe.ZAdd(addKeyPrefix(KeyCommentScoreZSet, pidStr), redisZ(time.Now().Unix(), commentId))
	_, err = pipe.Exec()
	return
}

// DeleteFatherComment 删除一条父评论与其所有的子评论
func DeleteFatherComment(pid, commentId int64) (err error) {
	pidStr := stvI64toa(pid)
	pipe := rdb.Pipeline()
	pipe.ZRem(addKeyPrefix(KeyCommentTimeZSet, pidStr), commentId)
	pipe.ZRem(addKeyPrefix(KeyCommentScoreZSet, pidStr), commentId)
	pipe.LTrim(addKeyPrefix(KeyCommentFather, stvI64toa(commentId)), 1, 0)
	_, err = pipe.Exec()
	return
}

// CreateChildComment 创建子评论信息
func CreateChildComment(fCommentId, cCommentId int64) (err error) {
	return rdb.LPush(addKeyPrefix(KeyCommentFather, stvI64toa(fCommentId)), cCommentId).Err()
}

// DeleteChildComment 删除子评论
func DeleteChildComment(fCommentId, cCommentId int64) (err error) {
	return rdb.LRem(addKeyPrefix(KeyCommentFather, stvI64toa(fCommentId)), 0, cCommentId).Err()
}

// GetFatherCommentId 获取该帖子的所有父评论id
func GetFatherCommentId(pid int64, key string) (fList []string, err error) {
	return rdb.ZRevRange(addKeyPrefix(key, stvI64toa(pid)), 0, -1).Result()
}

// GetChildCommentId 获取该父评论的所有子评论
func GetChildCommentId(fCommentId int64) ([]string, error) {
	return rdb.LRange(addKeyPrefix(KeyCommentFather, stvI64toa(fCommentId)), 0, -1).Result()
}

// GetAllCommentId 获取该帖子的所有父评论的所有子评论
func GetAllCommentId(fList []string) (cList [][]string, err error) {
	pipe := rdb.Pipeline()
	// lpush 最先进入的最晚出来，对应时间，最早发布的最晚被遍历
	for _, fid := range fList {
		pipe.LRange(addKeyPrefix(KeyCommentFather, fid), 0, -1)
	}
	cmders, err := pipe.Exec()
	if err != nil {
		return
	}
	cList = make([][]string, 0, len(fList))
	for _, cmder := range cmders {
		ids := cmder.(*redis.StringSliceCmd).Val()
		cList = append(cList, ids)
	}
	return
}

// -------------------------Comment Favorite
// CreateFFavorite 创建父评论点赞
func CreateFFavorite(pid, uid, to_uid, commentId int64) (err error) {
	pipe := rdb.TxPipeline()
	pipe.HSet(addKeyPrefix(KeyCommentFavorite, stvI64toa(commentId)), stvI64toa(uid), to_uid)
	pipe.ZIncrBy(addKeyPrefix(KeyCommentScoreZSet, stvI64toa(pid)), OneFavoriteScore, stvI64toa(commentId))
	_, err = pipe.Exec()
	return
}

// DeleteFFavorite 移除父评论点赞
func DeleteFFavorite(pid, uid, commentId int64) (err error) {
	pipe := rdb.Pipeline()
	pipe.HDel(addKeyPrefix(KeyCommentFavorite, stvI64toa(commentId)), stvI64toa(uid))
	pipe.ZIncrBy(addKeyPrefix(KeyCommentScoreZSet, stvI64toa(pid)), -OneFavoriteScore, stvI64toa(commentId))
	_, err = pipe.Exec()
	return
}

// CreateCFavorite 创建子评论点赞
func CreateCFavorite(uid, to_uid, commentId int64) error {
	return rdb.HSet(addKeyPrefix(KeyCommentFavorite, stvI64toa(commentId)), stvI64toa(uid), to_uid).Err()
}

// DeleteCFavorite 移除子评论点赞
func DeleteCFavorite(uid, commentId int64) error {
	return rdb.HDel(addKeyPrefix(KeyCommentFavorite, stvI64toa(commentId)), stvI64toa(uid)).Err()
}

// GetFavorites 获取点赞数
func GetFavorites(commentId int64) int64 {
	return rdb.HLen(addKeyPrefix(KeyCommentFavorite, stvI64toa(commentId))).Val()
}

// GetFavoriteList 获取点赞数集合
func GetFavoriteList(commentId []string) (favorites []int64, err error) {
	pipe := rdb.Pipeline()
	for _, id := range commentId {
		pipe.HLen(addKeyPrefix(KeyCommentFavorite, id)).Val()
	}
	cmders, err := pipe.Exec()
	if err != nil {
		return
	}
	favorites = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		favorite := cmder.(*redis.IntCmd).Val()
		favorites = append(favorites, favorite)
	}
	return
}
