package ktresolver

import (
	"net"
	"strconv"

	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/ahaostudy/kitextool/log"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosResolver(conf *ktconf.Resolver) discovery.Resolver {
	if len(conf.Address) == 0 {
		r, err := resolver.NewDefaultNacosResolver()
		if err != nil {
			log.Fatalf("service resolver failed: %s", err.Error())
		}
		return r
	}

	host, portStr, err := net.SplitHostPort(conf.Address[0])
	if err != nil {
		panic(err)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	port, err := strconv.ParseUint(portStr, 10, 64)
	if err != nil {
		panic(err)
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port),
	}
	cc := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		Username:            conf.Username,
		Password:            conf.Password,
	}
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return resolver.NewNacosResolver(cli)
}
