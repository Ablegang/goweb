package jobs

import (
	"goweb/pkg/dingrobot"
)

type QuoteNearCloseNotice struct {
}

func (job *QuoteNearCloseNotice) GetName() string {
	return "收盘操作提醒"
}

func (job *QuoteNearCloseNotice) GetTime() []string {
	return []string{
		"1 54 14 * * mon-fri",
	}
}

func (job *QuoteNearCloseNotice) GetHandler() func() {
	return func() {
		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      GetRobotToken(),
			Md:      "# 临近收盘，特此提醒，抓紧操作 \n @所有人",
			Title:   job.GetName(),
			At:      []string{},
			IsAtAll: false,
		})
	}
}
