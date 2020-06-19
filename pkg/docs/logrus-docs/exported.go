package logrus_docs

import (
	"context"
	"io"
	"time"
)

var (
	// 默认初始化的 logs，但只能在包内使用
	// 这个默认的 logs 主要用于一些直接调用了包函数的操作
	std = New()
)

// 返回 std ，这里返回的指针，所以多次使用此方法取 std，拿到的都是同一个 logs
// 看起来不是线程安全的，但实际上 logs 结构体本身就用了 mutex，所有方法都是线程安全的
// 除非设置了 std.SetNoLock()
func StandardLogger() *Logger {
	return std
}

// 给 std 设置 output
func SetOutput(out io.Writer) {
	std.SetOutput(out)
}

// 给 std 设置 Formatter
func SetFormatter(formatter Formatter) {
	std.SetFormatter(formatter)
}

// 给 std 设置 ReportCaller
func SetReportCaller(include bool) {
	std.SetReportCaller(include)
}

// 给 std 设置最小日志级别
func SetLevel(level Level) {
	std.SetLevel(level)
}

// 获取 std 的最小日志级别
func GetLevel() Level {
	return std.GetLevel()
}

// 判断 std 是否支持该日志级别
func IsLevelEnabled(level Level) bool {
	return std.IsLevelEnabled(level)
}

// 给 std 添加 hook
func AddHook(hook Hook) {
	std.AddHook(hook)
}

// 用默认的 error key 来创建一个带 Error field 的 entry
func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

// 返回一个指定了 context 的 entry
func WithContext(ctx context.Context) *Entry {
	return std.WithContext(ctx)
}

// 返回一个指定了 field 的 entry
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// 同上，但可指定多个 field
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// 返回一个指定了 time 的 entry
func WithTime(t time.Time) *Entry {
	return std.WithTime(t)
}

// 调用默认 logs 的 trace
func Trace(args ...interface{}) {
	std.Trace(args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}

func Print(args ...interface{}) {
	std.Print(args...)
}

func Info(args ...interface{}) {
	std.Info(args...)
}

func Warn(args ...interface{}) {
	std.Warn(args...)
}

func Warning(args ...interface{}) {
	std.Warning(args...)
}

func Error(args ...interface{}) {
	std.Error(args...)
}

func Panic(args ...interface{}) {
	std.Panic(args...)
}

func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}
// Traceln logs a message at level Trace on the standard logs.
func Traceln(args ...interface{}) {
	std.Traceln(args...)
}

func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

func Println(args ...interface{}) {
	std.Println(args...)
}

func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

func Warnln(args ...interface{}) {
	std.Warnln(args...)
}

func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}
