package ktlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	ktserver "github.com/aiagt/kitextool/suite/server"

	"github.com/aiagt/kitextool/utils"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	kitexslog "github.com/kitex-contrib/obs-opentelemetry/logging/slog"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogOption struct {
	ktserver.EmptyOption
	Logger klog.FullLogger
}

// WithLogger set the Logger through global config
func WithLogger(logger klog.FullLogger) ktserver.Option {
	return &LogOption{Logger: logger}
}

func (o *LogOption) Apply(_ *ktserver.KitexToolSuite, conf *ktconf.ServerConf) {
	klog.SetLogger(o.Logger)
	o.SetLogger(conf)
}

func (o *LogOption) SetLogger(conf *ktconf.ServerConf) {
	if conf.Log.EnableFile != nil && !*conf.Log.EnableFile {
		return
	}

	confLog := conf.Log
	utils.SetDefault(&confLog.FileName, logPath(conf.Server.Name))
	utils.SetDefault(&confLog.MaxSize, DefaultMaxSize)
	utils.SetDefault(&confLog.MaxAge, DefaultMaxAge)
	utils.SetDefault(&confLog.MaxBackups, DefaultMaxBackups)
	utils.SetDefault(&confLog.MaxSize, DefaultMaxSize)
	utils.SetDefault(
		&confLog.FlushInterval,
		utils.Ternary(ktconf.GetEnv() == ktconf.EnvProd, DefaultDevFlushInterval, DefaultProdFlushInterval),
	)

	klog.SetLevel(KLogLevel(confLog.LogLevel()))

	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   confLog.FileName,
			MaxSize:    confLog.MaxSize,
			MaxBackups: confLog.MaxBackups,
			MaxAge:     confLog.MaxAge,
		}),
		FlushInterval: time.Duration(confLog.FlushInterval) * time.Second,
	}

	output := io.MultiWriter(os.Stdout, asyncWriter)
	klog.SetOutput(output)
	server.RegisterShutdownHook(func() {
		_ = asyncWriter.Sync()
	})
}

func NewLogrusLogger(opts ...kitexlogrus.Option) *kitexlogrus.Logger {
	return kitexlogrus.NewLogger(opts...)
}

func NewZapLogger(opts ...kitexzap.Option) *kitexzap.Logger {
	return kitexzap.NewLogger(opts...)
}

func NewSlogLogger(opts ...kitexslog.Option) *kitexslog.Logger {
	return kitexslog.NewLogger(opts...)
}

const (
	DefaultMaxSize           = 10
	DefaultMaxAge            = 3
	DefaultMaxBackups        = 50
	DefaultDevFlushInterval  = 5
	DefaultProdFlushInterval = 60
)

func logPath(svc string) string {
	fileName := fmt.Sprintf("%s.log", svc)
	return filepath.Join("log", fileName)
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
