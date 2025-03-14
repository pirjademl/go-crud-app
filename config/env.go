package config

import (
	"fmt"
	"os"
)

type EnvConfig struct {
	MYSQL_PASSWORD string
	MYSQL_USERNAME string
	REDIS_PASSWORD string
	REDIS_URL      string
}

func (config *EnvConfig) LoadEnv() *EnvConfig {
	envConfig := EnvConfig{
		REDIS_URL:      os.Getenv("REDIS_URL"),
		MYSQL_PASSWORD: os.Getenv("MYSQL_PASSWORD"),
		MYSQL_USERNAME: os.Getenv("MYSQL_USERNAME"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
	}
	fmt.Println(&envConfig)
	return &envConfig
}
