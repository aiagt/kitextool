package ktconf

type Default struct {
	Server       Server   `yaml:"server"`
	Log          Log      `yaml:"log"`
	Registry     Registry `yaml:"registry"`
	ConfigCenter Center   `yaml:"config_center"`
	DB           DB       `yaml:"db"`
	Redis        Redis    `yaml:"redis"`
	Rabbitmq     Rabbitmq `yaml:"rabbitmq"`
}

func (d *Default) ParseDefault(data string) error {
	return Parse([]byte(data), d)
}

func (d *Default) GetDefault() *Default {
	return d
}

type Server struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Version string `yaml:"version"`
}

type Log struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"file_name"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type Center struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}

type Registry struct {
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type DB struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Rabbitmq struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"vhost"`
}

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

func (l Log) LogLevel() LogLevel {
	switch l.Level {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "notice":
		return LevelNotice
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}
