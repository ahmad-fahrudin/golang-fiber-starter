package migration

import (
	"fmt"
	"log"

	"golang-fiber-starter-kit/internal/model"
	"golang-fiber-starter-kit/internal/platform"

	"golang.org/x/crypto/bcrypt"
)

// SeedUsers inserts sample users if users table is empty.
func SeedUsers() error {
	db := platform.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		log.Println("users already seeded, skipping")
		return nil
	}

	users := []model.User{
		{Name: "Admin User", Email: "admin@gmail.com", Password: "password123"},
		{Name: "Regular User", Email: "user@gmail.com", Password: "password123"},
	}

	for i := range users {
		hashed, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		users[i].Password = string(hashed)
		if err := db.Create(&users[i]).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", users[i].Email, err)
		}
	}

	log.Println("users seeded successfully")
	return nil
}
