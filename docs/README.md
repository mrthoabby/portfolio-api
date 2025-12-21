# API Documentation

This directory contains the complete API documentation for the Portfolio API.

## Files

- **[openapi.yaml](./openapi.yaml)** - OpenAPI 3.0.3 specification
  - Import into Postman, Swagger UI, Insomnia, or any OpenAPI-compatible tool
  - Contains complete API schema, endpoints, request/response models, and examples

- **[POSTMAN_COLLECTION.md](./POSTMAN_COLLECTION.md)** - Guide for using the API with Postman
  - Instructions for importing the OpenAPI spec
  - Environment variables setup
  - Example tests

## Using the Documentation

### Postman

1. Open Postman
2. Click **Import**
3. Select **File** and choose `openapi.yaml`
4. Postman will automatically create a collection with all endpoints

### Swagger UI

1. Install Swagger UI:
   ```bash
   docker run -p 8081:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(pwd)/docs:/docs swaggerapi/swagger-ui
   ```
2. Open `http://localhost:8081` in your browser

### Insomnia

1. Open Insomnia
2. Go to **Application** > **Preferences** > **Data** > **Import Data**
3. Select **OpenAPI 3.0** and choose `openapi.yaml`

## API Base URL

- **Local Development**: `http://localhost:8080`
- **Production**: Update in `openapi.yaml` servers section

## Quick Reference

### Endpoints

- `GET /health` - Health check
- `GET /api/v1/profiles/{id}` - Get profile
- `GET /api/v1/profiles/{id}/skills` - Get skills
- `GET /api/v1/profiles/{id}/projects` - Get projects
- `GET /api/v1/profiles/{id}/certificates` - Get certificates
- `POST /api/v1/profiles/{id}/contacts` - Create contact

### Test Profile ID

Use this profile ID for testing:
```
123e4567-e89b-12d3-a456-426614174000
```



