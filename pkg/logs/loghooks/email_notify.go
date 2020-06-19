package loghooks

import (
	"github.com/sirupsen/logrus"
)

type emailNotify struct {
	// 一些需要初始化的配置...
}

// 检测是否实现了 Hook 接口
var _ logrus.Hook = &emailNotify{}

// 获取一个 emailNotify
func NewEmailNotify() *emailNotify {
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

// 运行 Hook 逻辑
func (e *emailNotify) Fire(entry *logrus.Entry) error {
	_, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}

	// 邮件通知逻辑...

	return nil
}
