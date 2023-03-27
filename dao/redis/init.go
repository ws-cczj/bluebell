package redis

import (
	"bluebell/settings"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	rdb                *redis.Client
	ErrVoteTimeExpired = errors.New("post vote time was expired")
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Conf.Rdb.Host, settings.Conf.Rdb.Port),
		Password: settings.Conf.Rdb.Password,
		DB:       settings.Conf.Rdb.Db,
		PoolSize: settings.Conf.Rdb.PoolSize, // -- 连接池大小
	})
	if _, err := rdb.Ping().Result(); err != nil {
		panic(fmt.Errorf("redis ping fail, err: %s", err))
	}
}

func Close() {
	_ = rdb.Close()
}
