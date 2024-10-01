package ktserver

import (
	"net"

	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

type KitexToolSuite struct {
	Conf    ktconf.Conf
	SvrOpts []server.Option

	opts []Option
}

func (s *KitexToolSuite) Options() []server.Option {
	conf := s.Conf.GetServerConf()
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
	// NOTE: add more default configuration items

	return suite
}

func NewKitexToolEmptySuite(conf ktconf.Conf, opts ...Option) *KitexToolSuite {
	suite := &KitexToolSuite{
		Conf: conf,
		opts: opts,
	}

	return suite
}
