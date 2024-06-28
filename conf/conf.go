package ktconf

var (
	dynamicConfigName = "config"
)

type Conf interface {
	ParseDefault(string) error // parse the default config from the string content
	LoadDefault(...string)     // parse the default config from multiple files
	GetDefault() *Default      // get the default config
}

type Callback func(conf *Default)

type ConfigCenter interface {
	Register(dest string, conf Conf, callbacks ...Callback)
}
