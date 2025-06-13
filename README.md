# Samudai Package

Common utility packages for Samudai Services.

## Packages

### Database (`db`)

Package for database connections and operations.

**Environment Variables:**

- `DATABASE_URL`: PostgreSQL database connection URL
- `MONGO_URL`: MongoDB connection URL
- `REDIS_URL`: Redis connection URL

### Logger (`logger`)

Structured logging package for consistent log formatting across services.

**Environment Variables:**

- `SERVICE_NAME`: Name of the service using the logger

### File Upload (`fileupload`)

Digital Ocean Spaces file upload integration.

**Environment Variables:**

- `SPACES_KEY`: Digital Ocean Spaces access key
- `SPACES_SECRET`: Digital Ocean Spaces secret key
- `ENDPOINT`: Digital Ocean Spaces endpoint URL
- `BUCKET_NAME`: Name of the storage bucket

### Requester (`requester`)

HTTP client package for making external API requests.

## Installation

```bash
go get github.com/Samudai/samudai-pkg
```

## Usage

Import the required packages in your Go code:

```go
import (
    "github.com/Samudai/samudai-pkg/db"
    "github.com/Samudai/samudai-pkg/logger"
    "github.com/Samudai/samudai-pkg/fileupload"
    "github.com/Samudai/samudai-pkg/requester"
)
```

## License

[MIT License](LICENSE)
