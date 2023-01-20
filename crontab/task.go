package crontab

const (
	Caps            = 3
	ExecOnceDay     = "0 0 0 * * *"
	ExecTwiceDay    = "0 0 2,14 * * *"
	ExecSecondHour  = "0 0 */2 * * *"
	ExecTest        = "*/5 * * * * *"
	TaskVoteExpired = "帖子投票时间过期任务"
	TaskPostDelete  = "帖子被删除清除数据任务"
	TaskMonitor     = "监视任务执行情况"
	TaskLogger      = "定时任务日志:"
)

type Task string
