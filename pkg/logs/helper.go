package logs

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Panicln(values ...interface{}) {
	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.PanicLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	log.Panicln(values)
}

func Println(values ...interface{}) {
	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	log.Println(values)
}

func Errorln(values ...interface{}) {
	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.ErrorLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	log.Errorln(values)
}
