package ktdb

import (
	ktconf "github.com/ahaostudy/kitextool/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDial() gorm.Dialector {
	klog.Info(ktconf.GlobalDefaultConf().MySQL.DSN)
	return mysql.Open(ktconf.GlobalDefaultConf().MySQL.DSN)
}
