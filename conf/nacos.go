package ktconf

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-nacos/nacos"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfigCenter struct {
	nacos.Client
}

func NewNacosConfigCenter(opts nacos.Options) *NacosConfigCenter {
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
	return &NacosConfigCenter{Client: client}
}

func (c *NacosConfigCenter) RegisterConfigCallback(dest string, conf Conf) {
	param, err := c.Client.ServerConfigParam(&nacos.ConfigParamConfig{
		Category:          dynamicConfigName,
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
		for _, callback := range callbacks {
			callback(conf.GetDefault())
		}
	}, nacos.GetUniqueID())
}
