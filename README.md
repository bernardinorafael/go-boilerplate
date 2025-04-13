# Go Boilerplate

A modern and well-structured Go project boilerplate.

## Features

- Modular architecture
- Observability (Prometheus + Grafana)
- JWT Authentication (Paseto coming soon) + Session management
- Custom error handling
- Graceful shutdown
- Email sending
- Cache
- Easy migration control

## How to create new modules?

```
modules/
└── module-name/         # Module name (auth, session, user, etc)
    ├── handler.go       # Manages HTTP requests
    ├── interface.go     # Defines the module's contracts/interfaces
    ├── service.go       # Implements business logic
    ├── entity.go        # Defines domain model (DDD based)
    └── repository.go    # Database operations
```

## Getting Started

1. Clone the repository

```bash
git clone https://github.com/yourusername/go-boilerplate.git
cd go-boilerplate
```

2. Build docker image

```bash
make docker-build
```

3. Up container

```bash
docker compose up -d
```

4. Up migrations

```bash
make migrate-up
```

5. Follow air logs

```bash
make air
```
