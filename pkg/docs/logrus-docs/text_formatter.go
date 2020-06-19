package logrus_docs

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

type TextFormatter struct {
	// 设置为 true 则在 TTY 上会显示颜色
	ForceColors bool

	// 设置为 true 则强制禁用 TTY 颜色
	DisableColors bool

	// 强制引用所有值
	ForceQuote bool

	// 禁用所有值的引用，优先级低于 ForceQuote
	DisableQuote bool

	// 根据环境变量 CLICOLOR 和 CLICOLOR_FORCE 给颜色 - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// 禁用时间戳记录，当输出重定向到已添加了时间戳的日志时很有用
	DisableTimestamp bool

	// 启用记录 TTY 的完整时间戳，而不仅仅是自开始执行以来经过的时间
	FullTimestamp bool

	// 打印时的时间格式
	TimestampFormat string

	// 默认 fields 是排序的后再记录的，但对于一些 logs 非常频繁且没有使用 json formatter 的情况，可能需要禁用排序
	DisableSorting bool

	// 可以指定一个排序的自定义逻辑，如果未指定，则默认使用 sort.Strings
	SortingFunc func([]string)

	// 禁用 level 文本自动截断
	DisableLevelTruncation bool

	// 是否补齐 level 文本，主要作用是让日志变得整齐
	PadLevelText bool

	// 如果为 true，则 QuoteEmptyFields 会将空字段括在引号中
	QuoteEmptyFields bool

	// 是否将 logs 输出到 terminal
	isTerminal bool

	// FieldMap 允许用户自定义默认字段的键名称
	// As an example:
	// formatter := &TextFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	// 激活 ReportCaller 时，用户可以设置 Call​​erPrettyfier 来修改数据中函数和文件键的内容
	// 如果任何返回值是空字符串，对应的键将从字段中删除
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	terminalInitOnce sync.Once

	// 级别文本的最大长度，在初始化时动态生成
	// 根据 AllLevels 计算出来的
	// 这个长度主要是用于补全，见 PadLevelText
	levelTextMaxLength int
}

func (f *TextFormatter) init(entry *Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}

	// 设置 levelTextMaxLength
	for _, level := range AllLevels {
		levelTextLength := utf8.RuneCount([]byte(level.String()))
		if levelTextLength > f.levelTextMaxLength {
			f.levelTextMaxLength = levelTextLength
		}
	}
}

func (f *TextFormatter) isColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		switch force, ok := os.LookupEnv("CLICOLOR_FORCE"); {
		case ok && force != "0":
			isColored = true
		case ok && force == "0", os.Getenv("CLICOLOR") == "0":
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}

// 格式化日志条目
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var funcVal, fileVal string

	fixedKeys := make([]string, 0, 4+len(data))
	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyTime))
	}
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLevel))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyMsg))
	}
	if entry.err != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLogrusError))
	}
	if entry.HasCaller() {
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}

		if funcVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyFunc))
		}
		if fileVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyFile))
		}
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		} else {
			if !f.isColored() {
				fixedKeys = append(fixedKeys, keys...)
				f.SortingFunc(fixedKeys)
			} else {
				f.SortingFunc(keys)
			}
		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if f.isColored() {
		f.printColored(b, entry, keys, data, timestampFormat)
	} else {

		for _, key := range fixedKeys {
			var value interface{}
			switch {
			case key == f.FieldMap.resolve(FieldKeyTime):
				value = entry.Time.Format(timestampFormat)
			case key == f.FieldMap.resolve(FieldKeyLevel):
				value = entry.Level.String()
			case key == f.FieldMap.resolve(FieldKeyMsg):
				value = entry.Message
			case key == f.FieldMap.resolve(FieldKeyLogrusError):
				value = entry.err
			case key == f.FieldMap.resolve(FieldKeyFunc) && entry.HasCaller():
				value = funcVal
			case key == f.FieldMap.resolve(FieldKeyFile) && entry.HasCaller():
				value = fileVal
			default:
				value = data[key]
			}
			f.appendKeyValue(b, key, value)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, data Fields, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case DebugLevel, TraceLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation && !f.PadLevelText {
		levelText = levelText[0:4]
	}
	if f.PadLevelText {
		// Generates the format string used in the next line, for example "%-6s" or "%-7s".
		// Based on the max level text length.
		formatString := "%-" + strconv.Itoa(f.levelTextMaxLength) + "s"
		// Formats the level text by appending spaces up to the max length, for example:
		// 	- "INFO   "
		//	- "WARNING"
		levelText = fmt.Sprintf(formatString, levelText)
	}

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus-docs text_formatter the same as the stdlib logs package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""
	if entry.HasCaller() {
		funcVal := fmt.Sprintf("%s()", entry.Caller.Function)
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)

		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}

		if fileVal == "" {
			caller = funcVal
		} else if funcVal == "" {
			caller = fileVal
		} else {
			caller = fileVal + " " + funcVal
		}
	}

	switch {
	case f.DisableTimestamp:
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m%s %-44s ", levelColor, levelText, caller, entry.Message)
	case !f.FullTimestamp:
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]%s %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
	default:
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%s %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), caller, entry.Message)
	}
	for _, k := range keys {
		v := data[k]
		_, _ = fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
