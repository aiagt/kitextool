package ktconf

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

type Parser struct{}

func (Parser) Decode(kind vo.ConfigType, data string, config interface{}) error {
	switch kind {
	case vo.YAML, vo.JSON:
		// since YAML is a superset of JSON, it can parse JSON using a YAML parser
		return Parse([]byte(data), config)
	default:
		return fmt.Errorf("unsupported config data type %s", kind)
	}
}

var defaultParser = Parser{}

func DefaultParser() Parser {
	return Parser{}
}

func Parse(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}
