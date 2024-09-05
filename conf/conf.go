package ktconf

import (
	"os"

	"github.com/aiagt/kitextool/log"
)

var dynamicConfigName = "config"

type Conf interface {
	ParseServerConf(string) error // parse the default config from the string content
	GetServerConf() *ServerConf   // get the default config
	ParseClientConf(string) error
	GetClientConf() *ClientConf
}

type Callback func(conf Conf)

type ConfigCenter interface {
	Init(conf *CenterConf)
	RegisterCallbacks(callbacks ...Callback)
	Register(dest string, conf Conf)
}

func LoadFiles(conf Conf, files ...string) {
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Warnf("read config file failed: %s", err.Error())
			continue
		}

		err = ParseConf(content, conf)
		if err != nil {
			log.Fatalf("parse config file failed: %s", err.Error())
		}
	}
}

type CenterConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}

func ApplyDynamicConfig(center ConfigCenter, centerConf *CenterConf, dest string, conf Conf) {
	logger := func(c Conf) {
		log.Infof("config changed: %+v", c)
	}
	callbacks := []Callback{logger}

	center.Init(centerConf)
	center.RegisterCallbacks(callbacks...)
	center.Register(dest, conf)
}
