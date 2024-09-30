package ktconf

import (
	"os"

	"github.com/aiagt/kitextool/log"
)

type Conf interface {
	ParseServerConf(string) error // parse the default config from the string content
	GetServerConf() *ServerConf   // get the default config
	ParseClientConf(string) error
	GetClientConf(name string) *ClientConf
}

func LoadFiles(conf Conf, files ...string) {
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Warnf("read config file failed: %s", err.Error())
			continue
		}

		err = ParseConf(string(content), conf)
		if err != nil {
			log.Fatalf("parse config file failed: %s", err.Error())
		}
	}
}
