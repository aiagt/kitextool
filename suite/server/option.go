package ktserver

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/log"
	ktlog "github.com/ahaostudy/kitextool/option/server/log"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/transport"
)

type Option interface {
	Apply(s *KitexToolSuite, conf *ktconf.ServerConf)
	Callback() ServerCallback
}

type ServerCallback func(conf *ktconf.ServerConf)

func confCallback(callback ServerCallback) ktconf.Callback {
	return func(conf ktconf.Conf) {
		callback(conf.GetServerConf())
	}
}

type EmptyOption struct{}

func (o EmptyOption) Apply(s *KitexToolSuite, conf *ktconf.ServerConf) {}

func (o EmptyOption) Callback() ServerCallback {
	return func(conf *ktconf.ServerConf) {}
}

type LogOption struct {
	EmptyOption
	opts []ktlog.LoggerOption
}

func (o LogOption) Apply(s *KitexToolSuite, conf *ktconf.ServerConf) {
	ktlog.SetLogger(conf, o.opts...)
}

// WithLogger set the logger through global config
func WithLogger(opts ...ktlog.LoggerOption) Option {
	return LogOption{opts: opts}
}

type ConfigOption struct {
	EmptyOption
	center ktconf.ConfigCenter
}

func (o *ConfigOption) Apply(s *KitexToolSuite, conf *ktconf.ServerConf) {
	var callbacks []ktconf.Callback
	for _, opt := range s.opts {
		callbacks = append(callbacks, confCallback(opt.Callback()))
	}
	ktconf.ApplyDynamicConfig(o.center, &conf.GetServerConf().ConfigCenter, conf.Server.Name, s.Conf)
}

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ktconf.ConfigCenter) Option {
	return &ConfigOption{center: center}
}

type TransportOption struct {
	EmptyOption
	protocol transport.Protocol
}

func (o *TransportOption) Apply(s *KitexToolSuite, conf *ktconf.ServerConf) {
	switch o.protocol {
	case transport.TTHeader:
	case transport.Framed, transport.TTHeaderFramed:
		s.SvrOpts = append(s.SvrOpts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))
	case transport.GRPC:
		s.SvrOpts = append(s.SvrOpts, server.WithMetaHandler(transmeta.ServerHTTP2Handler))
	default:
		log.Warnf("unsupported transport protocol: %v, please set it via Kitex option", o.protocol)
	}
}

// WithTransport set up the transport protocol and automatically add meta handler
func WithTransport(protocol transport.Protocol) Option {
	return &TransportOption{protocol: protocol}
}
