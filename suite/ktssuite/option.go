package ktssuite

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	ktlog "github.com/ahaostudy/kitextool/log"
	"github.com/cloudwego/kitex/pkg/klog"
)

type Option interface {
	Apply(s *KitexToolSuite, conf *ktconf.Default)
	OnChange() ktconf.Callback
}

type EmptyOption struct{}

func (o EmptyOption) Apply(s *KitexToolSuite, conf *ktconf.Default) {}

func (o EmptyOption) OnChange() ktconf.Callback {
	return func(conf *ktconf.Default) {}
}

type LogOption struct {
	EmptyOption
	opts []ktlog.LoggerOption
}

func (o LogOption) Apply(s *KitexToolSuite, conf *ktconf.Default) {
	ktlog.SetLogger(conf, o.opts...)
}

type ConfigOption struct {
	EmptyOption
	center ktconf.ConfigCenter
}

func (o *ConfigOption) Apply(s *KitexToolSuite, conf *ktconf.Default) {
	logger := func(conf *ktconf.Default) {
		klog.Infof("config changed: %+v\n", s.Conf)
	}
	s.callbacks = append(s.callbacks, logger)
	for _, opt := range s.opts {
		s.callbacks = append(s.callbacks, opt.OnChange())
	}
	o.center.Register(conf.Server.Name, s.Conf)
}

// WithLogger set the logger through global config
func WithLogger(opts ...ktlog.LoggerOption) Option {
	return LogOption{opts: opts}
}

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ktconf.ConfigCenter) Option {
	return &ConfigOption{center: center}
}
