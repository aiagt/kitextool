package ktdb

import (
	"errors"
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/suite/ktssuite"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	globalDB *gorm.DB
)

func GetDB() (*gorm.DB, error) {
	if globalDB == nil {
		return nil, errors.New("the global db is not initialized")
	}
	return globalDB, nil
}

func DB() *gorm.DB {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}
	return db
}

type GormDial func(dsn string) gorm.Dialector

type DBOption func(conf *ktconf.Default, gconf *gorm.Config)

type Option struct {
	dial     GormDial
	gormConf *gorm.Config
	opts     []DBOption
}

func (o *Option) Apply(s *ktssuite.KitexToolSuite, conf *ktconf.Default) {
	o.gormConf = &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(GormLogLevel(conf.Log.LogLevel())),
	}
	for _, opt := range o.opts {
		opt(conf, o.gormConf)
	}
	o.reconnect(conf)
}

func (o *Option) OnChange() ktconf.Callback {
	return func(conf *ktconf.Default) {
		o.reconnect(conf)
	}
}

func (o *Option) reconnect(conf *ktconf.Default) {
	klog.Infof("connecting to database with dsn: %s\n", conf.DB.DSN)

	if conf.DB.DSN == "" {
		klog.Error("failed to connect to database: dsn is empty")
		return
	}
	db, err := gorm.Open(
		o.dial(conf.DB.DSN),
		o.gormConf,
	)
	if err != nil {
		klog.Errorf("failed to connect to database: %s\n", err.Error())
		return
	}
	globalDB = db
}

func WithDB(dial GormDial, opts ...DBOption) ktssuite.Option {
	return &Option{dial: dial, opts: opts}
}

func WithGormConf(gormConf *gorm.Config) DBOption {
	return func(conf *ktconf.Default, gconf *gorm.Config) {
		if gormConf.Logger == nil {
			gormConf.Logger = gconf.Logger
		}
		*gconf = *gormConf
	}
}

func GormLogLevel(level ktconf.LogLevel) logger.LogLevel {
	switch level {
	case ktconf.LevelTrace, ktconf.LevelDebug, ktconf.LevelInfo:
		return logger.Info
	case ktconf.LevelNotice, ktconf.LevelWarn:
		return logger.Warn
	case ktconf.LevelError, ktconf.LevelFatal:
		return logger.Error
	default:
		return logger.Warn
	}
}
