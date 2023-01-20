package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"time"

	"go.uber.org/zap"
)

// CrontabDeleteComment 定实清除剩余redis评论信息
func CrontabDeleteComment(preT, now time.Time) error {
	// TODO 比对这段时间中的post状态为4的帖子,如果redis中这种帖子还存在
	pids, err := mysql.CrontabPostDelete(preT, now)
	if err != nil {
		zap.L().Error("Crontab mysql CrontabPostDelete method err", zap.Error(err))
		return err
	}
	// 删除mysql中没有父评论的子评论
	if err = mysql.CrontabDeleteComment(); err != nil {
		zap.L().Error("Crontab mysql CheckComment method err", zap.Error(err))
		return err
	}
	// 循环删除帖子评论信息 如果数据为空不会返回错误!
	for _, pid := range pids {
		err = redis.PostCommentDelete(pid)
	}
	return err
}
