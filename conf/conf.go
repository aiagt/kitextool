package ktconf

import (
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
)

var dynamicConfigName = "config"

type Conf interface {
	ParseDefault(string) error // parse the default config from the string content
	GetDefault() *Default      // get the default config
}

type Callback func(conf *Default)

type ConfigCenter interface {
	Register(dest string, conf Conf, callbacks ...Callback)
}

func LoadFiles(conf Conf, files ...string) {
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			klog.Warnf("read config file failed: %s", err.Error())
			continue
		}
		err = conf.GetDefault().ParseDefault(string(content))
		if err != nil {
			panic(err)
		}
		err = Parse(content, conf)
		if err != nil {
			panic(err)
		}
	}
}
