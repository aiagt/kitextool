package ktssuite

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/option/log"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"net"
)

type KitexToolSuite struct {
	opts []Option
}

func (s *KitexToolSuite) Options() (opts []server.Option) {
	conf := ktconf.GlobalDefaultConf()
	for _, opt := range s.opts {
		opt.Apply(conf)
	}
	if conf.Server.Address != "" {
		addr, err := net.ResolveTCPAddr("tcp", conf.Server.Address)
		if err != nil {
			panic(err)
		}
		opts = append(opts, server.WithServiceAddr(addr))
	}
	if conf.Server.Name != "" {
		opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.Server.Name}))
	}
	return opts
}

func NewKitexToolSuite(opts ...Option) *KitexToolSuite {
	opts = append(opts, log.WithLogger())
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
