# CleanShort - URL Shortener API

A REST API service for creating and managing custom short links with user authentication, built with Go, Fiber, and PostgreSQL.

## Features

- User registration and JWT-based authentication
- Create, read, update, and delete short links
- Public redirect functionality with click tracking
- Rate limiting for security
- Input validation and error handling
- PostgreSQL database with GORM
- Structured logging and health checks

## Quick Start

### Prerequisites

- Go 1.23.3 or later
- PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd cleanshort
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Update the `.env` file with your database credentials and configuration:
```env
APP_ENV=development
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

DB_DSN=postgres://<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>?sslmode=disable

JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

RATE_LIMIT_AUTH=5
RATE_LIMIT_REDIRECT=200

LOG_LEVEL=info
```

4. Install dependencies:
```bash
go mod tidy
```

5. Create the database:
```sql
CREATE DATABASE shortener;
```

6. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Checks

- `GET /healthz` - Liveness check
- `GET /readyz` - Readiness check (includes database connectivity)

### Authentication

All authentication endpoints are rate-limited to 5 requests per minute per IP.

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "StrongP@ssw0rd"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "StrongP@ssw0rd"
}
```

**Response (200 OK):**
```json
{
  "access_token": "jwt-token",
  "expires_in": 900,
  "refresh_token": "opaque-token",
  "refresh_expires_in": 604800
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "opaque-token"
}
```

**Response (200 OK):**
```json
{
  "access_token": "new-jwt-token",
  "expires_in": 900
}
```

#### Logout
```http
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refresh_token": "opaque-token"
}
```

**Response (204 No Content)**

### Links Management

All link endpoints require authentication via `Authorization: Bearer <access_token>` header.

#### Create Link
```http
POST /api/v1/links
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "target_url": "https://example.com/article/123",
  "short_code": "my-article-123",
  "title": "Favorite Article",
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "short_code": "my-article-123",
  "short_url": "http://localhost:8080/my-article-123",
  "target_url": "https://example.com/article/123",
  "title": "Favorite Article",
  "is_active": true,
  "click_count": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### List Links
```http
GET /api/v1/links?limit=20&offset=0&query=article&active=true
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `limit` (optional): Number of results (1-100, default: 20)
- `offset` (optional): Pagination offset (default: 0)
- `query` (optional): Search in short_code and title
- `active` (optional): Filter by active status (true/false)

**Response (200 OK):**
```json
{
  "links": [...],
  "total": 42,
  "limit": 20,
  "offset": 0
}
```

#### Get Link
```http
GET /api/v1/links/{id}
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "short_code": "my-article-123",
  "short_url": "http://localhost:8080/my-article-123",
  "target_url": "https://example.com/article/123",
  "title": "Favorite Article",
  "is_active": true,
  "click_count": 5,
  "last_clicked_at": "2024-01-01T12:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Link
```http
PATCH /api/v1/links/{id}
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Updated Title",
  "target_url": "https://example.com/new-url",
  "is_active": false
}
```

**Response (200 OK):** Updated link object

#### Delete Link
```http
DELETE /api/v1/links/{id}
Authorization: Bearer <access_token>
```

**Response (204 No Content)**

### Public Redirect

Rate-limited to 200 requests per minute per IP.

```http
GET /{shortCode}
```

**Response:**
- `302 Found` with `Location` header if link is active
- `404 Not Found` if link doesn't exist or is inactive

## Error Responses

All errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "request_id": "req_abc123"
  }
}
```

**Common Error Codes:**
- `VALIDATION_ERROR` - Invalid request data
- `UNAUTHORIZED` - Authentication required or invalid
- `FORBIDDEN` - Access denied
- `CONFLICT` - Resource already exists
- `LINK_NOT_FOUND` - Short link not found
- `TOO_MANY_REQUESTS` - Rate limit exceeded
- `INTERNAL_ERROR` - Server error

## Rate Limiting

The API implements rate limiting with the following limits:

- **Authentication endpoints**: 5 requests per minute per IP
- **Redirect endpoint**: 200 requests per minute per IP

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Request limit
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Reset timestamp

## Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `email` (Text, Unique)
- `password` (Text, Hashed)
- `created_at`, `updated_at` (Timestamps)

### Links Table
- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key)
- `short_code` (VARCHAR(32), Unique)
- `target_url` (Text)
- `title` (Text, Nullable)
- `is_active` (Boolean)
- `click_count` (BigInt)
- `last_clicked_at` (Timestamp, Nullable)
- `created_at`, `updated_at` (Timestamps)

### Refresh Tokens Table
- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key)
- `token_hash` (Text, Unique)
- `revoked` (Boolean)
- `expires_at` (Timestamp)
- `created_at` (Timestamp)

## Security Features

- Password hashing with bcrypt
- JWT tokens with configurable expiration
- Refresh token rotation
- Rate limiting
- Input validation and sanitization
- CORS configuration
- Reserved short code protection

## Development

### Project Structure
```
├── config/          # Configuration management
├── controllers/     # HTTP handlers
├── database/        # Database connection and migrations
├── middleware/      # Authentication and rate limiting
├── models/          # Data models and DTOs
├── routes/          # Route definitions
├── services/        # Business logic
├── utils/           # Utility functions
├── main.go          # Application entry point
├── go.mod           # Go module definition
└── .env.example     # Environment configuration template
```

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o cleanshort main.go
```

## License

This project is licensed under the MIT License.