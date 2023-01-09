package service

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"strconv"
)

// 开启定时任务，将redis中投票时间已经过期的帖子转移到数据库中，并通过数据库进行查询

// CrontabPostExpired 定时处理过期帖子
func CrontabPostExpired(preT, now int64) error {
	// 1. 获取已经过期的帖子ids
	ids, err := redis.GetPostExpired(preT, now)
	if err != nil {
		return err
	}
	// 2. 根据已经过期的帖子ids查找到对应的票数
	tickets, err := redis.GetPostVotes(ids)
	if err != nil {
		return err
	}
	// 3. 将票数和过期状态更新到mysql 将过期的帖子进行删除
	var cid int64
	for i, pidStr := range ids {
		pid, _ := strconv.ParseInt(pidStr, 10, 64)
		err = mysql.UpdateCtbPost(pid, tickets[i])
		if err != nil {
			continue
		}
		cid, err = mysql.FindCidByPid(pid)
		if err != nil {
			continue
		}
		err = redis.DeletePost(pid, cid)
		if err != nil {
			continue
		}
	}
	return err
}
