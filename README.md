# LOAN BILLING SYSTEM
A RESTful API for managing loan billing schedules, payments, and delinquency
tracking. Built with Go, PostgreSQL, and Redis.

## Features
- Create a loan with weekly installments (flat interest rate 10% p.a., 50 weeks)
- Get outstanding balance at any point
- Check delinquency (missed 2 consecutive weekly payments)
- Make payments with idempotency support (Redis)
- Full API documentation via Swagger UI

## Tech Stack
- Language: Go 1.24+
- Framework: Gin
- Database: PostgreSQL 15
- Cache: Redis 7 (for idempotency keys)
- Containerization: Docker & docker-compose
- Documentation: Swagger (swaggo)

## Prerequisites
- Docker and Docker Compose installed on your machine.
- (Optional) Go 1.24+ for local development.

## Quick Start With Docker
1. Clone the repository
    ```bash
   git clone https://github.com/yourusername/loan-billing.git
   cd loan-billing
   ```

2. Run with docker-compose
   ```bash
   docker-compose up --build
    ```
   This starts:
   - PostgreSQL on localhost:5432
   - Redis on localhost:6379
   - The API on localhost:8080

3. Access Swagger UI:
  <br>Open your browser at http://localhost:8080/swagger/index.html to explore the API.</br>

## API ENDPOINTS
All endpoints are prefixed with ***/api/v1***.

1. #### Create a Loan
   <mark>**POST**</mark> /loans
   <br>Request body:
   ```json
   {
     "principal": 5000000,
     "interest_rate": 10,
     "term_weeks": 50,
     "start_date": "2026-02-18"
   }
   ```
   - principal: Loan amount (e.g., 5000000)
   - interest_rate: Annual interest rate (e.g., 10 for 10%)
   - term_weeks: Number of weeks (e.g., 50)
   - start_date: First due date (format YYYY-MM-DD)
   
   Response: 
   - 201 Created with the created loan object.

2. #### Get Outstanding Amount
   <mark>**GET**</mark> /loans/**{id}**/outstanding
   <br>Path parameter: id – Loan ID.
   <br>Response:
   ```json
   {
     "outstanding": 5500000
   }
   ```

3. #### Check Delinquency
   <mark>**GET**</mark> /loans/**{id}**/delinquent
   <br>Path parameter: id – Loan ID.
   <br>Response:
   ```json
   {
     "delinquent": false
   }
   ```

4. #### Make a Payment
   <mark>**POST**</mark> /loans/**{id}**/payments
   <br>Path parameter: id – Loan ID.
   <br>Headers: Idempotency-Key: <unique-string> (required)
   <br>Request body:
    ```json
   {
     "amount": 110000
   }
    ```
   - The amount must be a multiple of the weekly installment.
   - If any installments are overdue, the payment must cover **ALL** overdue weeks.

   ## Idempotency Key:
    - Use a new, unique key (e.g., a UUID v4) for each distinct payment operation.

    - If you need to retry the same payment request (e.g., due to a network timeout), reuse the same key – the server will detect the duplicate and not process it again.

    - Keys are stored in Redis for 24 hours. After expiration, the same key can be reused, but it's recommended to always generate a fresh one for clarity and to avoid accidental collisions.

   **Response**: 200 OK with success message.

## LOCAL DEVELOPMENT (WITHOUT DOCKER)

1. Install Go dependencies
   ```bash
   go mod download
   ```

2. Set up environment variables
   <br>Create a .env file in the project root:
   ```bash
   DB_HOST={your_db_host}
   DB_PORT={your_db_port}
   DB_USER={your_db_user}
   DB_PASSWORD={your_db_password}
   DB_NAME={your_db_name}
   REDIS_ADDR={your_redis_addr}
   PORT=8080
   ```
   **OR**

   Copy from .env.example for default value
   ```bash
   cp .env.example .env
   ```

3. Start PostgreSQL and Redis (using Docker)
   ### PostgreSQL
   ```bash
   docker run -d --name loan-postgres \
     -e POSTGRES_USER=your_db_user \
     -e POSTGRES_PASSWORD=your_db_password \
     -e POSTGRES_DB=your_db_name \
     -p 5432:5432 \
     postgres:15-alpine
   ```
   ### Redis
   ```bash
   docker run -d --name loan-redis \
     -p 6379:6379 \
     redis:7-alpine
   ```

4. Run database migrations
   ```bash
   psql -h localhost -U user -d db_name -f migrations/001_init.sql
   (Password will be prompted; use "password".)
   ```

5. Run the application
   ```bash
   go run cmd/api/main.go
   ```

6. Generate Swagger docs (optional, for development)
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   swag init -g cmd/api/main.go
   ```

## PROJECT STRUCTURE
```bash
.
├── cmd
│   └── api
│       └── main.go
├── config
│   └── config.go
├── docker-compose.yml
├── dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── loan
│   │   ├── handler
│   │   │   └── http
│   │   │       └── handler.go
│   │   ├── loan_repository.go
│   │   ├── loan_usecase.go
│   │   ├── repository
│   │   │   └── loan_repository.go
│   │   └── usecase
│   │       └── loan_usecase.go
│   └── payment
│       ├── handler
│       │   └── http
│       │       └── handler.go
│       ├── payment_repository.go
│       ├── payment_usecase.go
│       ├── repository
│       │   └── payment_repository.go
│       └── usecase
│           └── payment_usecase.go
├── migrations
│   └── 001_init.sql
├── models
│   ├── loan.go
│   └── payment.go
├── pkg
│   ├── idempotency
│   │   └── redis.go
│   ├── postgres
│   │   └── client.go
│   └── redis
│       └── client.go
└── README.md
```
<!-- 
## TESTING
Run unit tests:
   go test ./...

Integration tests (requiring real PostgreSQL/Redis) can be added separately. -->

ACKNOWLEDGEMENTS
----------------
- Gin Web Framework (github.com/gin-gonic/gin)
- Swaggo (github.com/swaggo/swag) for Swagger integration
- lib/pq PostgreSQL driver (github.com/lib/pq)
- go-redis Redis client (github.com/go-redis/redis)
- joho/godotenv for environment loading (github.com/joho/godotenv)