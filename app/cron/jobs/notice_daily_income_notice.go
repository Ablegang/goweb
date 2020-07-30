package jobs

import (
	"fmt"
	"github.com/chenhg5/collection"
	"github.com/sirupsen/logrus"
	show2 "goweb/app/handlers/show"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/dingrobot"
	"goweb/pkg/quotes"
)

type QuoteDailyIncomeNotice struct {
}

func (job *QuoteDailyIncomeNotice) GetName() string {
	return "每日盈亏"
}

func (job *QuoteDailyIncomeNotice) GetTime() []string {
	return []string{
		"* 01 15 * * mon-fri",
	}
}

func (job *QuoteDailyIncomeNotice) GetHandler() func() {

	var (
		TemplateHead = "# 每日总结：\n @" + GetAtMobile() + "\n"
		TemplateBody = "- %s 入选价：%v 现价：%v 入选至今收益：%v%%\n"
	)

	return func() {
		// 展示中标的列表
		qs := make([]show.Quote, 0)
		_ = models.Show().Where("status = ?", show2.DisplayState).Find(&qs)
		keys := collection.Collect(qs).Pluck("key").ToStringArray()

		// 取远程数据
		driver := quotes.New(quotes.WyResource)
		driver.SetKeys(keys)
		qMap, err := driver.GetMap()
		if err != nil {
			logrus.Errorln(err)
			return
		}

		// 生成通知 markdown
		md := TemplateHead
		for _, q := range qs {
			q.TodayPrice = qMap[q.Key].NowPrice
			// 更新当日价格
			_, _ = models.Show().Update(&q)
			sy := (q.TodayPrice - q.InitialPrice) / q.InitialPrice * 100
			md += fmt.Sprintf(TemplateBody, q.Name, q.InitialPrice, q.TodayPrice, sy)
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
