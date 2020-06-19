package loghooks

import (
	"goweb/pkg/logs/logrus"
	"os"
	"time"
)

type FileWrite struct {
	// 日志模式
	LogMode          bool
	Dir              string
	FileNameFormater string

}

// 检查是否已实现 Hook 接口
var _ logrus.Hook = &FileWrite{}

func NewFileWrite() *FileWrite {
	return &FileWrite{}
}

// 注册到所有级别
func (fileCut *FileWrite) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 包内全局变量初始化
var (
	f *os.File

	// 文件相对路径
	logsDir = "storage/logs/"
	logFile = logsDir + time.Now().Format("2006-01-02") + ".txt"
)

// 日志切割
func (fileCut *FileWrite) Fire(entry *logrus.Entry) error {

	if f == nil {
		// 首次启动，尚未打开过 File
		f = open(logFile)
	}

	// 触发 Hook 时的时间作为文件名，如果与全局变量不同，则以新文件名为准
	newFile := logsDir + time.Now().Format("2006-01-02") + ".txt"
	if logFile != newFile {
		// 日志文件名已经更新，先关闭老文件
		f.Close()

		// 新文件名替换给全局变量
		logFile = newFile

		f = open(logFile)
	}

	entry.Logger.SetOutput(f)

	return nil
}

// 打开文件
func open(filePath string) *os.File {
	flag := os.O_APPEND | os.O_RDWR | os.O_CREATE
	perm := os.FileMode(0777)
	f, err := os.OpenFile(filePath, flag, perm)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(logsDir, perm)
		}

		if os.IsPermission(err) {
			_ = os.Chmod(logsDir, perm)
			_ = os.Chmod(filePath, perm)
		}

		f, err = os.OpenFile(filePath, flag, perm)
		if err != nil {
			panic("日志文件无法创建，请检查 storage 目录的权限")
		}
	}

	return f
}
