package logx

import (
	"io"
	"log/slog"
	"strings"
)

var Log *slog.Logger

type Option func(c *config)

type config struct {
	AddSource   bool
	Level       slog.Level
	HandlerType string
	handler     slog.Handler
}

func WithSource(addSource bool) Option {
	return func(c *config) {
		c.AddSource = addSource
	}
}

func WithLevel(level string) Option {
	return func(c *config) {
		switch strings.ToLower(level) {
		case "info":
			c.Level = slog.LevelInfo
		case "debug":
			c.Level = slog.LevelDebug
		case "warn", "warning":
			c.Level = slog.LevelWarn
		case "error":
			c.Level = slog.LevelError
		}
	}
}

func WithHandlerType(handlerType string) Option {
	return func(c *config) {
		c.HandlerType = handlerType
	}
}

func WithHandler(handler slog.Handler) Option {
	return func(c *config) {
		c.handler = handler
	}
}

func newConfig() config {
	return config{
		AddSource:   false,
		Level:       slog.LevelError,
		HandlerType: "json",
		handler:     nil,
	}
}

func New(writer io.Writer, opts ...Option) {
	conf := newConfig()
	for _, opt := range opts {
		opt(&conf)
	}
	hopts := &slog.HandlerOptions{
		AddSource: conf.AddSource,
		Level:     conf.Level,
	}

	if conf.handler != nil {
		Log = slog.New(conf.handler)
	} else {
		switch conf.HandlerType {
		case "text":
			Log = slog.New(slog.NewTextHandler(writer, hopts))
		default:
			Log = slog.New(slog.NewJSONHandler(writer, hopts))
		}
	}
}
