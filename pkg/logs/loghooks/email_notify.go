package loghooks

import (
	"fmt"
	"goweb/pkg/log/logrus"
)

type emailNotify struct {
	// 一些需要初始化的配置...
}

func newEmailNotify() *emailNotify {
	return &emailNotify{}
}

// 要注册的日志级别
func (e *emailNotify) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (e *emailNotify) Fire(entry *logrus.Entry) error {
	b, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println()
}
