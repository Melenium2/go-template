package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
)

type LoggerOption func(*loggerConfig)

type loggerConfig struct {
	Level   string `env:"LOGGER_LEVEL" envDefault:"info"`
	PathLen uint8  `env:"LOGGER_SOURCE_LEN" envDefault:"3"`
}

func newLoggerConfig() loggerConfig {
	var cfg loggerConfig

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func (c *loggerConfig) apply(options ...LoggerOption) {
	for _, opt := range options {
		opt(c)
	}
}

// WithLevel sets custom logger level. As the main logger we use log/slog logger from
// std library. You can provide any logger level that supports by log/slog pkg.
// These are the INFO, ERROR, WARN, DEBUG levels. If the specified level can not be
// matched to levels from log/slog, the default level will be used.
//
// Default: INFO.
func WithLevel(level string) LoggerOption {
	return func(c *loggerConfig) {
		c.Level = level
	}
}

// WithSourceLen provides a choice of how detailed the source of the logger
// will be described. By default log/slog logger will print full source path,
// it can be a little bit long. For example, on a unix system the path starts
// from your user home directory (/Users/anyuser/dev/../../../logger/file.go).
// For example, we provide SoruceLen = 2, then our source path will be 'logger/file.go'.
//
// Default: 3.
func WithSourceLen(n uint8) LoggerOption {
	return func(c *loggerConfig) {
		c.PathLen = n
	}
}

// SetupLogger setup the default log/slog logger and overwrite it with our
// own updated copy. After call this function, you can access the custom version
// of the logger using the default log/slog package functions.
//
// Example:
//
//	 import (
//		"context"
//		"log/slog"
//
//		"pkg/logger/tools"
//	 )
//
//	 func main() {
//	    tools.SetupLogger() <- customizing default logger here.
//
//	    slog.Info("message")
//	    slog.WarnContext(context.Background(), "warn message")
//	    slog.Error("error message", "any", "key", "value", "here")
//	    slog.Debug("debug message", slog.Int64("key"))
//	 }
//
// More usage examples can be found here
// https://pkg.go.dev/log/slog.
//
// Benefits of use custom version of logger:
// - Custom logger level
// - Ability to print shortest version of the source
// - Auto conversion log time to UTC
func SetupLogger(options ...LoggerOption) {
	cfg := newLoggerConfig()

	cfg.apply(options...)

	opts := slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel(cfg.Level),
		ReplaceAttr: replaceAttrFunc(cfg),
	}

	customHandler := &customSlogHandler{handler: slog.NewJSONHandler(os.Stdout, &opts)}

	l := slog.New(customHandler)

	slog.SetDefault(l)
}

func logLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// replaceAttrFunc is called to rewrite record attributes and can be customized inside this function.
// For example, we can `blur` user sensitivity data.
func replaceAttrFunc(c loggerConfig) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			a = replaceSourceFunc(a, c.PathLen)
		}

		if a.Key == slog.TimeKey {
			a = replaceTimeFunc(a)
		}

		if a.Key == slog.LevelKey {
			a = replaceLevelFunc(a)
		}

		return a
	}
}

// replaceSourceFunc is used to replace the default log source with the shortest version, since
// the default version of the source is very large and unreadable.
func replaceSourceFunc(a slog.Attr, pathLen uint8) slog.Attr {
	v, ok := a.Value.Any().(*slog.Source)
	if !ok {
		return a
	}

	file := v.File

	// Split path to file by separated pieces and reverse it.
	spl := strings.Split(file, string(filepath.Separator))

	for i, j := 0, len(spl)-1; i < j; i, j = i+1, j-1 {
		spl[i], spl[j] = spl[j], spl[i]
	}

	var (
		length = min(len(spl), int(pathLen))
		path   = make([]string, 0, length)
	)

	// Create new source path with specified length.

	for i := length - 1; i >= 0; i-- {
		path = append(path, spl[i])
	}

	resource := filepath.Join(path...)

	a.Value = slog.StringValue(fmt.Sprintf("%s:%d", resource, v.Line))

	return a
}

func replaceTimeFunc(a slog.Attr) slog.Attr {
	v := a.Value.Time()

	a.Value = slog.TimeValue(v.UTC())

	return a
}

func replaceLevelFunc(a slog.Attr) slog.Attr {
	v := a.Value.String()

	a.Value = slog.StringValue(strings.ToLower(v))

	return a
}

// customSlogHandler can be used to add specific attributes to the logger record
// each time the logger is called. We use it to add specific parameters that can be presented
// within a context.
type customSlogHandler struct {
	handler slog.Handler
}

func (h *customSlogHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

func (h *customSlogHandler) Handle(ctx context.Context, rec slog.Record) error {
	return h.handler.Handle(ctx, rec)
}

func (h *customSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customSlogHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *customSlogHandler) WithGroup(group string) slog.Handler {
	return &customSlogHandler{handler: h.handler.WithGroup(group)}
}
