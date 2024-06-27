package ktregistry

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/suite/ktssuite"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
)

type Registry func(conf *ktconf.Registry) registry.Registry

type Option struct {
	registry Registry
}

func (o *Option) Apply(s *ktssuite.KitexToolSuite, conf *ktconf.Default) {
	s.SvrOpts = append(s.SvrOpts, server.WithRegistry(o.registry(&conf.Registry)))
}

func (o *Option) OnChange() ktconf.Callback {
	return func(conf *ktconf.Default) {
		klog.Warn("dynamic registration service is not supported yet")
	}
}

func WithRegistry(r Registry) ktssuite.Option {
	return &Option{registry: r}
}
