# Shorty - Production-Oriented URL Shortener

**A URL shortener built in Go to explore backend engineering, system design, reliability, and cloud infrastructure.**

Shorty started with a simple goal: shorten long URLs. But the real objective was to understand how backend systems are designed, deployed, secured, and scaled beyond basic CRUD operations.

Instead of focusing on unnecessary product features, Shorty focuses on the engineering challenges behind running a reliable service.

**Live Application:** https://shorty.saisathvik.site

**API:** https://api.shorty.saisathvik.site

*The AWS backend is deployed on demand to control infrastructure costs. The live API may be temporarily unavailable.*

---

## Tech Stack

| Layer | Technologies |
|---|---|
| Backend | Go, Gin |
| Frontend | React, TypeScript, Vite |
| Database | PostgreSQL, Neon, pgxpool |
| Caching | Redis, Upstash |
| Containerization | Docker, GHCR |
| Orchestration | Kubernetes, AWS EKS |
| Networking | AWS ALB, HTTPS, Ingress |
| Infrastructure | Terraform |
| Frontend Hosting | Cloudflare Pages |

## System Architecture

```text
                   Client / Browser
                          |
                          v
                  Cloudflare Pages
                    React + Vite
                          |
                          v
                    HTTPS Request
                          |
                          v
                   AWS Application
                    Load Balancer
                          |
                          v
                    EKS Ingress
                          |
                          v
                 Kubernetes Service
                          |
                          v
                     Go API Pods
                          |
               +----------+----------+
               |                     |
               v                     v
          Upstash Redis        Neon PostgreSQL
        Cache / Rate Limit      Source of Truth
```

The backend follows a layered architecture:

**Handler → Service → Repository → Database**

This keeps HTTP handling, business logic, and persistence responsibilities separate.

## Key Engineering Features

### URL Management
- Cryptographically generated short codes
- URL validation and collision handling
- Custom expiration, up to 30 days
- URL deactivation and deletion
- Expiration updates
- HTTP redirects

### Caching and Performance
- Redis Cache-Aside strategy
- Cache invalidation after mutations
- Negative caching for nonexistent URLs
- Cache stampede protection using Go `singleflight`
- Cache TTL aligned with URL expiration
- PostgreSQL as the source of truth

### Reliability
- Redis-backed rate limiting
- Database connection pooling
- Graceful server shutdown
- Kubernetes readiness and liveness probes
- Rolling deployments
- Horizontal Pod Autoscaler configured for CPU-based scaling

### Security
- Request body size limits
- CORS configuration
- HTTP security headers
- Non-root Docker containers
- Read-only container filesystem
- Dropped Linux capabilities
- Kubernetes NetworkPolicy

### Infrastructure
- Multi-stage Docker builds with a minimal `scratch` runtime
- Container images published to GitHub Container Registry
- AWS EKS deployment
- AWS Application Load Balancer with HTTPS
- Terraform-managed AWS infrastructure
- Cloudflare Pages frontend deployment

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/urls` | Create a shortened URL |
| `GET` | `/:shortCode` | Redirect to the original URL |
| `PUT` | `/api/urls/:shortCode/deactivate` | Deactivate a URL |
| `DELETE` | `/api/urls/:shortCode` | Delete a URL |
| `PATCH` | `/api/urls/:shortCode/expiration` | Update expiration |
| `GET` | `/health` | API health check |

## Project Structure

```text
shorty/
├── apps/
│   ├── api/
│   │   ├── cmd/server/
│   │   ├── internal/
│   │   │   ├── cache/
│   │   │   ├── database/
│   │   │   ├── handlers/
│   │   │   ├── middlewares/
│   │   │   ├── repository/
│   │   │   ├── services/
│   │   │   └── workers/
│   │   ├── migrations/
│   │   ├── tests/
│   │   └── Dockerfile
│   └── web/
├── infra/
│   └── terraform/
├── k8s/
└── docker-compose.yml
```

## What I Learned

Building Shorty helped me explore engineering challenges that aren't always visible when developing a simple application:

- Designing cache consistency and invalidation strategies
- Preventing redundant database queries during concurrent requests
- Implementing distributed rate limiting
- Managing database connections and application shutdown
- Understanding Kubernetes deployments, networking, and autoscaling
- Deploying and debugging containerized services on AWS
- Resolving ARM64/AMD64 image compatibility issues between Apple Silicon and AWS EC2
- Managing infrastructure lifecycle and cloud costs using Terraform

## What's Next?

Shorty is an ongoing backend engineering learning project.

The next phase focuses on:

- **Observability:** Structured logging, metrics, and monitoring
- **Load testing:** Measuring throughput and latency under concurrent traffic
- **Failure testing:** Understanding behavior during Redis, database, and application failures
- **Performance optimization:** Identifying bottlenecks using real measurements
- **Capacity planning:** Understanding resource requirements and scaling limits

---

**Shorty isn't about reinventing URL shortening. It's about understanding the engineering behind reliable backend systems.**

Built with Go, curiosity, and a lot of debugging.