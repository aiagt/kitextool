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
}

func (o LogOption) Apply(s *KitexToolSuite, conf *ktconf.Default) {
	ktlog.SetLogger(conf)
}

type ConfigOption struct {
	EmptyOption
	center ktconf.ConfigCenter
}

func (o *ConfigOption) Apply(s *KitexToolSuite, conf *ktconf.Default) {
	ktconf.RegisterCallback(func(conf *ktconf.Default) {
		klog.Infof("config changed: %+v\n", ktconf.GlobalConf())
	})
	o.center.RegisterConfigCallback(conf.Server.Name, ktconf.GlobalConf())
}

// WithLogger set the logger through global config
func WithLogger() Option {
	return LogOption{}
}

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ktconf.ConfigCenter) Option {
	return &ConfigOption{center: center}
}
