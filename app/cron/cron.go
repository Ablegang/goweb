package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"goweb/app/cron/jobs"
)

// 所有计时任务在这注册
func commands() []JobCommands {
	return []JobCommands{
		&jobs.QuoteNotice{},
	}
}

func Start() {
	// 使支持秒级别
	c := cron.New(cron.WithSeconds())
	for _, job := range commands() {
		fmt.Println(job.GetName(), "脚本开始执行")
		_, _ = c.AddFunc(job.GetTime(), job.GetHandler())
	}
	c.Start()
}

type JobCommands interface {
	GetTime() string
	GetName() string
	GetHandler() func()
}
