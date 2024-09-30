package ktresolver

import (
	ktconf "github.com/aiagt/kitextool/conf"
	ktclient "github.com/aiagt/kitextool/suite/client"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
)

type Resolver func(conf *ktconf.Resolver) []discovery.Resolver

type Option struct {
	resolver Resolver
}

func (o *Option) Apply(s *ktclient.KitexToolSuite, conf *ktconf.ClientConf) {
	for _, resolver := range o.resolver(&conf.Resolver) {
		s.CliOpts = append(s.CliOpts, client.WithResolver(resolver))
	}
}

func WithResolver(r Resolver) ktclient.Option {
	return &Option{resolver: r}
}
