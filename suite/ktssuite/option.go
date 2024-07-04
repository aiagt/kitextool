package ktssuite

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	ktlog "github.com/ahaostudy/kitextool/log"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/transport"
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

type TransportOption struct {
	EmptyOption
	protocol transport.Protocol
}

func (o *TransportOption) Apply(s *KitexToolSuite, conf *ktconf.Default) {
	switch o.protocol {
	case transport.TTHeader:
	case transport.Framed, transport.TTHeaderFramed:
		s.SvrOpts = append(s.SvrOpts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))
	case transport.GRPC:
		s.SvrOpts = append(s.SvrOpts, server.WithMetaHandler(transmeta.ServerHTTP2Handler))
	default:
		klog.Warnf("[KitexTool] unsupported transport protocol: %v, please set it via Kitex option", o.protocol)
	}
}

// WithLogger set the logger through global config
func WithLogger(opts ...ktlog.LoggerOption) Option {
	return LogOption{opts: opts}
}

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ktconf.ConfigCenter) Option {
	return &ConfigOption{center: center}
}

// WithTransport set up the transport protocol and automatically add meta handler
func WithTransport(protocol transport.Protocol) Option {
	return &TransportOption{protocol: protocol}
}
