package ktconf

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

func Parse(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

func ParseConf(data []byte, conf Conf) error {
	err := conf.ParseServerConf(string(data))
	if err != nil {
		return err
	}

	err = conf.ParseClientConf(string(data))
	if err != nil {
		return err
	}

	err = Parse(data, conf)
	if err != nil {
		return err
	}

	return nil
}

type NacosParser struct{}

func (NacosParser) Decode(kind vo.ConfigType, data string, config interface{}) error {
	switch kind {
	case vo.YAML, vo.JSON:
		// since YAML is a superset of JSON, it can parse JSON using a YAML parser
		return Parse([]byte(data), config)
	default:
		return fmt.Errorf("unsupported config data type %s", kind)
	}
}

func DefaultNacosParser() NacosParser {
	return NacosParser{}
}
