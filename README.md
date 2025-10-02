# Go Fiber Boilerplate with File Storage

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmad-fahrudin/golang-fiber-starter)](https://goreportcard.com/report/github.com/ahmad-fahrudin/golang-fiber-starter)

A boilerplate/starter project for quickly building RESTful APIs using Go, Fiber, and PostgreSQL with integrated file storage capabilities.

The app comes with many built-in features, such as authentication using JWT and Google OAuth2, request validation, unit and integration tests, docker support, API documentation, pagination, and flexible file storage (Local/MinIO).

## Quick Start

Clone the repo:

```bash
git clone https://github.com/ahmad-fahrudin/golang-fiber-starter.git
cd golang-fiber-starter
```

Install the dependencies:

```bash
go mod tidy
```

Set the environment variables:

```bash
cp .env.example .env
# open .env and modify the environment variables as needed
```

Run the application:

```bash
make start
```

## Table of Contents

- [Features](#features)
- [Commands](#commands)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [File Storage Feature](#file-storage-feature)
- [API Documentation](#api-documentation)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)
- [Validation](#validation)
- [Authentication](#authentication)
- [Authorization](#authorization)
- [Logging](#logging)

## Features

- **SQL database**: [PostgreSQL](https://www.postgresql.org) Object Relation Mapping using [Gorm](https://gorm.io)
- **Database migrations**: with [golang-migrate](https://github.com/golang-migrate/migrate)
- **File Storage**: Support for Local Storage and MinIO object storage
- **Validation**: request data validation using [Package validator](https://github.com/go-playground/validator)
- **Logging**: using [Logrus](https://github.com/sirupsen/logrus) and [Fiber-Logger](https://docs.gofiber.io/api/middleware/logger)
- **Testing**: unit and integration tests using [Testify](https://github.com/stretchr/testify) and formatted test output using [gotestsum](https://github.com/gotestyourself/gotestsum)
- **Error handling**: centralized error handling mechanism
- **API documentation**: with [Swag](https://github.com/swaggo/swag) and [Swagger](https://github.com/gofiber/swagger)
- **Sending email**: using [Gomail](https://github.com/go-gomail/gomail)
- **Environment variables**: using [Viper](https://github.com/spf13/viper)
- **Security**: set security HTTP headers using [Fiber-Helmet](https://docs.gofiber.io/api/middleware/helmet)
- **CORS**: Cross-Origin Resource-Sharing enabled using [Fiber-CORS](https://docs.gofiber.io/api/middleware/cors)
- **Compression**: gzip compression with [Fiber-Compress](https://docs.gofiber.io/api/middleware/compress)
- **Docker support**
- **Linting**: with [golangci-lint](https://golangci-lint.run)

## Commands

Running locally:

```bash
make start
```

Or running with live reload:

```bash
air
```

> [!NOTE]
> Make sure you have `Air` installed.\
> See 👉 [How to install Air](https://github.com/air-verse/air)

Testing:

```bash
# run all tests
make tests

# run all tests with gotestsum format
make testsum

# run test for the selected function name
make tests-TestUserModel
```

Docker:

```bash
# run docker container
make docker

# run all tests in a docker container
make docker-test
```

Linting:

```bash
# run lint
make lint
```

Swagger:

```bash
# generate the swagger documentation
make swagger
```

Migration:

```bash
# Create migration
make migration-<table-name>

# Example for table users
make migration-users
```

```bash
# run migration up in local
make migrate-up

# run migration down in local
make migrate-down

# run migration up in docker container
make migrate-docker-up

# run migration down all in docker container
make migrate-docker-down
```

Seeder:

```bash
# run all seeders
make seed-all

# list available seeders
make seed-list

# run specific seeder
make seed-User

# truncate table
make seed-truncate-users
```

## Environment Variables

The environment variables can be found and modified in the `.env` file. They come with these default values:

```bash
# server configuration
APP_ENV=dev
APP_HOST=0.0.0.0
APP_PORT=3000
APP_URL=http://localhost:3000

# database configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=toor
DB_NAME=golang_fiberdb
DB_PORT=5432

# JWT
JWT_SECRET=thisisasamplesecret
JWT_ACCESS_EXP_MINUTES=10000
JWT_REFRESH_EXP_DAYS=30
JWT_RESET_PASSWORD_EXP_MINUTES=10
JWT_VERIFY_EMAIL_EXP_MINUTES=10

# SMTP configuration
SMTP_HOST=email-server
SMTP_PORT=587
SMTP_USERNAME=email-server-username
SMTP_PASSWORD=email-server-password
EMAIL_FROM=support@yourapp.com

# OAuth2 configuration
GOOGLE_CLIENT_ID=yourapps.googleusercontent.com
GOOGLE_CLIENT_SECRET=thisisasamplesecret
REDIRECT_URL=http://localhost:3000/v1/auth/google-callback

# File Storage configuration
STORAGE_TYPE=local
STORAGE_LOCAL_PATH=./uploads
STORAGE_MAX_FILE_SIZE=10485760

# MinIO configuration (required when STORAGE_TYPE=minio)
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET_NAME=uploads
MINIO_USE_SSL=false
```

## Project Structure

```
src\
 |--config\         # Environment variables and configuration
 |--controller\     # Route controllers (controller layer)
 |--database\       # Database connection & migrations
 |--docs\           # Swagger files
 |--middleware\     # Custom fiber middlewares
 |--model\          # Database models (data layer)
 |--response\       # Response models
 |--router\         # Routes
 |--service\        # Business logic (service layer)
 |--utils\          # Utility classes and functions
 |--validation\     # Request data validation schemas
 |--main.go         # Fiber app
```

## File Storage Feature

### Overview

Fitur file storage ini menyediakan dua opsi untuk menyimpan file:
1. **Local Storage** - Menyimpan file di file system lokal
2. **MinIO** - Menyimpan file menggunakan MinIO object storage

### Storage Types

#### Local Storage
- Menggunakan `STORAGE_TYPE=local`
- File disimpan di direktori yang ditentukan oleh `STORAGE_LOCAL_PATH`
- File dapat diakses melalui endpoint `/uploads/{folder}/{filename}`

#### MinIO Storage
- Menggunakan `STORAGE_TYPE=minio`
- Memerlukan MinIO server yang berjalan
- File disimpan di MinIO bucket

### File Validation

#### Allowed Extensions
- Images: `.jpg`, `.jpeg`, `.png`, `.gif`
- Documents: `.pdf`, `.doc`, `.docx`, `.txt`

#### File Size Limit
- Maximum file size ditentukan oleh `STORAGE_MAX_FILE_SIZE` (default: 10MB)

### MinIO Setup

#### Installation
```bash
# Install MinIO
go install github.com/minio/minio@latest

# Or download binary from https://min.io/download#/windows
```

#### Running MinIO Server
```bash
# Start MinIO server
minio server ./data --console-address ":9001"

# Default credentials:
# Access Key: minioadmin
# Secret Key: minioadmin
# Console: http://localhost:9001
# API: http://localhost:9000
```

#### MinIO Configuration
1. Akses MinIO Console di `http://localhost:9001`
2. Login dengan credentials default
3. Buat bucket baru dengan nama yang sesuai dengan `MINIO_BUCKET_NAME`
4. Set bucket policy jika diperlukan

## API Documentation

To view the list of available APIs and their specifications, run the server and go to `http://localhost:3000/v1/docs` in your browser.

This documentation page is automatically generated using the [Swag](https://github.com/swaggo/swag) definitions written as comments in the controller files.

## API Endpoints

### Auth routes
`POST /v1/auth/register` - register\
`POST /v1/auth/login` - login\
`POST /v1/auth/logout` - logout\
`POST /v1/auth/refresh-tokens` - refresh auth tokens\
`POST /v1/auth/forgot-password` - send reset password email\
`POST /v1/auth/reset-password` - reset password\
`POST /v1/auth/send-verification-email` - send verification email\
`POST /v1/auth/verify-email` - verify email\
`GET /v1/auth/google` - login with google account

### User routes
`POST /v1/users` - create a user\
`GET /v1/users` - get all users\
`GET /v1/users/:userId` - get user\
`PATCH /v1/users/:userId` - update user\
`DELETE /v1/users/:userId` - delete user

### File routes
`POST /v1/files/upload` - upload file\
`DELETE /v1/files/delete` - delete file\
`GET /v1/files/info` - get file info\
`GET /v1/files/my-files` - get user's files

### File Upload API

#### Upload File
```
POST /v1/files/upload
```

**Headers:**
- `Authorization: Bearer {token}`
- `Content-Type: multipart/form-data`

**Form Data:**
- `file` (required): File yang akan diupload
- `folder` (optional): Nama folder tujuan (default: "general")

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "File uploaded successfully",
  "data": {
    "file_name": "image_20241002120000_abcd1234.jpg",
    "file_path": "general/image_20241002120000_abcd1234.jpg",
    "file_size": 1024000,
    "file_url": "/uploads/general/image_20241002120000_abcd1234.jpg"
  }
}
```

#### Delete File
```
DELETE /v1/files/delete?file_path={path}
```

**Headers:**
- `Authorization: Bearer {token}`

**Query Parameters:**
- `file_path` (required): Path file yang akan dihapus

#### Get File Info
```
GET /v1/files/info?file_path={path}
```

**Headers:**
- `Authorization: Bearer {token}`

**Query Parameters:**
- `file_path` (required): Path file

#### Get My Files
```
GET /v1/files/my-files
```

**Headers:**
- `Authorization: Bearer {token}`

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "Files retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "file_name": "image_20241002120000_abcd1234.jpg",
      "file_path": "general/image_20241002120000_abcd1234.jpg",
      "file_size": 1024000,
      "file_url": "/uploads/general/image_20241002120000_abcd1234.jpg",
      "content_type": "image/jpeg",
      "folder": "general",
      "uploaded_by": "user_uuid",
      "created_at": "2024-10-02T12:00:00Z",
      "updated_at": "2024-10-02T12:00:00Z"
    }
  ]
}
```

## Error Handling

The app includes a custom error handling mechanism, which can be found in the `src/utils/error.go` file.

The error handling process sends an error response in the following format:

```json
{
  "code": 404,
  "status": "error",
  "message": "Not found"
}
```

For example, if you are trying to retrieve a user from the database but the user is not found:

```go
func (s *userService) GetUserByID(c *fiber.Ctx, id string) {
	user := new(model.User)

	err := s.DB.WithContext(c.Context()).First(user, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}
}
```

## Validation

Request data is validated using [Package validator](https://github.com/go-playground/validator). Check the [documentation](https://pkg.go.dev/github.com/go-playground/validator/v10) for more details on how to write validations.

The validation schemas are defined in the `src/validation` directory and are used within the services:

```go
import (
	"app/src/model"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

func (s *userService) CreateUser(c *fiber.Ctx, req validation.CreateUser) (*model.User, error) {
	if err := s.Validate.Struct(&req); err != nil {
		return nil, err
	}
}
```

## Authentication

To require authentication for certain routes, you can use the `Auth` middleware.

```go
import (
	"app/src/controllers"
	m "app/src/middleware"
	"app/src/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, u services.UserService, t services.TokenService) {
  userController := controllers.NewUserController(u, t)
	app.Post("/users", m.Auth(u), userController.CreateUser)
}
```

These routes require a valid JWT access token in the Authorization request header using the Bearer schema.

**Generating Access Tokens**:

An access token can be generated by making a successful call to the register (`POST /v1/auth/register`) or login (`POST /v1/auth/login`) endpoints.

An access token is valid for the time specified in `JWT_ACCESS_EXP_MINUTES` environment variable.

**Refreshing Access Tokens**:

After the access token expires, a new access token can be generated by making a call to the refresh token endpoint (`POST /v1/auth/refresh-tokens`) and sending along a valid refresh token in the request body.

A refresh token is valid for the time specified in `JWT_REFRESH_EXP_DAYS` environment variable.

## Authorization

The `Auth` middleware can also be used to require certain rights/permissions to access a route.

```go
app.Post("/users", m.Auth(u, "manageUsers"), userController.CreateUser)
```

In the example above, an authenticated user can access this route only if that user has the `manageUsers` permission.

The permissions are role-based. You can view the permissions/rights of each role in the `src/config/roles.go` file.

## Logging

Import the logger from `src/utils/logrus.go`. It is using the [Logrus](https://github.com/sirupsen/logrus) logging library.

```go
import "app/src/utils"

utils.Log.Error('message');
utils.Log.Warn('message');
utils.Log.Info('message');
utils.Log.Debug('message');
```

> [!NOTE]
> API request information (request url, response code, timestamp, etc.) are also automatically logged (using [Fiber-Logger](https://docs.gofiber.io/api/middleware/logger)).

## Security Considerations

1. **File Validation**: Selalu validasi file extension dan MIME type
2. **File Size**: Batasi ukuran file yang dapat diupload
3. **Authentication**: Semua endpoint memerlukan authentication
4. **Path Traversal**: File path di-sanitize untuk mencegah path traversal attacks
5. **Storage Permissions**: Pastikan direktori storage memiliki permissions yang tepat

## Troubleshooting

### Common Issues

1. **Permission Denied (Local Storage)**
   - Pastikan direktori `STORAGE_LOCAL_PATH` memiliki write permissions
   - Pastikan user yang menjalankan aplikasi memiliki akses ke direktori

2. **MinIO Connection Failed**
   - Pastikan MinIO server berjalan
   - Periksa `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, dan `MINIO_SECRET_KEY`
   - Pastikan bucket sudah dibuat

3. **File Not Found**
   - Periksa path file di database
   - Pastikan file masih ada di storage

4. **Upload Failed**
   - Periksa ukuran file tidak melebihi limit
   - Pastikan file extension diizinkan
   - Periksa space disk yang tersedia
