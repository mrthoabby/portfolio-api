# Portfolio API

REST API for the [portfolio-ui](https://github.com/mrthoabby/portfolio-ui) project. Manages personal portfolio data: profile, skills, projects, certificates, and contact messages.

## Local Testing

```bash
# Build
docker build \
  --build-arg VERSION=local-test \
  --build-arg BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t portfolio-api:local-test \
  .

# Run (make sure PORT is set in .env file)
# Example: if PORT=3000 in .env, use -p 3000:3000
docker run -d \
  --name portfolio-api-test \
  --env-file .env \
  -p 3000:3000 \
  portfolio-api:local-test

# Test health check
curl http://localhost:3000/health

# View logs
docker logs -f portfolio-api-test
```

## Environment Variables

Required in `.env` file:
- `DATABASE_URL` - MongoDB connection string
- `DATABASE_NAME` - Database name
- `ALLOWED_ORIGINS` - Comma-separated CORS origins
- `PORT` - Server port
