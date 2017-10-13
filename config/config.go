package config

import (
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml"
)

var config = &Config{}

type HttpConfig struct {
	Port string `yaml:"port"`
}

type MysqlConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	HttpConfig  HttpConfig  `yaml:"http"`
	RedisConfig RedisConfig `yaml:"redis"`
	MysqlConfig MysqlConfig `yaml:"mysql"`
}

func load(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panic(err)
	}
	if err = yaml.Unmarshal(file, config); err != nil {
		log.Panic(err)
	}
}

func GetConfig() Config {
	return *config
}

func init() {
	load("./config.yaml")
}
