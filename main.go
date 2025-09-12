// @title Fiber Starter Kit API
// @version 1.0
// @description Starter kit dengan GORM + JWT + Fiber
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
package main

import (
	"fmt"
	"log"

	"golang-fiber-starter-kit/internal/config"
	"golang-fiber-starter-kit/internal/http/handlers"
	"golang-fiber-starter-kit/internal/http/routes"
	"golang-fiber-starter-kit/internal/platform"
	"golang-fiber-starter-kit/internal/repository"
	"golang-fiber-starter-kit/internal/service"

	_ "golang-fiber-starter-kit/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swagger "github.com/gofiber/swagger"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository()

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	// Secure default: do not use wildcard when AllowCredentials is true.
	// Use environment-configured origin or fallback to http://localhost:3000
	allowOrigin := config.AppConfig.Port
	if allowOrigin == "" {
		allowOrigin = "http://localhost:3000"
	} else {
		// if a port is specified, assume localhost with that port
		allowOrigin = "http://localhost:" + config.AppConfig.Port
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigin,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(app, authHandler, userHandler)

	// Swagger UI (serve at /swagger/index.html)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start server
	port := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Server starting on port %s", config.AppConfig.Port)
	log.Fatal(app.Listen(port))
}
