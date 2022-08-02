package config

import "os"

var C *Config

type DatabaseConfig struct {
	DatabaseFilePath string
}

type RelayConfig struct {
	RelayIP   string
	RelayPort string
}

type Config struct {
	Database DatabaseConfig
	Relay    RelayConfig
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
		Relay: RelayConfig{
			RelayIP:   getEnvValueString("RELAY_IP", "localhost"),
			RelayPort: getEnvValueString("RELAY_PORT", "32767"),
		},
	}
}
