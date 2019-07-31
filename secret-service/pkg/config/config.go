package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	varLogLevel = "log.level"
	varHTTPPort = "http.port"

	varPostgresHost              = "postgres.host"
	varPostgresPort              = "postgres.port"
	varPostgresUser              = "postgres.user"
	varPostgresPwd               = "postgres.pwd"
	varPostgresDatabase          = "postgres.database"
	varPostgresSSLMode           = "postgres.sslmode"
	varPostgresConnectionTimeout = "postgres.connection.timeout"
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

	c.v.SetDefault(varPostgresHost, "0.0.0.0")
	c.v.SetDefault(varPostgresPort, 5432)
	c.v.SetDefault(varPostgresUser, "postgres")
	c.v.SetDefault(varPostgresPwd, "abcd1234") // TODO handle pwd
	c.v.SetDefault(varPostgresDatabase, "postgres")
	c.v.SetDefault(varPostgresSSLMode, "disable")
	c.v.SetDefault(varPostgresConnectionTimeout, 5)
}

func (c *Config) GetLogLevel() string {
	return c.v.GetString(varLogLevel)
}

func (c *Config) GetHttpPort() int {
	return c.v.GetInt(varHTTPPort)
}

func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.GetPostgresHost(),
		c.GetPostgresPort(),
		c.GetPostgresUser(),
		c.GetPostgresPwd(),
		c.GetPostgresDatabase(),
		c.GetPostgresSSLMode(),
		c.GetPostgresConnectionTimeout(),
	)
}

func (c *Config) GetPostgresHost() string {
	return c.v.GetString(varPostgresHost)
}

func (c *Config) GetPostgresPort() int {
	return c.v.GetInt(varPostgresPort)
}

func (c *Config) GetPostgresUser() string {
	return c.v.GetString(varPostgresUser)
}

func (c *Config) GetPostgresPwd() string {
	return c.v.GetString(varPostgresPwd)
}

func (c *Config) GetPostgresDatabase() string {
	return c.v.GetString(varPostgresDatabase)
}

func (c *Config) GetPostgresSSLMode() string {
	return c.v.GetString(varPostgresSSLMode)
}

func (c *Config) GetPostgresConnectionTimeout() int64 {
	return c.v.GetInt64(varPostgresConnectionTimeout)
}
