package config

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type Config struct {
	KafkaConfig `ini:"kafka"`
	ESConfig    `ini:"es"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type ESConfig struct {
	Address       string `ini:"address"`
	Index         string `ini:"index"`
	MaxChanSize   int    `ini:"max_chan_size"`
	GoroutineSize int    `ini:"goroutine_size"`
}

func NewConfig() (*Config, error) {
	config := new(Config)
	err := ini.MapTo(config, "./config/config.ini")
	if err != nil {
		logrus.Error("load config file failed.", err.Error())
		return nil, err
	}
	return config, nil
}
