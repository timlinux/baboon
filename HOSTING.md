# Hosting Baboon

This guide covers all the ways you can host Baboon, from local development to production Kubernetes deployments.

## Table of Contents

- [Quick Start](#quick-start)
- [Hosting Options](#hosting-options)
  - [Local Development](#local-development)
  - [Docker](#docker)
  - [Docker Compose](#docker-compose)
  - [Kubernetes](#kubernetes)
  - [Cloud Platforms](#cloud-platforms)
- [Configuration](#configuration)
- [Session Affinity](#session-affinity)
- [Google AdSense](#google-adsense)
- [Reverse Proxy Setup](#reverse-proxy-setup)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

The fastest way to get Baboon running:

```bash
# Using pre-built binary
./baboon web -port 8787

# Using Docker
docker run -p 8787:8787 ghcr.io/timlinux/baboon:latest

# Using Nix
nix run github:timlinux/baboon -- web -port 8787
```

Then open http://localhost:8787 in your browser.

---

## Hosting Options

### Local Development

For development, run the backend and frontend separately for hot-reloading:

```bash
# Terminal 1: Start the backend
make server

# Terminal 2: Start the web frontend with hot-reload
make web-dev
```

Or run everything together:

```bash
# Build and serve the production bundle
make web-serve
```

### Docker

#### Basic Docker Run

```bash
docker run -d \
  --name baboon \
  -p 8787:8787 \
  ghcr.io/timlinux/baboon:latest
```

#### With AdSense

```bash
docker run -d \
  --name baboon \
  -p 8787:8787 \
  ghcr.io/timlinux/baboon:latest \
  web -port 8787 -adsense ca-pub-YOUR_PUBLISHER_ID
```

#### Building Your Own Image

```bash
# Build the image
docker build -t baboon:local .

# Run it
docker run -p 8787:8787 baboon:local
```

### Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  baboon:
    image: ghcr.io/timlinux/baboon:latest
    ports:
      - "8787:8787"
    command: web -port 8787 -adsense ${ADSENSE_KEY:-}
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8787/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  # Optional: nginx reverse proxy with SSL
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    depends_on:
      - baboon
```

Run with:

```bash
# Without AdSense
docker-compose up -d

# With AdSense
ADSENSE_KEY=ca-pub-YOUR_ID docker-compose up -d
```

### Kubernetes

Baboon includes production-ready Kubernetes manifests with Kustomize.

#### Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured
- Optional: Ingress controller (NGINX, Traefik, or AWS ALB)

#### Development Deployment

```bash
# Deploy to development
kubectl apply -k k8s/overlays/development

# Check status
kubectl -n baboon get pods
kubectl -n baboon get svc
```

#### Production Deployment

1. **Configure your domain** in `k8s/overlays/production/ingress.yaml`:

```yaml
spec:
  rules:
    - host: baboon.yourdomain.com  # Change this
```

2. **Deploy**:

```bash
kubectl apply -k k8s/overlays/production
```

3. **Verify**:

```bash
kubectl -n baboon get all
kubectl -n baboon get ingress
```

#### Important: Session Affinity

Baboon stores game sessions in memory, so users must connect to the same pod throughout their session. The Kubernetes manifests include:

- **Service**: `sessionAffinity: ClientIP` with 3-hour timeout
- **Ingress**: Cookie-based sticky sessions (`BABOON_AFFINITY` cookie)

See [Session Affinity](#session-affinity) for details.

### Cloud Platforms

#### AWS (ECS/Fargate)

```bash
# Create task definition
aws ecs register-task-definition \
  --family baboon \
  --container-definitions '[{
    "name": "baboon",
    "image": "ghcr.io/timlinux/baboon:latest",
    "portMappings": [{"containerPort": 8787}],
    "healthCheck": {
      "command": ["CMD-SHELL", "wget -q --spider http://localhost:8787/api/health || exit 1"]
    }
  }]'

# Create service with ALB and sticky sessions
aws ecs create-service \
  --service-name baboon \
  --task-definition baboon \
  --load-balancers targetGroupArn=YOUR_TG_ARN,containerName=baboon,containerPort=8787
```

**ALB Target Group Settings**:
- Stickiness: Enabled
- Stickiness type: Application-based cookie
- Cookie name: `BABOON_AFFINITY`
- Stickiness duration: 3 hours

#### Google Cloud Run

```bash
gcloud run deploy baboon \
  --image ghcr.io/timlinux/baboon:latest \
  --port 8787 \
  --session-affinity \
  --allow-unauthenticated
```

Note: Cloud Run has limited session affinity support. For production, consider GKE.

#### DigitalOcean App Platform

Create `app.yaml`:

```yaml
name: baboon
services:
  - name: web
    image:
      registry_type: GHCR
      registry: timlinux
      repository: baboon
      tag: latest
    http_port: 8787
    instance_count: 1
    instance_size_slug: basic-xxs
    routes:
      - path: /
    health_check:
      http_path: /api/health
```

Deploy:

```bash
doctl apps create --spec app.yaml
```

#### Fly.io

Create `fly.toml`:

```toml
app = "baboon"
primary_region = "ams"

[build]
  image = "ghcr.io/timlinux/baboon:latest"

[http_service]
  internal_port = 8787
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true

[[services.http_checks]]
  interval = 10000
  timeout = 2000
  path = "/api/health"
```

Deploy:

```bash
fly launch
fly deploy
```

#### Railway

Connect your GitHub repo and Railway will auto-detect the Dockerfile. Set:

- Port: 8787
- Health check: `/api/health`

#### Render

Create a new Web Service:

- Environment: Docker
- Docker Command: `web -port 8787`
- Health Check Path: `/api/health`

---

## Configuration

### Command Line Flags

```bash
baboon web [options]

Options:
  -port int      Port for the web server (default 8787)
  -adsense str   Google AdSense publisher ID or "preview"
  -dir string    Directory containing built web frontend (default "web/dist")
  -p             Enable punctuation mode by default
```

### Environment Variables

While Baboon primarily uses command-line flags, you can wrap them in environment variables:

```bash
# In your entrypoint script or docker-compose
BABOON_PORT=${PORT:-8787}
BABOON_ADSENSE=${ADSENSE_KEY:-}

exec ./baboon web -port $BABOON_PORT -adsense "$BABOON_ADSENSE"
```

---

## Session Affinity

### Why It's Required

Baboon stores game sessions in memory on the server. When a user starts a typing session, they must continue hitting the same server instance. Without session affinity:

- Users lose their progress mid-game
- Statistics aren't recorded properly
- The experience breaks randomly

### Implementation by Platform

| Platform | Method | Configuration |
|----------|--------|---------------|
| Kubernetes (NGINX Ingress) | Cookie | `nginx.ingress.kubernetes.io/affinity: cookie` |
| Kubernetes (Traefik) | Cookie | `traefik.ingress.kubernetes.io/affinity: "true"` |
| AWS ALB | Cookie | Target group stickiness enabled |
| Google Cloud | Session affinity | `--session-affinity` flag |
| HAProxy | Cookie | `cookie SERVERID insert indirect nocache` |
| Nginx | ip_hash or sticky | `ip_hash;` or `sticky cookie` |

### Cookie Configuration

For cookie-based affinity:

- **Cookie name**: `BABOON_AFFINITY`
- **Duration**: 3 hours (10800 seconds)
- **Path**: `/`
- **Secure**: Yes (for HTTPS)
- **HttpOnly**: Yes

---

## Google AdSense

### Setup

1. **Get your publisher ID** from [Google AdSense](https://www.google.com/adsense/)

2. **Run with AdSense enabled**:
   ```bash
   baboon web -adsense ca-pub-YOUR_PUBLISHER_ID
   ```

3. **Preview ad placement** (without real ads):
   ```bash
   baboon web -adsense preview
   ```

### Ad Targeting

Baboon includes meta keywords for relevant ad targeting:

- Typing practice, typing tutor, typing speed
- Mechanical keyboards, ergonomic keyboards, split keyboards
- Keycaps, Cherry MX, Gateron, Kailh, Choc switches
- QMK, VIA, keyboard enthusiast terminology

### Ad Placement

Ads appear at the bottom of the typing screen, below the WPM bar. The ad slot is a 728x90 leaderboard format.

---

## Reverse Proxy Setup

### Nginx

```nginx
upstream baboon {
    ip_hash;  # Session affinity
    server 127.0.0.1:8787;
}

server {
    listen 80;
    server_name baboon.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name baboon.yourdomain.com;

    ssl_certificate /etc/ssl/certs/baboon.crt;
    ssl_certificate_key /etc/ssl/private/baboon.key;

    location / {
        proxy_pass http://baboon;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (if needed in future)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Health check endpoint
    location /api/health {
        proxy_pass http://baboon;
        access_log off;
    }
}
```

### Caddy

```caddyfile
baboon.yourdomain.com {
    reverse_proxy localhost:8787 {
        lb_policy ip_hash
        health_uri /api/health
        health_interval 30s
    }
}
```

### Traefik

```yaml
# traefik.yml
http:
  routers:
    baboon:
      rule: "Host(`baboon.yourdomain.com`)"
      service: baboon
      tls:
        certResolver: letsencrypt

  services:
    baboon:
      loadBalancer:
        sticky:
          cookie:
            name: BABOON_AFFINITY
            secure: true
            httpOnly: true
        servers:
          - url: "http://localhost:8787"
        healthCheck:
          path: /api/health
          interval: 30s
```

---

## Monitoring

### Health Check Endpoint

```bash
curl http://localhost:8787/api/health
```

Response:
```json
{
  "status": "healthy",
  "active_sessions": 5
}
```

### Prometheus Metrics

For production monitoring, consider adding a metrics endpoint. Example with middleware:

```go
// Add to your deployment
curl http://localhost:8787/metrics
```

### Recommended Alerts

- Health check failures
- High memory usage (sessions are in-memory)
- Response time > 500ms
- Error rate > 1%

---

## Troubleshooting

### Common Issues

#### "Session not found" errors

**Cause**: Session affinity not configured properly.

**Solution**: Ensure sticky sessions are enabled on your load balancer/ingress.

#### High memory usage

**Cause**: Many active sessions stored in memory.

**Solution**:
- Sessions are cleaned up automatically after inactivity
- Scale horizontally with proper session affinity
- Consider session timeout tuning

#### Ads not showing

**Cause**: AdSense not approved or blocked.

**Solution**:
- Verify your publisher ID is correct
- Check browser console for AdSense errors
- Ensure your domain is approved in AdSense
- Test with `-adsense preview` first

#### CORS errors

**Cause**: API and frontend on different origins.

**Solution**: Baboon serves both from the same origin. If you're splitting them, configure CORS headers.

### Debug Mode

For troubleshooting, check active sessions:

```bash
curl http://localhost:8787/api/sessions
```

### Logs

Baboon logs to stdout. In production:

```bash
# Docker
docker logs baboon

# Kubernetes
kubectl -n baboon logs -l app=baboon

# Systemd
journalctl -u baboon
```

---

## Support

- **Issues**: https://github.com/timlinux/baboon/issues
- **Discussions**: https://github.com/timlinux/baboon/discussions
- **Documentation**: https://timlinux.github.io/baboon/

---

Made with ♥ by [Kartoza](https://kartoza.com) | [Donate](https://github.com/sponsors/kartoza) | [GitHub](https://github.com/timlinux/baboon)
