# go-ecommerce-api

A REST API built in Go providing **products** and **orders** services with a clean layered architecture.

---

## Table of Contents

| S.No | Title                                    | Link                                                                |
| ---- | ---------------------------------------- | ------------------------------------------------------------------- |
| 1    | Services Overview                        | [→ Services Overview](#services-overview)                           |
| 2    | Project Setup                            | [→ Project Setup](#project-setup)                                   |
| 3    | Adding a New Table                       | [→ Adding a New Table](#adding-a-new-table)                         |
| 4    | Clean Layered Architecture               | [→ Clean Layered Architecture](#clean-layered-architecture)         |
| 5    | HTTP Server Setup                        | [→ HTTP Server Setup](#http-server-setup)                           |
| 6    | Structured Logging                       | [→ Structured Logging](#structured-logging)                         |
| 7    | The `internal/` Folder                   | [→ The internal/ Folder](#the-internal-folder)                      |
| 8    | Pointers — When to Use `*` and `&`       | [→ Pointers](#pointers--when-to-use--and-)                          |
| 9    | Handler & Service Layers for `/products` | [→ Handler & Service Layers](#handler--service-layers-for-products) |
| 10   | Products Service Deep Dive               | [→ Products Service Deep Dive](#products-service-deep-dive)         |
| 11   | Database with SQLC                       | [→ Database with SQLC](#database-with-sqlc)                         |
| 12   | Migrations with Goose                    | [→ Migrations with Goose](#migrations-with-goose)                   |
| 13   | Running Postgres Locally                 | [→ Running Postgres Locally](#running-postgres-locally)             |
| 14   | Connecting to Postgres with Go           | [→ Connecting to Postgres with Go](#connecting-to-postgres-with-go) |
| 15   | GET `/products/:id`                      | [→ GET /products/:id](#get-productsid)                              |
| 16   | POST `/orders` — Order Creation          | [→ POST /orders](#post-orders--order-creation)                      |
| 17   | GET `/orders` & GET `/orders/:id`        | [→ GET /orders](#get-orders--get-ordersid)                          |
| 18   | POST `/products`                         | [→ POST /products](#post-products)                                  |
| 19   | PUT `/products`                          | [→ PUT /products](#put-products)                                    |
| 20   | SQLC Query Parameters (`$1`, `$2`)       | [→ SQLC Query Parameters](#sqlc-query-parameters-1-2)               |
| 21   | Dependency Injection                     | [→ Dependency Injection](#dependency-injection)                     |

---

## Services Overview

This API provides 2 services:

- **[products](https://github.com/Prakash-Ravichandran/go-ecommerce-api/tree/main/internal/orders)** — `GET /products`, `POST /products`, `PUT /products`
- **[orders](https://github.com/Prakash-Ravichandran/go-ecommerce-api/tree/main/internal/products)** — `GET /orders`, `POST /orders`

---

## Project Setup

```bash
git clone https://github.com/Prakash-Ravichandran/go-ecommerce-api.git
```

**Dependencies:**

- [Docker](https://docs.docker.com/desktop/setup/install/windows-install/#start-docker-desktop)
- [Goose](https://github.com/pressly/goose)
- [SQLC](https://docs.sqlc.dev/en/latest/)

### Running the API

```bash
cd go-ecommerce-api
docker compose up
```

```bash
cd cmd
go run .
```

---

## Adding a New Table

**1. Create a new migration file:**

```bash
goose -s create create_products sql
# result: 00001_create_products.sql
```

**2. Run the migrations:**

```bash
goose up
# result: table created in DB
```

**3. Generate SQLC code:**

```bash
sqlc generate
# result: generates models, interfaces, and Go code for DB operations
```

---

## Clean Layered Architecture

| Layer          | Responsibility                                 | Examples                            |
| -------------- | ---------------------------------------------- | ----------------------------------- |
| **Transport**  | Handles incoming requests & outgoing responses | HTTP, gRPC, REST, GraphQL           |
| **Service**    | Contains business logic & orchestration        | Use cases, validators, transformers |
| **Repository** | Manages data access & persistence              | DB queries, ORM, external APIs      |

- **Transport** depends on both Service and Repository
- **Service** depends only on Repository

> Flow: `Transport → Service → Repository`

---

## HTTP Server Setup

- [Add structs for HTTP server](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/890847824557e7566fbfef66b78dcb09dd348170)
- Request flow: `user → handler GET /products → service getProducts → DB SELECT * FROM products`
- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/7b9d8aab0df2ce17917d7367501be29e87dde9fe)

---

## Structured Logging

Instead of `fmt.Println`, use structured logging with levels (`info`, `error`) and metadata.

```go
slog.SetDefault(logger) // replaces r.Use(middleware.Logger)
```

- [Go structured logging blog](https://go.dev/blog/slog)

---

## The `internal/` Folder

The `internal/` directory restricts access — external packages cannot import anything inside it.

- [Reference: Golang internal packages](https://www.bytesizego.com/blog/golang-internal-package)

---

## Pointers — When to Use `*` and `&`

`NewHandler` returns `*handler` so that:

- **Memory is efficient** — the handler struct is not copied on every call
- **State is shared** — all routes using the handler share the same logger, DB connection, etc.

```go
func NewHandler(s Service) *handler {
    return &handler{
        service: s,
    }
}
```

- [Golang pointers (video)](https://youtu.be/2XEQsJLsLN0?si=bAQUEkC2mONMBqKk)

---

## Handler & Service Layers for `/products`

- [Products handler commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/8124dc9ac98b430a8989e7486d82f39d347ca01b)
- [Custom JSON write package](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/0ef7ebf458d1c180b5556add763fd28952ba4b56)
- [Products service commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/7f6250ec98787c6ef7c04c6075b467bfaebd75dc)

---

## Products Service Deep Dive

### What is a service?

The service doesn't care about HTTP, JSON, or headers. It only enforces the rules of the app (e.g. only list products that are active).

### What is a context?

`context` is Go's way of managing request lifecycles — cancellation, timeouts, and tracing.

> If a user closes their browser mid-request, `ctx` signals cancellation. The service stops its DB query immediately, saving server resources.

### Why two versions of `ListProducts`?

| Version | Speaks         | Responsibility                                                               |
| ------- | -------------- | ---------------------------------------------------------------------------- |
| Handler | HTTP           | Extracts data from `http.Request`, calls service, formats JSON response      |
| Service | Business Logic | Talks to DB or cache, returns data or error — knows nothing about `w` or `r` |

### Data flow (request journey)

```
Router → Handler (peels HTTP layer) → Service (business logic) → Handler (formats & responds)
```

### Why use an interface for the service?

```go
type Service interface {
    ListProducts(ctx context.Context) error
}
```

- **Mocking for tests** — swap in a `mockService` without needing a real DB
- **Flexibility** — swap out the service implementation without touching the handler

---

## Database with SQLC

- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/ddc8d6f77562725027ab08df4fa48b10c87ee588)

---

## Migrations with Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Migrations track the history of changes to the database schema (adding tables, altering columns, etc.).

| Command      | Effect                       |
| ------------ | ---------------------------- |
| `goose up`   | Apply pending migrations     |
| `goose down` | Roll back the last migration |

**Create a migration file:**

```bash
goose -s create create_products sql
# result: 00001_create_products.sql
```

- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/6118520a8a395e417d6c7f3cdffa04defc045acc)

---

## Running Postgres Locally

```bash
cd go-ecommerce-api
docker compose up
docker compose down
```

---

## Connecting to Postgres with Go

- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/76452c5060fc45c099fe85b3278a4ffb99c051fe)

---

## GET `/products/:id`

- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/ed5c65ac1289fcab77b5ae8cb2cc0117e1c87c83)

---

## POST `/orders` — Order Creation

- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/a17278195b210527a6091b377c8151a4b5c77372)
- Uses [database transactions](https://github.com/Prakash-Ravichandran/go-ecommerce-api/blob/main/internal/orders/service.go#L55)

**SQLC generated models:**

```go
type Order struct {
    ID         int64              `json:"id"`
    CustomerID int64              `json:"customer_id"`
    CreatedAt  pgtype.Timestamptz `json:"created_at"`
}

type OrderItem struct {
    ID         int64 `json:"id"`
    OrderID    int64 `json:"order_id"`
    ProductID  int64 `json:"product_id"`
    Quantity   int32 `json:"quantity"`
    PriceCents int32 `json:"price_cents"`
}
```

### Request & response examples

**✅ Success — create an order (201)**

```json
// Request
{ "customerId": 33, "items": [{ "productId": 1, "quantity": 4 }, { "productId": 2, "quantity": 6 }] }

// Response
{ "id": 18, "customer_id": 67, "created_at": "2026-05-05T11:30:00.363912+05:30" }
```

**❌ Product not found**

```json
{
  "customerId": 67,
  "items": [
    { "productId": 1, "quantity": 4 },
    { "productId": 1000, "quantity": 6 }
  ]
}
```

**❌ Not enough stock**

```json
{
  "customerId": 67,
  "items": [
    { "productId": 1, "quantity": 4 },
    { "productId": 2, "quantity": 10000 }
  ]
}
```

**❌ Customer ID is required**

```json
{ "items": [{ "productId": 1, "quantity": 4 }] }
```

**❌ At least one item is required**

```json
{ "customerId": 67 }
```

---

## GET `/orders` & GET `/orders/:id`

- [GET /orders with active service](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/3cc16c84d1ee84444f170ad9a8ca07924bb9111a)
- [GET /orders with dummy service](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/fafee51fb3a56cc80eb9113af3cc232a2e3e40fc)
- [GET /orders/:id](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/49b2bd8c56f3d86ab0063f920597df5e9f069714) — e.g. `http://127.0.0.1:8080/orders/1`

---

## POST `/products`

> Send a unique product on each request.

```json
// Request
{ "id": 101, "name": "Wireless Mechanical Keyboard", "price_in_cents": 8900, "quantity": 50, "created_at": "2026-05-06T10:00:00Z" }

// Response (201)
{ "id": 101, "name": "Wireless Mechanical Keyboard", "price_in_cents": 8900, "quantity": 50, "created_at": "2026-05-06T10:00:00Z" }
```

---

## PUT `/products`

```json
// Request
{ "id": 17, "price_in_cents": 120 }

// Response
{ "id": 17, "name": "Omen", "price_in_cents": 120, "quantity": 10, "created_at": "2026-05-06T07:00:21.578824+05:30" }
```

---

## SQLC Query Parameters (`$1`, `$2`)

```sql
-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 AND customer_id = $2;
```

SQLC maps positional parameters to Go function arguments:

```go
func (q *Queries) GetOrder(ctx context.Context, id int64, customerID int64) (Order, error)
```

`$1` → `id`, `$2` → `customerID`

---

## Dependency Injection

By accepting an **interface** instead of a concrete struct, you gain:

- **Testability** — inject a `mockService` in tests without a real DB
- **Flexibility** — swap implementations without changing handler code

```go
type Service interface {
    ListProducts(ctx context.Context) error
}

func NewHandler(s Service) *handler {
    return &handler{service: s}
}
```
