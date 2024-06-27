package ktlog

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

const (
	DefaultFileName   = "log/kitex.log"
	DefaultMaxSize    = 10
	DefaultMaxAge     = 3
	DefaultMaxBackups = 50
)

func SetLogger(conf *ktconf.Default) {
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
	klog.SetLevel(KlogLevel(confLog.LogLevel()))
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

func KlogLevel(level ktconf.LogLevel) klog.Level {
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
