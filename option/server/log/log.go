package ktlog

import (
	"io"
	"os"
	"time"

	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	kitexslog "github.com/kitex-contrib/obs-opentelemetry/logging/slog"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultMaxSize    = 10
	DefaultMaxAge     = 3
	DefaultMaxBackups = 50
)

type LoggerOption func(*ktconf.ServerConf)

var (
	WithLogrus LoggerOption = func(_ *ktconf.ServerConf) {
		logger := kitexlogrus.NewLogger()
		klog.SetLogger(logger)
	}
	WithZap LoggerOption = func(_ *ktconf.ServerConf) {
		logger := kitexzap.NewLogger()
		klog.SetLogger(logger)
	}
	WithSlog LoggerOption = func(_ *ktconf.ServerConf) {
		logger := kitexslog.NewLogger()
		klog.SetLogger(logger)
	}
)

func SetLogger(conf *ktconf.ServerConf, opts ...LoggerOption) {
	for _, opt := range opts {
		opt(conf)
	}

	confLog := conf.Log
	if confLog.FileName == "" {
		confLog.FileName = filepathOption(conf)
	}

	if confLog.MaxSize == 0 {
		confLog.MaxSize = DefaultMaxSize
	}

	if confLog.MaxAge == 0 {
		confLog.MaxAge = DefaultMaxAge
	}

	if confLog.MaxBackups == 0 {
		confLog.MaxBackups = DefaultMaxBackups
	}

	klog.SetLevel(KLogLevel(confLog.LogLevel()))
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   confLog.FileName,
			MaxSize:    confLog.MaxSize,
			MaxBackups: confLog.MaxBackups,
			MaxAge:     confLog.MaxAge,
		}),
		FlushInterval: time.Minute,
	}

	output := io.MultiWriter(os.Stdout, asyncWriter)
	klog.SetOutput(output)
	server.RegisterShutdownHook(func() {
		_ = asyncWriter.Sync()
	})
}

func KLogLevel(level ktconf.LogLevel) klog.Level {
	switch level {
	case ktconf.LevelTrace:
		return klog.LevelTrace
	case ktconf.LevelDebug:
		return klog.LevelDebug
	case ktconf.LevelInfo:
		return klog.LevelInfo
	case ktconf.LevelNotice:
		return klog.LevelNotice
	case ktconf.LevelWarn:
		return klog.LevelWarn
	case ktconf.LevelError:
		return klog.LevelError
	case ktconf.LevelFatal:
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
