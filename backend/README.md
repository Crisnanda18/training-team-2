# SecureTask Backend API

Golang REST API for the SecureTask training project.

## ⚠️ Security Warning
This backend contains **intentional security vulnerabilities** for training purposes. DO NOT use in production!

## Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL running (use `docker-compose up -d` from project root)

### Installation

```bash
# Install dependencies
go mod download

# Run the server
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Tasks
- `GET /api/tasks` - Get user's tasks (requires auth)
- `POST /api/tasks` - Create new task (requires auth)
- `PUT /api/tasks/:id` - Update task (requires auth)
- `DELETE /api/tasks/:id` - Delete task ⚠️ NO AUTH
- `GET /api/tasks/search?q=term` - Search tasks ⚠️ NO AUTH

### Users
- `GET /api/users/me` - Get current user (requires auth)
- `PUT /api/users/:id/profile` - Update profile ⚠️ NO AUTH
- `GET /api/admin/users` - Get all users ⚠️ NO AUTH

## Default Users

- **Admin**: admin@example.com / admin123
- **User**: user@example.com / password123

## Known Vulnerabilities (For Training)

This backend contains the following intentional vulnerabilities:

1. **SQL Injection** - Search endpoint uses string concatenation
2. **Missing Authentication** - Several endpoints lack auth middleware
3. **Missing Authorization** - No role-based access control
4. **Hardcoded Credentials** - Secrets in source code and .env file
5. **Plain Text Passwords** - No password hashing
6. **Sensitive Data Exposure** - Passwords returned in API responses
7. **No Input Validation** - XSS-prone fields not sanitized
8. **No Rate Limiting** - Login endpoint vulnerable to brute force

## Technology Stack

- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (with vulnerabilities)
