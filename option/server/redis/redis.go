package ktrdb

import (
	"context"
	"fmt"

	"github.com/aiagt/kitextool/log"

	ktconf "github.com/aiagt/kitextool/conf"
	ktserver "github.com/aiagt/kitextool/suite/server"
	"github.com/redis/go-redis/v9"
)

var (
	rdbs           map[string]*redis.Client
	defaultRDBName string
)

func GetRDB(name string) (*redis.Client, error) {
	if rdb, ok := rdbs[name]; !ok {
		return nil, fmt.Errorf("the rdb-%s is not exists", name)
	} else {
		return rdb, nil
	}
}

func RDB() *redis.Client {
	rdb, err := GetRDB(defaultRDBName)
	if err != nil {
		panic(err)
	}

	return rdb
}

func SetDefaultRDBName(name string) {
	defaultRDBName = name
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

		confRedises = []*ktconf.Redis{conf.Redis}
	}

	rdbs = make(map[string]*redis.Client, len(confRedises))

	for _, redisConf := range confRedises {
		rdb, err := connect(redisConf)
		if err != nil {
			log.Fatalf("failed to connect redis-%s: %s", redisConf.Name, err.Error())
		}

		rdbs[redisConf.Name] = rdb
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
