package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Runmode      string `yaml:"runmode"`
	Runenv       string `yaml:"runenv"`
	LoadLocation string `yaml:"loadLocation"`
	DB
	GRPC
}

//获取配置结构体
func NewConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "open config file error")
	}
	cfg := &Config{}
	//解析配置  静态配置  暂时不考虑动态配置
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal error")
	}
	return cfg, nil
}

type GRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DB struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	Debug    bool   `yaml:"debug"`
}
