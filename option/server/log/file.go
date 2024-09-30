package ktlog

import (
	"fmt"
	ktconf "github.com/aiagt/kitextool/conf"
	"path/filepath"
)

type FilepathOption func(conf *ktconf.ServerConf) string

var filepathOption FilepathOption = func(conf *ktconf.ServerConf) string {
	fileName := fmt.Sprintf("%s.log", conf.Server.Name)
	return filepath.Join("log", fileName)
}

// WithFilepath dynamically set the logger output location through configuration
func WithFilepath(opt FilepathOption) LoggerOption {
	return func(conf *ktconf.ServerConf) {
		filepathOption = opt
	}
}
