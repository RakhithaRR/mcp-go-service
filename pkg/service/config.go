package service

import (
	"fmt"
	"os"
	"sync"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server Server `mapstructure:"server"`
	Http   Http   `mapstructure:"http"`
}

type Server struct {
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	KeyPath  string `mapstructure:"keyPath"`
	CertPath string `mapstructure:"certPath"`
	Secure   bool   `mapstructure:"secure"`
}

type Http struct {
	Insecure        bool `mapstructure:"insecure"`
	MaxIdleConns    int  `mapstructure:"maxIdleConns"`
	IdleConnTimeout int  `mapstructure:"idleConnTimeout"`
}

var (
	config     *Config
	onceConfig sync.Once
	errConfig  error
)

func InitConfig() (*Config, error) {
	onceConfig.Do(func() {
		data, err := os.ReadFile("config.toml")
		if err != nil {
			logger.Error("Failed to read config file", "error", err)
			errConfig = err
			return
		}
		err = toml.Unmarshal(data, &config)
		if err != nil {
			logger.Error("Failed to unmarshal config file", "error", err)
			errConfig = err
			return
		}
		err = validateConfig()
		if err != nil {
			logger.Error("Invalid config file", "error", err)
			errConfig = err
			return
		}
	})
	return config, errConfig
}

func GetConfig() *Config {
	if config == nil {
		return nil
	}
	return config
}

func validateConfig() error {
	if config.Server.Port == 0 {
		return fmt.Errorf("server port is not set")
	}
	if config.Server.Host == "" {
		return fmt.Errorf("server host is not set")
	}
	if config.Server.KeyPath == "" {
		return fmt.Errorf("server key is not set")
	}
	if config.Server.CertPath == "" {
		return fmt.Errorf("server cert is not set")
	}
	return nil
}
