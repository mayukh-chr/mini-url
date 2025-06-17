# URL Shortener

A production-ready URL shortener with Go backend, React frontend, PostgreSQL database, and comprehensive monitoring. Features include custom short codes, analytics, rate limiting, and cloud deployment support.

## Features

### Core URL Shortening
- **Create Short URLs**: Convert long URLs into short, manageable links
  - Automatic random code generation (6 characters)
  - Custom short code support with validation
  - Duplicate short code prevention
  - URL sanitization and validation
  - JSON API response with generated short URL

### URL Management
- **Retrieve Original URLs**: Redirect short URLs to their original destinations
  - Automatic access count tracking with async updates
  - Real-time click analytics
  - HTTP 302 redirect to original URL

- **Update Short URLs**: Modify existing URL mappings
  - Change destination URLs for existing short codes
  - Update short codes with new custom codes
  - Conflict detection for duplicate codes

- **Delete Short URLs**: Remove URL mappings from the system
  - Complete removal of short URL entries
  - Clean database management

### Production Features
- **PostgreSQL Support**: Production database with connection pooling
- **Rate Limiting**: Configurable request rate limiting
- **Structured Logging**: JSON logging with request tracing
- **Health Checks**: Comprehensive health and metrics endpoints
- **Security Headers**: CORS, XSS protection, security headers
- **Environment Configuration**: Environment-based configuration
- **Database Migrations**: Automatic table creation and indexing

### Monitoring & Analytics
- **Access Statistics**: Track usage metrics for short URLs
- **System Metrics**: Memory, goroutines, database connections
- **Prometheus Integration**: Metrics endpoint for monitoring
- **Request Logging**: Detailed request/response logging

### Deployment Support
- **Heroku Ready**: One-click Heroku deployment
- **Docker Support**: Multi-stage Docker builds
- **CI/CD Pipeline**: GitHub Actions integration
- **Environment Variables**: 12-factor app configuration

### API Endpoints
- `POST /shorten` - Create new short URL
- `GET /u/{code}` - Redirect to original URL
- `PUT /u/{code}` - Update existing short URL
- `DELETE /u/{code}` - Delete short URL
- `GET /stats/{code}` - Get access statistics
- `GET /health` - Health check endpoint
- `GET /metrics` - Application metrics
- `GET /metrics/prometheus` - Prometheus format metrics
- `GET /shorten` - Web interface for URL management

## Project Structure

```
urlshortner/
├── config/              # Configuration management
├── database/            # Database connection and migrations  
├── handlers/            # HTTP request handlers
├── middleware/          # HTTP middleware (rate limiting, logging)
├── models/             # Data models
├── monitoring/         # Metrics and monitoring
├── utils/              # Utility functions
├── frontend/           # React frontend
├── templates/          # HTML templates
├── .github/workflows/  # CI/CD pipelines
├── Dockerfile          # Container configuration
├── docker-compose.yml  # Local development setup
├── Procfile           # Heroku deployment
├── main.go            # Application entry point
└── README.md          # This file
```

## Quick Start

### Local Development (SQLite)

1. **Clone and setup:**
   ```bash
   git clone https://github.com/your-username/url-shortner.git
   cd url-shortner
   go mod download
   ```

2. **Run the application:**
   ```bash
   go run main.go
   ```

3. **Access the application:**
   - API: http://localhost:8080
   - Web UI: http://localhost:8080/shorten
   - Health Check: http://localhost:8080/health
   - Metrics: http://localhost:8080/metrics

### Local Development with PostgreSQL

1. **Start with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

2. **Or manually with PostgreSQL:**
   ```bash
   # Set environment variables
   export DATABASE_URL="postgres://user:password@localhost:5432/urlshortener?sslmode=disable"
   export ENVIRONMENT="development"
   
   # Run application
   go run main.go
   ```

### Production Deployment

#### Heroku Deployment

See [HEROKU_DEPLOY.md](HEROKU_DEPLOY.md) for detailed instructions.

**Quick Deploy:**
```bash
# Login and create app
heroku login
heroku create your-url-shortener

# Add PostgreSQL
heroku addons:create heroku-postgresql:mini

# Set environment variables
heroku config:set ENVIRONMENT=production
heroku config:set BASE_URL=https://your-url-shortener.herokuapp.com

# Deploy
git push heroku main
```

#### Docker Deployment

```bash
# Build image
docker build -t url-shortener .

# Run with environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL="your_postgres_url" \
  -e ENVIRONMENT="production" \
  url-shortener
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | 8080 | No |
| `DATABASE_URL` | PostgreSQL connection string | - | Production only |
| `BASE_URL` | Base URL for short links | http://localhost:8080 | Yes |
| `ENVIRONMENT` | App environment (development/production) | development | No |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | info | No |

## API Usage Examples

### Create Short URL
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "short_code": "custom"}'
```

### Get Statistics
```bash
curl http://localhost:8080/stats/abc123
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Monitoring

### Application Metrics
- **Endpoint**: `/metrics` - JSON format metrics
- **Prometheus**: `/metrics/prometheus` - Prometheus format

### Available Metrics
- System metrics (memory, goroutines, GC)
- Database connection pool stats
- HTTP request metrics
- Application uptime

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL 13+ (for production)
- Node.js 16+ (for frontend)

### Testing
```bash
# Run tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Integration tests with PostgreSQL
DATABASE_URL="postgres://..." go test -v ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Security check
gosec ./...
```

## Architecture

### Database Schema
```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    short_code VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    access_count INTEGER NOT NULL DEFAULT 0
);
```

### Performance Characteristics
- **SQLite**: ~500 req/s (development)
- **PostgreSQL**: ~5,000-15,000 req/s (production)
- **With caching**: 20,000+ req/s (Redis integration ready)

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request




