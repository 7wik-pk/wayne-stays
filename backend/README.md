# Wayne Stays Backend

Golang backend API for Wayne Stays luxury accommodation platform.

## Tech Stack

- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM
- **Storage**: Local filesystem (dev) → S3 (production)
- **Payments**: Stripe (planned)

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16 (or use Docker)

### Installation

1. **Start PostgreSQL**:
   ```bash
   docker-compose up -d postgres
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
GET    /api/v1/accommodations      # List all
GET    /api/v1/accommodations/:id  # Get one
POST   /api/v1/accommodations      # Create (admin)
PUT    /api/v1/accommodations/:id  # Update (admin)
DELETE /api/v1/accommodations/:id  # Delete (admin)
```

### Static Images
```
GET /static/accommodations/:property/:image.jpg
```

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── database.go          # DB connection
│   │   └── seed.go              # Data seeding
│   ├── models/
│   │   └── accommodation.go     # Data models
│   └── storage/
│       ├── storage.go           # Storage interface
│       ├── local.go             # Local filesystem
│       ├── s3.go                # AWS S3 (TODO)
│       └── factory.go           # Storage factory
├── .env.example                 # Environment template
├── docker-compose.yml           # PostgreSQL + MinIO
└── go.mod                       # Dependencies
```

## Environment Variables

See `.env.example` for all configuration options.

## Database

The application automatically:
1. Connects to PostgreSQL
2. Runs migrations (creates tables)
3. Seeds initial data (5 accommodations)

### Manual Database Operations

```bash
# Connect to PostgreSQL
docker exec -it wayne-stays-postgres psql -U wayne_stays -d wayne_stays_dev

# View accommodations
SELECT title, price_start, location FROM accommodations;

# Reset database
docker-compose down -v
docker-compose up -d postgres
go run cmd/api/main.go
```

## Development

### Adding New Features

1. **Models**: Add to `internal/models/`
2. **Routes**: Add to `cmd/api/main.go` or create handlers
3. **Database**: Update migration in `database.Migrate()`

### Storage Backends

**Local (Development)**:
- Images served from `../img/`
- No cloud costs

**S3 (Production)**:
- Switch `STORAGE_TYPE=s3` in `.env`
- Configure AWS credentials
- Implement `internal/storage/s3.go`

## Next Steps

- [ ] Build accommodation CRUD handlers
- [ ] Add authentication & authorization
- [ ] Implement S3 storage backend
- [ ] Add Stripe payment integration
- [ ] Create booking system
- [ ] Add user management
