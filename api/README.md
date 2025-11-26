# SlimJSON HTTP API Documentation

REST API documentation for SlimJSON daemon mode.

## Overview

SlimJSON can run as an HTTP daemon to provide JSON compression as a service. This is useful for:

- Microservices architecture
- API gateway integration
- CI/CD pipeline processing
- Real-time data compression

## Starting the Daemon

```bash
# Default port 8080
slimjson -d

# Custom port
slimjson -d -port 3000

# With custom config file
slimjson -d -c /path/to/.slimjson
```

## API Endpoints

### Health Check

Check if the service is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok",
  "version": "1.0"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

### List Profiles

Get all available compression profiles (built-in and custom).

**Endpoint:** `GET /profiles`

**Response:**
```json
{
  "builtin": [
    "light",
    "medium",
    "aggressive",
    "ai-optimized"
  ],
  "custom": [
    "my-custom-profile",
    "api-response"
  ]
}
```

**Example:**
```bash
curl http://localhost:8080/profiles
```

---

### Compress JSON

Compress JSON data using a specified profile or default settings.

**Endpoint:** `POST /slim`

**Query Parameters:**
- `profile` (optional): Profile name to use for compression

**Request Headers:**
- `Content-Type: application/json`

**Request Body:**
Any valid JSON object or array.

**Response:**
Compressed JSON object.

**Examples:**

#### Default Compression

```bash
curl -X POST http://localhost:8080/slim \
  -H "Content-Type: application/json" \
  -d '{
    "users": [
      {"id": 1, "name": "Alice", "email": "alice@example.com"},
      {"id": 2, "name": "Bob", "email": "bob@example.com"}
    ],
    "prices": [19.999, 29.123, 39.456]
  }'
```

Response:
```json
{
  "prices": [20, 29, 39],
  "users": [
    {"id": 1, "name": "Alice", "email": "alice@example.com"},
    {"id": 2, "name": "Bob", "email": "bob@example.com"}
  ]
}
```

#### With Profile

```bash
curl -X POST 'http://localhost:8080/slim?profile=medium' \
  -H "Content-Type: application/json" \
  -d @data.json
```

#### With Custom Profile

```bash
curl -X POST 'http://localhost:8080/slim?profile=my-custom-profile' \
  -H "Content-Type: application/json" \
  -d @data.json
```

## Built-in Profiles

### Light
- **MaxDepth:** 10
- **MaxListLength:** 20
- **Use Case:** Preserve most data, only limit depth and arrays
- **Reduction:** ~20-30%

### Medium
- **MaxDepth:** 5
- **MaxListLength:** 10
- **Use Case:** Balanced compression for general use
- **Reduction:** ~40-60%

### Aggressive
- **MaxDepth:** 3
- **MaxListLength:** 5
- **BlockList:** description, summary, comment, notes, bio, readme
- **Use Case:** Maximum reduction, remove verbose fields
- **Reduction:** ~60-80%

### AI-Optimized
- **MaxDepth:** 4
- **MaxListLength:** 8
- **BlockList:** avatar_url, gravatar_id, url, html_url, *_url
- **Use Case:** Optimized for LLM contexts, remove URLs
- **Reduction:** ~50-70%

## Error Responses

### 400 Bad Request

Invalid JSON or unknown profile.

```json
"Invalid JSON: unexpected end of JSON input"
```

or

```json
"Unknown profile: nonexistent"
```

### 405 Method Not Allowed

Only POST method is supported for `/slim` endpoint.

```json
"Method not allowed"
```

## OpenAPI Specification

Full OpenAPI 3.0 specification is available in [`swagger.yaml`](swagger.yaml).

### Viewing the Spec

You can view the API specification using:

- **Swagger UI:** https://editor.swagger.io/ (paste the swagger.yaml content)
- **Redoc:** https://redocly.github.io/redoc/ (paste the swagger.yaml URL)
- **VS Code:** Install "OpenAPI (Swagger) Editor" extension

## Integration Examples

### cURL

```bash
# Health check
curl http://localhost:8080/health

# List profiles
curl http://localhost:8080/profiles

# Compress with medium profile
curl -X POST 'http://localhost:8080/slim?profile=medium' \
  -H "Content-Type: application/json" \
  -d '{"data": "value"}'
```

### Python

```python
import requests

# Compress JSON
response = requests.post(
    'http://localhost:8080/slim?profile=medium',
    json={'users': [{'id': 1, 'name': 'Alice'}]}
)

compressed = response.json()
print(compressed)
```

### JavaScript/Node.js

```javascript
const fetch = require('node-fetch');

async function compressJSON(data, profile = 'medium') {
  const response = await fetch(
    `http://localhost:8080/slim?profile=${profile}`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    }
  );
  
  return await response.json();
}

// Usage
const data = { users: [{ id: 1, name: 'Alice' }] };
const compressed = await compressJSON(data, 'medium');
console.log(compressed);
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func compressJSON(data interface{}, profile string) (interface{}, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    url := fmt.Sprintf("http://localhost:8080/slim?profile=%s", profile)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return result, nil
}
```

## Docker Deployment

Run SlimJSON daemon in Docker:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o slimjson ./cmd/slimjson

FROM alpine:latest
COPY --from=builder /app/slimjson /usr/local/bin/
EXPOSE 8080
CMD ["slimjson", "-d", "-port", "8080"]
```

Build and run:

```bash
docker build -t slimjson .
docker run -p 8080:8080 slimjson
```

## Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: slimjson
spec:
  replicas: 3
  selector:
    matchLabels:
      app: slimjson
  template:
    metadata:
      labels:
        app: slimjson
    spec:
      containers:
      - name: slimjson
        image: slimjson:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: slimjson
spec:
  selector:
    app: slimjson
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Monitoring

### Health Check

The `/health` endpoint can be used for:
- Kubernetes liveness/readiness probes
- Load balancer health checks
- Monitoring systems (Prometheus, etc.)

### Metrics

For production deployments, consider adding:
- Request/response logging
- Prometheus metrics endpoint
- Distributed tracing (OpenTelemetry)

## Security Considerations

1. **Rate Limiting:** Implement rate limiting to prevent abuse
2. **Authentication:** Add API key or JWT authentication for production
3. **HTTPS:** Use TLS/SSL in production environments
4. **Input Validation:** The API validates JSON but consider additional validation
5. **Resource Limits:** Set appropriate memory and CPU limits

## Performance

- **Throughput:** ~1000-5000 requests/second (depending on JSON size)
- **Latency:** ~1-10ms per request (depending on JSON complexity)
- **Memory:** ~10-50MB per instance
- **CPU:** Minimal, scales linearly with request rate

## Troubleshooting

### Service won't start

```bash
# Check if port is already in use
lsof -i :8080

# Try different port
slimjson -d -port 3000
```

### Profile not found

```bash
# List available profiles
curl http://localhost:8080/profiles

# Check config file
cat .slimjson
```

### Invalid JSON errors

Ensure request has proper Content-Type header:
```bash
curl -X POST http://localhost:8080/slim \
  -H "Content-Type: application/json" \
  -d '{"valid": "json"}'
```

## Support

- **GitHub Issues:** https://github.com/tradik/slimjson/issues
- **Documentation:** https://github.com/tradik/slimjson
- **Examples:** https://github.com/tradik/slimjson/blob/main/LIBRARY_EXAMPLES.md
