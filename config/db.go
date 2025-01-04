package config

import (
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/utils"
)

type DBConfig struct {
	User     string
	Password string
	Driver   string
	Name     string
	Host     string
	Port     string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		User:     utils.GetEnv("DB_USER"),
		Password: utils.GetEnv("DB_PASSWORD"),
		Name:     utils.GetEnv("DB_NAME"),
		Host:     utils.GetEnvWithDefault("DB_HOST", "localhost"),
		Port:     utils.GetEnvWithDefault("DB_PORT", "3306"),
	}
}
