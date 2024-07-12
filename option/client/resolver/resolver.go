package ktresolver

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	ktclient "github.com/ahaostudy/kitextool/suite/client"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
)

type Resolver func(conf *ktconf.Resolver) discovery.Resolver

type Option struct {
	resolver Resolver
}

func (o *Option) Apply(s *ktclient.KitexToolSuite, conf *ktconf.ClientConf) {
	s.CliOpts = append(s.CliOpts, client.WithResolver(o.resolver(&conf.Resolver)))
}

func WithResolver(r Resolver) ktclient.Option {
	return &Option{resolver: r}
}
