package routes

import (
	"golang-fiber-starter-kit/internal/http/handlers"
	"golang-fiber-starter-kit/internal/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) {
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", middleware.JWTMiddleware(), authHandler.Logout)
	auth.Get("/me", middleware.JWTMiddleware(), authHandler.Me)

	// User routes (protected)
	users := api.Group("/users", middleware.JWTMiddleware())
	users.Get("/", userHandler.GetUsers)
	users.Post("/pagination", userHandler.GetUsersPaginated)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}
