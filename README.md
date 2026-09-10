# Electronics E-Commerce API

A RESTful e-commerce backend built with **Go, PostgreSQL, Redis, and Docker**.

This project focuses on real-world backend engineering concepts beyond basic CRUD, including authentication, authorization, transactions, caching, inventory management, checkout, payments, and review moderation.

---

## ✨ Features

### Authentication & Users

* Registration and login
* JWT access & refresh tokens
* Logout and token revocation
* Customer / Admin roles
* Role-based authorization

### Products & Catalog

* Products and product variants
* SKU, pricing, and discounts
* Product images
* Brands
* Nested categories
* Product attributes and variant attributes

### Commerce

* Inventory management
* Shopping cart
* Redis cart caching
* Wishlist
* Addresses
* Coupons
* Orders and checkout
* Mock payments
* Product reviews and moderation

### Admin

* Product and inventory management
* Order management
* Coupon management
* Review moderation
* Role-protected admin endpoints

---

## 🧠 Architecture

The application follows a layered architecture:

```text
Client
  ↓
Middleware
  ↓
Handlers
  ↓
Services
  ↓
Repositories
  ↓
PostgreSQL
```

Redis is used alongside PostgreSQL for cart caching.

Each layer has a clear responsibility:

* **Handlers** — HTTP requests, responses, and DTOs
* **Services** — business logic and transactions
* **Repositories** — database queries
* **Models** — application/database entities
* **DTOs** — API request and response contracts
* **Middleware** — authentication and role authorization

Dependency injection is used to connect the application components.

---

## 🛠️ Tech Stack

| Technology            | Purpose               |
| --------------------- | --------------------- |
| **Go**                | Backend               |
| **net/http**          | HTTP server & routing |
| **PostgreSQL**        | Primary database      |
| **Redis**             | Cart caching          |
| **Docker**            | Containerization      |
| **Docker Compose**    | Local infrastructure  |
| **JWT**               | Authentication        |
| **Swagger / OpenAPI** | API documentation     |
| **SQL Migrations**    | Database versioning   |

---

## ⚡ Redis Caching

Redis is currently used for shopping cart caching using a **cache-aside** strategy.

```text
Read Cart
   ↓
Redis
 ┌─┴──────────┐
Hit         Miss
 ↓             ↓
Return      PostgreSQL
              ↓
         Store in Redis
```

Cart mutations invalidate the relevant cache, keeping PostgreSQL as the source of truth.

---

## 💾 Transactions

Transactions are used where multiple database operations need to remain consistent.

Important transactional operations include:

* Inventory updates
* Order checkout
* Inventory reservation
* Payment / order state changes

Repositories use a `DBTX` abstraction so the same repository code can work with both `*sql.DB` and `*sql.Tx`.

---

## 🔐 Authentication

The API uses JWT access tokens with database-backed refresh tokens.

Refresh tokens support:

* Expiration
* Revocation
* Replacement tracking
* Hashed token storage

Admin endpoints require both authentication and the `admin` role.

```http
Authorization: Bearer <access_token>
```

---

## 🗃️ Database

The project uses **23 SQL migrations** to build the PostgreSQL schema.

The main entities include:

```text
Users
Products
Categories
Brands
Product Variants
Inventory
Cart
Wishlist
Coupons
Orders
Payments
Reviews
```

Foreign keys, unique constraints, and check constraints are used to enforce data integrity at the database level.

---

## 📁 Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── docs/
├── internal/
│   ├── app/
│   ├── apperror/
│   ├── cache/
│   ├── config/
│   ├── db/
│   ├── dto/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── redis/
│   ├── repositories/
│   ├── routes/
│   ├── services/
│   └── utils/
├── migrations/
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## 🚀 Getting Started

### Requirements

* Go
* Docker
* Docker Compose
* Git

### Clone

```bash
git clone https://github.com/ariiiiph/ecommerce
cd Ecommerce
```

### Environment

Create your environment file according to `.env.example`.

```bash
cp .env.example .env
```

### Run

```bash
docker compose up --build
```

### Stop

```bash
docker compose down
```

---

## 📖 API Documentation

Swagger / OpenAPI documentation is included in the repository.

Once the application is running:

```text
/swagger/
```

The `docs/` directory contains the generated OpenAPI documentation.

---

## ⭐ Example: Review Moderation

New reviews start as `pending`.

```text
pending ──→ approved
       └──→ rejected
```

Only approved reviews are publicly visible.

Admins can moderate reviews through:

```http
PATCH /api/admin/reviews/{id}/approve
PATCH /api/admin/reviews/{id}/reject
```

---

## 🎯 Project Goal

This project was built to practice designing a larger Go backend and making practical engineering decisions around:

* Layered architecture
* Database design
* Transactions
* Caching
* Authentication & authorization
* Dependency injection
* API design
* Data integrity

The focus was not simply on making endpoints work, but on understanding **why different backend components belong where they do**.

---

## 📌 Status

**Core e-commerce backend complete.**

Possible future improvements include Stripe integration, automated tests, CI/CD, rate limiting, observability, advanced search, and production deployment.

---

## License

This project is for educational and portfolio purposes.
