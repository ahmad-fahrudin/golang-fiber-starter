package platform

import (
	"fmt"
	"log"

	"golang-fiber-starter-kit/internal/config"
	"golang-fiber-starter-kit/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() error {
	var err error

	dsn := config.GetDatabaseURL()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func MigrateDatabase() error {
	err := DB.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
