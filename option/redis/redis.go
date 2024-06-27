package ktrdb

import (
	"context"
	"errors"

	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/suite/ktssuite"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

var globalRDB *redis.Client

func GetRDB() (*redis.Client, error) {
	if globalRDB == nil {
		return nil, errors.New("the global rdb is not initialized")
	}
	return globalRDB, nil
}

func RDB() *redis.Client {
	rdb, err := GetRDB()
	if err != nil {
		panic(err)
	}
	return rdb
}

type Option struct{}

func (o Option) Apply(s *ktssuite.KitexToolSuite, conf *ktconf.Default) {
	reconnect(&conf.Redis)
}

func (o Option) OnChange() ktconf.Callback {
	return func(conf *ktconf.Default) {
		reconnect(&conf.Redis)
	}
}

func reconnect(conf *ktconf.Redis) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		klog.Warnf("failed to connect to redis: %s\n", err.Error())
		return
	}
	globalRDB = rdb
}

func WithRedis() ktssuite.Option {
	return Option{}
}
