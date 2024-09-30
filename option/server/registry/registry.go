package ktregistry

import (
	ktconf "github.com/aiagt/kitextool/conf"
	ktserver "github.com/aiagt/kitextool/suite/server"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
)

type Registry func(conf *ktconf.Registry) []registry.Registry

type Option struct {
	ktserver.EmptyOption
	registry Registry
}

func (o *Option) Apply(s *ktserver.KitexToolSuite, conf *ktconf.ServerConf) {
	for _, r := range o.registry(&conf.Registry) {
		s.SvrOpts = append(s.SvrOpts, server.WithRegistry(r))
	}
}

func WithRegistry(r Registry) ktserver.Option {
	return &Option{registry: r}
}
