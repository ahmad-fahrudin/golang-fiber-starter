package main

import (
	"log"

	migration "golang-fiber-starter-kit/database/seeders"
	"golang-fiber-starter-kit/internal/config"
	"golang-fiber-starter-kit/internal/platform"
)

func main() {
	// load config
	if err := config.LoadConfig(); err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// connect to database
	if err := platform.ConnectDatabase(); err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// run seeders
	if err := migration.SeedUsers(); err != nil {
		log.Fatal("seeder failed: ", err)
	}

	log.Println("seeding completed")
}
