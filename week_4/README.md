# Week 4 - Provisioning and Interaction via RESTful API

## Overview

This week implements a production-ready RESTful API that provides access to the CloudNativePG-based database platform built in Week 3. Instead of manually creating YAML manifests, users can now provision and manage PostgreSQL databases through HTTP endpoints.

## What Was Built

### Core Features ✅

- **RESTful API Endpoints**
  - `POST /databases` - Create a new database cluster
  - `GET /databases` - List all database clusters
  - `GET /databases/{name}` - Get database details + connection info
  - `DELETE /databases/{name}` - Delete a database cluster (idempotent)

- **API Specification**
  - Full OpenAPI 3.0 documentation → [`docs/openapi.yaml`](docs/openapi.yaml)
  - Comprehensive schema definitions
  - Error response documentation

- **Unit Tests**
  - Handler tests with mocked dependencies
  - Table-driven test patterns
  - Coverage of success and error paths

- **Containerization**
  - Multi-stage Dockerfile with distroless runtime
  - Final image size: ~44MB
  - Security: runs as non-root user

- **Kubernetes Deployment**
  - Deployed to `paas-control-plane` namespace
  - RBAC: ClusterRole for cross-namespace access
  - Health checks (liveness + readiness probes)
  - Resource limits configured

- **Architecture Documentation**
  - System architecture diagram
  - Database creation sequence diagram
  - Error handling flowchart
  - See: [`docs/architecture.md`](docs/architecture.md)

## Architecture

The API acts as a **control plane** that sits above the CloudNativePG operator:

```
User → REST API → Kubernetes API → CNPG Operator → PostgreSQL Pods
```

For detailed architecture diagrams, see [`docs/architecture.md`](docs/architecture.md).

### Key Components

| Component | Purpose | Namespace |
|-----------|---------|-----------|
| **PaaS API** | REST interface for database management | `paas-control-plane` |
| **ServiceAccount** | Identity with RBAC permissions | `paas-control-plane` |
| **ClusterRole** | Permissions to manage Cluster CRs | cluster-wide |
| **Service** | Stable endpoint (paas-api:80) | `paas-control-plane` |
| **CNPG Operator** | Reconciles Cluster CRs into PostgreSQL | `cnpg-system` |
| **Database Pods** | Actual PostgreSQL instances | `default` |

## Quick Start

### Prerequisites

- SKE cluster from Week 3 with CNPG operator installed
- `kubectl` configured with cluster access
- Docker or Podman for building images

### 1. Build the Container Image

```bash
cd api/
podman build -t paas-api:latest .
```

### 2. Push to Registry

```bash
# Tag for your registry
podman tag paas-api:latest registry.onstackit.cloud/YOUR_PROJECT/paas-api:v1

# Push
podman push registry.onstackit.cloud/YOUR_PROJECT/paas-api:v1
```

### 3. Deploy to Kubernetes

```bash
cd ../manifests/
kubectl apply -f deployment.yaml
```

### 4. Verify Deployment

```bash
# Check pod status
kubectl get pods -n paas-control-plane

# Check logs
kubectl logs -n paas-control-plane -l app=paas-api
```

### 5. Test the API

```bash
# Port-forward to access locally
kubectl port-forward -n paas-control-plane svc/paas-api 8080:80

# Create a database
curl -X POST http://localhost:8080/databases \
  -H "Content-Type: application/json" \
  -d '{"name":"test-db","instances":3,"storage":"1Gi"}'

# List databases
curl http://localhost:8080/databases

# Get database details
curl http://localhost:8080/databases/test-db

# Delete database
curl -X DELETE http://localhost:8080/databases/test-db
```

## Project Structure

```
week_4/
├── README.md                   # This file
├── api/
│   ├── main.go                 # Entry point
│   ├── app/
│   │   └── app.go              # App context & interfaces
│   ├── handlers/
│   │   ├── database.go         # HTTP handlers
│   │   └── database_test.go    # Handler unit tests
│   ├── services/
│   │   └── database.go         # Business logic
│   ├── models/
│   │   └── database.go         # Request/response models
│   ├── config/
│   │   └── kubernetes.go       # K8s client setup
│   ├── Dockerfile              # Multi-stage build
│   ├── .dockerignore           # Build exclusions
│   ├── go.mod                  # Dependencies
│   └── go.sum
├── manifests/
│   └── deployment.yaml         # K8s manifests
└── docs/
    ├── openapi.yaml            # API specification
    └── architecture.md         # Architecture diagrams
```

## Testing

### Run Unit Tests

```bash
cd api/
go test ./handlers/... -v -cover
```

**Coverage:**
- ✅ CreateDatabase (success, already exists, invalid body, invalid config, unauthorized, internal error)
- ✅ ListDatabases (success, empty list, internal error)
- ✅ GetDatabase (success, not found, internal error)
- ✅ DeleteDatabase (success, idempotent delete, internal error)


## Technology Stack

| Layer | Technology |
|-------|-----------|
| **Language** | Go 1.23 |
| **HTTP Framework** | Echo v4 |
| **K8s Client** | controller-runtime |
| **Container Runtime** | Podman / Docker |
| **Base Image** | gcr.io/distroless/static-debian12 |
| **Orchestration** | Kubernetes (SKE) |

## Security

- **Non-root container**: Runs as UID 65532 (distroless nonroot user)
- **RBAC**: ServiceAccount with minimal required permissions
- **Static binary**: No runtime dependencies (CGO_ENABLED=0)
- **No shell**: Distroless image has no shell or package manager
- **Resource limits**: CPU and memory limits prevent resource exhaustion

## Performance

- **Image size**: 44.5 MB (vs ~800MB with full Go toolchain)
- **Startup time**: ~1-2 seconds
- **Build cache**: Subsequent builds ~5-10x faster with BuildKit mounts

## Troubleshooting

### Pod not starting

```bash
kubectl describe pod -n paas-control-plane -l app=paas-api
kubectl logs -n paas-control-plane -l app=paas-api
```

### API returns 401 Unauthorized when creating databases

Check RBAC permissions:
```bash
kubectl get clusterrole cnpg-cluster-manager -o yaml
kubectl get clusterrolebinding paas-api-cluster-binding -o yaml
```

### Database creation fails

Verify CNPG operator is running:
```bash
kubectl get pods -n cnpg-system
kubectl logs -n cnpg-system -l app.kubernetes.io/name=cloudnative-pg
```

## Week 4 Requirements Checklist

### Core Requirements ✅

- [x] API Development for Product Instances
  - [x] Create endpoint
  - [x] Delete endpoint
  - [x] List endpoint
  - [x] Get endpoint with connection info
- [x] API Specification (OpenAPI)
- [x] Unit Tests for each endpoint
- [x] Docker Container Image
- [x] Upload to STACKIT Container Registry
- [x] Provisioning via SKE
- [x] Understanding the Creation Process
  - [x] Flowchart for creation
  - [x] Understanding RESTful API basics

### Bonus Items

- [ ] Automated API Deployment (GitOps integration)
- [ ] Auto-Scaling and Performance Tests
  - [ ] HPA configuration
  - [ ] Load testing
- [ ] Update Functionality (PATCH endpoint)
