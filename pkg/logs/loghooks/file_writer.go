package loghooks

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type FileWriter struct {
	// 日志模式
	// daily 指文件按日切分，默认
	// single 指日志存放在单一文件内
	LogMode string

	// 日志文件存放目录
	// 默认为 storage/logs/
	// 注意这里最后一定要加上斜杠
	Dir string

	// 日志文件名格式
	// single 时，此参数就是日志文件名
	// daily 时，此参数则为 time.Now().Format 的参数
	// 默认为 2006-01-02.txt
	FileNameFormater string

	// 每一条日志的格式
	// 默认是 JSONFormatter
	// 自定义这个值时，请参考 logrus 的 Formatter
	EntryFormatter logrus.Formatter

	// 要将当前 Hook 应用到 logrus 的哪些级别里
	// 注意，logrus Std 默认的 Level 是 InfoLevel，所以想要开启 Debug 和 Trace ，需要自己设置下
	HookLevels []logrus.Level

	// 默认权限，默认是 0777
	// 这个一定要设置合理，否则无法写入
	Perm os.FileMode
}

// 检查是否已实现 Hook 接口
var _ logrus.Hook = &FileWriter{}

func NewFileWriter() *FileWriter {
	return &FileWriter{
		"daily",
		"storage/logs/custom/",
		"2006-01-02.txt",
		&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano, // 含纳秒
		},
		logrus.AllLevels,
		os.FileMode(0777),
	}
}

// 注册到所有级别
func (writer *FileWriter) Levels() []logrus.Level {
	return writer.HookLevels
}

// 日志切割
func (writer *FileWriter) Fire(entry *logrus.Entry) error {

	// 建立目录
	_ = os.MkdirAll(writer.Dir, writer.Perm)

	// 目录授权
	_ = os.Chmod(writer.Dir, writer.Perm)

	// 日志文件
	f, err := writer.open()
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入日志
	// 自己实现写入，而不使用 logrus 的 output ，有以下两个原因
	// 1、阅读 logrus 的源码可以发现，当 Hook 的 Fire 被执行时，entry.logger 处于 Lock 状态，所以在 Hook 内 SetOutput 会引发死锁
	// 2、就算设置了 output ，也难以处理 f.Close
	b, _ := writer.EntryFormatter.Format(entry)
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// 打开文件
func (writer *FileWriter) open() (*os.File, error) {

	// 根据日志模式取日志文件名
	fileName := ""
	if writer.LogMode == "daily" {
		fileName = time.Now().Format(writer.FileNameFormater)
	} else if writer.LogMode == "single" {
		fileName = writer.FileNameFormater
	} else {
		return nil, errors.New("错误的 LogMode")
	}

	// 打开文件
	f, err := os.OpenFile(writer.Dir+fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, writer.Perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
