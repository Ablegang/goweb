package syslog

import (
	"fmt"
	"log/syslog"
	"os"

	"goweb/pkg/docs/logrus-docs"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

// Creates a hook to be added to an instance of logs. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { logs.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*SyslogHook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &SyslogHook{w, network, raddr}, err
}

func (hook *SyslogHook) Fire(entry *logrus_docs.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus_docs.PanicLevel:
		return hook.Writer.Crit(line)
	case logrus_docs.FatalLevel:
		return hook.Writer.Crit(line)
	case logrus_docs.ErrorLevel:
		return hook.Writer.Err(line)
	case logrus_docs.WarnLevel:
		return hook.Writer.Warning(line)
	case logrus_docs.InfoLevel:
		return hook.Writer.Info(line)
	case logrus_docs.DebugLevel, logrus_docs.TraceLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

func (hook *SyslogHook) Levels() []logrus_docs.Level {
	return logrus_docs.AllLevels
}
