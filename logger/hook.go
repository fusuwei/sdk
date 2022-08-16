package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

const (
	maximumCallerDepth int = 25
)

type Caller struct {
	EndString string
}

func (c *Caller) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c *Caller) Fire(entry *logrus.Entry) error {
	file, funcName, line := getCaller(c.EndString)
	entry.Data["file"] = file
	entry.Data["func"] = fmt.Sprintf("%s:[%d]", funcName, line)
	return nil
}

func getCaller(endString string) (string, string, int) {
	var (
		callerDepth int
		fileName    string
		funcName    string
		line        int
	)
	pcs := make([]uintptr, maximumCallerDepth)
	_ = runtime.Callers(0, pcs)
	for i := 8; i < maximumCallerDepth; i++ {
		funcName := runtime.FuncForPC(pcs[i]).Name()
		if strings.Contains(funcName, endString) {
			callerDepth = i + 2
			break
		}
	}
	frames, ok := runtime.CallersFrames(pcs[callerDepth:]).Next()
	if !ok {
		return fileName, funcName, line
	}
	funcs := strings.Split(frames.Function, "/")
	if len(funcs) < 5 {
		return fileName, funcName, line
	}
	funcName = strings.Join(funcs[len(funcs)-4:], "/")

	files := strings.Split(frames.File, "/")
	fileName = strings.Join(files[len(funcs)-5:], "/")
	line = frames.Line
	return fileName, funcName, line
}
