package ktdb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDial() GormDial {
	return func(dsn string) gorm.Dialector {
		return mysql.Open(dsn)
	}
}
