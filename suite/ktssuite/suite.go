package ktssuite

import (
	"net"

	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

type KitexToolSuite struct {
	opts    []Option
	SvrOpts []server.Option
}

func (s *KitexToolSuite) Options() []server.Option {
	conf := ktconf.GlobalDefaultConf()
	for _, opt := range s.opts {
		opt.Apply(s, conf)
		ktconf.RegisterCallback(opt.OnChange())
	}
	if conf.Server.Address != "" {
		addr, err := net.ResolveTCPAddr("tcp", conf.Server.Address)
		if err != nil {
			panic(err)
		}
		s.SvrOpts = append(s.SvrOpts, server.WithServiceAddr(addr))
	}
	if conf.Server.Name != "" {
		s.SvrOpts = append(s.SvrOpts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.Server.Name}))
	}
	return s.SvrOpts
}

func NewKitexToolSuite(opts ...Option) *KitexToolSuite {
	opts = append(opts, WithLogger())
	suite := &KitexToolSuite{
		opts: opts,
	}
	return suite
}

func NewKitexToolEmptySuite(opts ...Option) *KitexToolSuite {
	suite := &KitexToolSuite{
		opts: opts,
	}
	return suite
}
