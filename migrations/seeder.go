package main

import (
	"log"

	"golang-fiber-starter-kit/internal/config"
	"golang-fiber-starter-kit/internal/model"
	"golang-fiber-starter-kit/internal/platform"
	"golang-fiber-starter-kit/pkg"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := platform.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := platform.MigrateDatabase(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed users
	if err := seedUsers(); err != nil {
		log.Fatal("Failed to seed users:", err)
	}

	log.Println("Database seeded successfully!")
}

func seedUsers() error {
	db := platform.GetDB()

	// Check if users already exist
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already exist, skipping seeding")
		return nil
	}

	// Sample users to seed
	users := []model.User{
		{
			Name:  "Admin User",
			Email: "admin@example.com",
		},
		{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
	}

	// Hash passwords and create users
	for i := range users {
		hashedPassword, err := pkg.HashPassword("password123")
		if err != nil {
			return err
		}
		users[i].Password = hashedPassword

		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Created %d users", len(users))
	return nil
}
