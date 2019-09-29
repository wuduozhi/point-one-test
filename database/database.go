package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *gorm.DB
var err error

func InitMysql(dbUrl string) error {
	DB, err = gorm.Open("mysql", dbUrl)
	if err != nil {
		log.WithError(err).Error("connect mysql meet error")
	}
	return err

}
