package jobs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"goweb/app/models"
	"goweb/app/models/show"
	"goweb/pkg/dingrobot"
	"goweb/pkg/quotes"
	"math"
	"os"
	"time"
)

type QuoteNotice struct {
}

func (job *QuoteNotice) GetName() string {
	return "涨跌幅提醒，每波动一个百分点都会通知"
}

func (job *QuoteNotice) GetTime() string {
	return "*/5 * * * * ?"
}

func (job *QuoteNotice) GetHandler() func() {
	var (
		// 通用配置
		AtMobile    = os.Getenv("AT_MOBILE")
		RobotToken  = os.Getenv("QUOTES_DING_ACCESS_TOKEN")
		CommonTitle = "涨跌幅通知"
		ZTemplate   = "%s \n - 涨幅 %s%% 现价：%s \n @%s"
		DTemplate   = "%s \n - 跌幅 %s%% 现价：%s \n @%s"
	)

	return func() {
		driver := quotes.New(quotes.WyResource)
		driver.SetKeys(job.getKeys())
		data, err := driver.GetQuotes()
		if err != nil {
			logrus.Errorln(err)
			return
		}

		needToNotice := job.getNeeds(data)
		fmt.Println(needToNotice)
		for _, n := range needToNotice {
			template := ZTemplate
			if n.Percent < 0 {
				template = DTemplate
			}
			go dingrobot.Markdown(&dingrobot.MarkdownParams{
				Ac:      RobotToken,
				Md:      fmt.Sprintf(template, n.Name, n.PercentStr, n.NowPriceStr, AtMobile),
				Title:   CommonTitle,
				At:      []string{AtMobile},
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

// 取数据库的 keys
func (job *QuoteNotice) getKeys() (keys []string) {
	model := show.Quote{}
	_ = models.Show().Table(model.TableName()).Cols("key").Find(&keys)
	return
}

// 取需要通知的 quote
func (job *QuoteNotice) getNeeds(data []quotes.QuoteData) (res []quotes.QuoteData) {
	qs := make([]show.Notice, 0)
	_ = models.Show().Where("created_at > ?", time.Now().Format("2006-01-02")).Find(&qs)
	fmt.Println(qs)
	// 遍历所有标的
	for _, d := range data {
		last := job.getLastNotice(qs, d)
		// 取百分点
		per := int64(math.Floor(d.Percent * 100))
		if per != last.Per {
			res = append(res, d)
		}
	}

	return
}

// 在今日所有通知里取出某标的最后的一条通知
func (job *QuoteNotice) getLastNotice(qs []show.Notice, d quotes.QuoteData) (last show.Notice) {
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
