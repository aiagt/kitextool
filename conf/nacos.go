package ktconf

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-nacos/nacos"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ConfigCenterNacos struct {
	nacos.Client
}

func NewConfigCenterNacos(opts nacos.Options) *ConfigCenterNacos {
	conf := GlobalDefaultConf().ConfigCenter
	if opts.Address == "" {
		opts.Address = conf.Host
	}
	if opts.Port == 0 {
		opts.Port = uint64(conf.Port)
	}
	if opts.NamespaceID == "" {
		opts.NamespaceID = conf.Key
	}
	if opts.ConfigParser == nil {
		opts.ConfigParser = DefaultParser()
	}
	client, err := nacos.NewClient(opts)
	if err != nil {
		panic(err)
	}
	return &ConfigCenterNacos{Client: client}
}

func (c *ConfigCenterNacos) RegisterConfigCallback(dest string, conf Conf) {
	param, err := c.Client.ServerConfigParam(&nacos.ConfigParamConfig{
		Category:          limiterConfigName,
		ServerServiceName: dest,
	})
	if err != nil {
		panic(err)
	}
	c.Client.RegisterConfigCallback(param, func(data string, parser nacos.ConfigParser) {
		var err error
		err = conf.ParseDefault(data)
		if err != nil {
			klog.Errorf("conf parse error, %s", err.Error())
			return
		}
		err = parser.Decode(vo.YAML, data, conf)
		if err != nil {
			klog.Errorf("conf parse error, %s", err.Error())
			return
		}
		klog.Infof("conf changed: %+v\n", conf)
	}, nacos.GetUniqueID())
}
