# Packs Service

A **Go microservice** for managing **packs** using **Echo, PostgreSQL (Bun), and OpenAPI validation**.

---

## Features

✅ REST API with Echo and OpenAPI validation  
✅ **All handler signatures, request and response types are generated automatically
from [`api/openapi.yaml`](api/openapi.yaml)** via codegen  
✅ Clean layered architecture (handlers, repository, migration, config)  
✅ PostgreSQL with Bun ORM and auto migrations  
✅ Structured logging with `slog`  
✅ Request ID, timeout and other middlewares  
✅ Structured error handling  
✅ Easily testable architecture

---

## Prerequisites

- `make`
- `docker` and `docker compose`

---

## Run with Docker (Recommended)

No local Go or PostgreSQL installation required.

```bash
make dev
```

This will:

✅ Start the **PostgreSQL container**  
✅ Build the **Go application**  
✅ Run database migrations  
✅ Start the API server

The service will be accessible at:

```
http://localhost:8080
```

---

## Run Locally (Using Local Go and Postgres)

### Prerequisites

- **Go 1.22+**
- **PostgreSQL 14+**
- `make`
- **Docker**

### Steps

- Create a PostgreSQL database and run the db:

```bash
make run-db
```

- Copy  `.env.example` to `.env`:

- Start the server:

```bash
make dev-local
```

The service will be accessible at:

```
http://localhost:8080
```

---

## How to run tests

```bash
make test
```

or:

```bash
go test ./...
```

---

## Other Available Makefile Commands

```make
lint:       # Run golangci-lint
generate:   # Generates code: mocks, openapi...
```

---

## OpenAPI Specification

The **API is defined in [`api/openapi.yaml`](api/openapi.yaml)**.

✅ **All handler signatures, request and response types are generated automatically via codegen from this file.**  
✅ Ensures strict request/response validation and consistency across API changes.  
✅ Use this file to understand the full API surface or generate client SDKs for integration.

---

## API Usage

### Create Pack

```bash
curl --location 'localhost:8080/packs' \
--header 'Content-Type: application/json' \
--data '{
    "size": 250
}'
```

---

### Get Packs

```bash
curl --location 'localhost:8080/packs'
```

---

### Update Pack

```bash
curl --location --request PATCH 'localhost:8080/packs/<pack id>' \
--header 'Content-Type: application/json' \
--data '{
        "size": 500
}'
```

---

### Delete Pack

```bash
curl --location --request DELETE 'localhost:8080/packs/<pack id>'
```

---

### Calculate Packs and Items

```bash
curl --location 'localhost:8080/packs/calculate' \
--header 'Content-Type: application/json' \
--data '
{
    "items": 1000
}
'
```
