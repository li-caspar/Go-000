package demo

import (
	"app/internal/demo/errorscode"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Name string
}

func NewConfig(name string) Config {
	c := Config{
		Name: name,
	}
	err := c.Init()
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Config) Init() error {
	if err := c.initConfig(); err != nil {
		return err
	}
	c.watchConfig()
	return nil
}

//初始化配置
func (c *Config) initConfig() error {
	if c.Name == "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("configs")
		viper.SetConfigName(AppName)
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(AppName)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(errorscode.ErrConfigFail, fmt.Sprintf("viper readinconfig err:%s", err))
	}
	return nil
}

//监听配置变化
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})
}
