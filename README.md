# Go Boilerplate

A modern and well-structured Go project boilerplate with modular architecture and observability.

## 🚀 Features

- **Modular Architecture**: Domain-based organization with clear separation of responsibilities
- **Observability**: Prometheus + Grafana for monitoring
- **JWT Authentication**: Authentication system with JWT tokens
- **Custom Error Handling**: Robust error handling system
- **Graceful Shutdown**: Graceful server shutdown
- **Email Sending**: Resend integration for email sending
- **Cache**: Redis-based caching system
- **Migration Control**: Automated database migrations
- **Rate Limiting**: Protection against brute force attacks
- **Metrics**: HTTP and business metrics collection

## 📋 Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make (optional but recommended)

## 🏗️ Project Structure

```
go-boilerplate/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── common/
│   │   └── dto/                 # Data Transfer Objects
│   ├── config/
│   │   └── config.go           # Application configurations
│   ├── domain/                  # Application domains
│   │   ├── category/           # Categories module
│   │   ├── code/               # Verification codes module
│   │   ├── product/            # Products module
│   │   └── user/               # Users module
│   └── infra/                  # Infrastructure
│       ├── container/          # Dependency injection
│       ├── database/           # Database configuration
│       ├── http/               # HTTP handlers and middlewares
│       └── logger/             # Logging configuration
├── pkg/                        # Reusable packages
│   ├── cache/                  # Caching system
│   ├── crypto/                 # Cryptography utilities
│   ├── fault/                  # Error handling
│   ├── httputil/               # HTTP utilities
│   ├── mail/                   # Email system
│   ├── metric/                 # Metrics
│   ├── pagination/             # Pagination
│   ├── retry/                  # Retry mechanism
│   ├── server/                 # Server configuration
│   ├── strutil/                # String utilities
│   ├── token/                  # Token generation and validation
│   └── uid/                    # Unique ID generation
├── grafana/                    # Grafana dashboards
├── docker-compose.yml          # Container configuration
├── Dockerfile                  # Docker image
├── Makefile                    # Automation commands
└── README.md                   # This file
```

## 🚀 How to Run the Project

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/go-boilerplate.git
cd go-boilerplate
```

### 2. Configure environment variables

Copy the `.env.sample` file to `.env` and configure the variables as needed:

```bash
cp .env.sample .env
```

### 3. Run with Docker (Recommended)

```bash
# Build Docker image
make docker-build

# Start containers
make compose-up

# Run migrations
make migrate-up

# Follow logs
make air
```

### 4. Run locally

```bash
# Install dependencies
go mod tidy

# Run migrations
make migrate-up

# Run the application
make run
```

## 🔧 Available Commands

### Docker

```bash
make docker-build      # Build Docker image
make compose-up        # Start containers
make compose-stop      # Stop containers
make compose-down      # Remove containers
```

### Migrations

```bash
make migrate-up        # Apply all migrations
make migrate-down      # Revert all migrations
make migrate-next      # Apply next migration
make migrate-prev      # Revert last migration
make migrate name=name # Create new migration
```

### Tests and Development

```bash
make test             # Run tests with coverage
make tidy             # Organize dependencies
make run              # Run application locally
```

### Mocks

```bash
make install-mockgen  # Install mockgen
make mock             # Generate mocks for all domains
```

### Container Access

```bash
make redis            # Access Redis CLI
make psql             # Access PostgreSQL
```

## 📡 Available Endpoints

### Authentication (`/api/v1/auth`)

| Method | Endpoint                     | Description                  | Authentication |
| ------ | ---------------------------- | ---------------------------- | -------------- |
| POST   | `/api/v1/auth/register`      | Register new user            | ❌             |
| POST   | `/api/v1/auth/login`         | Login (sends code via email) | ❌             |
| POST   | `/api/v1/auth/code/{userId}` | Verify access code           | ❌             |

### Users (`/api/v1/users`)

| Method | Endpoint           | Description          | Authentication |
| ------ | ------------------ | -------------------- | -------------- |
| GET    | `/api/v1/users/me` | Get logged user data | ✅             |

### Products (`/api/v1/products`)

| Method | Endpoint                       | Description                     | Authentication |
| ------ | ------------------------------ | ------------------------------- | -------------- |
| GET    | `/api/v1/products`             | List products (with pagination) | ✅             |
| POST   | `/api/v1/products`             | Create new product              | ✅             |
| GET    | `/api/v1/products/{productId}` | Get product by ID               | ✅             |
| PATCH  | `/api/v1/products/{productId}` | Update product                  | ✅             |
| DELETE | `/api/v1/products/{productId}` | Delete product                  | ✅             |

### Categories (`/api/v1/categories`)

| Method | Endpoint                  | Description                       | Authentication |
| ------ | ------------------------- | --------------------------------- | -------------- |
| GET    | `/api/v1/categories`      | List categories (with pagination) | ✅             |
| POST   | `/api/v1/categories`      | Create new category               | ✅             |
| GET    | `/api/v1/categories/{id}` | Get category by ID                | ✅             |
| DELETE | `/api/v1/categories/{id}` | Delete category                   | ✅             |

### Metrics

| Method | Endpoint   | Description         |
| ------ | ---------- | ------------------- |
| GET    | `/metrics` | Prometheus endpoint |

## 🔐 Authentication

The application uses JWT for authentication. For protected endpoints, include the header:

```
Authorization: Bearer <your-jwt-token>
```

### Authentication Flow

1. **Register**: `POST /api/v1/auth/register`
2. **Login**: `POST /api/v1/auth/login` (sends code via email)
3. **Verification**: `POST /api/v1/auth/code/{userId}` (returns JWT token)

## 📊 Observability

### Prometheus

- Endpoint: `http://localhost:8080/metrics`
- Collected metrics:
  - HTTP requests by method, path and status
  - Cache hits and misses
  - Errors by service and action

### Grafana

- URL: `http://localhost:3000`
- Username: `admin`
- Password: `grafana`

## 🧪 How to Generate Mocks

### 1. Install mockgen

```bash
make install-mockgen
```

### 2. Generate mocks for all domains

```bash
make mock
```

This will generate mocks for:

- Services from all domains
- Repositories from all domains

Mocks will be created in `__mocks/domain/{domain}/`.

### 3. Use mocks in tests

```go
import (
    "github.com/bernardinorafael/go-boilerplate/__mocks/domain/usermock"
)

func TestUserService(t *testing.T) {
    mockRepo := usermock.NewMockRepository(t)
    // ... use the mock
}
```

## 🏗️ How to Create New Modules

To create a new module, follow the structure:

```
internal/domain/
└── new-module/
    ├── handler.go       # Manages HTTP requests
    ├── interface.go     # Defines module contracts/interfaces
    ├── service.go       # Implements business logic
    ├── entity.go        # Defines domain model (DDD based)
    └── repository.go    # Database operations
```

### Implementation example:

1. **interface.go**: Define `Service` and `Repository` interfaces
2. **entity.go**: Define the domain entity
3. **repository.go**: Implement database operations
4. **service.go**: Implement business logic
5. **handler.go**: Register HTTP routes

## 📄 License

This project is under the MIT license. See the `LICENSE` file for more details.
