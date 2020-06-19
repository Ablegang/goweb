package logrus_docs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	// 对象池，主要为了重用 bytes.Buffer 实例
	bufferPool *sync.Pool

	// 包名，在首次使用时缓存
	logrusPackage string

	minimumCallerDepth int

	// sync.Once 结构体有一个 Do 方法，此方法能保证传入的函数在并发环境下只被执行一次，内部通过互斥锁实现
	callerInitOnce sync.Once
)

// 主要用于得到 trace stack（执行栈，如某文件、某行、某方法...）
const (
	maximumCallerDepth int = 25 // 最大堆栈帧
	knownLogrusFrames  int = 4  // Logrus 所处的堆栈帧
)

func init() {
	// 缓存池，负责字节缓冲
	// Entry.logs 被调用时用到，主要为了 bytes.Buffer 的重用
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// 指定 runtime.Caller 的 skip
	minimumCallerDepth = 1
}

// WithError 时指定一个默认的 key 传给 WithField 方法
var ErrorKey = "error"

// Entry 实例是 Logrus 记录的具体结构，包含所有通过 WithField 注册的字段
// 会在最终调用 Trace、Debug 之类的方法时记录
// 使用时，应尽可能将记录内容抽象为 Entry 以实现重用
// logrus-docs 的设计思想是将日志抽象为对象来管理，Entry 就是对这一层抽象的封装
type Entry struct {
	Logger *Logger

	// 包含所有自定义的 field
	Data Fields

	// 实例创建时的时间
	Time time.Time

	// 错误级别
	Level Level

	// 调用链 trace stack 含包名
	Caller *runtime.Frame

	// 要记录的内容文本内容
	Message string

	// 调用 entry.logs 时会用到，字节缓冲区
	Buffer *bytes.Buffer

	// 用户自定义的 Context
	Context context.Context

	// 包含的错误
	err string
}

// new entry
func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// 默认 Data 大小，注意：make 返回的是值，new 返回的是指针
		Data: make(Fields, 6),
	}
}

// 返回 bytes 形式的 Entry
func (entry *Entry) Bytes() ([]byte, error) {
	return entry.Logger.Formatter.Format(entry)
}

// 返回 string 形式的 Entry
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Bytes()
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// 给 entry 实例添加一个 error field
// 会从当前 entry 派生一个新的 entry
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// 给 entry 添加一个 context
// 会从当前 entry 派生一个新的 entry
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	// map 是引用类型，所以需要 copy （其实可以直接用 copy 来写）
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: dataCopy, Time: entry.Time, err: entry.err, Context: ctx}
}

// 给 entry 添加一个 field
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// 给 entry 添加多个 field
// 会从当前 entry 派生一个新的 entry
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	fieldErr := entry.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = entry.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}

	// 这里都会返回一个新的 entry 实例，而其属性，除了特定的内容外，都继承于原 entry
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, err: fieldErr, Context: entry.Context}
}

// 给 entry 添加一个 time
// 会派生一个新 entry
func (entry *Entry) WithTime(t time.Time) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	//
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: dataCopy, Time: t, err: entry.err, Context: entry.Context}
}

// 取包名
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// 用于检索第一个 logrus-docs 上游的函数名称
// 比如在 main 里面调用了 logrus-docs 的函数，则拿到 main
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// 是否有 Caller
func (entry Entry) HasCaller() (has bool) {
	return entry.Logger != nil &&
		entry.Logger.ReportCaller &&
		entry.Caller != nil
}

// 此函数未使用 *Entry 来作为接收者，避免多 goroutine 情况下出现问题
func (entry Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer

	// Default to now, but allow users to override if they want.
	//
	// We don't have to worry about polluting future calls to Entry#logs()
	// with this assignment because this function is declared with a
	// non-pointer receiver.
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = level
	entry.Message = msg
	entry.Logger.mu.Lock()
	if entry.Logger.ReportCaller {
		entry.Caller = getCaller()
	}
	entry.Logger.mu.Unlock()

	entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	// To avoid Entry#logs() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

// 执行 hooks
func (entry *Entry) fireHooks() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	err := entry.Logger.Hooks.Fire(entry.Level, entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}

// 写日志 （内容是整个 entry 格式化）
func (entry *Entry) write() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}
	if _, err = entry.Logger.Out.Write(serialized); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to logs, %v\n", err)
	}
}

// 写日志（不单是写，还会先执行 hooks 等操作）
func (entry *Entry) Log(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) Trace(args ...interface{}) {
	entry.Log(TraceLevel, args...)
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DebugLevel, args...)
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(InfoLevel, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WarnLevel, args...)
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FatalLevel, args...)
	// Fatal 还会执行 Exit
	entry.Logger.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	entry.Log(PanicLevel, args...)
	// Panic 还会执行 panic
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.Logf(TraceLevel, format, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(InfoLevel, format, args...)
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WarnLevel, format, args...)
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FatalLevel, format, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Logf(PanicLevel, format, args...)
}

// Entry Println family functions

func (entry *Entry) Logln(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.Log(level, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Traceln(args ...interface{}) {
	entry.Logln(TraceLevel, args...)
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Logln(DebugLevel, args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Logln(InfoLevel, args...)
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Logln(WarnLevel, args...)
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Logln(ErrorLevel, args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Logln(FatalLevel, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Logln(PanicLevel, args...)
}

func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	// 去掉换行字符
	return msg[:len(msg)-1]
}
