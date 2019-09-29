package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type config struct {
	Mysql struct {
		DBUrl string `yaml:"DBUrl"`
	}
}

var Config = config{}

func InitConfig() error {
	fileByte, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		log.WithError(err).Error("load conf meet error")
	}
	err = yaml.Unmarshal(fileByte, &Config)
	if err != nil {
		log.WithError(err).Error("unmarshal meet error")
	}
	return err
}
