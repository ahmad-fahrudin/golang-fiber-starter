package main

import (
	"app/src/config"
	"app/src/database"
	"app/src/middleware"
	"app/src/router"
	"app/src/utils"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"gorm.io/gorm"
)

// @title go-fiber-boilerplate API documentation
// @version 1.0.0
// @license.name MIT
// @license.url https://github.com/indrayyana/go-fiber-boilerplate/blob/main/LICENSE
// @host localhost:3000
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Example Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		runSeeder()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := setupFiberApp()
	db := setupDatabase()
	defer closeDatabase(db)
	setupRoutes(app, db)

	address := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)

	// Start server and handle graceful shutdown
	serverErrors := make(chan error, 1)
	go startServer(app, address, serverErrors)
	handleGracefulShutdown(ctx, app, serverErrors)
}

func setupFiberApp() *fiber.App {
	app := fiber.New(config.FiberConfig())

	// Middleware setup
	app.Use("/v1/auth", middleware.LimiterConfig())
	app.Use(middleware.LoggerConfig())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.RecoverConfig())

	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	// Add any additional database setup if needed
	return db
}

func setupRoutes(app *fiber.App, db *gorm.DB) {
	router.Routes(app, db)
	app.Use(utils.NotFoundHandler)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	if err := app.Listen(address); err != nil {
		errs <- fmt.Errorf("error starting server: %w", err)
	}
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("Error closing database connection: %v", err)
	} else {
		utils.Log.Info("Database connection closed successfully")
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		utils.Log.Fatalf("Server error: %v", err)
	case <-quit:
		utils.Log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			utils.Log.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		utils.Log.Info("Server exiting due to context cancellation")
	}

	utils.Log.Info("Server exited")
}

func runSeeder() {
	command := "all"
	if len(os.Args) > 2 {
		command = os.Args[2]
	}

	sc := database.NewSeederConfig()

	switch command {
	case "all":
		if err := sc.RunAllSeeders(); err != nil {
			fmt.Printf("Error running all seeders: %v\n", err)
			os.Exit(1)
		}
	case "list":
		sc.ListAvailableSeeders()
	case "run":
		if len(os.Args) < 4 {
			fmt.Println("Error: seeder name required for 'run' command")
			printSeederUsage()
			os.Exit(1)
		}
		seederName := strings.Join(os.Args[3:], " ")
		if err := sc.RunSpecificSeeder(seederName); err != nil {
			fmt.Printf("Error running seeder %s: %v\n", seederName, err)
			os.Exit(1)
		}
	case "refresh":
		if len(os.Args) < 5 {
			fmt.Println("Error: seeder name and table name required for 'refresh' command")
			printSeederUsage()
			os.Exit(1)
		}
		seederName := os.Args[3]
		tableName := os.Args[4]
		if err := sc.RefreshSeeder(seederName, tableName); err != nil {
			fmt.Printf("Error refreshing seeder %s: %v\n", seederName, err)
			os.Exit(1)
		}
	case "truncate":
		if len(os.Args) < 4 {
			fmt.Println("Error: table name required for 'truncate' command")
			printSeederUsage()
			os.Exit(1)
		}
		tableName := os.Args[3]
		if err := sc.TruncateTable(tableName); err != nil {
			fmt.Printf("Error truncating table %s: %v\n", tableName, err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown seeder command: %s\n", command)
		printSeederUsage()
		os.Exit(1)
	}
}

func printSeederUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run src/main.go --seed all                    - Run all seeders")
	fmt.Println("  go run src/main.go --seed list                   - List all available seeders")
	fmt.Println("  go run src/main.go --seed run <seeder_name>      - Run a specific seeder")
	fmt.Println("  go run src/main.go --seed refresh <seeder_name> <table_name> - Truncate table and run seeder")
	fmt.Println("  go run src/main.go --seed truncate <table_name>  - Truncate a table")
}
