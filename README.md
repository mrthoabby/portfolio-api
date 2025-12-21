# Portfolio API

REST API for the [portfolio-ui](https://github.com/mrthoabby/portfolio-ui) project. Manages personal portfolio data: profile, skills, projects, certificates, and contact messages.

## Deploy (Docker)

```bash
# Build
docker build -t portfolio-api .

# Run
docker run -p 8080:8080 \
  -e DATABASE_URL=mongodb://your-mongo-host:27017 \
  -e DATABASE_NAME=portfolio_db \
  -e ALLOWED_ORIGINS=https://yoursite.com \
  portfolio-api
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | ✅ | MongoDB connection string |
| `DATABASE_NAME` | ✅ | Database name |
| `ALLOWED_ORIGINS` | ✅ | Comma-separated CORS origins |
| `PORT` | ❌ | Server port (default: 8080) |
| `ENV` | ❌ | Environment (default: development) |

## Local Development

Requirements: Go 1.21+, MongoDB 4.4+

```bash
# Install dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your MongoDB connection

# Start MongoDB (if using Docker)
docker run -d -p 27017:27017 --name mongodb mongo:7

# Seed database
mongosh portfolio_db < mongodb-seed/insert.js

# Run
go run cmd/api/main.go
```

## Verify

```bash
curl http://localhost:8080/health
# {"status":"ok","database":"connected"}
```

## Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/profiles/{id}` | Get profile |
| GET | `/api/v1/profiles/{id}/skills` | Get skills |
| GET | `/api/v1/profiles/{id}/projects` | Get projects |
| GET | `/api/v1/profiles/{id}/certificates` | Get certificates |
| POST | `/api/v1/profiles/{id}/contacts` | Send contact message |
| POST | `/api/v1/profiles/{id}/questions` | Submit a question |

## Security Features

This API includes production-ready security measures:

| Feature | Description |
|---------|-------------|
| **Rate Limiting** | 100 req/min global, 5 req/min for contacts, 10 req/min for questions |
| **Body Size Limit** | Maximum 1MB request body |
| **Security Headers** | X-Content-Type-Options, X-Frame-Options, CSP, etc. |
| **Input Sanitization** | HTML/XSS protection on all inputs |
| **UUID Validation** | Profile IDs validated as proper UUIDs |
| **Request ID** | Unique ID per request for tracing |
| **Panic Recovery** | Graceful handling of unexpected errors |
| **Request Timeouts** | Read: 10s, Write: 30s, Idle: 120s |

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `BAD_REQUEST` | 400 | Invalid request format |
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `NOT_FOUND` | 404 | Resource not found |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `PAYLOAD_TOO_LARGE` | 413 | Request body too large |
| `INTERNAL_ERROR` | 500 | Server error |

## Documentation

- [OpenAPI Specification](./docs/openapi.yaml)
- [Postman Collection](./docs/POSTMAN_COLLECTION.md)

## Testing

```bash
go test ./... -cover
```

---

MIT License - Copyright © 2024 Danis Abadía
