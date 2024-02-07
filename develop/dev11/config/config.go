package config

import (
	"github.com/spf13/viper"
	"time"
)

type ServerConfig struct {
	Port              string
	MaxHeaderBytes    int
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
}

const (
	maxHeaderBytes    = 1 << 20
	readHeaderTimeout = 10 * time.Second
	writeTimeout      = 10 * time.Second
)

type Config struct {
	ServConf *ServerConfig
}

func NewConfig() *Config {
	return &Config{
		ServConf: &ServerConfig{
			Port:              viper.GetString("app.port"),
			MaxHeaderBytes:    maxHeaderBytes,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
		},
	}
}

func InitConfig(path, nameConfig string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(nameConfig)

	return viper.ReadInConfig()
}
