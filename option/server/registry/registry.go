package ktregistry

import (
	ktconf "github.com/aiagt/kitextool/conf"
	ktserver "github.com/aiagt/kitextool/suite/server"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
)

type Registry func(conf *ktconf.Registry) registry.Registry

type Option struct {
	ktserver.EmptyOption
	registry Registry
}

func (o *Option) Apply(s *ktserver.KitexToolSuite, conf *ktconf.ServerConf) {
	s.SvrOpts = append(s.SvrOpts, server.WithRegistry(o.registry(&conf.Registry)))
}

func WithRegistry(r Registry) ktserver.Option {
	return &Option{registry: r}
}
