package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// CreatePost 创建帖子的时间和初始分数
func CreatePost(pid int64) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(addKeyPrefix(KeyPostTimeZest), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid,
	})
	pipeline.ZAdd(addKeyPrefix(KeyPostScoreZest), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid,
	})
	_, err = pipeline.Exec()
	return
}
