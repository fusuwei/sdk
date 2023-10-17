package logx

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

func FileStdout(filepath, fileName string) io.Writer {
	writer := &lumberjack.Logger{
		Filename:   path.Join(filepath, fileName),
		MaxSize:    500,
		MaxBackups: 7,
		MaxAge:     28,
		LocalTime:  true,
		Compress:   true,
	}
	return writer
}

func FileAndTerminalStdout(path, fileName string) io.Writer {
	file := FileStdout(path, fileName)
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	return io.MultiWriter(writers...)
}
