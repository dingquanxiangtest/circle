package config

import (
	"git.internal.yunify.com/qxp/misc/client"
	"git.internal.yunify.com/qxp/misc/kafka"
	"git.internal.yunify.com/qxp/misc/logger"
	mongo2 "git.internal.yunify.com/qxp/misc/mongo"
	"git.internal.yunify.com/qxp/misc/mysql2"
	"git.internal.yunify.com/qxp/misc/redis2"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf 全局配置文件
var Conf *Config

// DefaultPath 默认配置路径
var DefaultPath = "./configs/config.yml"

// Config 配置文件
type Config struct {
	InternalNet client.Config `yaml:"internalNet"`
	ProcessPort string        `yaml:"processPort"`
	Port        string        `yaml:"port"`
	Model       string        `yaml:"model"`
	Log         logger.Config `yaml:"log"`
	Mysql       mysql2.Config `yaml:"mysql"`
	Mongo       mongo2.Config `yaml:"mongo"`
	Service     Service       `yaml:"service"`
	Redis       redis2.Config `yaml:"redis"`
	Kafka       kafka.Config  `yaml:"kafka"`
}

// Service service config
type Service struct {
	DB string `yaml:"db"`
}

// NewConfig 获取配置配置
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}
