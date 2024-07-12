package ktrdb

import (
	"context"
	"fmt"
	"github.com/ahaostudy/kitextool/log"

	ktconf "github.com/ahaostudy/kitextool/conf"
	ktserver "github.com/ahaostudy/kitextool/suite/server"
	"github.com/redis/go-redis/v9"
)

var rdbs []*redis.Client

func GetRDB(id int) (*redis.Client, error) {
	if len(rdbs) <= id {
		return nil, fmt.Errorf("the rdb-%d is not exists", id)
	}
	return rdbs[id], nil
}

func RDB() *redis.Client {
	rdb, err := GetRDB(0)
	if err != nil {
		panic(err)
	}
	return rdb
}

type Option struct {
	ktserver.EmptyOption
}

func (o Option) Apply(s *ktserver.KitexToolSuite, conf *ktconf.ServerConf) {
	confRedises := conf.Redises
	if len(confRedises) == 0 {
		if conf.Redis == nil {
			log.Fatalf("the redis config is empty")
		}
		confRedises = append(confRedises, *conf.Redis)
	}
	rdbs = make([]*redis.Client, len(confRedises))
	for i, redisConf := range confRedises {
		rdb, err := connect(&redisConf)
		if err != nil {
			log.Fatalf("failed to connect redis-%d: %s", i, err.Error())
		}
		rdbs[i] = rdb
	}
}

func connect(conf *ktconf.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return rdb, nil
}

func WithRedis() ktserver.Option {
	return Option{}
}
