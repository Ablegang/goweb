package jobs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"goweb/pkg/dingrobot"
	"goweb/pkg/quotes"
)

type QuoteCommonNotice struct {
}

func (job *QuoteCommonNotice) GetName() string {
	return "通用提醒"
}

func (job *QuoteCommonNotice) GetTime() []string {
	return []string{
		"* 16,26,30 9 * * mon-fri",
		"* 31 11 * * mon-fri",
		"* 1 13,15 * * mon-fri",
	}
}

func (job *QuoteCommonNotice) GetHandler() func() {

	var (
		TemplateHead = "# 监控：\n @" + GetAtMobile() + "\n"
		TemplateBody = "- %s 涨跌幅：%s%% 现价：%s\n"
	)

	return func() {
		driver := quotes.New(quotes.WyResource)
		driver.SetKeys(GetKeys())
		data, err := driver.GetQuotes()
		if err != nil {
			logrus.Errorln(err)
			return
		}

		md := TemplateHead
		for _, v := range data {
			// 解析数据
			md += fmt.Sprintf(TemplateBody, v.Name, v.PercentStr, v.NowPriceStr)
		}

		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      GetRobotToken(),
			Md:      md,
			Title:   job.GetName(),
			At:      []string{GetAtMobile()},
			IsAtAll: false,
		})
	}
}
