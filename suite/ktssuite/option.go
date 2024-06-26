package ktssuite

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
)

type Option interface {
	Apply(conf *ktconf.Default)
	OnChange(conf *ktconf.Default)
}
