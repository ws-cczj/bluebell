package crontab

import (
	"bluebell/logger"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Crontab struct {
	*cron.Cron
	Enable   bool
	PreTime  map[cron.EntryID]time.Time
	EntryIds map[Task]cron.EntryID
}

var (
	ctab    *Crontab
	ctbOnce sync.Once
)

// NewCrontabInstance 获取定时任务
func NewCrontabInstance() *Crontab {
	ctbOnce.Do(func() {
		ctab = new(Crontab)
		ctab.Cron = cron.New(cron.WithLogger(logger.Log{Logger: zap.L()}), cron.WithSeconds())
		ctab.EntryIds = make(map[Task]cron.EntryID, Caps)
		ctab.PreTime = make(map[cron.EntryID]time.Time, Caps)
		Exec(ctab)
	})
	return ctab
}

// RunAll 开始执行所有任务
func (c *Crontab) RunAll() {
	if c.Enable {
		return
	}
	c.Enable = true
	c.Cron.Start()
}

// StopAll 停止所有任务的执行
func (c *Crontab) StopAll() {
	if !c.Enable {
		return
	}
	c.Cron.Stop()
	c.Enable = false
}

// RemoveTask 移除任务
func (c *Crontab) RemoveTask(t Task) {
	if !c.Enable {
		return
	}
	entryIds := c.EntryIds
	if id, ok := entryIds[t]; ok {
		c.Cron.Remove(id)
		delete(entryIds, t)
		delete(c.PreTime, id)
	}
}
