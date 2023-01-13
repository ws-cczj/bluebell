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
	})
	return ctab
}

// GetKeysAndVals 获取当前任务列表中所有的任务信息
func (c Crontab) GetKeysAndVals() (keys []Task, vals []cron.EntryID) {
	if c.Cron != ctab.Cron || !ctab.Enable {
		return
	}
	cnt := 0
	n := len(ctab.EntryIds)
	keys = make([]Task, n, n)
	vals = make([]cron.EntryID, n, n)
	for s, id := range ctab.EntryIds {
		keys[cnt] = s
		vals[cnt] = id
		cnt++
	}
	return
}

// RunAll 开始执行所有任务
func (c Crontab) RunAll() {
	if c.Cron != ctab.Cron || ctab.Enable {
		return
	}
	ctab.Enable = true
	ctab.Cron.Start()
}

// StopAll 停止所有任务的执行
func (c Crontab) StopAll() {
	if c.Cron != ctab.Cron || !ctab.Enable {
		return
	}
	ctab.Cron.Stop()
	ctab.Enable = false
}

// RemoveTask 移除任务
func (c Crontab) RemoveTask(t Task) {
	if c.Cron != ctab.Cron || !ctab.Enable {
		return
	}
	entryIds := ctab.EntryIds
	if id, ok := entryIds[t]; ok {
		ctab.Cron.Remove(id)
		delete(entryIds, t)
		delete(ctab.PreTime, id)
	}
}
