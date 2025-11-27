# Docker Deployment Guide

Complete guide for running SlimJSON in Docker containers.

## Quick Start

### Pull and Run

```bash
# Pull latest image
docker pull ghcr.io/tradik/slimjson:latest

# Run daemon on port 8080
docker run -d -p 8080:8080 --name slimjson ghcr.io/tradik/slimjson:latest

# Test health
curl http://localhost:8080/health
```

### Build Locally

```bash
# Build image
docker build -t slimjson:local .

# Run
docker run -d -p 8080:8080 --name slimjson slimjson:local
```

## Usage Modes

### 1. Daemon Mode (Default)

Run as HTTP API service:

```bash
# Default port 8080
docker run -d -p 8080:8080 ghcr.io/tradik/slimjson:latest

# Custom port
docker run -d -p 3000:3000 ghcr.io/tradik/slimjson:latest -d -port 3000

# With custom config file
docker run -d -p 8080:8080 \
  -v $(pwd)/.slimjson:/app/.slimjson:ro \
  ghcr.io/tradik/slimjson:latest -d -c /app/.slimjson
```

### 2. CLI Mode

Process JSON files:

```bash
# Process file
docker run --rm -i ghcr.io/tradik/slimjson:latest \
  -profile medium < input.json > output.json

# With volume mount
docker run --rm \
  -v $(pwd)/data:/data \
  ghcr.io/tradik/slimjson:latest \
  -profile medium /data/input.json > /data/output.json

# Pretty print
docker run --rm -i ghcr.io/tradik/slimjson:latest \
  -profile medium -pretty < input.json
```

## Docker Compose

### Basic Setup

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  slimjson:
    image: ghcr.io/tradik/slimjson:latest
    container_name: slimjson
    ports:
      - "8080:8080"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

Run:

```bash
docker-compose up -d
```

### With Custom Config

```yaml
version: '3.8'

services:
  slimjson:
    image: ghcr.io/tradik/slimjson:latest
    container_name: slimjson
    ports:
      - "8080:8080"
    volumes:
      - ./config/.slimjson:/app/.slimjson:ro
    command: ["-d", "-c", "/app/.slimjson", "-port", "8080"]
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

### Production Setup

```yaml
version: '3.8'

services:
  slimjson:
    image: ghcr.io/tradik/slimjson:latest
    container_name: slimjson
    ports:
      - "8080:8080"
    volumes:
      - ./config/.slimjson:/app/.slimjson:ro
    environment:
      - TZ=UTC
    command: ["-d", "-c", "/app/.slimjson", "-port", "8080"]
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 256M
        reservations:
          cpus: '0.5'
          memory: 128M
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Kubernetes Deployment

### Basic Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: slimjson
  labels:
    app: slimjson
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
        image: ghcr.io/tradik/slimjson:latest
        ports:
        - containerPort: 8080
          name: http
        args: ["-d", "-port", "8080"]
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
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
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
    protocol: TCP
  type: LoadBalancer
```

### With ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: slimjson-config
data:
  .slimjson: |
    [production]
    depth=5
    list-len=20
    strip-empty=true
    decimal-places=2
    deduplicate=true
    
    [staging]
    depth=10
    list-len=50
    strip-empty=true
---
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
        image: ghcr.io/tradik/slimjson:latest
        ports:
        - containerPort: 8080
        args: ["-d", "-c", "/config/.slimjson", "-port", "8080"]
        volumeMounts:
        - name: config
          mountPath: /config
          readOnly: true
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
      volumes:
      - name: config
        configMap:
          name: slimjson-config
```

## Environment Variables

Currently, SlimJSON is configured via command-line flags. For Docker, pass them via `CMD`:

```bash
docker run -d -p 8080:8080 \
  ghcr.io/tradik/slimjson:latest \
  -d -port 8080 -c /app/.slimjson
```

## Health Checks

### Docker Health Check

Built-in health check runs every 30 seconds:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

Check status:

```bash
docker ps
# Look for "healthy" status

docker inspect --format='{{.State.Health.Status}}' slimjson
```

### Manual Health Check

```bash
# Inside container
docker exec slimjson wget -qO- http://localhost:8080/health

# From host
curl http://localhost:8080/health
```

## Logging

### View Logs

```bash
# Follow logs
docker logs -f slimjson

# Last 100 lines
docker logs --tail 100 slimjson

# With timestamps
docker logs -t slimjson
```

### Log Configuration

In docker-compose.yml:

```yaml
services:
  slimjson:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Networking

### Bridge Network

```bash
# Create network
docker network create slimjson-net

# Run with network
docker run -d \
  --network slimjson-net \
  --name slimjson \
  -p 8080:8080 \
  ghcr.io/tradik/slimjson:latest
```

### Host Network

```bash
# Use host network (Linux only)
docker run -d \
  --network host \
  --name slimjson \
  ghcr.io/tradik/slimjson:latest \
  -d -port 8080
```

## Volumes

### Mount Config File

```bash
docker run -d -p 8080:8080 \
  -v $(pwd)/.slimjson:/app/.slimjson:ro \
  ghcr.io/tradik/slimjson:latest \
  -d -c /app/.slimjson
```

### Process Files

```bash
docker run --rm \
  -v $(pwd)/data:/data \
  ghcr.io/tradik/slimjson:latest \
  -profile medium /data/input.json > /data/output.json
```

## Security

### Run as Non-Root

Container runs as user `slimjson` (UID 1000) by default:

```dockerfile
USER slimjson
```

### Read-Only Root Filesystem

```bash
docker run -d -p 8080:8080 \
  --read-only \
  --tmpfs /tmp \
  ghcr.io/tradik/slimjson:latest
```

### Drop Capabilities

```bash
docker run -d -p 8080:8080 \
  --cap-drop=ALL \
  ghcr.io/tradik/slimjson:latest
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker logs slimjson

# Check if port is available
lsof -i :8080

# Run interactively
docker run -it --rm ghcr.io/tradik/slimjson:latest /bin/sh
```

### Health Check Failing

```bash
# Check health status
docker inspect --format='{{json .State.Health}}' slimjson | jq

# Test health endpoint manually
docker exec slimjson wget -qO- http://localhost:8080/health
```

### Permission Issues

```bash
# Check user
docker exec slimjson id

# Check file permissions
docker exec slimjson ls -la /app
```

## Performance Tuning

### Resource Limits

```bash
docker run -d -p 8080:8080 \
  --memory="256m" \
  --cpus="1" \
  ghcr.io/tradik/slimjson:latest
```

### Multiple Instances

```bash
# Run 3 instances with different ports
for i in {1..3}; do
  docker run -d \
    --name slimjson-$i \
    -p $((8080+i)):8080 \
    ghcr.io/tradik/slimjson:latest
done
```

## Examples

### API Gateway Integration

```bash
# Run behind nginx
docker network create api-net

docker run -d \
  --network api-net \
  --name slimjson \
  ghcr.io/tradik/slimjson:latest

docker run -d \
  --network api-net \
  --name nginx \
  -p 80:80 \
  -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro \
  nginx:alpine
```

### CI/CD Pipeline

```bash
# In CI/CD script
docker run --rm -i ghcr.io/tradik/slimjson:latest \
  -profile aggressive < large-api-response.json > optimized.json
```

### Development

```bash
# Mount source code for development
docker run -it --rm \
  -v $(pwd):/app \
  -w /app \
  golang:1.25-alpine \
  sh -c "go build -o slimjson ./cmd/slimjson && ./slimjson -d"
```

## Maintenance

### Update Image

```bash
# Pull latest
docker pull ghcr.io/tradik/slimjson:latest

# Recreate container
docker stop slimjson
docker rm slimjson
docker run -d -p 8080:8080 --name slimjson ghcr.io/tradik/slimjson:latest
```

### Cleanup

```bash
# Remove container
docker stop slimjson
docker rm slimjson

# Remove image
docker rmi ghcr.io/tradik/slimjson:latest

# Prune unused resources
docker system prune -a
```

## References

- [Dockerfile](../Dockerfile)
- [API Documentation](../api/README.md)
- [GitHub Container Registry](https://github.com/tradik/slimjson/pkgs/container/slimjson)
