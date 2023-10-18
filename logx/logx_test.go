package logx

import "testing"

func Test_test(t *testing.T) {
	writer := FileAndTerminalStdout("./", "log.log")
	New(writer, WithLevel("debug"))

	Log.Info("hello", "标题", "路多辛的博客", "2312", "11")
}
