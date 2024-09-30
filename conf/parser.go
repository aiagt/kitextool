package ktconf

import (
	"gopkg.in/yaml.v2"
)

func Parse(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

func ParseConf(data string, conf Conf) error {
	err := conf.ParseServerConf(data)
	if err != nil {
		return err
	}

	err = conf.ParseClientConf(data)
	if err != nil {
		return err
	}

	err = Parse([]byte(data), conf)
	if err != nil {
		return err
	}

	return nil
}
