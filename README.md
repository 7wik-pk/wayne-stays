# Wayne Stays

A fictional luxury accommodation booking platform set in the DC Comics universe. Book high-end properties across Gotham City, from Wayne Tower penthouses to private island retreats.

**Tech Stack:** React + TypeScript (Frontend) | Golang + PostgreSQL (Backend)

---

## 🚀 Quick Start

```bash
# 1. Start infrastructure
docker compose up -d postgres

# 2. Start backend (terminal 1)
cd backend
cp .env.example .env   # first time only
go run cmd/api/main.go

# 3. Start frontend (terminal 2)
cd frontend
npm install
npm run dev
```

**Access:**
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **API Docs**: See [API.md](API.md)

---

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.23+
- Node.js 20.19+ (required by Vite 7)

---

## 📚 Documentation

- **[Development Guide](DEVELOPMENT.md)** - Complete setup and troubleshooting
- **[API Reference](API.md)** - Backend API documentation
- **[Backend Architecture](backend/ARCHITECTURE.md)** - Architecture decisions
- **[Payment Architecture](backend/docs/PAYMENT_GATEWAY_STRATEGY.md)** - Multi-gateway payment design
- **[Image Reference](DC_CITIES_REFERENCE.md)** - Real-world location mappings

---

## 🏗️ Project Structure

```
wayne-stays/
├── docker-compose.yml          # Infrastructure (PostgreSQL, Redis, MinIO)
├── frontend/                   # React + TypeScript + Vite
├── backend/                    # Golang + Gin + GORM + PostgreSQL
│   ├── docs/                   # Architecture & design decisions
│   └── internal/payment/       # Multi-gateway payment abstraction
├── img/                        # Static accommodation & merch images (served by backend)
├── scripts/                    # Utility scripts
└── DEVELOPMENT.md              # Setup guide
```

---
## Database Schema (ERD diagram)

draw.io link to the ERD diagram:
https://drive.google.com/file/d/1BJ7NUKANT-lhKWoc-eQ8d7itU7xDwLrn/view?usp=sharing

---

## ✨ Features

### Current (MVP)
- ✅ Browse luxury accommodations
- ✅ View property details with image galleries
- ✅ Realistic pricing ($950 - $3,200/night)
- ✅ Category filtering (urban-luxury, island, historic, countryside)
- ✅ RESTful API with JSONB support
- ✅ Optimized image serving
- ✅ CORS-enabled for SPA

### Planned
- 🔜 User authentication & authorization
- 🔜 Booking system with date selection
- 🔜 Multi-gateway payments (Stripe, Razorpay, Juspay)
- 🔜 Saved payment methods with nicknames
- 🔜 Seasonal & dynamic pricing
- 🔜 Real-time availability calendar
- 🔜 Reviews & ratings
- 🔜 Admin dashboard

---

## 🏨 Available Properties

| Property | Location | Category | Price/Night |
|----------|----------|----------|-------------|
| Wayne Tower Sky Suites | Financial District | Urban Luxury | $1,850 |
| Miagani Island Sanctuary | Gotham Harbor | Private Island | $3,200 |
| The Clocktower Penthouse | Burnley, Old Gotham | Historic | $1,400 |
| Arkham District Gentleman's Club | Arkham | Historic | $950 |
| Kane County Vineyard Estate | Gotham Outskirts | Countryside | $2,100 |

---

## 🛠️ Technology Stack

### Frontend
- React 19 with TypeScript
- React Router 7 (client-side routing)
- Vite (build tool)
- Vitest + Testing Library (tests)
- Vanilla CSS (custom design system)
- Google Fonts (Cinzel, Lato)

### Backend
- Go 1.23
- Gin (HTTP framework)
- GORM (ORM)
- goose (versioned SQL migrations)
- PostgreSQL 16 (JSONB, UUID, full-text search)
- Docker Compose (infrastructure)

### Architecture
- **Current**: Layered (Handler → Service → Repository)
- **Future**: DDD + CQRS + Event-Driven (see [ARCHITECTURE.md](backend/ARCHITECTURE.md))

---

## 🧪 API Examples

```bash
# List all accommodations
curl http://localhost:8080/api/v1/accommodations

# Filter by category
curl http://localhost:8080/api/v1/accommodations?category=historic

# Get by slug
curl http://localhost:8080/api/v1/accommodations/slug/wayne-tower-sky-suites

# Merchandise, optionally by category
curl http://localhost:8080/api/v1/merchandise?category=apparel
```

---

## ✅ Tests

```bash
cd backend && go test ./...    # spins up a throwaway Postgres; no Docker needed
cd frontend && npm test        # Vitest
```

---

## 🤝 Contributing

This is a learning project exploring:
- React vs Vue.js concepts (see [react_vue_map.md](frontend/react_vue_map.md))
- Golang backend development
- RESTful API design
- PostgreSQL with JSONB
- Architecture evolution (Layered → DDD)

Feel free to:
- Open issues for bugs or suggestions
- Submit PRs for improvements
- Use this as a learning resource

---

## 📝 License

MIT License - Free to use for learning and commercial purposes

---

## 🦇 About

Wayne Stays is a fictional subsidiary of Wayne Enterprises (DC Comics) offering luxury accommodations. This project serves as a full-stack learning platform for modern web development practices.

**Happy coding!**