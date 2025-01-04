package config

import (
	"log"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	DB      DBConfig
	Port    string
	Timeout time.Duration
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	appTimeout, err := time.ParseDuration(utils.GetEnvWithDefault("APP_TIMEOUT", "30s"))

	return &Config{
		DB:      LoadDBConfig(),
		Port:    utils.GetEnvWithDefault("APP_PORT", "8800"),
		Timeout: appTimeout,
	}
}
