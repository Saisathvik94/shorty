# Shorty — High-Performance URL Shortener

A full-stack URL shortening platform built with **Go, Gin, PostgreSQL, and Redis**, designed to explore production-oriented backend engineering, caching, concurrency, performance optimization, and cloud infrastructure.

Shorty converts long URLs into compact links while supporting expiration, link management, Redis-backed caching, and rate limiting.

**Frontend:** https://shorty.saisathvik.site  
**API:** https://api.shorty.saisathvik.site

*Cloud infrastructure may be temporarily unavailable when the AWS environment is scaled down to control costs.*

## Performance

**39,518 successful redirects/second — measured locally with k6.**

| Metric | Result |
|---|---|
| Total requests | 1,185,571 |
| Successful redirects | 1,185,571 |
| Failed requests | 0 |
| Throughput | 39,518 RPS |
| Average response time | 1.20 ms |
| P95 response time | 2.06 ms |
| P90 response time | 1.75 ms |
| Concurrent virtual users | 50 |
| Duration | 30 seconds |

### Benchmark Environment

- **Machine:** Apple MacBook Air M5, 16 GB RAM
- **Backend:** Go + Gin running locally
- **Cache:** Native Redis running on macOS
- **Database:** PostgreSQL on Neon
- **Load generator:** Grafana k6, running on the same machine
- **Workload:** Repeated HTTP GET requests to a valid, warmed short URL
- **Rate limiting:** Disabled only in a dedicated local benchmark mode
- **Response:** HTTP 302 redirect, without following the destination URL

These are **local warmed-cache benchmark results**, not production throughput guarantees or measurements of the deployed AWS infrastructure.

### Docker Redis vs Native Redis

We investigated Redis as a potential performance bottleneck using Go's `pprof` CPU profiler and Redis connection-pool statistics.

| Metric | Docker Redis | Native Redis |
|---|---:|---:|
| Throughput | 22,367 RPS | 39,518 RPS |
| Total requests | 671,038 | 1,185,571 |
| Average latency | 2.14 ms | 1.20 ms |
| P95 latency | 3.44 ms | 2.06 ms |
| Error rate | 0% | 0% |
| Test duration | 30 seconds | 30 seconds |
| Virtual users | 50 | 50 |

Native Redis achieved approximately **76.7% higher throughput** in these local runs.

CPU profiling showed substantial activity in operating-system calls and Redis-related socket operations. Redis pool statistics reported approximately 671K pool hits, 50 misses, and zero timeouts, suggesting effective connection reuse rather than pool-acquisition failure.

This comparison highlights the cost of the local Docker Desktop networking and virtualization path. It does not imply that Redis containers have equivalent overhead on native Linux infrastructure.

### High-Load Stress Testing

A separate test attempted a configured arrival rate of **500,000 requests/second**.

With native Redis:

| Metric | Result |
|---|---:|
| Configured target | 500,000 RPS |
| Achieved throughput | 32,511 RPS |
| Completed requests | 162,775 |
| Successful completed requests | 100% |
| Dropped iterations | 1,717,134 |
| P95 latency | 40.93 ms |

The target was **not achieved**. The test demonstrated substantial load-generator or shared-machine contention under extreme concurrency.

Dropped k6 iterations are scheduled requests that never started, and are not counted as backend HTTP failures.

The next scaling milestone is benchmarking Shorty with load generation running on separate machines and profiling individual infrastructure components.

## Architecture

```text
                  Client / Browser
                         |
                         v
                  React Frontend
                         |
                         v
                    Go + Gin API
                         |
              +----------+----------+
              |                     |
              v                     v
       Redis Rate Limiter      URL Service
                                    |
                              Cache-Aside
                                    |
                         +----------+----------+
                         |                     |
                         v                     v
                    Redis Cache          PostgreSQL
                         |               (NeonDB)
                         |
                         v
                  HTTP 302 Redirect
```

On a warmed cache hit, Shorty resolves the destination from Redis without querying PostgreSQL.

PostgreSQL remains the source of truth for URL mappings and lifecycle state.

## Features

- Cryptographically generated six-character short codes
- URL validation for HTTP and HTTPS destinations
- PostgreSQL persistence using `pgxpool`
- Redis cache-aside implementation with TTL management
- Cache invalidation when links are modified or disabled
- Negative caching for missing URLs
- Go `singleflight` to mitigate concurrent cache misses
- URL expiration with a configurable maximum lifetime
- Link deactivation and deletion
- Redis-backed IP rate limiting for read and write operations
- Fail-open rate limiting when Redis encounters an error
- Request body size limits and HTTP security headers
- Periodic expired-link cleanup worker
- Graceful HTTP server shutdown
- Local profiling with Go `pprof`
- Kubernetes deployment with health probes and horizontal pod autoscaling

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React, TypeScript, Vite |
| Backend | Go, Gin |
| Database | PostgreSQL, Neon |
| Cache | Redis, go-redis |
| Database driver | pgx |
| Rate limiting | redis_rate |
| Load testing | Grafana k6 |
| Profiling | Go pprof |
| Containerization | Docker |
| Container registry | GitHub Container Registry |
| Orchestration | Kubernetes, Amazon EKS |
| Infrastructure as code | Terraform |
| Networking | AWS Application Load Balancer |
| Frontend hosting | Cloudflare Pages |

## API

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/urls` | Create a shortened URL |
| GET | `/:shortCode` | Redirect to original URL |
| PUT | `/api/urls/:shortCode/deactivate` | Deactivate a link |
| DELETE | `/api/urls/:shortCode` | Delete a link |
| PATCH | `/api/urls/:shortCode/expiration` | Update expiration |
| GET | `/health` | Health check |

### Create a short URL

```bash
curl -X POST http://localhost:3000/api/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "expires_at": "2026-11-01T18:00:00Z"
  }'
```

### Resolve a short URL

```bash
curl -i http://localhost:3000/YOUR_SHORT_CODE
```

A valid active link returns an HTTP 302 response with a `Location` header pointing to the original URL.

## Run Locally

### Requirements

- Go
- PostgreSQL database
- Redis
- Node.js or Bun for the frontend
- Grafana k6 for benchmarking (optional)

Clone the repository:

```bash
git clone https://github.com/Saisathvik94/Shorty.git
cd Shorty
```

Configure your backend environment in `apps/api/.env`:

```env
PORT=3000
DATABASE_URL=your_postgresql_connection_string
REDIS_URL=redis://localhost:6379

RATE_LIMIT_PER_MINUTE=10
RATE_LIMIT_REDIRECT_PER_MINUTE=100
MAX_BODY_SIZE=16384

APP_ENV=development
BENCHMARK_MODE=false
```

Set the database URL to your own PostgreSQL instance and apply the migrations included in the project.

Start Redis, then run the API:

```bash
cd apps/api
go run ./cmd/server
```

Verify that the server is running:

```bash
curl http://localhost:3000/health
```

Expected response:

```text
Shorty API
```

## Load Testing

Benchmark mode is intended exclusively for trusted local development environments.

```env
APP_ENV=local
BENCHMARK_MODE=true
REDIS_URL=redis://localhost:6379
```

When enabled locally, benchmark mode bypasses the redirect-specific rate limiter, allowing measurements of the cached redirect path without HTTP 429 rejection distorting throughput results.

Production environments must always keep the rate limiter enabled.

### Example k6 configuration

```javascript
import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 50,
  duration: "30s",
  maxRedirects: 0,
  discardResponseBodies: true,
};

export default function () {
  const res = http.get(
    "http://localhost:3000/YOUR_SHORT_CODE",
    { redirects: 0 }
  );

  check(res, {
    "redirect successful": (r) =>
      r.status === 302 && Boolean(r.headers.Location),
  });
}
```

Replace `YOUR_SHORT_CODE` with an active URL and warm the cache before benchmarking.

Run:

```bash
k6 run tests/load/redirect.js
```

This test uses a fixed number of virtual users. It measures completed HTTP requests, rather than guaranteeing a specific arrival rate.

### CPU profiling

Shorty exposes a separate pprof server at `127.0.0.1:6060` in local development.

Capture a CPU profile while a benchmark is running:

```bash
go tool pprof -seconds 10 \
  http://127.0.0.1:6060/debug/pprof/profile
```

Inspect hotspots using:

```text
(pprof) top 15
```

Profiling endpoints should never be exposed publicly.

## Deployment

Shorty has been deployed using:

- A multi-stage Docker build
- GitHub Container Registry for container images
- Amazon EKS for Kubernetes orchestration
- Kubernetes Deployment, Service, and Ingress resources
- Application Load Balancer for external traffic
- Horizontal Pod Autoscaler for CPU-based scaling
- Terraform for infrastructure provisioning and teardown
- Cloudflare Pages for frontend hosting

The backend is designed to remain stateless across API replicas, using shared PostgreSQL and Redis services.

## Engineering Learnings

Shorty was built to explore practical backend and infrastructure challenges beyond basic CRUD operations.

Key areas explored include cache-aside design, expiration and invalidation, rate limiting, cache stampede prevention, connection pooling, graceful shutdown, container image compatibility across architectures, Kubernetes deployment, distributed load testing, and CPU profiling.

The performance experiments reinforced a central lesson:

**Throughput, latency, and load-generator capacity must be measured independently. A configured benchmark target is not the same as achieved backend throughput.**

## Future Improvements

- Benchmark with independent load-generator machines
- Reproduce benchmarks across repeated runs
- Investigate HTTP and Redis socket overhead
- Profile allocations and garbage collection
- Measure cold-cache and database performance
- Add detailed request and cache observability
- Evaluate scaling across multiple Go API replicas
- Measure production-like latency with remote cache infrastructure

---

Built as a hands-on project for learning Go backend engineering, performance analysis, and cloud-native system design.
