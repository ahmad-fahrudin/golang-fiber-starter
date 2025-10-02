package service

import (
	"app/src/config"
	"app/src/model"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

// StorageService interface untuk file storage
type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string, userID *uuid.UUID) (*FileUploadResult, error)
	DeleteFile(ctx context.Context, filePath string) error
	GetFileURL(filePath string) string
	ValidateFile(file *multipart.FileHeader) error
	GetFileByPath(filePath string) (*model.File, error)
	GetFilesByUser(userID uuid.UUID) ([]model.File, error)
}

// FileUploadResult result dari upload file
type FileUploadResult struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`
	FileURL  string `json:"file_url"`
}

// NewStorageService membuat instance StorageService berdasarkan konfigurasi
func NewStorageService(db *gorm.DB) StorageService {
	switch config.StorageType {
	case "minio":
		return NewMinIOStorageService(db)
	default:
		return NewLocalStorageService(db)
	}
}

// LocalStorageService implementasi storage untuk local file system
type LocalStorageService struct {
	basePath string
	db       *gorm.DB
}

// NewLocalStorageService membuat instance LocalStorageService
func NewLocalStorageService(db *gorm.DB) *LocalStorageService {
	return &LocalStorageService{
		basePath: config.StorageLocalPath,
		db:       db,
	}
}

// UploadFile upload file ke local storage
func (s *LocalStorageService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string, userID *uuid.UUID) (*FileUploadResult, error) {
	if err := s.ValidateFile(file); err != nil {
		return nil, err
	}

	// Generate unique filename
	fileName := s.generateFileName(file.Filename)

	// Create full path
	fullPath := filepath.Join(s.basePath, folder)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Save file
	filePath := filepath.Join(fullPath, fileName)
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Save file info to database
	relativePath := filepath.Join(folder, fileName)
	fileRecord := &model.File{
		FileName:    fileName,
		FilePath:    relativePath,
		FileSize:    file.Size,
		FileURL:     s.GetFileURL(relativePath),
		ContentType: file.Header.Get("Content-Type"),
		Folder:      folder,
		UploadedBy:  userID,
	}

	if err := s.db.Create(fileRecord).Error; err != nil {
		// Clean up uploaded file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// Return result
	return &FileUploadResult{
		FileName: fileName,
		FilePath: relativePath,
		FileSize: file.Size,
		FileURL:  s.GetFileURL(relativePath),
	}, nil
}

// DeleteFile menghapus file dari local storage
func (s *LocalStorageService) DeleteFile(ctx context.Context, filePath string) error {
	// Get file record from database
	fileRecord, err := s.GetFileByPath(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Delete physical file
	fullPath := filepath.Join(s.basePath, filePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete file record from database
	if err := s.db.Delete(fileRecord).Error; err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// GetFileURL mendapatkan URL file untuk local storage
func (s *LocalStorageService) GetFileURL(filePath string) string {
	return fmt.Sprintf("/uploads/%s", strings.ReplaceAll(filePath, "\\", "/"))
}

// ValidateFile validasi file yang diupload
func (s *LocalStorageService) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > config.StorageMaxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes", config.StorageMaxFileSize)
	}

	// Validate file extension
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt"}
	ext := strings.ToLower(filepath.Ext(file.Filename))

	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return nil
		}
	}

	return fmt.Errorf("file extension %s is not allowed", ext)
}

// generateFileName generate nama file yang unik
func (s *LocalStorageService) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Format("20060102150405")
	uuid := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s_%s%s", nameWithoutExt, timestamp, uuid, ext)
}

// GetFileByPath mendapatkan file berdasarkan path
func (s *LocalStorageService) GetFileByPath(filePath string) (*model.File, error) {
	var file model.File
	if err := s.db.Where("file_path = ?", filePath).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFilesByUser mendapatkan files berdasarkan user ID
func (s *LocalStorageService) GetFilesByUser(userID uuid.UUID) ([]model.File, error) {
	var files []model.File
	if err := s.db.Where("uploaded_by = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// MinIOStorageService implementasi storage untuk MinIO
type MinIOStorageService struct {
	client     *minio.Client
	bucketName string
	db         *gorm.DB
}

// NewMinIOStorageService membuat instance MinIOStorageService
func NewMinIOStorageService(db *gorm.DB) *MinIOStorageService {
	// Initialize MinIO client
	client, err := minio.New(config.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIOAccessKey, config.MinIOSecretKey, ""),
		Secure: config.MinIOUseSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize MinIO client: %v", err))
	}

	service := &MinIOStorageService{
		client:     client,
		bucketName: config.MinIOBucketName,
		db:         db,
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.MinIOBucketName)
	if err != nil {
		panic(fmt.Sprintf("Failed to check bucket existence: %v", err))
	}

	if !exists {
		err = client.MakeBucket(ctx, config.MinIOBucketName, minio.MakeBucketOptions{})
		if err != nil {
			panic(fmt.Sprintf("Failed to create bucket: %v", err))
		}
	}

	return service
}

// UploadFile upload file ke MinIO
func (s *MinIOStorageService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string, userID *uuid.UUID) (*FileUploadResult, error) {
	if err := s.ValidateFile(file); err != nil {
		return nil, err
	}

	// Generate unique filename
	fileName := s.generateFileName(file.Filename)
	objectName := fmt.Sprintf("%s/%s", folder, fileName)

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Upload to MinIO
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Save file info to database
	fileRecord := &model.File{
		FileName:    fileName,
		FilePath:    objectName,
		FileSize:    file.Size,
		FileURL:     s.GetFileURL(objectName),
		ContentType: contentType,
		Folder:      folder,
		UploadedBy:  userID,
	}

	if err := s.db.Create(fileRecord).Error; err != nil {
		// Clean up uploaded file if database save fails
		s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	return &FileUploadResult{
		FileName: fileName,
		FilePath: objectName,
		FileSize: file.Size,
		FileURL:  s.GetFileURL(objectName),
	}, nil
}

// DeleteFile menghapus file dari MinIO
func (s *MinIOStorageService) DeleteFile(ctx context.Context, filePath string) error {
	// Get file record from database
	fileRecord, err := s.GetFileByPath(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Delete file from MinIO
	err = s.client.RemoveObject(ctx, s.bucketName, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	// Delete file record from database
	if err := s.db.Delete(fileRecord).Error; err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// GetFileURL mendapatkan URL file untuk MinIO
func (s *MinIOStorageService) GetFileURL(filePath string) string {
	protocol := "http"
	if config.MinIOUseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, config.MinIOEndpoint, s.bucketName, filePath)
}

// ValidateFile validasi file yang diupload untuk MinIO
func (s *MinIOStorageService) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > config.StorageMaxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes", config.StorageMaxFileSize)
	}

	// Validate file extension
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt"}
	ext := strings.ToLower(filepath.Ext(file.Filename))

	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return nil
		}
	}

	return fmt.Errorf("file extension %s is not allowed", ext)
}

// generateFileName generate nama file yang unik untuk MinIO
func (s *MinIOStorageService) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Format("20060102150405")
	uuid := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s_%s%s", nameWithoutExt, timestamp, uuid, ext)
}

// GetFileByPath mendapatkan file berdasarkan path
func (s *MinIOStorageService) GetFileByPath(filePath string) (*model.File, error) {
	var file model.File
	if err := s.db.Where("file_path = ?", filePath).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFilesByUser mendapatkan files berdasarkan user ID
func (s *MinIOStorageService) GetFilesByUser(userID uuid.UUID) ([]model.File, error) {
	var files []model.File
	if err := s.db.Where("uploaded_by = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
