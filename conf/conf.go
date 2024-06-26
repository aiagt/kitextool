package ktconf

var (
	dynamicConfigName = "config"
	globalConf        Conf
	callbacks         []Callback
)

type Conf interface {
	ParseDefault(string) error // parse the default config from the string content
	LoadDefault(...string)     // parse the default config from multiple files
	GetDefault() *Default      // get the default config
}

type ConfigCenter interface {
	RegisterConfigCallback(dest string, conf Conf, callbacks []Callback)
}

type Callback func(conf *Default)

func RegisterGlobalConf(conf Conf) {
	globalConf = conf
}

func GlobalConf() Conf {
	if globalConf == nil {
		globalConf = &Default{}
	}
	return globalConf
}

func GlobalDefaultConf() *Default {
	return GlobalConf().GetDefault()
}

type Option struct {
	center ConfigCenter
}

func (o *Option) Apply(conf *Default) {
	o.center.RegisterConfigCallback(GlobalDefaultConf().Server.Name, GlobalConf(), callbacks)
}

func (o *Option) OnChange(conf *Default) {

}

// WithDynamicConfig dynamically fetch config from the config center
func WithDynamicConfig(center ConfigCenter) *Option {
	return &Option{center: center}
}
