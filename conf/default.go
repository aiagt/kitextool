package ktconf

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Default struct {
	Server       Server   `yaml:"server"`
	Log          Log      `yaml:"log"`
	Registry     Registry `yaml:"registry"`
	ConfigCenter Center   `yaml:"config_center"`
	MySQL        MySQL    `yaml:"mysql"`
	Redis        Redis    `yaml:"redis"`
	Rabbitmq     Rabbitmq `yaml:"rabbitmq"`
}

func (d *Default) ParseDefault(data string) error {
	return yaml.Unmarshal([]byte(data), d)
}

func (d *Default) LoadDefault(files ...string) {
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			klog.Warnf("read config file failed: %s", err.Error())
			continue
		}
		if d == nil {
			*d = Default{}
		}
		err = d.ParseDefault(string(content))
		if err != nil {
			panic(err)
		}
	}
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
	RegistryAddress []string `yaml:"registry_address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
}

type MySQL struct {
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
