package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kubearmor/koach/koach/config"
	"github.com/kubearmor/koach/koach/model"
)

var DB *gorm.DB

func InitDatabase(databaseConfig config.DatabaseConfig) error {
	db, err := gorm.Open(sqlite.Open(databaseConfig.DatabaseFilePath), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	return nil
}

func MigrateDatabase() error {
	err := DB.AutoMigrate(
		&model.Observability{},
		&model.FileAccess{},
		&model.NetworkCall{},
		&model.ProcessSpawn{},
	)

	return err
}
