package logrus

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Logger struct {

	// Out 需要实现 io.Writer，用于指定 logs 写入的地方以及写入逻辑
	// 通常 Out 会被分配到一个文件，不分配的话，默认是 os.Stderr
	// 当然也可以指定为 sentry 、 kafka 等
	Out io.Writer

	// log 实例的 hooks，hooks 会在特定的日志级别触发
	// 也可以在特定的 entry 上触发，触发条件，由用户自己设置
	Hooks LevelHooks

	// 所有 entry 在被输出之前，都会先格式化
	// 有两个选项，TextFormatter 和 JSONFormatter，TextFormatter 是默认的
	// 开发环境下，TTY 是 attached 时，控制台输出的 log 是彩色的（色彩代码不会写入文件）
	// Formatter 也可以自己实现，这一点可以看 README
	Formatter Formatter

	// 是否记录 caller 信息，默认否
	ReportCaller bool

	// 最低日志级别
	Level Level

	// 用于保证多个 goroutine 能安全写日志
	mu MutexWrap

	// 可重复使用的空 entry（不自定义 entry，就会使用空 entry）
	// 用对象池来使 entry 可重用
	entryPool sync.Pool

	// 退出 app 的方法实现，若不指定，默认使用 os.Exit(1)
	ExitFunc exitFunc
}

// 定义 exit 方法类型，传参为 int，表示 code
type exitFunc func(int)

// 互斥锁
type MutexWrap struct {
	lock     sync.Mutex
	// 默认为 false，即开启并发互斥锁功能
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

// 禁用互斥锁
func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// log 构建函数
// 默认构建的 log 实例，formatter 是 TextFormatter，Hooks 是 LevelHooks，level 是 info
// 如果需要自定义，可以直接给返回值的属性赋值
func New() *Logger {
	return &Logger{
		Out:          os.Stderr,
		Formatter:    new(TextFormatter),
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}

// Entry 构建函数
func (logger *Logger) newEntry() *Entry {
	// 空 Entry 通过对象池重用
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

// 释放 Entry
func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{} // 置空
	logger.entryPool.Put(entry) // put 回对象池
}

// 从空 entry 得到一个新的 Entry 实例，并添加 field
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry() // 先得到实例（从对象池）
	defer logger.releaseEntry(entry) // 函数结束时释放实例
	return entry.WithField(key, value) // 返回一个新的 entry
}

// 同上，区别是这里分配的是多个 field
func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// 同上，区别是这里分配一个 error field
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// 同上，区别是这里分配一个 context 给 entry
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithContext(ctx)
}

// 同上，区别是这里分配了一个 time 给 entry
func (logger *Logger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}

// 下列 log、print 系列方法都是以空 entry 实例来执行

// 用空 entry 来 Logf
func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logf(level, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(TraceLevel, format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

// 空 entry 来 Printf
func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.Logf(FatalLevel, format, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.Logf(PanicLevel, format, args...)
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Log(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *Logger) Print(args ...interface{}) {
	entry := logger.newEntry()
	entry.Print(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}

func (logger *Logger) Logln(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logln(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Traceln(args ...interface{}) {
	logger.Logln(TraceLevel, args...)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.Logln(DebugLevel, args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.Logln(InfoLevel, args...)
}

func (logger *Logger) Println(args ...interface{}) {
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.Logln(WarnLevel, args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	logger.Warnln(args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.Logln(ErrorLevel, args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.Logln(FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.Logln(PanicLevel, args...)
}

// exit 方法会先执行所有已注册的 Handler
func (logger *Logger) Exit(code int) {
	runHandlers() // 这个 handler 是专门用于 exit 的，所以命名为 exitHandler
	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}

	// 虽然可以自定义 ExitFunc，但也还是会执行所有已注册的 ExitHandler
	logger.ExitFunc(code)
}

// 如果使用附加模式写入文件，这种时候可以安全地读写，就可以用下面的方法禁用 Lock
func (logger *Logger) SetNoLock() {
	logger.mu.Disable()
}

// 返回当前 log 的 level
// 从指针位置返回 level 值
func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

// 设置 log 的 level
// 通过指针将值存储到该内存地址
func (logger *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

// 取 log 的 level
func (logger *Logger) GetLevel() Level {
	return logger.level()
}

// 给 log 添加 hook
func (logger *Logger) AddHook(hook Hook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
}

// 判断日志级别是否被允许
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}

// 设置 Formatter
func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
}

// 设置 output
func (logger *Logger) SetOutput(output io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Out = output
}

// 设置 ReportCaller
func (logger *Logger) SetReportCaller(reportCaller bool) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.ReportCaller = reportCaller
}

// 替换 Hooks，并返回老的 Hooks
func (logger *Logger) ReplaceHooks(hooks LevelHooks) LevelHooks {
	logger.mu.Lock()
	oldHooks := logger.Hooks
	logger.Hooks = hooks
	logger.mu.Unlock()
	return oldHooks
}
