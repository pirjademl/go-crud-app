package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type EnvConfig struct {
	MYSQL_PASSWORD string
	MYSQL_USERNAME string
	REDIS_PASSWORD string
	REDIS_URL      string
}

func (config *EnvConfig) LoadEnv() *EnvConfig {
	file, err := os.Open(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Split into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove optional quotes around the value
		value = strings.Trim(value, `"'`)

		// Set the environment variable
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading .env file: %v\n", err)
	}

	envConfig := EnvConfig{
		REDIS_URL:      os.Getenv("REDIS_URL"),
		MYSQL_PASSWORD: os.Getenv("MYSQL_PASSWORD"),
		MYSQL_USERNAME: os.Getenv("MYSQL_USERNAME"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
	}

	fmt.Printf("Config: %+v\n", envConfig)
	return &envConfig
}
