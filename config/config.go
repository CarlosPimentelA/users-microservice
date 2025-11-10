package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB_CONNECTION                string
	DB_NAME                      string
	DB_COLLECTION_USERS          string
	DB_COLLECTION_REFRESH_TOKENS string
	JWT_SECRET_KEY               string
	REFRESH_TOKEN_CONFIG         RefreshTokenConfig
}

type RefreshTokenConfig struct {
	EXPIRY_TIME time.Duration
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DB_CONNECTION:                os.Getenv("DB_CONNECTION"),
		DB_NAME:                      os.Getenv("DB_NAME"),
		DB_COLLECTION_USERS:          os.Getenv("DB_COLLECTION_USERS"),
		JWT_SECRET_KEY:               os.Getenv("JWT_SECRET_KEY"),
		DB_COLLECTION_REFRESH_TOKENS: os.Getenv("DB_COLLECTION_REFRESH_TOKENS"),
		REFRESH_TOKEN_CONFIG: RefreshTokenConfig{
			EXPIRY_TIME: 24 * 7 * time.Hour,
		},
	}
	if config.DB_CONNECTION == "" {
		return nil, fmt.Errorf("variable de entorno DB_CONNECTION no configurada")
	}
	return config, nil
}
