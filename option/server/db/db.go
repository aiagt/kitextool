package ktdb

import (
	"context"
	"errors"
	"fmt"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/aiagt/kitextool/log"
	ktserver "github.com/aiagt/kitextool/suite/server"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbs           map[string]*gorm.DB
	defaultDBName string
)

func GetDB(name string) (*gorm.DB, error) {
	if db, ok := dbs[name]; ok {
		return nil, fmt.Errorf("the db-%s is not exists", name)
	} else {
		return db, nil
	}
}

func GetDBCtx(ctx context.Context, name string) (*gorm.DB, error) {
	db, err := GetDB(name)
	if err != nil {
		return nil, err
	}

	return db.WithContext(ctx), nil
}

func DB() *gorm.DB {
	db, err := GetDB(defaultDBName)
	if err != nil {
		panic(err)
	}

	return db
}

func DBCtx(ctx context.Context) *gorm.DB {
	return DB().WithContext(ctx)
}

func SetDefaultDBName(name string) {
	defaultDBName = name
}

type Option struct {
	ktserver.EmptyOption
	dial     GormDial
	gormConf *gorm.Config
	opts     []DBOption
}

func (o *Option) Apply(s *ktserver.KitexToolSuite, conf *ktconf.ServerConf) {
	o.gormConf = &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(GormLogLevel(conf.Log.LogLevel())),
	}
	for _, opt := range o.opts {
		opt(conf, o.gormConf)
	}

	confDBs := conf.DBs
	if len(confDBs) == 0 {
		if conf.DB == nil {
			log.Fatalf("the database config is empty")
		}

		confDBs = []*ktconf.DB{{DSN: conf.DB.DSN}}
	}

	dbs = make(map[string]*gorm.DB, len(confDBs))

	for _, confDB := range confDBs {
		db, err := o.connect(confDB.DSN)
		if err != nil {
			log.Fatalf("failed to connect to database with DSN: %s", confDB.DSN)
		}

		dbs[confDB.Name] = db
	}
}

func (o *Option) connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn is empty")
	}

	db, err := gorm.Open(o.dial(dsn), o.gormConf)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type GormDial func(dsn string) gorm.Dialector

type DBOption func(conf *ktconf.ServerConf, gormConf *gorm.Config)

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

func WithDB(dial GormDial, opts ...DBOption) ktserver.Option {
	return &Option{dial: dial, opts: opts}
}

func WithGormConf(gormConf *gorm.Config) DBOption {
	return func(conf *ktconf.ServerConf, gconf *gorm.Config) {
		if gormConf.Logger == nil {
			gormConf.Logger = gconf.Logger
		}

		*gconf = *gormConf
	}
}
