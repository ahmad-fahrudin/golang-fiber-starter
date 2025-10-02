package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// FileController struct
type FileController struct {
	storageService service.StorageService
}

// NewFileController membuat instance FileController
func NewFileController(storageService service.StorageService) *FileController {
	return &FileController{
		storageService: storageService,
	}
}

// FileUploadRequest request untuk upload file
type FileUploadRequest struct {
	Folder string `form:"folder" json:"folder"`
}

// UploadFile godoc
// @Summary Upload file
// @Description Upload file ke storage (local atau MinIO)
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param folder formData string false "Folder destination"
// @Router /files/upload [post]
func (fc *FileController) UploadFile(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	userClaims := c.Locals("user")
	var userID *uuid.UUID
	if userClaims != nil {
		if claims, ok := userClaims.(*model.User); ok {
			userID = &claims.ID
		}
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	// Get folder from form data
	folder := c.FormValue("folder", "general")

	// Upload file
	result, err := fc.storageService.UploadFile(c.Context(), file, folder, userID)
	if err != nil {
		utils.Log.Errorf("Failed to upload file: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload file")
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "File uploaded successfully",
		Data:    result,
	})
}

// DeleteFile godoc
// @Summary Delete file
// @Description Delete file dari storage (local atau MinIO)
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file_path query string true "File path to delete"
// @Router /files/delete [delete]
func (fc *FileController) DeleteFile(c *fiber.Ctx) error {
	filePath := c.Query("file_path")
	if filePath == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File path is required")
	}

	// Delete file
	err := fc.storageService.DeleteFile(c.Context(), filePath)
	if err != nil {
		utils.Log.Errorf("Failed to delete file: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete file")
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "File deleted successfully",
	})
}

// GetFileInfo godoc
// @Summary Get file info
// @Description Get file information and URL
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file_path query string true "File path"
// @Router /files/info [get]
func (fc *FileController) GetFileInfo(c *fiber.Ctx) error {
	filePath := c.Query("file_path")
	if filePath == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File path is required")
	}

	fileURL := fc.storageService.GetFileURL(filePath)

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "File info retrieved successfully",
		Data: map[string]string{
			"file_path": filePath,
			"file_url":  fileURL,
		},
	})
}

// GetMyFiles godoc
// @Summary Get user files
// @Description Get files uploaded by current user
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /files/my-files [get]
func (fc *FileController) GetMyFiles(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	userClaims := c.Locals("user")
	if userClaims == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	user, ok := userClaims.(*model.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user data")
	}

	// Get user files
	files, err := fc.storageService.GetFilesByUser(user.ID)
	if err != nil {
		utils.Log.Errorf("Failed to get user files: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get files")
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Files retrieved successfully",
		Data:    files,
	})
}
