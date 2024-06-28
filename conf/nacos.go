package ktconf

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-nacos/nacos"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sync"
)

type NacosConfigCenter struct {
	opts   nacos.Options
	client nacos.Client
	once   sync.Once
}

func NewNacosConfigCenter(opts nacos.Options) *NacosConfigCenter {
	c := &NacosConfigCenter{opts: opts}
	return c
}

func (c *NacosConfigCenter) Client() nacos.Client {
	if c.client == nil {
		panic("the nacos client is not initialized")
	}
	return c.client
}

func (c *NacosConfigCenter) InitClient(conf *Center) {
	c.once.Do(func() {
		opts := c.opts
		if conf != nil {
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
		}

		client, err := nacos.NewClient(opts)
		if err != nil {
			panic(err)
		}
		c.client = client
	})
}

func (c *NacosConfigCenter) Register(dest string, conf Conf, callbacks ...Callback) {
	c.InitClient(&conf.GetDefault().ConfigCenter)

	param, err := c.Client().ServerConfigParam(&nacos.ConfigParamConfig{
		Category:          dynamicConfigName,
		ServerServiceName: dest,
	})
	if err != nil {
		panic(err)
	}
	c.Client().RegisterConfigCallback(param, func(data string, parser nacos.ConfigParser) {
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
