package logs

import (
	"errors"
	"github.com/sirupsen/logrus"
	"goweb/pkg/dingrobot"
	"io"
	"os"
	"time"
)

type CustomFileWriter struct {
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

	// 默认权限，默认是 0777
	// 这个一定要设置合理，否则无法写入
	Perm os.FileMode

	// 是否启用钉钉推送
	IsDingRobot bool
}

// 检验是否实现
var _ io.Writer = &CustomFileWriter{}

func NewCustomFileWriter() *CustomFileWriter {
	return &CustomFileWriter{
		"daily",
		"storage/logs/ginStd/",
		"2006-01-02.txt",
		os.FileMode(0777),
		false,
	}
}

// 日志切割
func (writer *CustomFileWriter) Write(p []byte) (n int, err error) {

	// 建立目录
	_ = os.MkdirAll(writer.Dir, writer.Perm)

	// 目录授权
	_ = os.Chmod(writer.Dir, writer.Perm)

	// 日志文件
	var f *os.File
	f, err = writer.open()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// 写入日志
	n, err = f.Write(p)
	if err != nil {
		return 0, err
	}

	// robot
	if writer.IsDingRobot {
		robot := dingrobot.NewRobot(os.Getenv("LOG_DING_ACCESS_TOKEN"))
		md := "# PROD Recover 告警：\n" + "```text\n" + string(p) + "\n```"
		msg := dingrobot.NewMessageBuilder(dingrobot.TypeMarkdown).Markdown("PROD 接口告警", md).Build()
		err = robot.SendMessage(msg)
		if err != nil {
			logrus.Println(err)
		}
	}

	return n, nil
}

// 打开文件
func (writer *CustomFileWriter) open() (*os.File, error) {

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
