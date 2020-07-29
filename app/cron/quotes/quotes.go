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
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
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
