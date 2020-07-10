package quotes

import (
	"encoding/json"
	"fmt"
	"goweb/pkg/dingrobot"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	// 通用配置
	AtMobile    = os.Getenv("AT_MOBILE")
	Quotes      = os.Getenv("QUOTES")
	RobotToken  = os.Getenv("QUOTES_DING_ACCESS_TOKEN")
	CommonTitle = "行情播报"
	SpiderApi   = "http://api.money.126.net/data/feed/"
)

// 涨跌幅通知
func ListenQuotesNotice() {
	var (
		Z5Template = "涨超5%%：%s \n - 涨幅 %s%% 现价：%s \n @%s"
		Z2Template = "涨超2%%：%s \n - 涨幅 %s%% 现价：%s \n @%s"
		D5Template = "跌超5%%：%s \n - 跌幅 %s%% 现价：%s \n @%s"
		D2Template = "跌超2%%：%s \n - 跌幅 %s%% 现价：%s \n @%s"

		// 已推送列表（只要在列表内，同一事件就不再推送）
		zNoticed2 = make(map[string]int, 0)
		zNoticed5 = make(map[string]int, 0)
		dNoticed2 = make(map[string]int, 0)
		dNoticed5 = make(map[string]int, 0)

		TikeDua   = 5 * time.Second
		DeleteDua = 10 * time.Minute
	)

	// 5 秒取数据
	tick := time.NewTicker(TikeDua)
	for true {
		// 取数据
		u := GetQuotes(Quotes)
		if len(u) == 0 {
			<-tick.C
			continue
		}

		// 解析数据
		for _, v := range u {
			name, percent, percentStr, _, nowStr := FormatQuotesCoreData(v)
			// 涨超 5%
			if _, ok := zNoticed5[name]; percent*100 >= 5 && !ok {
				dingrobot.Markdown(&dingrobot.MarkdownParams{
					Ac:      RobotToken,
					Md:      fmt.Sprintf(Z5Template, name, percentStr, nowStr, AtMobile),
					Title:   CommonTitle,
					At:      []string{AtMobile},
					IsAtAll: false,
				})
				zNoticed5[name] = 1
				// 先发 5 的，接着就算大于 2 ，也不会再发了
				zNoticed2[name] = 1

				// 到了静默时间后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(DeleteDua)
					<-timer.C
					delete(zNoticed5, name)
					delete(zNoticed2, name)
				}()
			}

			// 涨超 2%
			if _, ok := zNoticed2[name]; percent*100 >= 2 && !ok {
				dingrobot.Markdown(&dingrobot.MarkdownParams{
					Ac:      RobotToken,
					Md:      fmt.Sprintf(Z2Template, name, percentStr, nowStr, AtMobile),
					Title:   CommonTitle,
					At:      []string{AtMobile},
					IsAtAll: false,
				})
				zNoticed2[name] = 1

				go func() {
					timer := time.NewTimer(DeleteDua)
					<-timer.C
					delete(zNoticed2, name)
				}()
			}

			// 跌超 5%
			if _, ok := dNoticed5[name]; percent*100 <= -5 && !ok {
				dingrobot.Markdown(&dingrobot.MarkdownParams{
					Ac:      RobotToken,
					Md:      fmt.Sprintf(D5Template, name, percentStr, nowStr, AtMobile),
					Title:   CommonTitle,
					At:      []string{AtMobile},
					IsAtAll: false,
				})
				dNoticed5[name] = 1
				dNoticed2[name] = 1
				// 10 分钟后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(DeleteDua)
					<-timer.C
					delete(dNoticed5, name)
					delete(dNoticed2, name)
				}()
			}

			// 跌超 2%
			if _, ok := dNoticed2[name]; percent*100 <= -2 && !ok {
				dingrobot.Markdown(&dingrobot.MarkdownParams{
					Ac:      RobotToken,
					Md:      fmt.Sprintf(D2Template, name, percentStr, nowStr, AtMobile),
					Title:   CommonTitle,
					At:      []string{AtMobile},
					IsAtAll: false,
				})
				dNoticed2[name] = 1
				// 10 分钟后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(DeleteDua)
					<-timer.C
					delete(dNoticed2, name)
				}()
			}
		}

		<-tick.C
	}
}

// 自选通知
func ListenQuotesCommonPush() {
	var (
		TikeDua      = 15 * time.Minute
		TemplateHead = "# 监控：\n @" + AtMobile + "\n"
		TemplateBody = "- %s 涨跌幅：%s%% 现价：%s\n"
	)

	tick := time.NewTicker(TikeDua)
	for true {
		// 取数据
		u := GetQuotes(Quotes)
		if len(u) == 0 {
			<-tick.C
			continue
		}

		md := TemplateHead
		for _, v := range u {
			// 解析数据
			name, _, percentStr, _, nowStr := FormatQuotesCoreData(v)
			md += fmt.Sprintf(TemplateBody, name, percentStr, nowStr)
		}

		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      RobotToken,
			Md:      md,
			Title:   CommonTitle,
			At:      []string{AtMobile},
			IsAtAll: false,
		})

		<-tick.C
	}
}

// 收盘前通知
func NearCloseNotice() {
	// 每分钟
	tick := time.NewTicker(time.Minute)
	for true {
		now := time.Now()
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			<-tick.C
			continue
		}
		if now.Hour() != 14 || now.Minute() != 54 {
			<-tick.C
			continue
		}

		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      RobotToken,
			Md:      "# 临近收盘，特此提醒，抓紧操作 \n @所有人",
			Title:   CommonTitle,
			At:      []string{},
			IsAtAll: true,
		})

		<-tick.C
	}
}

// 开盘前通知
func NearOpenNotice() {
	// 每分钟
	tick := time.NewTicker(time.Minute)
	for true {
		now := time.Now()
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			<-tick.C
			continue
		}
		if now.Hour() != 9 || now.Minute() != 28 {
			<-tick.C
			continue
		}

		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      RobotToken,
			Md:      "# 临近开盘，特此提醒，切忌追涨杀跌，只在尾盘操作！ \n @所有人",
			Title:   CommonTitle,
			At:      []string{},
			IsAtAll: true,
		})

		<-tick.C
	}
}

// 取行情数据
func GetQuotes(quotesList string) map[string]map[string]interface{} {
	u := make(map[string]map[string]interface{})

	// 非盘中处理
	now := time.Now()
	morning := (now.Hour() >= 9) && (now.Hour() <= 11)
	afternoon := now.Hour() >= 13 && now.Hour() <= 14
	if !morning && !afternoon {
		return u
	}
	if now.Weekday() == time.Sunday || now.Weekday() == time.Saturday {
		return u
	}

	url := SpiderApi + quotesList
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body[21:len(body)-2], &u)
	return u
}

// 对数据进行格式处理
// 如果用结构体来绑定，就不用封装这个方法了
func FormatQuotesCoreData(data map[string]interface{}) (string, float64, string, float64, string) {
	name, _ := data["name"].(string)
	percent, _ := data["percent"].(float64)
	percentStr := strconv.FormatFloat(percent*100, 'f', 6, 64)
	now, _ := data["price"].(float64)
	nowStr := strconv.FormatFloat(now, 'f', 6, 64)
	return name, percent, percentStr, now, nowStr
}
