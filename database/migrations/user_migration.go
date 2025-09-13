package migration

import (
	"fmt"
	"log"

	"golang-fiber-starter-kit/internal/model"
	"golang-fiber-starter-kit/internal/platform"
)

// MigrateUsers runs GORM AutoMigrate for the users table
func MigrateUsers() error {
	db := platform.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("failed to migrate users: %w", err)
	}

	log.Println("users migration completed")
	return nil
}

// RollbackUsers drops the users table
func RollbackUsers() error {
	db := platform.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	if err := db.Migrator().DropTable(&model.User{}); err != nil {
		return fmt.Errorf("failed to rollback users: %w", err)
	}

	log.Println("users rollback completed")
	return nil
}
