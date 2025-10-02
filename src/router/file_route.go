package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// FileRoutes setup routes untuk file operations
func FileRoutes(api fiber.Router, db *gorm.DB) {
	// Initialize validator
	validate := validation.Validator()

	// Initialize services
	storageService := service.NewStorageService(db)
	userService := service.NewUserService(db, validate)

	// Initialize controllers
	fileController := controller.NewFileController(storageService)

	// File routes
	files := api.Group("/files")

	// Protected routes (require authentication)
	files.Post("/upload", middleware.Auth(userService), fileController.UploadFile)
	files.Delete("/delete", middleware.Auth(userService), fileController.DeleteFile)
	files.Get("/info", middleware.Auth(userService), fileController.GetFileInfo)
	files.Get("/my-files", middleware.Auth(userService), fileController.GetMyFiles)
}
