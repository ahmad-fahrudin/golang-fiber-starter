# Golang Fiber Starter Kit

A modern, production-ready REST API starter kit built with Go Fiber framework, featuring JWT authentication, GORM ORM, and PostgreSQL database.

## Features

- 🚀 **Go Fiber** - Express-inspired web framework
- 🔐 **JWT Authentication** - Secure token-based authentication
- 🗄️ **GORM** - Go ORM with database migrations and seeding
- 🐘 **PostgreSQL** - Robust relational database
- 🏗️ **Clean Architecture** - Well-organized project structure
- 📝 **CRUD Operations** - Complete user management
- 🔒 **Middleware** - Authentication and authorization
- 🌍 **CORS** - Cross-origin resource sharing
- 📊 **Logging** - Request logging middleware

## Project Structure

```
golang-fiber-starter-kit/
├─ internal/                # Private application code
│  ├─ config/               # Configuration management
│  ├─ http/
│  │  ├─ middleware/        # Custom Fiber middlewares
│  │  ├─ routes/            # Route registration
│  │  └─ handlers/          # HTTP request handlers
│  ├─ repository/           # Database access layer (GORM)
│  ├─ service/              # Business logic layer
│  ├─ model/                # Data models and DTOs
│  └─ platform/             # Database connections, logger, etc.
├─ pkg/                     # Public utility packages
├─ migrations/              # Database migrations and seeders
├─ docs/                    # API documentation
├─ .env.example            # Environment variables template
├─ go.mod                  # Go module file
└─ main.go                 # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd golang-fiber-starter-kit
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Setup environment variables:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` file with your database credentials:
   ```env
   PORT=3000
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=fiber_starter
   DB_SSLMODE=disable
   JWT_SECRET=your_very_strong_jwt_secret
   ```

4. **Create PostgreSQL database:**
   ```sql
   CREATE DATABASE fiber_starter;
   ```

5. **Run database migrations and seeders:**
   ```bash
   go run migrations/seeder.go
   ```

6. **Start the server:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:3000`

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | User login | No |
| POST | `/api/v1/auth/logout` | User logout | Yes |
| GET | `/api/v1/auth/me` | Get current user | Yes |

### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users` | Get all users (paginated) | Yes |
| GET | `/api/v1/users/:id` | Get user by ID | Yes |
| PUT | `/api/v1/users/:id` | Update user | Yes |
| DELETE | `/api/v1/users/:id` | Delete user | Yes |

### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/health` | Server health check | No |

## API Usage Examples

### Register a new user
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get current user (with JWT token)
```bash
curl -X GET http://localhost:3000/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get all users
```bash
curl -X GET "http://localhost:3000/api/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Default Users

After running the seeder, you can use these default accounts:

- **Admin:** admin@example.com / password123
- **User 1:** john@example.com / password123
- **User 2:** jane@example.com / password123

## Development

### Project Architecture

This project follows Clean Architecture principles:

- **Handlers** handle HTTP requests and responses
- **Services** contain business logic
- **Repositories** handle data access
- **Models** define data structures

### Adding New Features

1. Define models in `internal/model/`
2. Create repository interfaces and implementations in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Create HTTP handlers in `internal/http/handlers/`
5. Register routes in `internal/http/routes/`

### Running in Development

```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 3000 |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | |
| DB_NAME | Database name | fiber_starter |
| DB_SSLMODE | SSL mode | disable |
| JWT_SECRET | JWT signing secret | your-secret-key |

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request