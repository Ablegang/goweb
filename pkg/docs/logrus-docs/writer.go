// logrus 的 io.Writer 实例
// 比如要在 gin 中接入 logrus ，可以直接将 logger.Writer 返回的实例传递给 gin

package logrus_docs

import (
	"bufio"
	"io"
	"runtime"
)

// 以 InfoLevel 的级别获取 writer
func (logger *Logger) Writer() *io.PipeWriter {
	return logger.WriterLevel(InfoLevel)
}

// 此方法返回一个 io.Writer，用于写入任意给定日志级别的文本到 logs
// writer 的每一行写入都会以指定的 formatter 作为格式，且会触发相关 Hook
// writer 是 io.Pipe 的一部分，写入完成后，writer 需要由调用者关闭
func (logger *Logger) WriterLevel(level Level) *io.PipeWriter {
	return NewEntry(logger).WriterLevel(level)
}

// 以 InfoLeve 的级别获取 writer
func (entry *Entry) Writer() *io.PipeWriter {
	return entry.WriterLevel(InfoLevel)
}

// 根据指定的 level，返回一个 writer
func (entry *Entry) WriterLevel(level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...interface{})

	switch level {
	case TraceLevel:
		printFunc = entry.Trace
	case DebugLevel:
		printFunc = entry.Debug
	case InfoLevel:
		printFunc = entry.Info
	case WarnLevel:
		printFunc = entry.Warn
	case ErrorLevel:
		printFunc = entry.Error
	case FatalLevel:
		printFunc = entry.Fatal
	case PanicLevel:
		printFunc = entry.Panic
	default:
		printFunc = entry.Print
	}

	// 并发写
	go entry.writerScanner(reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

// 逐行写入
func (entry *Entry) writerScanner(reader *io.PipeReader, printFunc func(args ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		printFunc(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		entry.Errorf("Error while reading from Writer: %s", err)
	}
	_ = reader.Close()
}

// 写入结束后释放资源
func writerFinalizer(writer *io.PipeWriter) {
	_ = writer.Close()
}
