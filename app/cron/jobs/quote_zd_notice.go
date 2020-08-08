package jobs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/dingrobot"
	"goweb/pkg/quotes"
	"math"
	"time"
)

type QuoteZdNotice struct {
}

func (job *QuoteZdNotice) GetName() string {
	return "涨跌幅提醒"
}

func (job *QuoteZdNotice) GetTime() []string {
	return []string{
		"*/5 30-59 9 * * mon-fri",
		"*/5 0-30 11 * * mon-fri",
		"*/5 * 10,13,14 * * mon-fri",
	}
}

func (job *QuoteZdNotice) GetHandler() func() {

	var (
		ZTemplate = "%s \n - 涨幅 %s 现价：%s \n @%s"
		DTemplate = "%s \n - 跌幅 %s 现价：%s \n @%s"
	)

	return func() {
		driver := quotes.New(quotes.WyResource)
		driver.SetKeys(GetKeys())
		data, err := driver.GetQuotes()
		if err != nil {
			logrus.Errorln(err)
			return
		}

		needToNotice := job.getNeeds(data)
		for _, n := range needToNotice {
			template := ZTemplate
			if n.Percent < 0 {
				template = DTemplate
			}
			dingrobot.Markdown(&dingrobot.MarkdownParams{
				Ac:      GetRobotToken(),
				Md:      fmt.Sprintf(template, n.Name, n.PercentStr, n.NowPriceStr, GetAtMobile()),
				Title:   job.GetName(),
				At:      []string{GetAtMobile()},
				IsAtAll: false,
			})
			_, _ = models.Show().Insert(&show.Notice{
				Key:       n.Key,
				Per:       int64(math.Floor(n.Percent * 100)),
				CreatedAt: time.Now(),
			})
		}
	}
}

// 取需要通知的 quote
func (job *QuoteZdNotice) getNeeds(data []quotes.QuoteData) (res []quotes.QuoteData) {
	qs := make([]show.Notice, 0)
	_ = models.Show().Where("created_at > ?", time.Now()).Find(&qs)
	// 遍历所有标的
	for _, d := range data {
		last := job.getLastNotice(qs, d)
		per := job.getPerFlag(d)
		if per != last.Per {
			res = append(res, d)
		}
	}

	return
}

// 获取百分点标志位
func (job *QuoteZdNotice) getPerFlag(d quotes.QuoteData) int64 {
	// 取百分点
	var per int64
	if d.Percent > -0.005 && d.Percent < 0.005 {
		// -0.5% 到 0.5% 之间只需通知一次
		per = 0
	} else if d.Percent < 0 {
		// -0.5% 之后，每波动 1% 都将通知
		per = int64(math.Floor(d.Percent * 100))
	} else {
		// 0.5% 之后，每波动 1% 都将通知
		per = int64(math.Ceil(d.Percent * 100))
	}
	return per
}

// 在今日所有通知里取出某标的最后的一条通知
func (job *QuoteZdNotice) getLastNotice(qs []show.Notice, d quotes.QuoteData) (last show.Notice) {
	// 遍历所有今日的 notice
	for _, q := range qs {
		// 如果是同一标的
		if d.Key == q.Key {
			// 取最晚的那条
			if &last == nil {
				last = q
			} else {
				if q.CreatedAt.After(last.CreatedAt) {
					last = q
				}
			}
		}
	}
	return
}
