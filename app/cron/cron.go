package cron

import (
	"github.com/robfig/cron/v3"
	"goweb/app/cron/jobs"
)

// 所有计时任务在这注册
func commands() []JobCommands {
	return []JobCommands{
		&jobs.QuoteZdNotice{},
		&jobs.QuoteCommonNotice{},
		&jobs.QuoteNearCloseNotice{},
		&jobs.QuoteNearOpenNotice{},
		//&jobs.QuoteDailyIncomeNotice{},
	}
}

// 计时任务启动
func Start() {
	// 使支持秒级别
	c := cron.New(cron.WithSeconds())
	commands := commands()
	for _, job := range commands {
		times := job.GetTime()
		for _, time := range times {
			_, _ = c.AddFunc(time, job.GetHandler())
		}
	}
	// TODO:建立 recover 机制
	c.Start()
}

// 计时任务接口
type JobCommands interface {
	GetTime() []string
	GetName() string
	GetHandler() func()
}
