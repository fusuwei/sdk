package logger

import (
	"github.com/FuSuwei/sdk/utils"
	"github.com/sirupsen/logrus"
	"io"
	"log"
)

const TimestampFormat = "2006-01-02 15:04:05"

type Config struct {
	Path      string
	Name      string
	Level     string
	Formatter logrus.Formatter
}

type Logger struct {
	*logrus.Logger
	io.Writer
	Config
}

func New(writer io.Writer, config Config) *Logger {
	l := &Logger{
		logrus.New(),
		writer,
		config,
	}
	if l.Name == "" {
		l.Name = "log"
	}
	if v, ok := utils.MakeDir(l.Path); ok {
		l.Path = v
	} else {
		log.Fatal("创建文件失败")
	}
	if l.Config.Formatter != nil {
		l.SetFormatter(l.Config.Formatter)
	} else {
		l.SetFormatter(GetJsonFormatter())
	}
	var level logrus.Level
	if l.Config.Level != "" {
		if le, err := logrus.ParseLevel(l.Config.Level); err != nil {
			level = logrus.ErrorLevel
		} else {
			level = le
		}
	} else {
		level = logrus.ErrorLevel
	}
	l.Logger.SetLevel(level)
	l.Logger.SetOutput(l.Writer)
	return l
}

func (l *Logger) AddHook(hook logrus.Hook) {
	l.Logger.AddHook(hook)
}

func (l *Logger) Panic(msg string) {
	l.WithFields(logrus.Fields{"name": l.Name}).Panic(msg)
}

func (l *Logger) Fatal(msg string) {
	l.WithFields(logrus.Fields{"name": l.Name}).Fatal(msg)
}

func (l *Logger) Error(msg string) {
	l.WithFields(logrus.Fields{"name": l.Name}).Error(msg)
}

func (l *Logger) Info(msg string) {
	l.WithFields(logrus.Fields{"name": l.Name}).Info(msg)
}

func (l *Logger) Debug(msg string) {
	l.WithFields(logrus.Fields{"name": l.Name}).Debug(msg)
}

func (l *Logger) Warn(msg string) {
	l.WithFields(logrus.Fields{"name": l.Name}).Warn(msg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{"name": l.Name}).Errorf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{"name": l.Name}).Infof(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{"name": l.Name}).Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{"name": l.Name}).Warnf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{"name": l.Name}).Panicf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.WithFields(logrus.Fields{"name": l.Name}).Fatalf(format, args...)
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.Logger.WithFields(fields)
	return l
}
