package ktcenter

import (
	"fmt"
	ktconf "github.com/aiagt/kitextool/conf"
)

type ConfigType string

const (
	PROPERTIES ConfigType = "properties"
	XML        ConfigType = "xml"
	JSON       ConfigType = "json"
	TEXT       ConfigType = "text"
	HTML       ConfigType = "html"
	YAML       ConfigType = "yaml"
)

type Parser[T any] struct {
	Translate func(T) ConfigType
}

func NewParser[T any](t func(T) ConfigType) *Parser[T] {
	return &Parser[T]{Translate: t}
}

func NewStringParser[T ~string]() *Parser[T] {
	return &Parser[T]{Translate: TranslateString[T]}
}

func TranslateString[T ~string](t T) ConfigType {
	return ConfigType(t)
}

func (p *Parser[T]) Decode(kind T, data string, config interface{}) error {
	k := p.Translate(kind)

	switch k {
	case YAML, JSON:
		// since YAML is a superset of JSON, it can parse JSON using a YAML parser
		return ktconf.Parse([]byte(data), config)
	default:
		return fmt.Errorf("unsupported config data type %s", k)
	}
}
