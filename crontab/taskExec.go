package crontab

import (
	"bluebell/service"
	"time"

	"go.uber.org/zap"
)

// MysqlTask mysql任务
func MysqlTask() {
	id, _ := ctab.Cron.AddFunc(ExecOnceDay, func() {
		id := ctab.EntryIds[TaskVoteExpired]
		preT := ctab.PreTime[id]
		now := time.Now()
		err := service.CrontabPostExpired(preT.Unix()-1, now.Unix())
		if err != nil {
			zap.L().Error("crontab service CrontabPostExpired method err",
				zap.Int("crontab ID", int(id)),
				zap.Error(err))
			ctab.PreTime[id] = preT
			return
		}
		ctab.PreTime[id] = now
	})
	ctab.PreTime[id] = time.Now()
	ctab.EntryIds[TaskVoteExpired] = id
	zap.L().Debug(TaskVoteExpired + "开启!")
}

// MonitorTask 监视执行任务
func MonitorTask() {
	id, _ := ctab.Cron.AddFunc(ExecSecondHour, func() {
		for task, id := range ctab.EntryIds {
			entry := ctab.Entry(id)
			zap.L().Info(TaskLogger,
				zap.String("[执行]:", string(task)),
				zap.Int("[Id]:", int(id)),
				zap.Time("[上次执行时间]:", ctab.PreTime[id]),
				zap.Time("[下次执行时间]:", entry.Next),
			)
		}
		ctab.PreTime[ctab.EntryIds[TaskMonitor]] = time.Now()
	})
	ctab.PreTime[id] = time.Now()
	ctab.EntryIds[TaskMonitor] = id
	zap.L().Debug(TaskMonitor + "开启!")
}
