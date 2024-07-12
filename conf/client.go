package ktconf

type ClientConf struct {
	Resolver Resolver `yaml:"resolver"`
}

func (c *ClientConf) ParseClientConf(data string) error {
	return Parse([]byte(data), c)
}

func (c *ClientConf) GetClientConf() *ClientConf {
	return c
}

func (c *ClientConf) ParseServerConf(data string) error {
	return nil
}

func (c *ClientConf) GetServerConf() *ServerConf {
	return nil
}

type Resolver struct {
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
