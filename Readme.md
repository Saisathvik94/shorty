# Advanced URL Shortner

A production-oriented URL shortener built with Go, PostgreSQL, Redis, and Gin.

## Current Progress

### Phase 1 — Core URL Shortener ✅

* [x] Go + Gin HTTP API
* [x] URL creation
* [x] Short code generation using `crypto/rand`
* [x] PostgreSQL persistence
* [x] URL redirection
* [x] URL expiration
* [x] URL deactivation
* [x] URL deletion
* [x] Update URL expiration
* [x] Repository → Service → Handler architecture
* [x] PostgreSQL connection pooling with `pgxpool`
* [x] Graceful HTTP server shutdown
* [x] Database migrations with `golang-migrate`

### Phase 2 — Caching ✅

* [x] Redis integration
* [x] Cache-Aside pattern
* [x] Redis TTL
* [x] TTL tied to URL expiration
* [x] Cache invalidation after updates/deactivation/deletion
* [x] Negative caching for missing URLs
* [x] Cache stampede protection with `singleflight`
* [x] PostgreSQL remains the source of truth

### Phase 3 — Concurrency & Background Work 🚧

Currently implementing:

* [x] Goroutines fundamentals
* [x] `context.Context`
* [x] Channels
* [x] `time.Ticker`
* [x] `sync.WaitGroup` fundamentals
* [ ] Expired URL cleanup worker
* [ ] Batch processing
* [ ] Retry and backoff
* [ ] Worker pool
* [ ] Job queue
* [ ] Idempotent background jobs
* [ ] Background job failure handling

## Architecture

```text
Client
  ↓
Gin
  ↓
Handler
  ↓
Service
  ├── Redis
  └── Repository
        ↓
    PostgreSQL

Background Workers
        ↓
    Repository
        ↓
    PostgreSQL
```

## Tech Stack

* Go
* Gin
* PostgreSQL
* pgx / pgxpool
* Redis
* golang-migrate
* Docker
* Neon PostgreSQL

## Learning Goal

Shorty is being built incrementally to understand real backend engineering concepts such as:

* Caching
* Concurrency
* Background processing
* Rate limiting
* Reliability
* Observability
* Distributed systems
* Production infrastructure
* Scaling
