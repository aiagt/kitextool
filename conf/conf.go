package ktconf

var (
	dynamicConfigName = "config"
	globalConf        Conf
	callbacks         []Callback
)

type Callback func(conf *Default)

func RegisterCallback(callback Callback) {
	callbacks = append(callbacks, callback)
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

func RegisterGlobalConf(conf Conf) {
	globalConf = conf
}

type Conf interface {
	ParseDefault(string) error // parse the default config from the string content
	LoadDefault(...string)     // parse the default config from multiple files
	GetDefault() *Default      // get the default config
}

type ConfigCenter interface {
	RegisterConfigCallback(dest string, conf Conf)
}
