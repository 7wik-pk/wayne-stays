# Wayne Stays - Development Guide

Complete guide to get Wayne Stays running locally.

---

## 🚀 Quick Start

```bash
# 1. Clone the repository
git clone <your-repo-url>
cd wayne-stays

# 2. Start everything (Database + Backend + Frontend)
docker compose up -d

# 3. Access the app
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

**Access:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Docs: [API.md](API.md)

## Running Frontend and Backend in development environments

### Start the database (& other services to be added in the future)
```bash
docker compose up -d postgres
```

For the next 2 steps, you can either run both backend and frontend on 2 separate terminal windows, or add an ampersand (&) to the end of each command to run them in the background. In the latter case, remember to run cd .. after each command to return to the root directory.

### Start the backend
```bash
cd backend
cp .env.example .env   # first time only
go run cmd/api/main.go
```

### Start the frontend
```bash
cd frontend
npm run dev
```

---

## 📋 Prerequisites

- **Docker** & **Docker Compose** (for PostgreSQL)
- **Go** 1.23+ (for backend)
- **Node.js** 20.19+ (for frontend; required by Vite 7)
- **Git**

---

## 🛠️ Infrastructure Setup

### PostgreSQL Database

```bash
# Start database
docker compose up -d postgres

# Check status
docker compose ps

# View logs
docker compose logs postgres -f

# Stop database
docker compose down

# Reset database (removes all data!)
docker compose down -v
docker compose up -d postgres
```

**Connection Details:**
- Host: `localhost`
- Port: `5432`
- Database: `wayne_stays_dev`
- User: `wayne_stays`
- Password: `gotham_secure_2024`

**Connect via CLI:**
```bash
docker exec -it wayne-stays-postgres psql -U wayne_stays -d wayne_stays_dev
```

---

### Future Services (Currently Commented)

**Redis** (for caching):
```bash
# Uncomment redis service in docker-compose.yml
docker compose up -d redis
```

**MinIO** (for local S3 testing):
```bash
# Uncomment minio service in docker-compose.yml
docker compose up -d minio
# Console: http://localhost:9001
# Credentials: minioadmin / minioadmin
```

---

## 🔧 Backend Setup

### First Time Setup

```bash
cd backend

# Install dependencies
go mod download

# Copy environment config
cp .env.example .env

# (Optional) Edit .env if needed
# vim .env
```

### Running the Backend

```bash
# Development mode (auto-connects to Docker PostgreSQL)
go run cmd/api/main.go

# The backend will:
# 1. Connect to PostgreSQL
# 2. Apply any pending SQL migrations (internal/database/migrations/)
# 3. Seed any empty tables (5 accommodations, 10 merch items)
# 4. Start HTTP server on :8080
```

### Database Migrations

Schema changes are versioned SQL files in `backend/internal/database/migrations/`,
applied in order by [goose](https://github.com/pressly/goose) on startup. Goose
records applied versions in the `goose_db_version` table, so each file runs once.

To add a migration, create the next numbered file:

```bash
cd backend
go run github.com/pressly/goose/v3/cmd/goose@v3.26.0 -s -dir internal/database/migrations create add_bookings sql
```

Then fill in its `-- +goose Up` and `-- +goose Down` sections. Never edit a
migration that has already been applied anywhere; add a new one instead.

### Backend Environment Variables

Key variables in `backend/.env`:

```bash
# Server
PORT=8080
GIN_MODE=debug

# Database (matches docker-compose.yml)
DB_HOST=localhost
DB_PORT=5432
DB_USER=wayne_stays
DB_PASSWORD=gotham_secure_2024
DB_NAME=wayne_stays_dev

# Storage
STORAGE_TYPE=local
LOCAL_STORAGE_PATH=../img
LOCAL_STORAGE_URL=http://localhost:8080/static

# CORS (allow frontend)
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

---

## 🎨 Frontend Setup

### First Time Setup

```bash
cd frontend

# Install dependencies
npm install
```

### Running the Frontend

```bash
# Development mode (with hot reload)
npm run dev

# Access at http://localhost:5173
```

### Frontend Environment Variables

Create `frontend/.env` if needed:

```bash
VITE_API_URL=http://localhost:8080/api/v1
```

---

## 📁 Project Structure

```
wayne-stays/
├── docker-compose.yml          # Infrastructure (PostgreSQL, Redis, MinIO)
├── img/                        # Static images (served by backend)
│   ├── accommodations/
│   │   ├── wayne_tower_sky_suites/
│   │   ├── miagani_island_sanctuary/
│   │   └── ...
│   └── merch/
├── backend/                    # Golang API
│   ├── cmd/api/main.go         # Entry point
│   ├── internal/               # Application code
│   │   ├── accommodation/      # Feature modules (handler/service/repository)
│   │   ├── merchandise/
│   │   ├── server/             # Router: middleware + route registration
│   │   ├── models/             # Data models
│   │   ├── config/             # Configuration
│   │   ├── database/           # Connection, migrations/, seed data
│   │   ├── storage/            # File storage
│   │   └── testutil/           # Throwaway Postgres for tests
│   ├── .env                    # Backend config
│   ├── go.mod                  # Go dependencies
│   └── README.md
├── frontend/                   # React + TypeScript
│   ├── src/
│   │   ├── api/                # Typed API client + response types
│   │   ├── hooks/              # Custom hooks (useApi)
│   │   ├── layouts/            # Page frames (navbar + footer)
│   │   ├── pages/              # One component per route
│   │   ├── components/         # Reusable components
│   │   ├── router.tsx          # Route table
│   │   └── main.tsx            # Entry point
│   ├── package.json            # NPM dependencies
│   └── README.md
├── scripts/                    # Utility scripts
│   └── optimize-images.sh      # Image optimization
└── README.md                   # Project overview
```

---

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./...          # all tests
go test -short ./...   # skip the database-backed tests
```

Database-backed tests (migrations, seeding, HTTP endpoints) run against a real
PostgreSQL 16 started by [embedded-postgres](https://github.com/fergusstrange/embedded-postgres),
so neither Docker nor a local Postgres is needed. The first run downloads the
Postgres binaries into `~/.embedded-postgres-go`; later runs reuse them. Each
test gets its own empty database.

### Frontend Tests

```bash
cd frontend
npm test               # run once
npm run test:watch     # re-run on save
```

Tests use [Vitest](https://vitest.dev) with [Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
in a simulated browser (jsdom). `fetch` is stubbed with `src/test/mockFetch.ts`,
so the backend does not need to be running.

---

## 🗄️ Database Operations

### View Data

```bash
# Connect to database
docker exec -it wayne-stays-postgres psql -U wayne_stays -d wayne_stays_dev

# List tables
\dt

# View accommodations
SELECT title, price_start, location FROM accommodations;

# Exit
\q
```

### Reset Database

```bash
# Stop and remove all data
docker compose down -v

# Start fresh
docker compose up -d postgres

# Restart backend to re-apply migrations and seed
cd backend
go run cmd/api/main.go
```

---

## 🐛 Troubleshooting

### Backend won't start

**"connection refused" error:**
```bash
# Check if PostgreSQL is running
docker compose ps

# Start it if needed
docker compose up -d postgres

# Wait for health check
docker compose logs postgres
```

**"database does not exist":**
```bash
# Reset database
docker compose down -v
docker compose up -d postgres
```

### Frontend API errors

**CORS errors:**
- Check `CORS_ALLOWED_ORIGINS` in `backend/.env`
- Ensure it includes `http://localhost:5173`

**Connection refused:**
- Ensure backend is running on port 8080
- Check `VITE_API_URL` in frontend

### Port conflicts

**Port 5432 (PostgreSQL) already in use:**
```yaml
# Edit docker-compose.yml
ports:
  - "5433:5432"  # Use different host port

# Update backend/.env
DB_PORT=5433
```

**Port 8080 (Backend) already in use:**
```bash
# Edit backend/.env
PORT=8081

# Update frontend API_URL
```

---

## 🚢 Production Deployment

### Building for Production

**Backend:**
```bash
cd backend
go build -o wayne-stays-api cmd/api/main.go
./wayne-stays-api
```

**Frontend:**
```bash
cd frontend
npm run build
# Output in dist/
```

### Environment Setup

- Switch `STORAGE_TYPE=s3` in backend
- Configure AWS credentials
- Update `CORS_ALLOWED_ORIGINS` to production domain
- Set `GIN_MODE=release`

---

## 📚 Documentation

- **[Backend Architecture](backend/ARCHITECTURE.md)** - Architecture decisions and evolution
- **[API Reference](API.md)** - Complete API documentation
- **[Image Reference](DC_CITIES_REFERENCE.md)** - Image sourcing guide

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

MIT License - Feel free to use this project for learning and commercial purposes

---

## 🆘 Getting Help

- Check existing [Issues](https://github.com/your-repo/issues)
- Read the [Architecture Guide](backend/ARCHITECTURE.md)
- Review [API Documentation](API.md)

---

**Happy coding! 🦇**
