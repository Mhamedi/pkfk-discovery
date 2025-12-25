# API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "admin123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "...",
  "expires_in": 3600
}
```

## Endpoints

### Adapters

- `GET /adapters` - List adapters
- `POST /adapters` - Create adapter
- `GET /adapters/{id}` - Get adapter
- `PUT /adapters/{id}` - Update adapter
- `DELETE /adapters/{id}` - Delete adapter
- `POST /adapters/{id}/publish` - Publish adapter

### Adapter Studio

- `POST /studio/drafts` - Create draft
- `GET /studio/drafts/{id}` - Get draft
- `PUT /studio/drafts/{id}` - Update draft
- `POST /studio/drafts/{id}/probe` - Probe connection
- `POST /studio/drafts/{id}/validate` - Validate adapter
- `POST /studio/drafts/{id}/optimize` - Optimize with AI
- `POST /studio/drafts/{id}/publish` - Publish draft

### Connections

- `GET /connections` - List connections
- `POST /connections` - Create connection
- `GET /connections/{id}` - Get connection
- `PUT /connections/{id}` - Update connection
- `DELETE /connections/{id}` - Delete connection
- `POST /connections/{id}/test` - Test connection

### Scans

- `GET /scans` - List scans
- `POST /scans` - Create scan
- `GET /scans/{id}` - Get scan
- `GET /scans/{id}/results` - Get scan results
- `DELETE /scans/{id}` - Cancel scan

### AI Providers

- `GET /ai/providers` - List providers
- `POST /ai/providers` - Create provider
- `GET /ai/providers/{id}` - Get provider
- `PUT /ai/providers/{id}` - Update provider
- `DELETE /ai/providers/{id}` - Delete provider

### Administration

- `GET /admin/users` - List users
- `POST /admin/users` - Create user
- `GET /admin/users/{id}` - Get user
- `PUT /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user

- `GET /admin/audit` - List audit logs
- `GET /admin/audit/{id}` - Get audit log
- `GET /admin/audit/export` - Export audit logs

### Settings

- `GET /settings` - Get settings
- `PUT /settings` - Update settings

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

Common status codes:
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Rate Limiting

Rate limits are applied per user/IP:
- Default: 100 requests per minute
- Burst: 20 requests per second

## OpenAPI Specification

The complete OpenAPI 3.1 specification is generated from code annotations and available at:

```
GET /api/v1/openapi.json
```

