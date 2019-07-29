package config

import (
	"github.com/spf13/viper"
)

const (
	varLogLevel = "log.level"
	varHTTPPort = "http.port"
)

type Config struct {
	v *viper.Viper
}

func New() *Config {
	c := &Config{
		v: viper.New(),
	}
	c.setDefaults()
	return c
}

func (c *Config) setDefaults() {
	c.v.SetDefault(varLogLevel, "debug")
	c.v.SetDefault(varHTTPPort, 8080)
}

func (c *Config) GetLogLevel() string {
	return c.v.GetString(varLogLevel)
}

func (c *Config) GetHttpPort() int {
	return c.v.GetInt(varHTTPPort)
}
