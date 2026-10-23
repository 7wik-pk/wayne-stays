# Wayne Stays API Documentation

Base URL: `http://localhost:8080/api/v1`

---

## Accommodations

### Get All Accommodations
**GET** `/accommodations`

Returns all active accommodations with full image URLs.

**Query Parameters:**
- `category` (optional): Filter by category
  - Valid values: `urban-luxury`, `island`, `historic`, `countryside`

**Example Requests:**
```bash
# Get all accommodations
curl http://localhost:8080/api/v1/accommodations

# Filter by category
curl http://localhost:8080/api/v1/accommodations?category=historic
```

**Response:**
```json
{
  "count": 5,
  "data": [
    {
      "id": "uuid",
      "slug": "wayne-tower-sky-suites",
      "title": "Wayne Tower Sky Suites",
      "description": "Executive suites in the iconic Wayne Tower...",
      "location": "Financial District, Gotham City",
      "priceStart": 1850,
      "capacity": 4,
      "category": "urban-luxury",
      "amenities": [
        "helipad",
        "concierge",
        "skylounge",
        "board-room-access",
        "panoramic-views"
      ],
      "images": {
        "hero": "http://localhost:8080/static/accommodations/wayne_tower_sky_suites/1.jpg",
        "gallery": [
          "http://localhost:8080/static/accommodations/wayne_tower_sky_suites/2.jpg",
          ...
        ]
      }
    }
  ]
}
```

---

### Get Accommodation by ID
**GET** `/accommodations/:id`

Returns a single accommodation by UUID. A malformed or unknown ID returns 404.

**Example:**
```bash
curl http://localhost:8080/api/v1/accommodations/caade5b2-8db2-433e-b256-09444c30ee01
```

**Response:**
```json
{
  "data": {
    "id": "caade5b2-8db2-433e-b256-09444c30ee01",
    "slug": "wayne-tower-sky-suites",
    "title": "Wayne Tower Sky Suites",
    ...
  }
}
```

**Error Response (404):**
```json
{
  "error": "Accommodation not found"
}
```

---

### Get Accommodation by Slug
**GET** `/accommodations/slug/:slug`

Returns a single accommodation by slug (URL-friendly identifier).

**Example:**
```bash
curl http://localhost:8080/api/v1/accommodations/slug/wayne-tower-sky-suites
```

**Response:**
```json
{
  "data": {
    "id": "caade5b2-8db2-433e-b256-09444c30ee01",
    "slug": "wayne-tower-sky-suites",
    "title": "Wayne Tower Sky Suites",
    ...
  }
}
```

---

## Available Accommodations

| Title | Slug | Category | Price/Night | Capacity |
|-------|------|----------|-------------|----------|
| Wayne Tower Sky Suites | `wayne-tower-sky-suites` | urban-luxury | $1,850 | 4 |
| Miagani Island Sanctuary | `miagani-island-sanctuary` | island | $3,200 | 8 |
| The Clocktower Penthouse | `clocktower-penthouse` | historic | $1,400 | 2 |
| Arkham District Gentleman's Club | `arkham-district-gentlemens-club` | historic | $950 | 2 |
| Kane County Vineyard Estate | `kane-county-vineyard-estate` | countryside | $2,100 | 6 |

---

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Unique identifier |
| `slug` | string | URL-friendly identifier |
| `title` | string | Property name |
| `description` | string | Detailed description |
| `location` | string | Geographic location |
| `priceStart` | number | Starting price per night (in dollars) |
| `capacity` | number | Maximum guest capacity |
| `category` | string | Property category |
| `amenities` | string[] | List of amenities |
| `images.hero` | string | Main image URL |
| `images.gallery` | string[] | Gallery image URLs |

---

## Categories

- **urban-luxury**: High-rise penthouses and executive suites in Gotham City
- **island**: Private island retreats and waterfront properties
- **historic**: Heritage buildings and landmarks with historical significance
- **countryside**: Wine estates and rural retreats outside the city

---

## Merchandise

### Get All Merchandise
**GET** `/merchandise`

Returns all shop items with full image URLs.

**Query Parameters:**
- `category` (optional, case-insensitive): Filter by category. `all` is the same as omitting it.
  - Valid values: `apparel`, `accessories`, `collectibles`, `home`, `lifestyle`

**Example:**
```bash
curl http://localhost:8080/api/v1/merchandise?category=apparel
```

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Gotham Skyline Silk Tie",
      "price": 185,
      "category": "apparel",
      "imageUrl": "http://localhost:8080/static/merch/apparel/tie.png",
      "isNew": false
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Unique identifier |
| `name` | string | Product name |
| `description` | string | Optional; omitted when empty |
| `price` | number | Price (in dollars) |
| `category` | string | Product category |
| `imageUrl` | string | Product image URL |
| `isNew` | boolean | Shown as "New Arrival" |

---

## Static Images

Images are served from `/static/accommodations/:property/:image.jpg` and `/static/merch/:category/:image.png`

**Example:**
```
http://localhost:8080/static/accommodations/wayne_tower_sky_suites/1.jpg
```

---

## Health Check

**GET** `/health`

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "wayne-stays-backend"
}
```

---

## CORS

The API has CORS enabled for:
- `http://localhost:5173` (Vite dev server)
- `http://localhost:3000` (Alternative dev port)

---

## Future Endpoints (Planned)

- `POST /bookings` - Create a booking
- `GET /bookings/:id` - Get booking details
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /payments` - Process payment via Stripe
- `GET /users/me` - Get current user profile
