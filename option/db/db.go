package ktdb

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/suite/ktssuite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	globalDB *gorm.DB
)

func DB() *gorm.DB {
	if globalDB == nil {
		panic("the global db is not initialized, please use the `WithDB` option to initialize the global db.")
	}
	return globalDB
}

type Option struct {
	dial gorm.Dialector
	opts []DBOption
}

func (o Option) Apply(conf *ktconf.Default) {
	gormConf := &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(LogLevel()),
	}
	for _, opt := range o.opts {
		opt(conf, gormConf)
	}
	var err error
	globalDB, err = gorm.Open(o.dial, gormConf)
	if err != nil {
		panic(err)
	}
}

func (o Option) OnChange(conf *ktconf.Default) {
}

func WithDB(dial gorm.Dialector, opts ...DBOption) ktssuite.Option {
	return Option{dial: dial, opts: opts}
}

type DBOption func(conf *ktconf.Default, gconf *gorm.Config)

func WithLogger(gormConf *gorm.Config) DBOption {
	return func(conf *ktconf.Default, gconf *gorm.Config) {
		if gormConf.Logger == nil {
			gormConf.Logger = gconf.Logger
		}
		*gconf = *gormConf
	}
}

func LogLevel() logger.LogLevel {
	level := ktconf.GlobalDefaultConf().Log.Level
	switch level {
	case "trace":
		return logger.Info
	case "debug":
		return logger.Info
	case "info":
		return logger.Info
	case "notice":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "fatal":
		return logger.Error
	default:
		return logger.Warn
	}
}
