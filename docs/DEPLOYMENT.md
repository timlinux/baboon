# Baboon Kubernetes Deployment Guide

This guide provides detailed instructions for deploying Baboon in a Kubernetes environment with load balancing and session affinity (sticky sessions).

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Prerequisites](#prerequisites)
3. [Docker Image](#docker-image)
4. [Kubernetes Manifests](#kubernetes-manifests)
5. [Session Affinity Configuration](#session-affinity-configuration)
6. [Ingress Configuration](#ingress-configuration)
7. [Persistent Storage](#persistent-storage)
8. [Health Checks](#health-checks)
9. [Scaling Considerations](#scaling-considerations)
10. [Monitoring](#monitoring)
11. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

```
                                    ┌─────────────────────────────────────────┐
                                    │           Kubernetes Cluster            │
                                    │                                         │
┌──────────┐    ┌──────────────┐   │  ┌─────────────────────────────────┐   │
│  Users   │───▶│   Ingress    │───┼─▶│         Service (ClusterIP)     │   │
│ (Browser)│    │  Controller  │   │  │    session affinity: ClientIP   │   │
└──────────┘    │  (nginx/     │   │  └───────────────┬─────────────────┘   │
                │   traefik)   │   │                  │                      │
                └──────────────┘   │    ┌─────────────┼─────────────┐        │
                      │            │    │             │             │        │
                      │            │    ▼             ▼             ▼        │
                 sticky sessions   │ ┌─────┐     ┌─────┐       ┌─────┐      │
                 (cookie-based)    │ │Pod 1│     │Pod 2│  ...  │Pod N│      │
                                   │ │:8787│     │:8787│       │:8787│      │
                                   │ └──┬──┘     └──┬──┘       └──┬──┘      │
                                   │    │           │             │          │
                                   │    └───────────┴─────────────┘          │
                                   │                │                        │
                                   │         ┌──────┴──────┐                 │
                                   │         │ PVC (stats) │                 │
                                   │         │  (optional) │                 │
                                   │         └─────────────┘                 │
                                   └─────────────────────────────────────────┘
```

### Why Sticky Sessions Are Required

Baboon maintains **in-memory game sessions** on each backend instance:

1. When a user starts a typing session, a session ID is created on a specific pod
2. All subsequent requests (keystrokes, space, backspace) must go to the **same pod**
3. If requests are load-balanced to different pods, the session will not be found (404 error)

**Session lifecycle:**
- `POST /api/sessions` → Creates session on Pod A, returns `session_id`
- `POST /api/sessions/{id}/keystroke` → Must reach Pod A
- `POST /api/sessions/{id}/space` → Must reach Pod A
- ...and so on until session ends

---

## Prerequisites

- Kubernetes cluster (1.19+)
- `kubectl` configured
- Ingress controller (nginx-ingress or Traefik recommended)
- Container registry access
- (Optional) Persistent volume provisioner for stats storage

---

## Docker Image

### Dockerfile

Create `Dockerfile` in the project root:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o baboon .

# Production stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/baboon .

# Copy web frontend (pre-built)
COPY --from=builder /app/web/dist ./web/dist

# Create non-root user
RUN adduser -D -u 1000 baboon
USER baboon

# Create config directory
RUN mkdir -p /home/baboon/.config/baboon

EXPOSE 8787

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8787/api/health || exit 1

ENTRYPOINT ["./baboon"]
CMD ["web", "-port", "8787"]
```

### Build and Push

```bash
# Build the web frontend first
cd web && npm ci && npm run build && cd ..

# Build Docker image
docker build -t your-registry/baboon:v1.0.0 .

# Push to registry
docker push your-registry/baboon:v1.0.0
```

---

## Kubernetes Manifests

### Namespace

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: baboon
  labels:
    app.kubernetes.io/name: baboon
```

### ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: baboon-config
  namespace: baboon
data:
  # Server configuration
  PORT: "8787"

  # Optional: Google AdSense publisher ID
  # ADSENSE_KEY: "ca-pub-xxxxxxxxxx"
```

### Secret (Optional - for AdSense)

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: baboon-secrets
  namespace: baboon
type: Opaque
stringData:
  adsense-key: "ca-pub-xxxxxxxxxx"
```

### Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: baboon
  namespace: baboon
  labels:
    app.kubernetes.io/name: baboon
    app.kubernetes.io/component: server
spec:
  replicas: 3
  selector:
    matchLabels:
      app.kubernetes.io/name: baboon
  template:
    metadata:
      labels:
        app.kubernetes.io/name: baboon
        app.kubernetes.io/component: server
    spec:
      containers:
        - name: baboon
          image: your-registry/baboon:v1.0.0
          imagePullPolicy: Always

          args:
            - "web"
            - "-port"
            - "8787"
            # Uncomment to enable AdSense:
            # - "-adsense"
            # - "$(ADSENSE_KEY)"

          ports:
            - name: http
              containerPort: 8787
              protocol: TCP

          env:
            - name: ADSENSE_KEY
              valueFrom:
                secretKeyRef:
                  name: baboon-secrets
                  key: adsense-key
                  optional: true

          resources:
            requests:
              cpu: "100m"
              memory: "64Mi"
            limits:
              cpu: "500m"
              memory: "256Mi"

          livenessProbe:
            httpGet:
              path: /api/health
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 3
            failureThreshold: 3

          readinessProbe:
            httpGet:
              path: /api/health
              port: http
            initialDelaySeconds: 3
            periodSeconds: 5
            timeoutSeconds: 2
            failureThreshold: 2

          securityContext:
            runAsNonRoot: true
            runAsUser: 1000
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false

          volumeMounts:
            - name: config
              mountPath: /home/baboon/.config/baboon
            - name: tmp
              mountPath: /tmp

      volumes:
        - name: config
          emptyDir: {}
        - name: tmp
          emptyDir: {}

      securityContext:
        fsGroup: 1000

      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app.kubernetes.io/name: baboon
                topologyKey: kubernetes.io/hostname
```

### Service with Session Affinity

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: baboon
  namespace: baboon
  labels:
    app.kubernetes.io/name: baboon
spec:
  type: ClusterIP
  selector:
    app.kubernetes.io/name: baboon
  ports:
    - name: http
      port: 80
      targetPort: http
      protocol: TCP

  # CRITICAL: Enable session affinity
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      # Session timeout: 3 hours (typing sessions can be long)
      timeoutSeconds: 10800
```

---

## Session Affinity Configuration

### Method 1: Service-Level ClientIP Affinity

The Service manifest above uses `sessionAffinity: ClientIP`. This works but has limitations:
- Based on source IP, which may change (NAT, proxies)
- All users behind same NAT share affinity

### Method 2: Ingress Cookie-Based Affinity (Recommended)

Cookie-based affinity is more reliable as it tracks individual browser sessions.

#### For NGINX Ingress Controller

```yaml
# ingress-nginx.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: baboon
  namespace: baboon
  annotations:
    kubernetes.io/ingress.class: nginx

    # STICKY SESSIONS - Cookie-based affinity
    nginx.ingress.kubernetes.io/affinity: "cookie"
    nginx.ingress.kubernetes.io/affinity-mode: "persistent"
    nginx.ingress.kubernetes.io/session-cookie-name: "BABOON_AFFINITY"
    nginx.ingress.kubernetes.io/session-cookie-expires: "10800"  # 3 hours
    nginx.ingress.kubernetes.io/session-cookie-max-age: "10800"
    nginx.ingress.kubernetes.io/session-cookie-change-on-failure: "true"
    nginx.ingress.kubernetes.io/session-cookie-samesite: "Lax"
    nginx.ingress.kubernetes.io/session-cookie-conditional-samesite-none: "true"

    # SSL redirect
    nginx.ingress.kubernetes.io/ssl-redirect: "true"

    # Websocket support (if needed in future)
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
spec:
  tls:
    - hosts:
        - baboon.example.com
      secretName: baboon-tls
  rules:
    - host: baboon.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: baboon
                port:
                  number: 80
```

#### For Traefik Ingress Controller

```yaml
# ingress-traefik.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: baboon
  namespace: baboon
  annotations:
    kubernetes.io/ingress.class: traefik

    # STICKY SESSIONS
    traefik.ingress.kubernetes.io/service.sticky.cookie: "true"
    traefik.ingress.kubernetes.io/service.sticky.cookie.name: "BABOON_AFFINITY"
    traefik.ingress.kubernetes.io/service.sticky.cookie.secure: "true"
    traefik.ingress.kubernetes.io/service.sticky.cookie.httponly: "true"
    traefik.ingress.kubernetes.io/service.sticky.cookie.samesite: "lax"
spec:
  tls:
    - hosts:
        - baboon.example.com
      secretName: baboon-tls
  rules:
    - host: baboon.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: baboon
                port:
                  number: 80
```

#### For AWS ALB Ingress Controller

```yaml
# ingress-alb.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: baboon
  namespace: baboon
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip

    # STICKY SESSIONS
    alb.ingress.kubernetes.io/target-group-attributes: >-
      stickiness.enabled=true,
      stickiness.lb_cookie.duration_seconds=10800,
      stickiness.type=lb_cookie

    # Health check
    alb.ingress.kubernetes.io/healthcheck-path: /api/health
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: "15"
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: "5"

    # SSL
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:region:account:certificate/xxx
spec:
  rules:
    - host: baboon.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: baboon
                port:
                  number: 80
```

---

## Persistent Storage

By default, each pod stores statistics in an emptyDir volume (ephemeral). For persistent stats across pod restarts:

### Option 1: Shared PVC (NFS/EFS)

```yaml
# pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: baboon-stats
  namespace: baboon
spec:
  accessModes:
    - ReadWriteMany  # Required for multi-pod access
  storageClassName: efs-sc  # Or your NFS storage class
  resources:
    requests:
      storage: 1Gi
```

Update the Deployment to use this PVC:

```yaml
# In deployment.yaml, update volumes:
volumes:
  - name: config
    persistentVolumeClaim:
      claimName: baboon-stats
```

### Option 2: Per-Pod Storage with StatefulSet

For isolated per-user stats, use a StatefulSet:

```yaml
# statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: baboon
  namespace: baboon
spec:
  serviceName: baboon
  replicas: 3
  selector:
    matchLabels:
      app.kubernetes.io/name: baboon
  template:
    # ... same as Deployment spec.template
  volumeClaimTemplates:
    - metadata:
        name: stats
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: standard
        resources:
          requests:
            storage: 100Mi
```

**Note:** With StatefulSet, sticky sessions are even more important as each pod has its own stats storage.

---

## Health Checks

### Kubernetes Probes

The deployment includes:

- **Liveness Probe**: Restarts pod if unhealthy
- **Readiness Probe**: Removes pod from service if not ready

Both use the `/api/health` endpoint which returns:

```json
{
  "status": "healthy",
  "active_sessions": 5
}
```

### External Health Monitoring

For external monitoring (e.g., uptime services):

```bash
# Simple health check
curl -f https://baboon.example.com/api/health

# Check with session count threshold
curl -s https://baboon.example.com/api/health | jq '.active_sessions < 1000'
```

---

## Scaling Considerations

### Horizontal Pod Autoscaler

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: baboon
  namespace: baboon
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: baboon
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
```

### Scaling Warnings

1. **Don't scale down aggressively**: Active sessions will be lost if pods are terminated
2. **Use PodDisruptionBudget** to ensure availability during updates:

```yaml
# pdb.yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: baboon
  namespace: baboon
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: baboon
```

### Session Draining

Before scaling down or updating, sessions should be drained. The application handles this gracefully:

1. Stop sending new sessions to the pod (remove from service)
2. Wait for active sessions to complete or timeout
3. Terminate the pod

Configure graceful termination:

```yaml
# In deployment.yaml
spec:
  template:
    spec:
      terminationGracePeriodSeconds: 300  # 5 minutes
```

---

## Monitoring

### Prometheus Metrics (Future)

The `/api/health` endpoint can be enhanced with Prometheus metrics:

```yaml
# servicemonitor.yaml (if using Prometheus Operator)
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: baboon
  namespace: baboon
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: baboon
  endpoints:
    - port: http
      path: /api/metrics
      interval: 30s
```

### Key Metrics to Monitor

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `baboon_active_sessions` | Current active sessions | > 100 per pod |
| `baboon_http_requests_total` | Total HTTP requests | Rate > 1000/s |
| `baboon_http_request_duration_seconds` | Request latency | p99 > 100ms |
| `baboon_session_duration_seconds` | Session duration | Avg > 30 min |

### Logging

```yaml
# In deployment.yaml, add logging sidecar if needed:
containers:
  - name: baboon
    # ... existing config

  - name: fluent-bit
    image: fluent/fluent-bit:latest
    volumeMounts:
      - name: logs
        mountPath: /var/log/baboon
```

---

## Troubleshooting

### Session Not Found (404)

**Symptom:** User gets 404 errors after starting a session.

**Causes:**
1. Sticky sessions not configured
2. Session cookie not being set/sent
3. Pod was terminated

**Debug:**
```bash
# Check if cookie is set
curl -v https://baboon.example.com/api/sessions -X POST

# Check ingress annotations
kubectl get ingress baboon -n baboon -o yaml | grep -A20 annotations

# Check service session affinity
kubectl get svc baboon -n baboon -o yaml | grep -A5 sessionAffinity
```

### Pod Keeps Restarting

**Check logs:**
```bash
kubectl logs -f deployment/baboon -n baboon

# Check previous container logs
kubectl logs deployment/baboon -n baboon --previous
```

**Check resources:**
```bash
kubectl top pods -n baboon
kubectl describe pod -l app.kubernetes.io/name=baboon -n baboon
```

### High Latency

**Check pod distribution:**
```bash
kubectl get pods -n baboon -o wide

# Check node resources
kubectl top nodes
```

**Check network policies:**
```bash
kubectl get networkpolicies -n baboon
```

### Session Affinity Not Working

**For NGINX Ingress:**
```bash
# Verify ingress controller config
kubectl exec -it $(kubectl get pods -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx -o jsonpath='{.items[0].metadata.name}') -n ingress-nginx -- cat /etc/nginx/nginx.conf | grep -A10 upstream
```

**For Traefik:**
```bash
# Check Traefik logs
kubectl logs -f deployment/traefik -n traefik
```

---

## Quick Start Commands

```bash
# Create namespace and deploy
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml  # If using AdSense
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress-nginx.yaml  # Or your ingress variant
kubectl apply -f hpa.yaml
kubectl apply -f pdb.yaml

# Verify deployment
kubectl get all -n baboon
kubectl get ingress -n baboon

# Test health endpoint
kubectl port-forward svc/baboon 8787:80 -n baboon
curl http://localhost:8787/api/health

# Watch pods
kubectl get pods -n baboon -w

# Check logs
kubectl logs -f deployment/baboon -n baboon

# Scale manually
kubectl scale deployment/baboon --replicas=5 -n baboon
```

---

## Complete Kustomize Setup

For easier management, use Kustomize:

```
k8s/
├── base/
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── hpa.yaml
│   └── pdb.yaml
└── overlays/
    ├── development/
    │   ├── kustomization.yaml
    │   └── ingress.yaml
    ├── staging/
    │   ├── kustomization.yaml
    │   └── ingress.yaml
    └── production/
        ├── kustomization.yaml
        ├── ingress.yaml
        └── resources.yaml
```

```yaml
# base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - deployment.yaml
  - service.yaml
  - hpa.yaml
  - pdb.yaml

# overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - ../../base
resources:
  - ingress.yaml
patchesStrategicMerge:
  - resources.yaml
images:
  - name: your-registry/baboon
    newTag: v1.0.0
```

Deploy with:
```bash
kubectl apply -k k8s/overlays/production/
```

---

## Support

For issues or questions:
- GitHub Issues: https://github.com/timlinux/baboon/issues
- Documentation: https://timlinux.github.io/baboon/

Made with ❤️ by [Kartoza](https://kartoza.com)
