package logx

import (
	"context"
	"io"
	"log/slog"
)

var logger *slog.Logger

type Option func(opt *slog.HandlerOptions)

func WithSource(ok bool) Option {
	return func(opt *slog.HandlerOptions) {
		opt.AddSource =
	}
}

func New(opt slog.HandlerOptions) {

}

func test(level slog.Level, addSource bool,writer io.Writer) {
	opts := slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	}

	loggerJSON := slog.New(slog.NewJSONHandler(writer, &opts))

	loggerJSON = loggerJSON.WithGroup("g1").With("k1", 1).WithGroup("g2").With("k2", 2)
	loggerJSON.Info("hello", "标题", "路多辛的博客", "2312", "11")
	loggerJSON.LogAttrs(context.Background(), slog.LevelDebug, "k1", slog.String("dsa", "sss"))

}
