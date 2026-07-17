# Restaurant Management Platform

A production-grade, multi-tenant restaurant management backend built with Go. Designed to handle the full lifecycle of restaurant operations — from table reservations and order management to real-time kitchen notifications and payment processing.

## Overview

This platform mirrors the architecture patterns used in enterprise fintech CRM systems, applying them to the restaurant domain. It is built for scale, reliability, and real-time responsiveness.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go (Golang) |
| HTTP Framework | Gin |
| ORM | GORM |
| SQL Query Layer | sqlc |
| Database | PostgreSQL |
| Workflow Engine | Temporal |
| Real-time Events | Pusher Channels |
| Authentication | JWT (golang-jwt/jwt v5) |
| Password Hashing | bcrypt |
| Containerization | Docker |
| Caching | Redis

## Architecture

### Multi-Tenancy
Each restaurant is an isolated tenant. All data is scoped by `tenant_id` at the database level. JWT tokens carry tenant context, enforced by middleware on every protected endpoint — preventing any cross-tenant data access.

### Layered Design
```
HTTP Request
    ↓
Route → Middleware (JWT Auth + Tenant Context)
    ↓
Handler (request parsing, response formatting)
    ↓
Module / Service (business logic)
    ↓
PersistenceDB (sqlc queries + GORM)
    ↓
PostgreSQL
```

### Durable Order Workflows (Temporal)
Order processing is implemented as a durable Temporal workflow, replacing fragile synchronous request chains with a fault-tolerant, long-running process:

```
CreateOrderWorkflow
    ↓
Activity: ValidateAndPriceItems   — validates menu items, snapshots prices
    ↓
Activity: PersistOrder            — atomic insert of order + items
    ↓
Activity: ConfirmOrder            — updates order status
    ↓
Activity: NotifyKitchen           — fires Pusher event to kitchen display
    ↓
Signal: kitchen-started           — kitchen acknowledges via API
    ↓
Activity: MarkPreparing
    ↓
Signal: kitchen-done              — kitchen signals completion via API
    ↓
Activity: MarkReady
    ↓
Activity: NotifyWaiter            — fires Pusher event to floor staff
    ↓
Signal: order-served              — waiter signals delivery via API
    ↓
Activity: MarkServed
```

If the server crashes mid-workflow, Temporal resumes from the last completed activity automatically.

### Real-time Notifications (Pusher)
Events are published to role-scoped channels:

| Event | Channel | Consumers |
|---|---|---|
| `order.created` | `kitchen-{tenantID}` | Kitchen display |
| `order.ready` | `floor-{tenantID}` | Waiter devices |
| `order.paid` | `manager-{tenantID}` | Manager dashboard |
| `table.status_changed` | `floor-{tenantID}` | Floor map |

## Features

### Auth
- Tenant registration
- Staff login with JWT
- Auth middleware with tenant context injection
- Role-based access (owner, manager, waiter, kitchen, cashier)

### Menu Management
- Multi-level hierarchy: Menu → Category → Item
- Item availability toggle
- Price management with soft validation

### Table Management
- Table CRUD with capacity tracking
- Status lifecycle: `available` → `occupied` → `reserved` → `inactive`
- Dedicated status update endpoint
- Filter by status

### Order Management
- Full order lifecycle via Temporal workflow
- Atomic order + items insert (price snapshot at order time)
- Signal-based kitchen coordination
- Predictable workflow ID via UUID reference ID
- Status transitions: `created` → `confirmed` → `preparing` → `ready` → `served` → `paid`

### Payment Processing
- Single-payment model
- Supports: Cash, Card, TeleBirr (mobile money)
- Two-step payment: `pending` → `completed`
- Auto-releases table to `available` on payment
- Itemized bill endpoint with menu item details

### Reservations
- Double-booking prevention (±2 hour window check)
- Status lifecycle: `pending` → `confirmed` → `completed` / `cancelled`
- Auto-updates table status on confirmation and cancellation
- Built with GORM (demonstrating both sqlc and GORM patterns)
- Pagination + status filtering

## API Endpoints

### Public
```
POST /api/v1/auth/login
POST /api/v1/tenant
```

### Protected (Bearer Token required)

**Tenants**
```
GET    /api/v1/tenants
GET    /api/v1/tenants/:tenantID
PUT    /api/v1/tenants/:tenantID
DELETE /api/v1/tenants/:tenantID
PATCH  /api/v1/tenants/:tenantID/restore
```

**Users**
```
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/:userID
PUT    /api/v1/users/:userID
DELETE /api/v1/users/:userID
PATCH  /api/v1/users/:userID/restore
```

**Menus**
```
POST   /api/v1/menus
GET    /api/v1/menus
GET    /api/v1/menus/:menuID
PUT    /api/v1/menus/:menuID
DELETE /api/v1/menus/:menuID
PATCH  /api/v1/menus/:menuID/restore
POST   /api/v1/menus/:menuID/categories
GET    /api/v1/menus/:menuID/categories
PUT    /api/v1/menus/:menuID/categories/:categoryID
DELETE /api/v1/menus/:menuID/categories/:categoryID
PATCH  /api/v1/menus/:menuID/categories/:categoryID/restore
POST   /api/v1/menus/:menuID/categories/:categoryID/items
GET    /api/v1/menus/:menuID/categories/:categoryID/items
PUT    /api/v1/menus/:menuID/categories/:categoryID/items/:ID
DELETE /api/v1/menus/:menuID/categories/:categoryID/items/:ID
PATCH  /api/v1/menus/:menuID/categories/:categoryID/items/:ID/restore
```

**Tables**
```
POST   /api/v1/tables
GET    /api/v1/tables
GET    /api/v1/tables/:tableID
PUT    /api/v1/tables/:tableID
DELETE /api/v1/tables/:tableID
PATCH  /api/v1/tables/:tableID/status
```

**Orders**
```
POST   /api/v1/tables/:tableID/orders
GET    /api/v1/orders
GET    /api/v1/orders/:referenceID/bill
PATCH  /api/v1/orders/:referenceID/kitchen-start
PATCH  /api/v1/orders/:referenceID/kitchen-done
PATCH  /api/v1/orders/:referenceID/served
POST   /api/v1/orders/:referenceID/pay
```

**Reservations**
```
POST   /api/v1/reservations
GET    /api/v1/reservations
GET    /api/v1/reservations/:reservationID
PATCH  /api/v1/reservations/:reservationID/confirm
PATCH  /api/v1/reservations/:reservationID/cancel
PATCH  /api/v1/reservations/:reservationID/complete
```

## Database Schema

```
tenants
users (tenant_id)
menus (tenant_id)
menu_categories (tenant_id, menu_id)
menu_items (tenant_id, category_id)
tables (tenant_id)
orders (tenant_id, table_id, user_id, reference_id)
order_items (tenant_id, order_id, menu_item_id)
payments (tenant_id, order_id)
reservations (tenant_id, table_id)
```

## Local Development

### Prerequisites
- Go 1.21+
- Docker
- golang-migrate CLI
- sqlc CLI

### Setup

```bash
# Clone the repository
git clone https://github.com/Asimro-B/restaurant-platform
cd restaurant-platform

# Start infrastructure
make infra-up

# Run migrations
make migrate-up

# Generate sqlc
make sqlc

# Start the server
go run cmd/api/main.go
```

### Environment Variables

```dotenv
SERVER_PORT=8000

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=restaurant_platform
DB_SSLMODE=disable

JWT_SECRET=your-secret-key

TEMPORAL_HOST=localhost:7233

PUSHER_APP_ID=your_app_id
PUSHER_KEY=your_key
PUSHER_SECRET=your_secret
PUSHER_CLUSTER=your_cluster

REDIS_HOST=localhost
REDIS_PORT=6379
```

### Makefile Commands

```bash
make infra-up          # Start PostgreSQL, Temporal, Redis
make infra-down        # Stop all infrastructure
make migrate-up        # Run pending migrations
make migrate-down      # Rollback last migration
make migration-create NAME=migration_name  # Create new migration
make sqlc              # Regenerate sqlc code
```

## Design Decisions

**Why Temporal for orders?** Orders are long-running processes that span multiple actors (waiter, kitchen, cashier) and can take 20-60 minutes. Temporal provides durable execution — if the server restarts mid-order, the workflow resumes from exactly where it left off. It also gives us a complete audit trail of every state transition.

**Why price snapshots on order_items?** Menu prices change over time. Storing `unit_price` on each order item at the time of ordering ensures historical orders always reflect the price the customer was actually charged — critical for financial accuracy.

**Why sqlc + GORM together?** sqlc is used for complex queries (joins, aggregations, performance-critical paths) where raw SQL gives full control. GORM is used for simpler CRUD and dynamic queries (reservations module) where the Go-first API is more readable. Both patterns are demonstrated intentionally.

**Why UUID reference IDs for orders?** Numeric auto-increment IDs are predictable and expose business data (e.g. order volume). UUID reference IDs are used as the Temporal workflow ID, making workflows addressable without exposing internal database IDs to clients.

**Why Redis caching?** Menu data and bills are read far more frequently than they change. Caching these at the module layer reduces DB load significantly — menu list responses drop from ~8ms to ~0.7ms on cache hit. Cache is invalidated on every 
write operation to keep data consistent.

**Why soft delete?** Orders and payments are financial records — permanently deleting them would break audit trails and historical reporting. Soft delete preserves the data while hiding it from normal queries. Restore endpoints allow recovery from accidental deletions.

**Why per-route RBAC?** Role checks are applied at the route level using explicit allowed-role lists rather than a hierarchy. This prevents privilege escalation bugs (e.g. a waiter processing payments) and makes permissions immediately visible from the route definition.
