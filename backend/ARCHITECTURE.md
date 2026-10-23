# Wayne Stays Backend Architecture

**Current Status:** Layered Architecture (Production-ready for MVP)  
**Last Updated:** 2026-01-28  
**Architecture Decision:** Pragmatic, scalable, maintainable

---

## 📋 Table of Contents

1. [Current Architecture](#current-architecture)
2. [Why This Architecture](#why-this-architecture)
3. [Future Evolution Path](#future-evolution-path)
4. [Migration Strategy](#migration-strategy)
5. [Technology Decisions](#technology-decisions)
6. [References & Resources](#references--resources)

---

## 🏗️ Current Architecture

### **Layered / N-Tier Architecture**

```
┌───────────────────────────────────────────────────┐
│                  HTTP Layer                       │
│  • Gin Router                                     │
│  • CORS Middleware                                │
│  • Request Validation                             │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│              Handler Layer                        │
│  internal/accommodation/handler.go                │
│  • HTTP Request/Response handling                 │
│  • Parameter extraction                           │
│  • Error responses                                │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│              Service Layer                        │
│  internal/accommodation/service.go                │
│  • Business logic                                 │
│  • Data transformation                            │
│  • URL building (storage integration)             │
│  • Price conversion (cents ↔ dollars)             │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│            Repository Layer                       │
│  internal/accommodation/repository.go             │
│  • Database operations (GORM)                     │
│  • Query building                                 │
│  • CRUD operations                                │
└─────────────────┬─────────────────────────────────┘
                  │
┌─────────────────▼─────────────────────────────────┐
│              Data Layer                           │
│  internal/models/accommodation.go                 │
│  • Anemic domain models (data containers)         │
│  • GORM annotations                               │
│  • PostgreSQL mapping                             │
└───────────────────────────────────────────────────┘

          ┌───────────────────────────┐
          │  Infrastructure Layer     │
          ├───────────────────────────┤
          │  • PostgreSQL (GORM)      │
          │  • Local Storage/S3       │
          │  • Config Management      │
          └───────────────────────────┘
```

---

### **Directory Structure**

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── accommodation/                 # Feature module
│   │   ├── handler.go                 # HTTP handlers
│   │   ├── service.go                 # Business logic
│   │   └── repository.go              # Data access
│   │
│   ├── models/                        # Domain models
│   │   └── accommodation.go           # Accommodation entity
│   │
│   ├── config/                        # Configuration
│   │   └── config.go                  # Env var management
│   │
│   ├── database/                      # Database setup
│   │   ├── database.go                # Connection & migration runner
│   │   ├── migrations/                # Versioned SQL (goose)
│   │   └── seed.go                    # Initial data
│   │
│   ├── storage/                       # File storage abstraction
│   │   ├── storage.go                 # Interface
│   │   ├── local.go                   # Local filesystem
│   │   └── factory.go                 # Storage factory (S3: future)
│   │
│   ├── merchandise/                   # Feature module (same layering)
│   ├── server/                        # Router & route registration
│   └── testutil/                      # Throwaway Postgres for tests
│
├── .env                               # Environment config
└── go.mod                             # Dependencies
```

---

### **Current Implementation Details**

#### **1. Handler Layer** (`handler.go`)
**Responsibility:** HTTP request/response handling

```go
type Handler struct {
    service *Service
}

func (h *Handler) GetAll(c *gin.Context) {
    // Extract query params
    category := c.Query("category")
    
    // Call service
    accommodations, err := h.service.GetByCategory(category)
    
    // Return JSON response
    c.JSON(200, gin.H{"data": accommodations})
}
```

**Characteristics:**
- ✅ Thin layer (no business logic)
- ✅ Parameter extraction
- ✅ Error handling & HTTP status codes
- ✅ JSON serialization

---

#### **2. Service Layer** (`service.go`)
**Responsibility:** Business logic, orchestration

```go
type Service struct {
    repo    *Repository
    storage storage.Storage
}

func (s *Service) GetAll() ([]AccommodationResponse, error) {
    // Fetch from repository
    accommodations, err := s.repo.FindAll()
    
    // Business logic: transform to API response
    return s.toResponseList(accommodations), nil
}

func (s *Service) toResponse(acc *Accommodation) AccommodationResponse {
    // Price conversion (cents → dollars)
    priceInDollars := float64(acc.PriceStart) / 100
    
    // Build full image URLs
    images := s.buildImageURLs(acc.Images)
    
    return AccommodationResponse{...}
}
```

**Characteristics:**
- ✅ All business logic lives here
- ✅ Data transformation
- ✅ Service orchestration
- ✅ External service integration (storage)

---

#### **3. Repository Layer** (`repository.go`)
**Responsibility:** Database operations

```go
type Repository struct {
    db *gorm.DB
}

func (r *Repository) FindAll() ([]Accommodation, error) {
    var accommodations []Accommodation
    err := r.db.Where("status = ?", "active").Find(&accommodations).Error
    return accommodations, err
}

func (r *Repository) FindBySlug(slug string) (*Accommodation, error) {
    var accommodation Accommodation
    err := r.db.Where("slug = ?", slug).First(&accommodation).Error
    return &accommodation, err
}
```

**Characteristics:**
- ✅ CRUD operations only
- ✅ Query building
- ✅ No business logic
- ✅ Database abstraction

---

#### **4. Model Layer** (`models/accommodation.go`)
**Responsibility:** Data structures

```go
type Accommodation struct {
    ID          string         `gorm:"type:uuid;primary_key"`
    Slug        string         `gorm:"uniqueIndex"`
    Title       string         `gorm:"not null"`
    PriceStart  int            `gorm:"not null"` // in cents
    Amenities   StringArray    `gorm:"type:jsonb"`
    Images      ImageSet       `gorm:"type:jsonb"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

**Characteristics:**
- ✅ Anemic models (just data)
- ✅ GORM annotations
- ✅ Database mapping
- ❌ No behavior/methods (by design)

---

## 🎯 Why This Architecture?

### **Decision Rationale**

| Factor | Consideration | Decision |
|--------|---------------|----------|
| **Complexity** | Simple CRUD domain | ✅ Layered is sufficient |
| **Team Size** | 1-3 developers | ✅ Easy to understand |
| **Timeline** | MVP in weeks | ✅ Fast to implement |
| **Scalability** | Medium scale (thousands of bookings) | ✅ Adequate performance |
| **Maintainability** | Clear separation of concerns | ✅ Easy to modify |
| **Testability** | Can mock each layer | ✅ Highly testable |

---

### **What We Optimized For**

1. **✅ Speed to Market**
   - Standard pattern developers know
   - Quick to implement features
   - Minimal learning curve

2. **✅ Maintainability**
   - Clear layer boundaries
   - Easy to locate code
   - Predictable structure

3. **✅ Testability**
   - Mock repositories easily
   - Test services in isolation
   - Integration tests straightforward

4. **✅ Scalability (Current Needs)**
   - Handles thousands of requests/second
   - PostgreSQL can scale vertically
   - Horizontal scaling possible with load balancers

---

### **What We're NOT Optimizing For (Yet)**

1. **Complex Business Logic**
   - Current: Simple CRUD
   - Future: Dynamic pricing, complex rules

2. **High Write Throughput**
   - Current: Moderate booking volume
   - Future: Thousands of concurrent bookings

3. **Eventual Consistency**
   - Current: Immediate consistency is fine
   - Future: May need event-driven updates

4. **Microservices**
   - Current: Monolith is simpler
   - Future: May split by bounded context

---

## 🚀 Future Evolution Path

### **When to Evolve the Architecture**

Consider evolving when you encounter:

- ❌ **Performance bottlenecks** in read queries
- ❌ **Complex pricing rules** that clutter services
- ❌ **Inventory conflicts** (double bookings)
- ❌ **Long transaction times** (blocking operations)
- ❌ **Team scaling issues** (conflicts in monolith)

---

### **Phase 1: Current (MVP) ✅**

**Timeline:** Months 0-6  
**Scale:** 0-1000 bookings/month  
**Architecture:** Layered (current)

**Features:**
- ✅ List accommodations
- ✅ View details (API)
- 🔜 Simple booking
- 🔜 Basic payment (Stripe)
- ✅ Static pricing

**Stack:**
- PostgreSQL (single instance)
- Local storage → S3
- Gin HTTP server
- GORM ORM

---

### **Phase 2: Add Complexity (Growth) 🔄**

**Timeline:** Months 6-12  
**Scale:** 1000-10,000 bookings/month  
**Architecture:** Layered + Domain Services

**New Features:**
- 🔄 Seasonal pricing
- 🔄 Basic discounts
- 🔄 Availability calendar
- 🔄 Multiple payment providers

**Architectural Changes:**

```diff
+ Add: Domain Services for complex logic
+ Add: Strategy Pattern for pricing
+ Add: Adapter Pattern for payments
+ Keep: Layered structure
```

**Example: Pricing Domain Service**

```go
// New: internal/domain/pricing/
type PricingService struct {
    baseRateRepo   Repository
    seasonalRules  []SeasonalRule
    discountEngine DiscountEngine
}

func (p *PricingService) CalculatePrice(
    accommodation Accommodation,
    dates DateRange,
    guests int,
) (Money, error) {
    basePrice := accommodation.PriceStart
    
    // Apply seasonal multiplier
    seasonal := p.applySeasonalRates(basePrice, dates)
    
    // Apply discounts
    final := p.discountEngine.Apply(seasonal, dates, guests)
    
    return final, nil
}
```

**Payment Adapter Pattern:**

```go
// Port (interface)
type PaymentProcessor interface {
    Charge(amount Money, method PaymentMethod) (string, error)
    Refund(transactionID string) error
}

// Adapters
type StripeAdapter struct { /* ... */ }
type PayPalAdapter struct { /* ... */ }
type SquareAdapter struct { /* ... */ }

// Service uses interface
type PaymentService struct {
    processor PaymentProcessor // Injected
}
```

---

### **Phase 3: Event-Driven + CQRS (Scale) 🎯**

**Timeline:** Months 12-24  
**Scale:** 10,000-100,000 bookings/month  
**Architecture:** Hybrid DDD + CQRS + Event-Driven

**New Features:**
- 🎯 Dynamic pricing (demand-based)
- 🎯 Real-time inventory
- 🎯 Complex booking workflows
- 🎯 Cancellation/refund sagas
- 🎯 Overbooking prevention

**Architectural Changes:**

```diff
+ Add: Event Store (EventStoreDB or PostgreSQL)
+ Add: Message Bus (NATS, RabbitMQ, Kafka)
+ Add: Read Models (CQRS)
+ Add: Domain Aggregates (DDD)
+ Add: Saga Orchestrator
+ Refactor: Separate read/write databases
```

**New Directory Structure:**

```
backend/
├── internal/
│   ├── domain/                      # DDD Core
│   │   ├── pricing/
│   │   │   ├── aggregate.go         # Pricing aggregate
│   │   │   ├── events.go            # PriceUpdated, etc.
│   │   │   └── rules.go             # Pricing rules
│   │   ├── booking/
│   │   │   ├── aggregate.go         # Booking aggregate
│   │   │   ├── events.go            # BookingCreated, etc.
│   │   │   └── saga.go              # Booking saga
│   │   └── inventory/
│   │       ├── aggregate.go         # Inventory aggregate
│   │       └── events.go            # RoomBlocked, Released
│   │
│   ├── application/                 # Use Cases
│   │   ├── commands/                # Write operations
│   │   │   ├── create_booking.go
│   │   │   └── cancel_booking.go
│   │   └── queries/                 # Read operations
│   │       ├── list_accommodations.go
│   │       └── check_availability.go
│   │
│   ├── infrastructure/
│   │   ├── eventstore/              # Event persistence
│   │   ├── messaging/               # Event bus
│   │   └── readmodels/              # Denormalized views
│   │
│   └── interfaces/                  # HTTP, gRPC, etc.
```

**Example: Booking Saga**

```go
type BookingSaga struct {
    pricingService   PricingService
    inventoryService InventoryService
    paymentService   PaymentService
    eventBus         EventBus
}

func (s *BookingSaga) Execute(cmd CreateBookingCommand) error {
    // Step 1: Calculate price
    price, err := s.pricingService.Calculate(...)
    if err != nil {
        return err
    }
    
    // Step 2: Block inventory
    if err := s.inventoryService.BlockRoom(...); err != nil {
        return err
    }
    
    // Step 3: Process payment
    txID, err := s.paymentService.Charge(price)
    if err != nil {
        // Compensate: Release room
        s.inventoryService.ReleaseRoom(...)
        return err
    }
    
    // Step 4: Confirm booking
    booking := NewBooking(...)
    s.bookingRepo.Save(booking)
    
    // Step 5: Emit event
    s.eventBus.Publish(BookingConfirmedEvent{...})
    
    return nil
}
```

**CQRS: Read vs Write**

```go
// Command Side (Write)
type CreateBookingCommand struct {
    AccommodationID string
    Dates           DateRange
    PaymentMethod   string
}

// Query Side (Read) - Denormalized
type AccommodationListingView struct {
    ID              string
    Title           string
    PricePerNight   float64      // Pre-calculated
    AvailableDates  []Date       // Pre-computed
    Rating          float64
    Images          []string
}

// Update read model via events
func (h *AccommodationViewHandler) On(event PriceUpdatedEvent) {
    h.readDB.UpdatePrice(event.AccommodationID, event.NewPrice)
}
```

---

### **Phase 4: Microservices (Optional) 🌐**

**Timeline:** 24+ months  
**Scale:** 100,000+ bookings/month  
**Architecture:** Event-Driven Microservices

**Split by Bounded Context:**

```
┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
│  Accommodation  │   │    Booking      │   │    Payment      │
│    Service      │   │    Service      │   │    Service      │
│  (Catalog)      │   │  (Reservations) │   │  (Transactions) │
└────────┬────────┘   └────────┬────────┘   └────────┬────────┘
         │                     │                      │
         └─────────────────────┴──────────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │     Event Bus       │
                    │   (NATS/Kafka)      │
                    └─────────────────────┘
```

**Only if:**
- ✅ Team > 10 developers
- ✅ Different scaling needs per service
- ✅ Independent deployment required
- ✅ Clear bounded contexts

---

## 📚 Migration Strategy

### **From Layered to DDD + CQRS**

#### **Step 1: Extract Domain Services**

Move complex logic from service layer to domain services:

```go
// Before (Service Layer)
func (s *Service) CalculatePrice(...) float64 {
    // Complex pricing logic here
}

// After (Domain Service)
// internal/domain/pricing/service.go
type PricingDomainService struct {
    rules []PricingRule
}

func (p *PricingDomainService) Calculate(...) Money {
    // Complex pricing logic here
}

// Service layer becomes thin orchestrator
func (s *Service) CalculatePrice(...) float64 {
    return s.pricingDomain.Calculate(...)
}
```

---

#### **Step 2: Introduce Events**

Add event publishing without changing core logic:

```go
// Add event bus
type EventBus interface {
    Publish(event DomainEvent)
}

// Update service to publish events
func (s *Service) CreateBooking(...) error {
    booking := s.repo.Create(...)
    
    // New: Publish event
    s.eventBus.Publish(BookingCreatedEvent{
        BookingID: booking.ID,
        // ...
    })
    
    return nil
}

// Start with async event handlers
type NotificationHandler struct {}

func (h *NotificationHandler) On(event BookingCreatedEvent) {
    // Send confirmation email
}
```

---

#### **Step 3: Add Read Models (CQRS)**

Create denormalized views for complex queries:

```go
// Write model (existing)
type Accommodation struct {
    ID    string
    Title string
    // ... all fields
}

// New: Read model (denormalized)
type AccommodationListView struct {
    ID            string
    Title         string
    LowestPrice   float64    // Pre-calculated
    AvgRating     float64    // Pre-calculated
    NextAvailable Date       // Pre-computed
}

// Update read model on events
func (h *ListViewUpdater) On(event BookingCreatedEvent) {
    // Update next available date
    h.readDB.UpdateAvailability(event.AccommodationID)
}
```

---

#### **Step 4: Create Aggregates**

Graduate from anemic models to rich aggregates:

```go
// Before (Anemic Model)
type Booking struct {
    ID     string
    Status string
    Total  int
}

// After (Rich Aggregate)
type BookingAggregate struct {
    id      BookingID
    status  BookingStatus
    pricing PricingInfo
    events  []DomainEvent
}

// Business logic in aggregate
func (b *BookingAggregate) Cancel() error {
    if !b.status.IsCancellable() {
        return ErrCannotCancel
    }
    
    b.status = StatusCancelled
    b.recordEvent(BookingCancelledEvent{...})
    
    return nil
}
```

---

#### **Step 5: Implement Event Sourcing (Optional)**

Store events instead of current state:

```go
type EventStore interface {
    Save(aggregateID string, events []DomainEvent) error
    Load(aggregateID string) ([]DomainEvent, error)
}

type InventoryAggregate struct {
    id     string
    events []DomainEvent
}

// Rebuild state from events
func (i *InventoryAggregate) Load(events []DomainEvent) {
    for _, event := range events {
        i.apply(event)
    }
}

func (i *InventoryAggregate) apply(event DomainEvent) {
    switch e := event.(type) {
    case RoomBlockedEvent:
        // Update internal state
    case RoomReleasedEvent:
        // Update internal state
    }
}
```

---

## 🛠️ Technology Decisions

### **Current Stack**

| Component | Technology | Reasoning |
|-----------|-----------|-----------|
| **Language** | Go 1.23 | Performance, simplicity, strong typing |
| **HTTP Framework** | Gin | Fast, popular, good middleware support |
| **ORM** | GORM | Productivity, PostgreSQL support |
| **Migrations** | goose (embedded SQL files) | Versioned, reviewable schema changes; runs on startup |
| **Testing** | go test + embedded-postgres | Real Postgres in tests without Docker |
| **Database** | PostgreSQL 16 | ACID, JSONB, scalability, pgvector (future AI) |
| **Storage** | Local → S3 | Development simplicity, production scalability |
| **Config** | godotenv | 12-factor app compliance |
| **CORS** | gin-contrib/cors | Frontend integration |

---

### **Future Stack (Phase 3+)**

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Event Store** | EventStoreDB or PostgreSQL | Event sourcing |
| **Message Bus** | NATS or RabbitMQ | Event-driven communication |
| **Cache** | Redis | Availability cache, session storage |
| **Search** | PostgreSQL full-text or Elasticsearch | Accommodation search |
| **Orchestration** | Temporal or custom Saga framework | Workflow management |
| **Payments** | Stripe SDK, PayPal REST | Multi-provider support |
| **Monitoring** | Prometheus + Grafana | Metrics, alerting |
| **Tracing** | OpenTelemetry | Distributed tracing |

---

## 📖 References & Resources

### **Architecture Patterns**

1. **Layered Architecture**
   - [Microsoft: N-Tier Architecture](https://docs.microsoft.com/en-us/azure/architecture/guide/architecture-styles/n-tier)
   - Martin Fowler: "Patterns of Enterprise Application Architecture"

2. **Domain-Driven Design**
   - Eric Evans: "Domain-Driven Design" (Blue Book)
   - Vaughn Vernon: "Implementing Domain-Driven Design" (Red Book)
   - [DDD Reference](https://www.domainlanguage.com/ddd/reference/)

3. **CQRS + Event Sourcing**
   - Greg Young: [CQRS Documents](https://cqrs.files.wordpress.com/2010/11/cqrs_documents.pdf)
   - Martin Fowler: [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
   - [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

4. **Saga Pattern**
   - [Microservices.io: Saga Pattern](https://microservices.io/patterns/data/saga.html)
   - Chris Richardson: "Microservices Patterns"

5. **Hexagonal Architecture**
   - Alistair Cockburn: [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

---

### **Go-Specific Resources**

1. **Go Best Practices**
   - [Uber Go Style Guide](https://github.com/uber-go/guide)
   - [Effective Go](https://golang.org/doc/effective_go)
   - [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

2. **Go DDD Examples**
   - [go-ddd](https://github.com/marcusolsson/goddd)
   - [go-clean-architecture](https://github.com/bxcodec/go-clean-arch)

3. **Event-Driven Go**
   - [watermill.io](https://watermill.io/) - Event-driven library
   - [EventStoreDB Go Client](https://github.com/EventStore/EventStore-Client-Go)

---

### **Case Studies (Similar Domains)**

1. **Booking Platforms**
   - [Airbnb Engineering Blog](https://medium.com/airbnb-engineering)
   - [Booking.com Tech Blog](https://blog.booking.com/)
   - [Expedia Group Tech Blog](https://medium.com/expedia-group-tech)

2. **E-commerce Platforms**
   - [Shopify Engineering](https://shopify.engineering/)
   - [Amazon Architecture](https://aws.amazon.com/architecture/)

---

## 🎯 Decision Log

| Date | Decision | Rationale | Status |
|------|----------|-----------|--------|
| 2026-01-28 | Adopt Layered Architecture | MVP speed, team familiarity, adequate for initial scale | ✅ Active |
| 2026-01-28 | Use PostgreSQL over MongoDB | ACID compliance (Stripe), relational needs, pgvector (AI) | ✅ Active |
| 2026-01-28 | Local storage → S3 pattern | Development simplicity, production cloud-native | ✅ Active |
| 2026-01-28 | Defer DDD/CQRS | Current complexity doesn't warrant it | ✅ Deferred |
| 2026-10-23 | Replace GORM AutoMigrate with goose SQL migrations | AutoMigrate can't rename, drop or backfill; bookings/payments need reviewable, ordered schema changes | ✅ Active |
| TBD | Introduce pricing domain service | When seasonal/dynamic pricing needed | 🔜 Planned |
| TBD | Add event bus | When async workflows required | 🔜 Planned |
| TBD | Implement CQRS | When read performance becomes bottleneck | 🔜 Planned |

---

## 🤝 Contributing

When making architectural decisions:

1. **Document here first** before implementing
2. **Discuss trade-offs** with the team
3. **Update decision log** with date and rationale
4. **Keep it pragmatic** - don't over-engineer

**Questions or proposals? Open an issue or PR!**

---

*Last updated: 2026-01-28 by AI Assistant*  
*Next review: When adding complex features (pricing, inventory, workflows)*
