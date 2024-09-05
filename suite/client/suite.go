package ktclient

import (
	ktconf "github.com/aiagt/kitextool/conf"
	"github.com/cloudwego/kitex/client"
)

type KitexToolSuite struct {
	Conf    *ktconf.ClientConf
	CliOpts []client.Option

	opts []Option
}

func (s *KitexToolSuite) Options() []client.Option {
	for _, opt := range s.opts {
		opt.Apply(s, s.Conf)
	}

	return s.CliOpts
}

func NewKitexToolSuite(conf *ktconf.ClientConf, opts ...Option) *KitexToolSuite {
	suite := NewKitexToolEmptySuite(conf, opts...)
	return suite
}

func NewKitexToolEmptySuite(conf *ktconf.ClientConf, opts ...Option) *KitexToolSuite {
	suite := &KitexToolSuite{
		Conf: conf,
		opts: opts,
	}

	return suite
}
