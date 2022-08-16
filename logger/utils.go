package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func GetJsonFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{
		TimestampFormat: TimestampFormat,
		FieldMap: logrus.FieldMap{
			"time": "@time",
			"msg":  "~msg",
		},
	}
}

func GetTextFormatter() logrus.Formatter {
	return &logrus.TextFormatter{
		TimestampFormat: TimestampFormat,
		FieldMap: logrus.FieldMap{
			"time": "@time",
			"msg":  "~msg",
		},
	}
}

func FileStdout(path, fileName string) io.Writer {
	p := path + "/" + fileName + ".%Y-%m-%d.log"
	linkPath := path + "/" + fileName + ".link.log"
	logf, err := rotatelogs.New(
		p,
		rotatelogs.WithLinkName(linkPath),
		rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationTime(time.Hour*12),
	)
	if err != nil {
		return os.Stdout
	}
	return logf
}

func FileAndTerminalStdout(path, fileName string) io.Writer {
	file := FileStdout(path, fileName)
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	return io.MultiWriter(writers...)
}
