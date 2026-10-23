# Wayne Stays Backend

Golang backend API for Wayne Stays luxury accommodation platform.

## Tech Stack

- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM
- **Storage**: Local filesystem (dev) → S3 (production)
- **Payments**: Stripe (planned)

## Documentation

- **[Architecture Overview](ARCHITECTURE.md)** - Current architecture, design decisions, and future evolution path
- **[API Documentation](../API.md)** - Complete API endpoint reference

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16 (or use Docker)

### Installation

**Quick Start:** See [DEVELOPMENT.md](../DEVELOPMENT.md) for complete setup guide.

1. **Start PostgreSQL** (from project root):
   ```bash
   cd ..  # Go to project root
   docker compose up -d postgres
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

4. **Run the server**:
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /health
```

### Accommodations
```
GET    /api/v1/accommodations             # List all (?category=historic)
GET    /api/v1/accommodations/:id         # Get one by UUID
GET    /api/v1/accommodations/slug/:slug  # Get one by slug
```

### Merchandise
```
GET    /api/v1/merchandise                # List all (?category=apparel)
```

Full request/response details: [API.md](../API.md)

### Static Images
```
GET /static/accommodations/:property/:image.jpg
```

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point: config → DB → migrate → seed → serve
├── internal/
│   ├── accommodation/           # handler.go / service.go / repository.go
│   ├── merchandise/             # handler.go / service.go / repository.go
│   ├── server/
│   │   └── server.go            # Router: middleware, static files, routes
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── database.go          # DB connection + migration runner
│   │   ├── migrations/          # Versioned SQL migrations (goose)
│   │   └── seed.go              # Data seeding
│   ├── models/                  # GORM models
│   ├── payment/
│   │   └── gateway.go           # Multi-gateway interfaces (design sketch, unused)
│   ├── storage/
│   │   ├── storage.go           # Storage interface
│   │   ├── local.go             # Local filesystem
│   │   └── factory.go           # Storage factory (S3 not yet implemented)
│   └── testutil/                # Throwaway Postgres for tests
├── .env.example                 # Environment template
├── schema.dbml                  # Full planned schema (dbdiagram.io)
└── go.mod                       # Dependencies
```

## Environment Variables

See `.env.example` for all configuration options.

## Database

The application automatically:
1. Connects to PostgreSQL
2. Applies pending SQL migrations from `internal/database/migrations/` (goose)
3. Seeds any empty tables (5 accommodations, 10 merch items)

See [DEVELOPMENT.md](../DEVELOPMENT.md#database-migrations) for adding a migration.

### Manual Database Operations

```bash
# Connect to PostgreSQL
docker exec -it wayne-stays-postgres psql -U wayne_stays -d wayne_stays_dev

# View accommodations
SELECT title, price_start, location FROM accommodations;

# Reset database (from the project root)
docker compose down -v
docker compose up -d postgres
cd backend && go run cmd/api/main.go
```

## Development

### Adding New Features

1. **Database**: Add a new SQL migration in `internal/database/migrations/`
2. **Models**: Add to `internal/models/`
3. **Feature module**: `internal/<feature>/` with handler, service and repository
4. **Routes**: Register the handler in `internal/server/server.go`
5. **Tests**: Endpoint tests go in `internal/server/server_test.go`

## Testing

```bash
go test ./...          # includes database-backed tests (no Docker needed)
go test -short ./...   # unit tests only
```

### Static Files & Images

When running in Docker:
- Images are served from `./img` (mounted to `/root/img` in container)
- Accessible via `/static/*` endpoint
- CORS is configured to allow frontend access from `http://localhost:3000`

**Note:** If you add new images locally, they show up instantly in the container due to volume mounting.

### Storage Backends

**Local (Development)**:
- Images served from `../img/`
- No cloud costs

**S3 (Production)**:
- Switch `STORAGE_TYPE=s3` in `.env`
- Configure AWS credentials
- Implement an S3 backend in `internal/storage/` (not yet written) and wire it into `factory.go`

## Next Steps

- [x] Build accommodation read endpoints ✅
- [x] Versioned SQL migrations ✅
- [x] Database-backed test suite ✅
- [ ] Add authentication & authorization
- [ ] Implement S3 storage backend
- [ ] Add Stripe payment integration
- [ ] Create booking system
- [ ] Add user management

