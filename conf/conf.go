package ktconf

const (
	limiterConfigName = "config"
)

type Conf interface {
	ParseDefault(string) error // parse the default config from the string content
	LoadDefault(...string)     // parse the default config from multiple files
	GetDefault() *Default      // get the default config
}

type ConfigCenter interface {
	RegisterConfigCallback(dest string, conf Conf)
}

var globalConf Conf

func RegisterGlobalConf(conf Conf) {
	globalConf = conf
}

func GlobalConf() Conf {
	return globalConf
}

func GlobalDefaultConf() *Default {
	return globalConf.GetDefault()
}
