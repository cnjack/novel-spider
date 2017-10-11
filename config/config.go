package config

import (
	"github.com/Unknwon/goconfig"
)

var ConfigPath = "./config.ini"

var config *goconfig.ConfigFile

type HttpConfig struct {
	Port string
}

type MysqlConfig struct {
	DSN string
}

var (
	h = &HttpConfig{}
	m = &MysqlConfig{}
)

func GetHttpConfig() *HttpConfig {
	return h
}

func GetMysqlConfig() *MysqlConfig {
	return m
}

func load() {
	var err error
	h.Port, err = config.GetValue("http", "port")
	if err != nil {
		h.Port = ":1314"
	}
	m.DSN, err = config.GetValue("mysql", "dsn")
	if err != nil {
		m.DSN = "root:root@tcp(127.0.0.1:3306)/novel?charset=utf8&parseTime=True&loc=Local"
	}
}

func init() {
	var err error
	config, err = goconfig.LoadConfigFile(ConfigPath)
	if err != nil {
		config = &goconfig.ConfigFile{}
	}
	load()
}
