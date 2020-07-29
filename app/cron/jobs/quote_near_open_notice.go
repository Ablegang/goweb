package jobs

import (
	"goweb/app/cron"
	"goweb/pkg/dingrobot"
)

type QuoteNearOpenNotice struct {
}

func (job *QuoteNearOpenNotice) GetName() string {
	return "开盘提醒"
}

func (job *QuoteNearOpenNotice) GetTime() []string {
	return []string{
		"* 28 9 * * mon-fri",
	}
}

func (job *QuoteNearOpenNotice) GetHandler() func() {
	return func() {
		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      cron.RobotToken,
			Md:      "# 临近开盘，特此提醒，切忌追涨杀跌，只在尾盘操作！ \n @所有人",
			Title:   job.GetName(),
			At:      []string{},
			IsAtAll: true,
		})
	}
}
