package ktssuite

import (
	"net"

	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

type KitexToolSuite struct {
	Conf    ktconf.Conf
	SvrOpts []server.Option

	opts      []Option
	callbacks []ktconf.Callback
}

func (s *KitexToolSuite) Options() []server.Option {
	conf := s.Conf.GetDefault()
	for _, opt := range s.opts {
		opt.Apply(s, conf)
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

func NewKitexToolSuite(conf ktconf.Conf, opts ...Option) *KitexToolSuite {
	suite := NewKitexToolEmptySuite(conf, opts...)
	suite.opts = append(suite.opts, WithLogger())
	return suite
}

func NewKitexToolEmptySuite(conf ktconf.Conf, opts ...Option) *KitexToolSuite {
	suite := &KitexToolSuite{
		Conf: conf,
		opts: opts,
	}
	return suite
}
