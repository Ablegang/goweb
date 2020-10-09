package logs

import (
	"errors"
	"goweb/pkg/dingrobot"
	"io"
	"os"
	"strings"
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

	// 路径校验
	if isInvalidDir(writer.Dir) {
		return 0, errors.New("错误的路径")
	}

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
		dingrobot.Markdown(&dingrobot.MarkdownParams{
			Ac:      os.Getenv("LOG_DING_ACCESS_TOKEN"),
			Md:      "# PROD Recover 告警：\n" + "```text\n" + string(p) + "\n```",
			Title:   "监控告警",
			At:      []string{"15868100475"},
			IsAtAll: false,
		})
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

	// 路径校验
	thePath := writer.Dir + fileName
	if isInvalidDir(thePath) {
		return nil, errors.New("错误的路径")
	}

	// 打开文件
	f, err := os.OpenFile(thePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, writer.Perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 是否非法路径
func isInvalidDir(dir string) bool {

	if strings.Contains(dir, "../") || strings.Contains(dir, "..\\") {
		return true
	}

	if strings.Index(dir, "storage/") != 0 && strings.Index(dir, "storage\\") != 0 {
		return true
	}

	return false
}
