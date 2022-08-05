package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type Config struct {
	Name string
	Path string
	Type string
	*viper.Viper
}

type Option interface {
	Apply(c *Config)
}

func New(path string, contentType string) *Config {
	c := &Config{
		Name: "config",
		Path: ".",
		Type: "json",
	}
	if path != "" {
		c.Path, c.Name = filepath.Split(path)
	}
	if contentType != "" {
		c.Type = contentType
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

func NewWithOptions(options ...Option) *Config {
	c := &Config{}
	for index := range options {
		options[index].Apply(c)
	}
	return c
}

func (c Config) Get(key string) interface{} {
	return c.Viper.Get(key)
}

func (c Config) GetString(key string) string {
	return c.Viper.GetString(key)
}

func (c Config) GetInt(key string) int {
	return c.Viper.GetInt(key)
}

func (c Config) GetBool(key string) bool {
	return c.Viper.GetBool(key)
}
