package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"strconv"

	"go.uber.org/zap"
)

// 开启定时任务，将redis中投票时间已经过期的帖子转移到数据库中，并通过数据库进行查询

// CrontabPostExpired 定时处理过期帖子
func CrontabPostExpired(preT, now int64) error {
	// 1. 获取已经过期的帖子ids
	ids, err := redis.GetPostExpired(preT, now)
	if err != nil {
		zap.L().Error("redis getPostExpired method err", zap.Error(err))
		return err
	}
	// 2. 根据已经过期的帖子ids查找到对应的票数
	tickets, err := redis.GetPostVotes(ids)
	if err != nil {
		zap.L().Error("redis getPostVotes method err", zap.Error(err))
		return err
	}
	// 3. 根据已过期的帖子ids更新redis中帖子的有关信息,不包括评论和社区、用户管理
	if err = redis.PostExpire(ids); err != nil {
		zap.L().Error("redis PostExpire method err", zap.Error(err))
		return err
	}
	// 4. 将票数和过期状态更新到mysql 将过期的帖子进行删除
	for i, pidStr := range ids {
		pid, _ := strconv.ParseInt(pidStr, 10, 64)
		err = mysql.UpdateCtbPost(pid, tickets[i])
	}
	return err
}
