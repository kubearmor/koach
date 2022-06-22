package config

import "os"

var C *Config

type DatabaseConfig struct {
	DatabaseFilePath string
}

type Config struct {
	Database DatabaseConfig
}

func getEnvValueString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

func InitConfig() {
	C = &Config{
		Database: DatabaseConfig{
			DatabaseFilePath: getEnvValueString("DATABASE_FILE_PATH", "./database/koach.db"),
		},
	}
}
