package config

import (
	"github.com/Unknwon/goconfig"
)

var ConfigPath = "./config.ini"

var config *goconfig.ConfigFile

type HttpConfig struct {
	Port       string
	AccessKey  string
	SecretKey  string
	BucketName string
	BucketUrl  string
}

type MysqlConfig struct {
	DSN string
}

type SpiderConfig struct {
	BrokerAddr  []string
	BrokerTopic string
	MaxProcess  int
}

var (
	h = &HttpConfig{}
	s = &SpiderConfig{}
	m = &MysqlConfig{}
)

func GetHttpConfig() *HttpConfig {
	return h
}

func GetSpiderConfig() *SpiderConfig {
	return s
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
	h.AccessKey, _ = config.GetValue("http", "accessKey")
	h.SecretKey, _ = config.GetValue("http", "secretKey")
	h.BucketName, _ = config.GetValue("http", "bucketName")
	h.BucketUrl, _ = config.GetValue("http", "bucketHost")
	s.BrokerAddr = config.MustValueArray("spider", "broker_addrs", ",")
	s.MaxProcess, err = config.Int("spider", "max_process")
	if err != nil {
		s.MaxProcess = 50
	}
	s.BrokerTopic, err = config.GetValue("spider", "broker_topic")
	if err != nil {
		s.BrokerTopic = "novel-task"
	}
	m.DSN, err = config.GetValue("mysql", "dsn")
	if err != nil {
		m.DSN = "root:root@tcp(127.0.0.1:3306)/novel?charset=utf8&parseTime=True&loc=Local"
	}
}

func ReloadConfig() {
	var err error
	config, err = goconfig.LoadConfigFile(ConfigPath)
	if err != nil {
		config = &goconfig.ConfigFile{}
	}
	load()
}

func init() {
	var err error
	config, err = goconfig.LoadConfigFile(ConfigPath)
	if err != nil {
		config = &goconfig.ConfigFile{}
	}
	load()
}
