package main

import (
	"context"
	"log"
	"time"

	ktcenter "github.com/aiagt/kitextool/conf/center"

	ktregistry "github.com/aiagt/kitextool/option/server/registry"

	ktconf "github.com/aiagt/kitextool/conf"
	echo "github.com/aiagt/kitextool/example/kitex_gen/echo/echoservice"
	ktdb "github.com/aiagt/kitextool/option/server/db"
	ktrdb "github.com/aiagt/kitextool/option/server/redis"
	ktserver "github.com/aiagt/kitextool/suite/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

// Configuration must be of type ktconf.Conf, you can directly use ktconf.ServerConf, which implements ktconf.Conf.
// Additionally, you can customize the structure of the configuration, but it should embed ktconf.ServerConf in its structure, example below:
// ```go
//
//	type Conf struct {
//	    ktconf.ServerConf
//	    File FileConf `yaml:"file"`
//	}
//
//	type FileConf struct {
//	    Path string `yaml:"path"`
//	}
//
// ```
var conf = new(ktconf.ServerConf)

func init() {
	// Load configuration content from a specified files
	ktconf.LoadFiles(conf, "conf.yaml", "server/conf.yaml")
}

func main() {
	svr := echo.NewServer(
		new(EchoServiceImpl),
		// Use the KitexTool suite in the server
		server.WithSuite(ktserver.NewKitexToolSuite(
			// Global configuration of KitexTool suite
			conf,
			// Use consul dynamic configuration
			ktserver.WithDynamicConfig(ktcenter.WithConsulConfigCenter(nil)),
			// Use consul as the registry
			ktregistry.WithRegistry(ktregistry.NewConsulRegistry()),
			// Introduce MySQL database into the project
			ktdb.WithDB(ktdb.NewMySQLDial()),
			// Introduce Redis into the project
			ktrdb.WithRedis(),
		)),
	)

	migrateUserModel()

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

type EchoServiceImpl struct{}

func (s *EchoServiceImpl) Echo(ctx context.Context, message string) (resp string, err error) {
	err = dbExample(ctx)
	if err != nil {
		return "", err
	}

	err = redisExample(ctx)
	if err != nil {
		return "", err
	}

	return message, nil
}

type User struct {
	ID   uint64 `gorm:"primary_key;auto_increment"`
	Name string
}

func migrateUserModel() {
	err := ktdb.DB().AutoMigrate(new(User))
	if err != nil {
		panic(err)
	}
}

func dbExample(ctx context.Context) error {
	// KitexTool automatically initializes global DB objects based on the configuration
	// You can access it directly via `ktdb.DB()`
	err := ktdb.DB().WithContext(ctx).Create(&User{Name: "user1"}).Error
	if err != nil {
		return err
	}

	user := new(User)

	err = ktdb.DB().WithContext(ctx).First(user).Error
	if err != nil {
		return err
	}

	klog.CtxInfof(ctx, "[DB] result: %+v", user)

	return nil
}

func redisExample(ctx context.Context) error {
	// KitexTool automatically initializes global redis clients based on the configuration
	// You can access it directly via `ktrdb.RDB()`
	err := ktrdb.RDB().Set(ctx, "user:1", "{\"name\":\"user1\"}", 5*time.Minute).Err()
	if err != nil {
		return err
	}

	result := ktrdb.RDB().Get(ctx, "user:1")
	if err = result.Err(); err != nil {
		return err
	}

	klog.CtxInfof(ctx, "[Redis] result: %s", result.String())

	return nil
}
