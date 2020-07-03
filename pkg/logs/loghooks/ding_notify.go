package loghooks

import (
	"github.com/sirupsen/logrus"
	"goweb/pkg/dingrobot"
	"time"
)

type dingNotify struct {
	// 一些需要初始化的配置...
	accessToken    string
	entryFormatter logrus.Formatter
}

// 检测是否实现了 Hook 接口
var _ logrus.Hook = &dingNotify{}

// 获取一个 emailNotify
func NewDingNotify(accessToken string) *dingNotify {
	return &dingNotify{
		accessToken: accessToken,
		entryFormatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano, // 含纳秒
		},
	}
}

// 要注册的日志级别
func (e *dingNotify) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// 运行 Hook 逻辑
func (e *dingNotify) Fire(entry *logrus.Entry) error {
	_, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}

	// 钉钉机器人通知逻辑...
	robot := dingrobot.NewRobot(e.accessToken)
	b, _ := e.entryFormatter.Format(entry)
	md := "# PROD Custom 告警：\n" + "```json\n" + string(b) + "\n```"
	msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("PROD 接口告警", md).Build()
	err = robot.SendMessage(msg)
	if err != nil {
		logrus.Println("钉钉群机器人推送失败", err)
	}

	return nil
}
