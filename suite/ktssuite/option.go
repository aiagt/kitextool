package ktssuite

import ktconf "github.com/ahaostudy/kitextool/conf"

type Option func(suite *KitexToolSuite)

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ktconf.ConfigCenter) Option {
	return func(suite *KitexToolSuite) {
		center.RegisterConfigCallback(ktconf.GlobalDefaultConf().Server.Name, ktconf.GlobalConf())
	}
}
