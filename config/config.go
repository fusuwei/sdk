package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Config struct {
	Name string
	Path string
	Type string
	*viper.Viper
}

type Option func(config *Config)

func New(configPath string) *Config {
	c := &Config{
		Name: "config",
		Path: ".",
		Type: "yaml",
	}
	if configPath != "" {
		c.Path, c.Name = filepath.Split(configPath)
		c.Type = strings.Trim(filepath.Ext(configPath), ".")
	}

	c.Viper = viper.GetViper()
	c.SetConfigName(c.Name)
	c.SetConfigType(c.Type)
	c.AddConfigPath(c.Path)
	err := c.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("read config failed: %v", err))
	}
	return c
}

func (c Config) Get(key string) interface{} {
	return c.Viper.Get(key)
}

func (c Config) GetDefault(key string, item interface{}) interface{} {
	if v := c.Viper.Get(key); v != nil {
		return v
	}
	return item
}

func (c Config) GetString(key string) string {
	return c.Viper.GetString(key)
}

func (c Config) GetDefaultString(key string, item string) string {
	if v := c.Viper.GetString(key); v != "" {
		return v
	}
	return item
}

func (c Config) GetInt(key string) int {
	return c.Viper.GetInt(key)
}

func (c Config) GetDefaultInt(key string, item int) int {
	if v := c.Viper.GetInt(key); v != 0 {
		return v
	}
	return item
}

func (c Config) GetBool(key string) bool {
	return c.Viper.GetBool(key)
}
