package crontab

import (
	"bluebell/pkg/logger"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Crontab struct {
	*cron.Cron
	PreTime  map[cron.EntryID]time.Time
	EntryIds map[Task]cron.EntryID
	enable   bool
}

var (
	ctab    *Crontab
	ctbOnce sync.Once
)

// NewCrontabInstance 获取定时任务
func NewCrontabInstance() *Crontab {
	ctbOnce.Do(func() {
		ctab = new(Crontab)
		ctab.Cron = cron.New(cron.WithLogger(&logger.Log{Logger: zap.L()}), cron.WithSeconds())
		ctab.EntryIds = make(map[Task]cron.EntryID, Caps)
		ctab.PreTime = make(map[cron.EntryID]time.Time, Caps)
		ctab.RunAll()
	})
	return ctab
}

// RunAll 开始执行所有任务
func (c *Crontab) RunAll() {
	if c.enable {
		return
	}
	c.enable = true
	Task(TaskPostDelete).PostTask()
	Task(TaskVoteExpired).CommentTask()
	Task(TaskMonitor).MonitorTask()
	c.Cron.Start()
}

// StopAll 停止所有任务的执行
func (c *Crontab) StopAll() {
	if c.enable {
		c.Cron.Stop()
	}
}

// RemoveTask 移除任务
func (c *Crontab) RemoveTask(t Task) {
	entryIds := c.EntryIds
	if id, ok := entryIds[t]; ok {
		c.Cron.Remove(id)
		delete(entryIds, t)
		delete(c.PreTime, id)
	}
}
