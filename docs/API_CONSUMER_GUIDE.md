# Portfolio API - Consumer Guide

> Complete documentation for AI agents and developers consuming this API.

---

## 🎯 API Overview

This is a **read-mostly** REST API for a personal portfolio website. It provides endpoints to retrieve profile information, skills, projects, certificates, and allows visitors to send contact messages.

**Base URL**: `https://your-api-domain.com`

---

## 🔐 Authentication & Security

### No Authentication Required
This is a **public API**. No API keys or tokens are needed.

### CORS Policy
- Only requests from **whitelisted origins** are allowed
- The server must have your domain in `ALLOWED_ORIGINS`
- Credentials are **NOT** allowed (`credentials: 'omit'`)

### Allowed Headers
```
Accept
Content-Type
X-Request-ID (optional, for tracing)
```

### Allowed Methods
```
GET, POST, OPTIONS
```

---

## 📊 Rate Limiting

| Scope | Limit |
|-------|-------|
| Global (all endpoints) | 100 requests / minute |
| Contact endpoint (POST) | 5 requests / minute |
| Questions endpoint (POST) | 10 requests / minute |

**When exceeded**: HTTP `429 Too Many Requests`

---

## 📦 Request Constraints

| Constraint | Value |
|------------|-------|
| Max body size | 1 MB |
| Content-Type for POST | `application/json` |
| Response format | JSON |

---

## 🛣️ Endpoints

### Health Check

```http
GET /health
```

**Response** `200 OK`:
```json
{
  "status": "ok",
  "database": "connected",
  "version": {
    "version": "v1.2.0",
    "gitCommit": "a1b2c3d",
    "buildDate": "2025-12-20T15:30:00Z"
  }
}
```

---

### Get Profile

```http
GET /api/v1/profiles/{profileId}
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| profileId | UUID string | Yes | The profile identifier |

**Response** `200 OK`:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "photoUrl": "https://example.com/photo.jpg",
  "title": "Full Stack Developer",
  "aboutMe": "Passionate developer with 5+ years of experience...",
  "firstExperienceDate": "2019-01-15T00:00:00Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-12-01T00:00:00Z"
}
```

**Response** `404 Not Found`:
```json
{
  "error": {
    "type": "NOT_FOUND",
    "message": "Profile not found"
  }
}
```

---

### Get Skills

```http
GET /api/v1/profiles/{profileId}/skills
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| profileId | UUID string | Yes | The profile identifier |

**Response** `200 OK`:
```json
{
  "skills": [
    {
      "id": "skill-001",
      "profileId": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Go",
      "category": "backend",
      "proficiency": "advanced"
    },
    {
      "id": "skill-002",
      "profileId": "550e8400-e29b-41d4-a716-446655440000",
      "name": "React",
      "category": "frontend",
      "proficiency": "advanced"
    }
  ]
}
```

**Skill Categories**:
- `backend`
- `frontend`
- `tools`
- `softSkills`

**Proficiency Levels**:
- `advanced`
- `occasional`
- `past`

---

### Get Projects

```http
GET /api/v1/profiles/{profileId}/projects
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| profileId | UUID string | Yes | The profile identifier |

**Response** `200 OK`:
```json
{
  "projects": [
    {
      "id": "project-001",
      "profileId": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Portfolio API",
      "description": "RESTful API for personal portfolio",
      "techStack": ["Go", "MongoDB", "Docker"],
      "githubUrl": "https://github.com/user/portfolio-api",
      "liveUrl": "https://api.portfolio.com",
      "imageDiagramUrl": "https://example.com/diagram.png",
      "visible": true,
      "createdAt": "2024-06-15T00:00:00Z"
    }
  ]
}
```

**Note**: Only projects with `visible: true` are returned.

---

### Get Certificates

```http
GET /api/v1/profiles/{profileId}/certificates
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| profileId | UUID string | Yes | The profile identifier |

**Response** `200 OK`:
```json
{
  "certificates": [
    {
      "id": "cert-001",
      "profileId": "550e8400-e29b-41d4-a716-446655440000",
      "name": "AWS Solutions Architect",
      "issuer": "Amazon Web Services",
      "credentialId": "ABC123XYZ",
      "credentialUrl": "https://aws.amazon.com/verify/ABC123XYZ",
      "skills": ["AWS", "Cloud Architecture", "DevOps"]
    }
  ]
}
```

---

### Send Contact Message

```http
POST /api/v1/profiles/{profileId}/contacts
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| profileId | UUID string | Yes | The profile identifier |

**Request Headers**:
```
Content-Type: application/json
Accept: application/json
```

**Request Body**:
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "message": "Hello! I'm interested in working together on a project."
}
```

**Validation Rules**:
| Field | Type | Required | Min | Max |
|-------|------|----------|-----|-----|
| name | string | Yes | 2 chars | 100 chars |
| email | string | Yes | Valid email format | - |
| message | string | Yes | 10 chars | 1000 chars |

**Response** `201 Created`:
```json
{
  "id": "contact-001",
  "message": "Contact message sent successfully",
  "contactedAt": "2024-12-20T15:30:00Z"
}
```

**Response** `400 Bad Request` (validation error):
```json
{
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": "Key: 'Request.Email' Error:Field validation for 'Email' failed on the 'email' tag"
  }
}
```

**Response** `429 Too Many Requests`:
```json
{
  "error": {
    "type": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests"
  }
}
```

---

### Submit a Question

```http
POST /api/v1/profiles/{profileId}/questions
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| profileId | UUID string | Yes | The profile identifier |

**Request Headers**:
```
Content-Type: application/json
Accept: application/json
```

**Request Body**:
```json
{
  "message": "What technologies do you recommend for building microservices?"
}
```

**Validation Rules**:
| Field | Type | Required | Min | Max |
|-------|------|----------|-----|-----|
| message | string | Yes | 5 chars | 500 chars |

**Note**: The client's IP address is automatically captured and stored with the question.

**Response** `201 Created`:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Question received successfully",
  "createdAt": "2024-12-20T15:30:00Z"
}
```

**Response** `400 Bad Request` (validation error):
```json
{
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": "Key: 'Request.Message' Error:Field validation for 'Message' failed on the 'min' tag"
  }
}
```

**Response** `429 Too Many Requests`:
```json
{
  "error": {
    "type": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests"
  }
}
```

---

## ❌ Error Response Format

All errors follow this structure:

```json
{
  "error": {
    "type": "ERROR_TYPE",
    "message": "Human readable message",
    "details": "Optional additional details"
  }
}
```

**Error Types**:
| Type | HTTP Code | Description |
|------|-----------|-------------|
| `BAD_REQUEST` | 400 | Invalid request format or missing required fields |
| `VALIDATION_ERROR` | 400 | Request body failed validation |
| `NOT_FOUND` | 404 | Resource not found |
| `PAYLOAD_TOO_LARGE` | 413 | Request body exceeds 1 MB |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## 💻 Code Examples

### TypeScript/JavaScript (fetch)

```typescript
const API_URL = 'https://your-api-domain.com';
const PROFILE_ID = 'your-profile-uuid';

// Helper function for API calls
async function apiCall<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'omit', // IMPORTANT: No credentials allowed
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error?.message || 'API Error');
  }

  return response.json();
}

// Get profile
const profile = await apiCall<Profile>(`/api/v1/profiles/${PROFILE_ID}`);

// Get all portfolio data in parallel
const [profile, skills, projects, certificates] = await Promise.all([
  apiCall<Profile>(`/api/v1/profiles/${PROFILE_ID}`),
  apiCall<SkillsResponse>(`/api/v1/profiles/${PROFILE_ID}/skills`),
  apiCall<ProjectsResponse>(`/api/v1/profiles/${PROFILE_ID}/projects`),
  apiCall<CertificatesResponse>(`/api/v1/profiles/${PROFILE_ID}/certificates`),
]);

// Send contact message
const contactResponse = await apiCall<ContactResponse>(
  `/api/v1/profiles/${PROFILE_ID}/contacts`,
  {
    method: 'POST',
    body: JSON.stringify({
      name: 'Jane Smith',
      email: 'jane@example.com',
      message: 'Hello! I would like to discuss a project opportunity.',
    }),
  }
);
```

### Python (requests)

```python
import requests

API_URL = "https://your-api-domain.com"
PROFILE_ID = "your-profile-uuid"

headers = {
    "Accept": "application/json",
    "Content-Type": "application/json"
}

# Get profile
response = requests.get(
    f"{API_URL}/api/v1/profiles/{PROFILE_ID}",
    headers=headers
)
profile = response.json()

# Get skills
response = requests.get(
    f"{API_URL}/api/v1/profiles/{PROFILE_ID}/skills",
    headers=headers
)
skills = response.json()["skills"]

# Send contact
response = requests.post(
    f"{API_URL}/api/v1/profiles/{PROFILE_ID}/contacts",
    headers=headers,
    json={
        "name": "Jane Smith",
        "email": "jane@example.com",
        "message": "Hello! I would like to discuss a project."
    }
)
contact = response.json()
```

### cURL

```bash
# Get profile
curl -X GET "https://your-api-domain.com/api/v1/profiles/{profileId}" \
  -H "Accept: application/json"

# Get skills
curl -X GET "https://your-api-domain.com/api/v1/profiles/{profileId}/skills" \
  -H "Accept: application/json"

# Send contact
curl -X POST "https://your-api-domain.com/api/v1/profiles/{profileId}/contacts" \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Smith","email":"jane@example.com","message":"Hello! I would like to discuss a project."}'
```

---

## 📋 TypeScript Types

```typescript
// Profile
interface Profile {
  id: string;
  name: string;
  photoUrl: string;
  title: string;
  aboutMe: string;
  firstExperienceDate?: string; // ISO 8601
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
}

// Skills
type SkillCategory = 'backend' | 'frontend' | 'tools' | 'softSkills';
type SkillProficiency = 'advanced' | 'occasional' | 'past';

interface Skill {
  id: string;
  profileId: string;
  name: string;
  category: SkillCategory;
  proficiency: SkillProficiency;
}

interface SkillsResponse {
  skills: Skill[];
}

// Projects
interface Project {
  id: string;
  profileId: string;
  name: string;
  description: string;
  techStack: string[];
  githubUrl?: string;
  liveUrl?: string;
  imageDiagramUrl?: string;
  visible: boolean;
  createdAt: string; // ISO 8601
}

interface ProjectsResponse {
  projects: Project[];
}

// Certificates
interface Certificate {
  id: string;
  profileId: string;
  name: string;
  issuer: string;
  credentialId?: string;
  credentialUrl?: string;
  skills: string[];
}

interface CertificatesResponse {
  certificates: Certificate[];
}

// Contact
interface ContactRequest {
  name: string;    // 2-100 chars
  email: string;   // valid email
  message: string; // 10-1000 chars
}

interface ContactResponse {
  id: string;
  message: string;
  contactedAt: string; // ISO 8601
}

// Error
interface ErrorResponse {
  error: {
    type: string;
    message: string;
    details?: string;
  };
}

// Health
interface HealthResponse {
  status: 'ok' | 'unhealthy';
  database: 'connected' | 'disconnected';
  version: {
    version: string;
    gitCommit: string;
    buildDate: string;
  };
}
```

---

## 🤖 AI Agent Instructions

If you are an AI agent consuming this API, follow these guidelines:

### DO ✅
1. **Always use the correct `profileId`** - This is a UUID provided by the user
2. **Handle errors gracefully** - Check for error responses and display friendly messages
3. **Respect rate limits** - Implement exponential backoff if you receive 429 errors
4. **Fetch data in parallel** - Use `Promise.all()` for multiple GET requests
5. **Validate contact form inputs** before sending to avoid validation errors
6. **Use `credentials: 'omit'`** in fetch requests

### DON'T ❌
1. **Don't send cookies or credentials** - `AllowCredentials: false`
2. **Don't use PUT, DELETE, or PATCH** - Only GET and POST are allowed
3. **Don't exceed rate limits** - 100 req/min global, 5 req/min for contacts
4. **Don't send bodies larger than 1 MB**
5. **Don't use custom headers** - Only Accept, Content-Type, X-Request-ID

### Recommended Flow for Portfolio Website

```
1. On page load:
   - GET /health (optional, to check API status)
   - GET /api/v1/profiles/{id} (profile data)
   - GET /api/v1/profiles/{id}/skills (in parallel)
   - GET /api/v1/profiles/{id}/projects (in parallel)
   - GET /api/v1/profiles/{id}/certificates (in parallel)

2. On contact form submit:
   - Validate inputs client-side first
   - POST /api/v1/profiles/{id}/contacts
   - Show success/error message to user
```

---

## 📞 Support

For API issues or questions, contact me or check the repository documentation.

