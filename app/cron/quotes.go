package cron

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goweb/pkg/dingrobot"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 涨跌幅通知
func ListenQuotesNotice() {
	// 已通知 map ，通知过后，10 分钟内不再通知
	zNoticed2 := make(map[string]int, 0)
	zNoticed5 := make(map[string]int, 0)
	dNoticed2 := make(map[string]int, 0)
	dNoticed5 := make(map[string]int, 0)

	// 5 秒取数据
	tick := time.NewTicker(5 * time.Second)
	for true {
		// 取数据
		u := GetQuotes()
		if len(u) == 0 {
			<-tick.C
			continue
		}
		for _, v := range u {
			// 解析数据
			name, _ := v["name"].(string)
			percent, _ := v["percent"].(float64)
			percentStr := strconv.FormatFloat(percent*100, 'f', 6, 64)
			now := v["price"].(float64)
			nowStr := strconv.FormatFloat(now, 'f', 6, 64)

			// 涨超 5%
			if _, ok := zNoticed5[name]; percent*100 >= 5 && !ok {
				robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
				md := "涨超5%：" + name + "\n" + "- " + percentStr + "%" + " " + nowStr + " \n"
				msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
				err := robot.SendMessage(msg)
				if err != nil {
					logrus.Errorln("钉钉行情推送失败", err)
				}
				zNoticed5[name] = 1
				zNoticed2[name] = 1
				// 10 分钟后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(10 * time.Minute)
					<-timer.C
					delete(zNoticed5, name)
					delete(zNoticed2, name)
				}()
			}

			// 涨超 2%
			if _, ok := zNoticed2[name]; percent*100 >= 2 && !ok {
				robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
				md := "涨超2%：" + name + "\n" + "- " + percentStr + "%" + " " + nowStr + " \n"
				msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
				err := robot.SendMessage(msg)
				if err != nil {
					logrus.Errorln("钉钉行情推送失败", err)
				}
				zNoticed2[name] = 1
				// 10 分钟后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(10 * time.Minute)
					<-timer.C
					delete(zNoticed2, name)
				}()
			}

			// 跌超 5%
			if _, ok := dNoticed5[name]; percent*100 <= -5 && !ok {
				robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
				md := "跌超5%：" + name + "\n" + "- " + percentStr + "%" + " " + nowStr + " \n"
				msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
				err := robot.SendMessage(msg)
				if err != nil {
					logrus.Errorln("钉钉行情推送失败", err)
				}
				dNoticed5[name] = 1
				dNoticed2[name] = 1
				// 10 分钟后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(10 * time.Minute)
					<-timer.C
					delete(dNoticed5, name)
					delete(dNoticed2, name)
				}()
			}

			// 跌超 2%
			if _, ok := dNoticed2[name]; percent*100 <= -2 && !ok {
				robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
				md := "跌超2%：" + name + "\n" + "- " + percentStr + "%" + " " + nowStr + " \n"
				msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
				err := robot.SendMessage(msg)
				if err != nil {
					logrus.Errorln("钉钉行情推送失败", err)
				}
				dNoticed2[name] = 1
				// 10 分钟后从 zNoticed 里去除
				go func() {
					timer := time.NewTimer(10 * time.Minute)
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
	// 15 分钟取一次数据
	tick := time.NewTicker(15 * time.Minute)
	for true {
		// 取数据
		u := GetQuotes()
		if len(u) == 0 {
			<-tick.C
			continue
		}
		robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
		md := "# 监控：\n"
		for _, v := range u {
			// 解析数据
			name, _ := v["name"].(string)
			percent, _ := v["percent"].(float64)
			percentStr := strconv.FormatFloat(percent*100, 'f', 6, 64)
			now := v["price"].(float64)
			nowStr := strconv.FormatFloat(now, 'f', 6, 64)
			md += "- " + name + " 涨跌幅：" + percentStr + "% 现价：" + nowStr + "\n"
		}
		msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
		err := robot.SendMessage(msg)
		if err != nil {
			logrus.Errorln("钉钉行情推送失败", err)
		}

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

		robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
		md := "临近收盘，特此提醒，抓紧操作"
		msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("市场监控", md).Build()
		err := robot.SendMessage(msg)
		if err != nil {
			logrus.Errorln("钉钉行情推送失败", err)
		}

		<-tick.C
	}
}

// 取行情数据
func GetQuotes() map[string]map[string]interface{} {
	u := make(map[string]map[string]interface{})
	now := time.Now()
	// 9 - 12
	morning := (now.Hour() >= 9) && (now.Hour() <= 11)
	// 13 - 15
	afternoon := now.Hour() >= 13 && now.Hour() <= 14
	// 非盘中
	if !morning && !afternoon {
		return u
	}

	// 非交易日
	if now.Weekday() == time.Sunday || now.Weekday() == time.Saturday {
		return u
	}

	url := "http://api.money.126.net/data/feed/" + os.Getenv("QUOTES")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body[21:len(body)-2], &u)
	return u
}
