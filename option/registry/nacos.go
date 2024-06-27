package ktregistry

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/registry"
	nacosregistry "github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"strconv"
)

func NewNacosRegistry() Registry {
	return func(conf *ktconf.Registry) registry.Registry {
		if len(conf.RegistryAddress) == 0 {
			panic("service registration failed, registry address is empty")
		}

		host, portstr, err := net.SplitHostPort(conf.RegistryAddress[0])
		if err != nil {
			panic(err)
		}
		if host == "" {
			host = "127.0.0.1"
		}
		port, err := strconv.ParseUint(portstr, 10, 64)
		if err != nil {
			panic(err)
		}

		var sc = []constant.ServerConfig{
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
		return nacosregistry.NewNacosRegistry(cli)
	}
}
