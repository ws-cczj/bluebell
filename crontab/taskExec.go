package crontab

import (
	"bluebell/service"
	"time"

	"go.uber.org/zap"
)

const (
	Caps            = 3
	ExecOnceDay     = "0 0 0 * * *"
	ExecTwiceDay    = "0 0 2,14 * * *"
	ExecSecondHour  = "0 0 */2 * * *"
	ExecTest        = "*/5 * * * * *"
	TaskVoteExpired = "帖子投票时间过期任务"
	TaskPostDelete  = "帖子被删除清除数据任务"
	TaskMonitor     = "监视任务执行情况"
)

type Task string

// PostTask  post任务 清除在Redis中已经过期的帖子缓存，将数据移入数据库中
func (t Task) PostTask() {
	id, _ := ctab.Cron.AddFunc(ExecTwiceDay, func() {
		id := ctab.EntryIds[t]
		preT := ctab.PreTime[id]
		now := time.Now()
		if err := service.NewCrontabPostInstance().ExpiredHandle(preT.Unix()-1, now.Unix()); err != nil {
			zap.L().Error("crontab service CrontabPostExpired method err",
				zap.Int("crontab ID", int(id)),
				zap.Error(err))
			ctab.PreTime[id] = preT
			return
		}
		ctab.PreTime[id] = now
	})
	ctab.PreTime[id] = time.Now()
	ctab.EntryIds[t] = id
	zap.L().Debug(string(t) + "开启!")
}

// CommentTask comment任务 清除被删除掉的视频后遗留的评论
func (t Task) CommentTask() {
	id, _ := ctab.Cron.AddFunc(ExecOnceDay, func() {
		id := ctab.EntryIds[t]
		preT := ctab.PreTime[id]
		now := time.Now()
		if err := service.NewCrontabCommentInstance().Clear(preT, now); err != nil {
			zap.L().Error("crontab service CrontabDeleteComment method err",
				zap.Int("crontab ID", int(id)),
				zap.Error(err))
			ctab.PreTime[id] = preT
			return
		}
		ctab.PreTime[id] = now
	})
	ctab.PreTime[id] = time.Now()
	ctab.EntryIds[t] = id
	zap.L().Debug(string(t) + "开启!")
}

// MonitorTask 监视执行任务
func (t Task) MonitorTask() {
	id, _ := ctab.Cron.AddFunc(ExecSecondHour, func() {
		for task, id := range ctab.EntryIds {
			entry := ctab.Entry(id)
			zap.L().Info("定时任务日志:",
				zap.String("[执行]:", string(task)),
				zap.Int("[Id]:", int(id)),
				zap.Time("[上次执行时间]:", ctab.PreTime[id]),
				zap.Time("[下次执行时间]:", entry.Next),
			)
		}
		ctab.PreTime[ctab.EntryIds[TaskMonitor]] = time.Now()
	})
	ctab.PreTime[id] = time.Now()
	ctab.EntryIds[t] = id
	zap.L().Debug(string(t) + "开启!")
}
