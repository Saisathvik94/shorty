# Advanced URL Shortner

# Shorty

**Shorty** is a production-oriented URL shortener built with Go.

The project is being developed incrementally to understand how real backend systems are designed, optimized, and operated under increasing traffic and failure conditions.

The goal isn't just to build a URL shortener. The goal is to use a relatively simple system to learn backend engineering concepts such as:

* API design
* PostgreSQL
* Connection pooling
* Caching
* Cache invalidation
* Concurrency
* Background workers
* Rate limiting
* Reliability
* Observability
* Distributed systems
* Production infrastructure
* Scaling

---

# Current Status

🚧 **Currently Building: Background Workers & Concurrency**

The core URL-shortening functionality and caching layer are implemented.

The current focus is building an **expired URL cleanup worker** and using it to understand Go's concurrency and background-processing model.

---

# Tech Stack

| Technology     | Purpose                                  |
| -------------- | ---------------------------------------- |
| Go             | Backend language                         |
| Gin            | HTTP framework                           |
| PostgreSQL     | Primary database                         |
| pgx / pgxpool  | PostgreSQL driver and connection pooling |
| Redis          | Caching                                  |
| golang-migrate | Database migrations                      |
| Docker         | Local infrastructure                     |
| Neon           | Hosted PostgreSQL                        |

---

# Architecture

```text
                         ┌───────────────┐
                         │    Client     │
                         └───────┬───────┘
                                 │
                                 ▼
                         ┌───────────────┐
                         │     Gin       │
                         │    Router     │
                         └───────┬───────┘
                                 │
                                 ▼
                         ┌───────────────┐
                         │    Handler    │
                         └───────┬───────┘
                                 │
                                 ▼
                         ┌───────────────┐
                         │    Service    │
                         └───────┬───────┘
                                 │
                    ┌────────────┴────────────┐
                    ▼                         ▼
             ┌─────────────┐          ┌──────────────┐
             │    Redis    │          │  Repository  │
             │    Cache    │          └──────┬───────┘
             └─────────────┘                 │
                                             ▼
                                      ┌──────────────┐
                                      │  PostgreSQL  │
                                      │   Source of  │
                                      │     Truth    │
                                      └──────────────┘


                  Background Processing

                    ┌─────────────────┐
                    │ Cleanup Worker  │
                    └────────┬────────┘
                             │
                             ▼
                       ┌─────────────┐
                       │ Repository  │
                       └──────┬──────┘
                              │
                              ▼
                       ┌─────────────┐
                       │ PostgreSQL  │
                       └─────────────┘
```

---

# Project Structure

```text
shorty/
├── apps/
│   └── api/
│       ├── cmd/
│       │   └── server/
│       │       └── main.go
│       │
│       ├── internal/
│       │   ├── cache/
│       │   │   └── redis.go
│       │   │
│       │   ├── database/
│       │   │   └── postgres.go
│       │   │
│       │   ├── handlers/
│       │   │   └── url_handler.go
│       │   │
│       │   ├── repository/
│       │   │   └── url_repository.go
│       │   │
│       │   ├── services/
│       │   │   └── shortener.go
│       │   │
│       │   └── workers/
│       │       └── cleanup.go
│       │
│       ├── migrations/
│       │   ├── 000001_create_urls.up.sql
│       │   └── 000001_create_urls.down.sql
│       │
│       ├── .env
│       ├── .env.example
│       ├── go.mod
│       └── go.sum
│
├── docker-compose.yml
├── .gitignore
└── README.md
```

---

# What Shorty Can Do

## URL Creation

The API accepts an original URL and an expiration time.

The service:

1. Parses the URL.
2. Validates the scheme.
3. Validates that the URL contains a hostname.
4. Validates the expiration timestamp.
5. Ensures the expiration is in the future.
6. Limits expiration to 30 days.
7. Generates a random short code.
8. Stores the URL in PostgreSQL.

### Short Code Generation

Short codes are generated using Go's cryptographically secure random generator.

The current short code length is **6 characters**.

The system also handles possible database uniqueness collisions by generating another code and retrying.

---

# PostgreSQL

PostgreSQL is the primary persistence layer and the **source of truth** for URL data.

Current schema:

```text
urls
├── id
├── short_code
├── original_url
├── created_at
├── expires_at
└── is_active
```

The database stores:

* Original URLs
* Short codes
* Creation time
* Expiration time
* Active/inactive state

Database migrations are managed using `golang-migrate`.

---

# Repository Layer

Database operations are isolated inside the repository layer.

Current repository responsibilities include:

* Creating URLs
* Retrieving URLs by short code
* Deactivating URLs
* Deleting URLs
* Updating expiration times
* Deactivating expired URLs

The service layer does not directly interact with PostgreSQL.

```text
Service
   ↓
Repository
   ↓
PostgreSQL
```

This keeps database-specific operations separate from business logic.

---

# Service Layer

The service layer contains URL business logic.

Responsibilities include:

* URL validation
* Expiration validation
* Short code generation
* Collision handling
* Cache interaction
* URL expiration checks
* Active/inactive validation

The service coordinates the cache and repository rather than exposing infrastructure details to the handlers.

---

# Handler Layer

Gin handlers are responsible for HTTP-specific concerns.

Current endpoints:

```text
POST   /api/urls
GET    /:shortCode
PATCH  /api/urls/:shortCode/deactivate
DELETE /api/urls/:shortCode
PUT    /api/urls/:shortCode/expiration
```

Handlers:

1. Receive HTTP requests.
2. Bind and validate request data.
3. Call the service layer.
4. Translate service errors into HTTP responses.
5. Return JSON responses or redirects.

---

# Connection Pooling

PostgreSQL connections are managed using `pgxpool`.

Instead of opening a new database connection for every request:

```text
Request
   ↓
Connection Pool
   ↓
Available PostgreSQL connection
   ↓
Query
   ↓
Connection returned to pool
```

This provides controlled connection reuse and prepares the application for handling concurrent requests.

Connection-pool tuning for high traffic is planned for the production and scaling work.

---

# Graceful Shutdown

The HTTP server supports graceful shutdown.

When the application receives a termination signal:

```text
Signal
  ↓
Stop accepting new requests
  ↓
Allow existing requests to finish
  ↓
Close Redis
  ↓
Close PostgreSQL pool
  ↓
Exit
```

This prevents the application from abruptly terminating while requests or infrastructure resources are still in use.

---

# Caching

Shorty uses Redis to reduce unnecessary PostgreSQL reads and improve URL lookup performance.

## Cache-Aside

The request flow is:

```text
GET /abc123
      ↓
    Redis
      │
   ┌──┴──┐
   │     │
 HIT    MISS
   │     │
   │     ▼
   │  PostgreSQL
   │     │
   │     ▼
   │   Redis
   │
   ▼
Original URL
```

On a cache hit, PostgreSQL is not queried.

On a cache miss, the database is queried and the result is placed into Redis.

---

## Redis TTL

Each cached URL has a TTL based on the URL's expiration time.

```text
URL expires at
       ↓
Calculate remaining lifetime
       ↓
Redis TTL
```

This prevents cached URL data from living longer than the URL itself.

---

## Cache Invalidation

The cache is invalidated whenever URL state changes.

```text
Deactivate URL
      ↓
PostgreSQL update
      ↓
Delete Redis key
```

The same approach is used when:

* A URL is deactivated
* A URL is deleted
* A URL's expiration is updated

This prevents stale URL state from remaining in Redis after writes.

---

## Negative Caching

Shorty caches missing URLs for a short period.

```text
GET /doesnotexist
        ↓
    Redis MISS
        ↓
   PostgreSQL
        ↓
   No matching row
        ↓
Redis: NOT_FOUND
```

The current negative-cache TTL is **30 seconds**.

This prevents repeated requests for the same nonexistent short code from repeatedly reaching PostgreSQL.

---

## Cache Stampede Protection

Shorty uses Go's `singleflight` package to prevent multiple concurrent requests for the same uncached URL from simultaneously querying PostgreSQL.

Without protection:

```text
Request 1 ──→ Redis MISS ──→ PostgreSQL
Request 2 ──→ Redis MISS ──→ PostgreSQL
Request 3 ──→ Redis MISS ──→ PostgreSQL
Request 4 ──→ Redis MISS ──→ PostgreSQL
```

With `singleflight`:

```text
Request 1 ──┐
Request 2 ──┤
Request 3 ──┼──→ singleflight ──→ PostgreSQL
Request 4 ──┘
                         │
                         ▼
                       Redis
                         │
                  shared result
```

This reduces duplicate database work during concurrent cache misses.

---

# What I Learned

## Backend Architecture

* [x] Separating HTTP handlers from business logic
* [x] Service layer responsibilities
* [x] Repository pattern
* [x] Dependency injection through constructors
* [x] Composition root in `main.go`
* [x] Keeping infrastructure concerns separated

---

## PostgreSQL & Database Engineering

* [x] PostgreSQL fundamentals
* [x] SQL `INSERT`, `SELECT`, `UPDATE`, and `DELETE`
* [x] Database migrations
* [x] PostgreSQL constraints
* [x] Unique constraint handling
* [x] PostgreSQL error codes
* [x] `pgx`
* [x] Connection pooling with `pgxpool`
* [x] `RowsAffected()`
* [x] Using PostgreSQL as the source of truth

---

## URL Shortening

* [x] URL validation
* [x] Expiration handling
* [x] Cryptographically secure random generation
* [x] Short-code collision handling
* [x] Active/inactive URL state

---

## Caching

* [x] Redis fundamentals
* [x] Cache-Aside pattern
* [x] Redis TTL
* [x] TTL based on business expiration
* [x] Cache invalidation
* [x] Negative caching
* [x] Cache stampede
* [x] Request coalescing with `singleflight`
* [x] Understanding Redis as a cache rather than the source of truth

---

## Go Concurrency

* [x] Goroutines
* [x] How goroutines differ from JavaScript async/promises
* [x] Why the Go process exits when `main()` returns
* [ ] `context.Context`
* [ ] Context cancellation
* [ ] Context propagation
* [x] Channels
* [ ] Blocking behavior of channels
* [ ] `time.Ticker`
* [x] `sync.WaitGroup`
* [ ] Basic background-worker architecture

---

# What's Next to Build

## Background Workers

**Currently building**

* [ ] Expired URL cleanup worker
* [ ] Periodic cleanup using `time.Ticker`
* [ ] Context-aware worker cancellation
* [ ] Batch processing
* [ ] Retry and backoff
* [ ] Worker pool
* [ ] Job queues
* [ ] Idempotent background jobs
* [ ] Background job failure handling

### First Worker

The first worker will periodically find URLs where:

```text
is_active = true
AND
expires_at <= NOW()
```

and deactivate them.

```text
Cleanup Worker
      ↓
Periodic trigger
      ↓
Repository
      ↓
PostgreSQL
      ↓
Deactivate expired URLs
      ↓
Return affected row count
```

---

# Upcoming Engineering Work

## Rate Limiting & Security

* [ ] Rate limiting
* [ ] Token bucket / sliding window concepts
* [ ] Redis-based rate limiting
* [ ] Abuse prevention
* [ ] Request validation
* [ ] URL security considerations
* [ ] Concurrency limits

---

## Analytics & Observability

* [ ] Click analytics
* [ ] Event processing
* [ ] Structured logging
* [ ] Metrics
* [ ] Prometheus
* [ ] Grafana
* [ ] Distributed tracing
* [ ] Health checks
* [ ] Error tracking
* [ ] Performance measurement

---

## Production & Distributed Systems

* [ ] Docker production setup
* [ ] Production configuration
* [ ] Load testing
* [ ] Database connection-pool tuning
* [ ] Redis reliability
* [ ] PostgreSQL performance
* [ ] AWS infrastructure
* [ ] Horizontal scaling
* [ ] Load balancing
* [ ] Stateless services
* [ ] Failure recovery
* [ ] Distributed systems concepts
* [ ] Queue-based architectures
* [ ] Kubernetes
* [ ] Kubernetes deployments
* [ ] ArgoCD / GitOps
* [ ] Production deployment
* [ ] Monitoring and alerting

---

# Engineering Principles

### PostgreSQL is the source of truth

Redis improves performance but does not own the data.

### Minimize unnecessary work

Caching, negative caching, and `singleflight` exist to prevent repeated or duplicate work.

### Failures should be expected

The project will progressively introduce failure scenarios rather than assuming infrastructure always works.

### Keep responsibilities separated

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
Database
```

Infrastructure such as Redis and background workers should have clear boundaries.

### Learn the primitive before the abstraction

Instead of immediately using a job framework, the project first implements basic Go concurrency concepts directly.

The goal is to understand what abstractions such as worker pools and job queues are actually solving.

---

# Project Goal

Shorty started as a simple URL shortener.

It is gradually becoming a **backend engineering laboratory** for understanding what happens when a simple service needs to handle:

```text
More traffic
     ↓
More concurrent requests
     ↓
More database pressure
     ↓
Caching
     ↓
Background processing
     ↓
Failures
     ↓
Observability
     ↓
Multiple instances
     ↓
Distributed systems
     ↓
Production infrastructure
```

The end goal is not simply to make Shorty work.

The goal is to understand **why each piece exists, what problem it solves, and what happens when it fails.**
