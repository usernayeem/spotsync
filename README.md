# SpotSync API

> Smart Parking & EV Charging Reservation — A centralized platform for airports and malls to manage parking zones and handle high-demand EV charging spot reservations.

**Live URL:** [https://spotsync-production.up.railway.app](https://spotsync-production.up.railway.app)  
**Health Check:** [https://spotsync-production.up.railway.app/ping](https://spotsync-production.up.railway.app/ping)  
**API Zones:** [https://spotsync-production.up.railway.app/api/v1/zones](https://spotsync-production.up.railway.app/api/v1/zones)  

---

## Features

- User registration and login with JWT authentication
- Role-based access control (`driver` / `admin`)
- Parking zone management (create, update, delete) — admin only
- Real-time available spot calculation per zone
- Concurrency-safe reservation system using database transactions and row-level locking
- Reservation cancellation with ownership enforcement

---

## Tech Stack

| Technology | Version |
| --- | --- |
| Go | 1.22+ |
| Echo | v4 |
| GORM | v1 |
| PostgreSQL | (NeonDB) |
| go-playground/validator | v10 |
| golang-jwt/jwt | v5 |
| bcrypt | golang.org/x/crypto |

---

## Architecture

This project follows **Clean Architecture** with strict layer separation:

```
main.go           → Wires all layers via Dependency Injection
├── models/       → GORM structs (database table definitions)
├── dto/          → Request and response data transfer objects
├── repository/   → All database operations (GORM queries, transactions)
├── service/      → Business logic (hashing, JWT generation, capacity rules)
├── handler/      → HTTP layer (bind, validate, call service, return JSON)
├── middlewares/  → JWT authentication and role enforcement
└── config/       → Database connection setup
```

**Dependency flow:** `Handler → Service → Repository → Database`

Handlers never touch the database directly. Raw GORM errors never reach the client.

### Concurrency — The EV Spot Bottleneck

When two drivers attempt to reserve the last available spot simultaneously, a naive implementation would allow both to succeed, exceeding zone capacity.

**Solution:** The `CreateReservationTx` repository method wraps the entire check-and-create flow inside a single GORM transaction with a `FOR UPDATE` row-level lock on the parking zone record. This serializes concurrent requests at the database level, guaranteeing that only one reservation is created when a zone is at full capacity.

---

## Local Setup

### Prerequisites
- Go 1.22+
- PostgreSQL database (local or NeonDB)
- [Air](https://github.com/air-verse/air) for hot-reload (optional)

### Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/usernayeem/spotsync.git
   cd spotsync
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and JWT secret
   ```

4. **Run the server**
   ```bash
   # With hot-reload (recommended for development)
   air

   # Without hot-reload
   go run main.go
   ```

The server starts on the port defined in your `.env` file (default: `8080`).

### Required `.env` Variables

| Variable | Description |
| --- | --- |
| `PORT` | Server port (e.g. `8080`) |
| `DATABASE_URL` | Full PostgreSQL connection string |
| `JWT_SECRET` | Secret key for signing JWT tokens |

---

## API Endpoints

### Authentication

| Method | Endpoint | Access | Description |
| --- | --- | --- | --- |
| POST | `/api/v1/auth/register` | Public | Register a new user |
| POST | `/api/v1/auth/login` | Public | Login and receive JWT token |

### Parking Zones

| Method | Endpoint | Access | Description |
| --- | --- | --- | --- |
| GET | `/api/v1/zones` | Public | Get all zones with available spots |
| GET | `/api/v1/zones/:id` | Public | Get a single zone by ID |
| POST | `/api/v1/zones` | Admin | Create a new parking zone |
| PATCH | `/api/v1/zones/:id` | Admin | Update an existing zone |
| DELETE | `/api/v1/zones/:id` | Admin | Delete a zone |

### Reservations

| Method | Endpoint | Access | Description |
| --- | --- | --- | --- |
| POST | `/api/v1/reservations` | Authenticated | Reserve a parking spot |
| GET | `/api/v1/reservations/my-reservations` | Authenticated | Get your own reservations |
| DELETE | `/api/v1/reservations/:id` | Authenticated | Cancel your reservation |
| GET | `/api/v1/reservations` | Admin | Get all reservations in the system |

---

## HTTP Status Codes

| Code | Usage |
| --- | --- |
| `200` | Successful GET, PATCH, DELETE |
| `201` | Successful POST (resource created) |
| `400` | Validation errors or invalid input |
| `401` | Missing or invalid JWT token |
| `403` | Valid token but insufficient permissions |
| `404` | Resource not found |
| `409` | Business logic conflict (e.g. zone is full) |
| `500` | Unexpected server error |
