package ktserver

import (
	ktconf "github.com/aiagt/kitextool/conf"
	ktcenter "github.com/aiagt/kitextool/conf/center"
)

type Option interface {
	Apply(s *KitexToolSuite, conf *ktconf.ServerConf)
	Callback() ServerCallback
}

type ServerCallback func(conf *ktconf.ServerConf)

func confCallback(callback ServerCallback) ktcenter.Callback {
	return func(conf ktconf.Conf) {
		callback(conf.GetServerConf())
	}
}

type EmptyOption struct{}

func (o EmptyOption) Apply(s *KitexToolSuite, conf *ktconf.ServerConf) {}

func (o EmptyOption) Callback() ServerCallback {
	return func(conf *ktconf.ServerConf) {}
}

type ConfigOption struct {
	EmptyOption
	center ktcenter.ConfigCenter
}

func (o *ConfigOption) Apply(s *KitexToolSuite, conf *ktconf.ServerConf) {
	var callbacks []ktcenter.Callback
	for _, opt := range s.opts {
		callbacks = append(callbacks, confCallback(opt.Callback()))
	}

	o.center.RegisterCallbacks(callbacks...)
	ktcenter.ApplyDynamicConfig(o.center, &conf.GetServerConf().ConfigCenter, conf.Server.Name, s.Conf)
}

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ktcenter.ConfigCenter) Option {
	return &ConfigOption{center: center}
}
