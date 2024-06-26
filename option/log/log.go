package log

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const (
	DefaultFileName   = "log/kitex.log"
	DefaultMaxSize    = 10
	DefaultMaxAge     = 3
	DefaultMaxBackups = 50
)

type Option struct {
}

func (o Option) Apply(conf *ktconf.Default) {
	confLog := conf.Log
	if confLog.FileName == "" {
		confLog.FileName = DefaultFileName
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

	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(LogLevel(confLog.Level))
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   confLog.FileName,
			MaxSize:    confLog.MaxSize,
			MaxBackups: confLog.MaxBackups,
			MaxAge:     confLog.MaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		_ = asyncWriter.Sync()
	})
}

func (o Option) OnChange(conf *ktconf.Default) {
}

// WithLogger set the logger through global config
func WithLogger() Option {
	return Option{}
}

func LogLevel(level string) klog.Level {
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
