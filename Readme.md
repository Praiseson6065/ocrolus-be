# Ocrolus Backend API

This is a RESTful API backend service built with Go, Gin, and PostgreSQL. The application provides endpoints for user management, authentication, and article operations.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Development Setup](#local-development-setup)
- [Running with Docker](#running-with-docker)
- [Project Structure](#project-structure)

## Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Docker and Docker Compose (for containerized deployment)
- Git

## Local Development Setup

### 1. Clone the repository

```bash
git clone <https://github.com/Praiseson6065/ocrolus-be.git>
cd ocrolus-be  
```

### 2. Set up environment variables

Create a `.env` file in the root directory with the following variables:

```
ENVIRONMENT=DEV
SERVER_PORT=:8000

# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_SERVER=localhost #for Docker POSTGRES_SERVER=db
POSTGRES_PORT=5432
POSTGRES_DB=ocrolus

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRE=24
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Set up the PostgreSQL database

Ensure PostgreSQL is installed and running on your machine. Create a database with the name specified in your `.env` file:

```bash
createdb ocrolus
```

### 5. Run the application

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8000`.

## Running with Docker

### 1. Set up environment variables

Create a `.env` file as described in the local setup section.

### 2. Build and run with Docker Compose

```bash
docker-compose up -d
```

This will:
- Start a PostgreSQL database container
- Build and run the backend API container
- Set up an Nginx container as a reverse proxy

The API will be accessible at `http://localhost:80`.

### 3. Stop the containers

```bash
docker-compose down
```

To remove volumes as well:

```bash
docker-compose down -v
```


## Project Structure

```
├── cmd/                  # Application entry points
│   ├── main.go           # Main application file
│   ├── router.go         # API routes setup
│   ├── server.go         # Server configuration
├── config/               # Configuration
│   ├── config.go
│   ├── env.go            # Environment variable handling
├── database/             # Database connection and repositories
│   ├── db.go
│   ├── db.article.go
│   ├── db.user.go
├── handlers/             # Request handlers
│   ├── article.go
│   ├── auth.go
│   ├── user.go
├── middleware/           # HTTP middleware
│   ├── cors.go
│   ├── jwt.go
│   ├── middleware.go
├── models/               # Data models
│   ├── article.go
│   ├── model.hooks.go
│   ├── recently-viewed.go
│   ├── user.go
├── nginx/                # Nginx configuration for proxy
│   ├── default.conf
│   ├── Dockerfile
├── util/                 # Utility functions
│   ├── auth.go
├── docker-compose.yaml   # Docker Compose configuration
├── Dockerfile            # Docker image definition
├── go.mod                # Go modules
├── go.sum                # Go modules checksums
```