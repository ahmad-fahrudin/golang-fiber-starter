package main

import (
	"log"

	migration "golang-fiber-starter-kit/database/migrations"
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

	// run migrations
	if err := migration.MigrateUsers(); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("migrations completed")
}
