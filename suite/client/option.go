package ktclient

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/log"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

type Option interface {
	Apply(s *KitexToolSuite, conf *ktconf.ClientConf)
}

type TransportOption struct {
	protocol transport.Protocol
}

func (o *TransportOption) Apply(s *KitexToolSuite, conf *ktconf.ClientConf) {
	switch o.protocol {
	case transport.TTHeader:
	case transport.Framed, transport.TTHeaderFramed:
		s.CliOpts = append(s.CliOpts, client.WithMetaHandler(transmeta.ClientTTHeaderHandler))
	case transport.GRPC:
		s.CliOpts = append(s.CliOpts, client.WithMetaHandler(transmeta.ClientHTTP2Handler))
	default:
		log.Warnf("unsupported transport protocol: %v, please set it via Kitex option", o.protocol)
		return
	}
	s.CliOpts = append(s.CliOpts, client.WithTransportProtocol(o.protocol))
}

// WithTransport set up the transport protocol and automatically add meta handler
func WithTransport(protocol transport.Protocol) Option {
	return &TransportOption{protocol: protocol}
}
