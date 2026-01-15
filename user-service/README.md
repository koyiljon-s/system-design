# User Service

A Go-based microservice for user management using Gin framework, PostgreSQL, and JWT authentication.

## Prerequisites

- Go 1.25.4 or later
- Docker and Docker Compose

## Setup

1. Clone the repository and navigate to the user-service directory:
   ```bash
   cd user-service
   ```

2. Copy the environment file and update the values:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your database credentials and JWT secret.

## Running the Application

### Using Docker Compose (Recommended)

This will start both the PostgreSQL database and the Go application:

```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`.


## API Endpoints

- `POST /api/register` - Register a new user
- `POST /api/login` - Login and get JWT token
- `GET /api/me` - Get current user info (requires auth)
- `PUT /api/me` - Update current user (requires auth)
- `DELETE /api/me` - Delete current user (requires auth)
- `GET /api/users/:id` - Get user by ID (requires auth, internal use)

## Environment Variables

- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - JWT signing secret
- `PORT` - Application port (default: 8080)