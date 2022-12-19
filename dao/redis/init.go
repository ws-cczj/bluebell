package redis

import (
	"bluebell/settings"
	"fmt"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var rdb *redis.Client

func InitRedis(cfg *settings.RedisConfig) error {
	red := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize, // -- 连接池大小
	})
	_, err := red.Ping().Result()
	if err != nil {
		zap.L().Error("redis ping is fail", zap.Error(err))
		return err
	}
	rdb = red
	return nil
}

func Close() {
	_ = rdb.Close()
}
