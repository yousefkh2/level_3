# Building a Cloud-Native Platform: From IaaS to PaaS

A 6-week build of a managed PostgreSQL Platform-as-a-Service on Kubernetes, starting from OpenStack infrastructure and progressing through API, UI, GitOps, and full-stack observability.

**Live URL:** `https://paas.null.stackit.run`

## Current Status (March 2026)

| Week | Topic | Status |
|------|-------|--------|
| 1 | IaaS with OpenStack (DevStack) | ✅ Complete |
| 2 | Kubernetes foundation on OpenStack (Terraform) | ✅ Complete |
| 3 | Managed PostgreSQL PaaS (CloudNativePG) | ✅ Complete |
| 4 | REST API control plane (Go + Echo + JWT + OpenAPI) | ✅ Complete |
| 5 | Web UI + Ingress + TLS + E2E tests | ✅ Complete |
| 6 | Observability (Prometheus + Grafana + Loki) | ✅ Complete |

## What Is Implemented

### Platform Capabilities
- **Managed PostgreSQL** provisioning on SKE via CloudNativePG Custom Resources.
- **REST API** (Go/Echo) for full CRUD + update (PATCH) of database instances.
- **JWT authentication** (`/auth/login`) protecting all database operations.
- **Audit logging** — every create/update/delete is recorded as a structured log and queryable via Loki.
- **Service event watcher** — background goroutine polls CNPG cluster status every 15 s, emits structured logs on state transitions (creating → healthy, scaling, etc.).
- **Prometheus metrics** — `api_requests_total`, `api_request_duration_seconds`, `api_requests_in_flight` exposed on `/metrics`.
- **HPA** — API deployment auto-scales 1→5 replicas based on CPU/memory utilisation.
- **Vue 3 web UI** for login, database CRUD, inline editing, connection info display, audit log viewer, and service event viewer.
- **NGINX Ingress** with TLS (cert-manager + Let's Encrypt):
  - `/api/*` → PaaS API service
  - `/*` → PaaS UI service
- **GitOps** — Argo CD app-of-apps pattern with ordered sync waves for the entire stack.
- **Observability stack** — Prometheus, Grafana, Loki + Promtail, all deployed via Helm through ArgoCD.

### API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | No | Health check |
| `GET` | `/metrics` | No | Prometheus metrics |
| `POST` | `/auth/login` | No | Authenticate, receive JWT |
| `POST` | `/databases` | JWT | Create a database cluster |
| `GET` | `/databases` | JWT | List all database clusters |
| `GET` | `/databases/:name` | JWT | Get database details + connection info |
| `PATCH` | `/databases/:name` | JWT | Update instances / storage |
| `DELETE` | `/databases/:name` | JWT | Delete a database cluster (idempotent) |
| `GET` | `/databases/:name/logs` | JWT | Audit logs for a database (via Loki) |
| `GET` | `/databases/:name/service-logs` | JWT | Service status events (via Loki) |
| `GET` | `/audit/logs` | JWT | Global audit log across all databases |

### GitOps Sync Order
All apps are managed in `gitops/apps/` via Argo CD sync waves:

| Wave | App | Source |
|------|-----|--------|
| 1 | `cnpg-operator` | Helm (CloudNativePG) |
| 2 | `databases` | Git — `week_3/manifests/databases/` |
| 3 | `paas-api` | Git — `week_4/manifests/` |
| 4 | `paas-ui` | Git — `week_5/manifests/` |
| 5 | `kube-prometheus-stack` | Helm (Prometheus + Grafana) |
| 5 | `loki-stack` | Helm (Loki + Promtail) |
| 6 | `paas-api-servicemonitor` | Git — `week_6/manifests/` |

Root app: `gitops/root-app.yaml`

## Repository Structure

```text
cloud/
├── README.md
├── gitops/
│   ├── root-app.yaml              # Argo CD root Application
│   └── apps/
│       ├── cnpg-operator.yaml     # wave 1 — CloudNativePG operator
│       ├── databases.yaml         # wave 2 — Database CRs
│       ├── paas-api.yaml          # wave 3 — API deployment + RBAC + HPA
│       ├── paas-ui.yaml           # wave 4 — UI deployment + service
│       ├── kube-prometheus-stack.yaml  # wave 5 — Prometheus + Grafana
│       ├── loki-stack.yaml        # wave 5 — Loki + Promtail
│       └── paas-api-servicemonitor.yaml # wave 6 — ServiceMonitor
├── week_1_2/                      # OpenStack + Kubernetes foundation
│   ├── *.tf                       # Terraform (VMs, networking, security groups)
│   ├── local.conf.example         # DevStack config template
│   └── README.md
├── week_3/                        # Managed PostgreSQL product on SKE
│   ├── infrastructure/            # Terraform for SKE + ArgoCD bootstrap
│   ├── manifests/databases/       # CNPG Cluster CRs
│   ├── docs/                      # Product documentation
│   └── scripts/                   # Demo & validation scripts
├── week_4/                        # API control plane
│   ├── api/                       # Go API source
│   │   ├── main.go                # Entry point, routes, graceful shutdown
│   │   ├── app/                   # App context, metrics, audit logging
│   │   ├── handlers/              # HTTP handlers + unit tests
│   │   ├── services/              # DB service, Loki adapter, status watcher
│   │   ├── middleware/            # JWT auth, observability middleware
│   │   ├── models/                # Request / response structs
│   │   ├── config/                # Kubernetes client setup
│   │   └── Dockerfile             # Multi-stage distroless build (~44 MB)
│   ├── manifests/                 # Deployment, Service, RBAC, HPA (Kustomize)
│   ├── tests/                     # k6 load tests + results
│   └── docs/                      # OpenAPI spec + architecture diagrams
├── week_5/                        # UI + Ingress + TLS
│   ├── ui/                        # Vue 3 + Vite SPA
│   │   ├── src/App.vue            # Single-file app (login, CRUD, logs viewer)
│   │   ├── e2e/                   # Playwright end-to-end tests
│   │   ├── nginx.conf             # Production NGINX config
│   │   └── Dockerfile             # Multi-stage Node → NGINX build
│   ├── manifests/                 # UI Deployment + Ingress (Kustomize)
│   └── cluster-issuer.yaml        # cert-manager Let's Encrypt ClusterIssuer
└── week_6/                        # Observability
    └── manifests/
        └── servicemonitor.yaml    # ServiceMonitor for PaaS API → Prometheus
```

## Week Highlights

### Week 1–2: Infrastructure Foundation
- DevStack / OpenStack installation and configuration.
- Terraform-based multi-node VM provisioning (3 Debian 12 VMs on OpenStack).
- Security group setup for inter-node communication, floating IPs, auto-generated Ansible inventory.
- Kubernetes the Hard Way walkthrough for foundational understanding.
- **Reference:** `week_1_2/README.md`

### Week 3: Managed Database Product
- STACKIT Kubernetes Engine (SKE) cluster provisioned via Terraform.
- CloudNativePG operator deployed; PostgreSQL clusters managed as Kubernetes CRs.
- Full lifecycle demo script (create, scale, connect, failover, delete).
- Argo CD bootstrapped via Helm in Terraform, root app-of-apps pattern established.
- **Reference:** `week_3/README.md`

### Week 4: API Control Plane
- Go / Echo REST API with full CRUD + PATCH for database clusters.
- JWT authentication middleware (HMAC-signed, 24 h expiry).
- Structured audit logging (create / update / delete events) for Loki ingestion.
- Background status watcher emitting service events on CNPG state transitions.
- Prometheus metrics via `promhttp` + custom counters / histograms / gauges.
- Observability middleware: per-request structured logging with method, path, status, duration.
- Loki query adapter — audit logs and service events retrievable per-database or globally.
- Multi-stage Dockerfile (distroless, non-root, ~44 MB).
- Kustomize manifests: Deployment, Service, ClusterRole, ClusterRoleBinding, ServiceAccount, HPA.
- k6 load tests (read-heavy + CRUD lifecycle scenarios) validating HPA scale-out.
- Unit tests with table-driven patterns and mocked dependencies.
- **Reference:** `week_4/README.md`

### Week 5: UI, Ingress, TLS
- Vue 3 / Vite single-page application served via NGINX.
- Login form → JWT stored in `localStorage` → all subsequent requests authenticated.
- Database dashboard: create, list, inline edit (PATCH), delete, view connection info.
- Audit log viewer and service event viewer per database (pulls from Loki via API).
- NGINX Ingress with regex-based path routing (`/api` → API, `/` → UI).
- TLS termination with cert-manager + Let's Encrypt (`ClusterIssuer` → `paas-tls-cert`).
- Playwright end-to-end tests covering the full user journey (login → create → view → edit → delete).
- **Reference:** `week_5/ui/`, `week_5/manifests/`

### Week 6: Observability
- **Prometheus** (via `kube-prometheus-stack` Helm chart): 7-day retention, 10 Gi persistent storage, scrapes all ServiceMonitors across namespaces.
- **Grafana**: bundled with kube-prometheus-stack, pre-configured with Prometheus as a data source, Loki data source addable for log exploration.
- **Loki + Promtail** (via `loki-stack` Helm chart): 10 Gi persistent storage, Promtail DaemonSet tailing container logs from every node.
- **ServiceMonitor** (`week_6/manifests/servicemonitor.yaml`): scrapes `/metrics` on the PaaS API every 15 s, wired to the `release: kube-prometheus-stack` label selector.
- End-to-end observability loop: API emits structured JSON logs → Promtail ships to Loki → API queries Loki to surface audit / service events in the UI; Prometheus scrapes `/metrics` → Grafana dashboards.
- **Reference:** `week_6/manifests/`, `gitops/apps/kube-prometheus-stack.yaml`, `gitops/apps/loki-stack.yaml`

## Technology Stack

| Layer | Technology |
|-------|-----------|
| **IaaS** | OpenStack (DevStack) |
| **Infrastructure as Code** | Terraform |
| **Kubernetes** | STACKIT Kubernetes Engine (SKE) |
| **Database Operator** | CloudNativePG |
| **API** | Go 1.23, Echo v4, controller-runtime |
| **Auth** | JWT (HMAC-SHA256) |
| **UI** | Vue 3, Vite |
| **Web Server** | NGINX (UI container + Ingress controller) |
| **TLS** | cert-manager + Let's Encrypt |
| **GitOps** | Argo CD (app-of-apps) |
| **Metrics** | Prometheus, Grafana |
| **Logging** | Loki, Promtail |
| **Testing** | Go `testing`, Playwright (E2E), k6 (load) |
| **Container Build** | Multi-stage Dockerfile, distroless base |
| **Registry** | STACKIT Container Registry |

## Quick Start (Current Flow)

### 1. Bootstrap Argo CD root app
```bash
kubectl apply -f gitops/root-app.yaml
```

### 2. Verify apps
```bash
kubectl get applications -n argocd
```

### 3. Verify control plane workloads
```bash
kubectl get pods -n paas-control-plane
kubectl get ingress -n paas-control-plane
```

Note: TLS issuer manifest lives at `week_5/cluster-issuer.yaml`. If not managed elsewhere, apply it manually before expecting TLS certificate issuance.

## Technology Stack
- Infrastructure: OpenStack, Terraform, STACKIT SKE
- Kubernetes: CRDs, operators, Argo CD GitOps
- Database: PostgreSQL via CloudNativePG
- API: Go, Echo, controller-runtime, OpenAPI
- UI: Vue 3, Vite, Axios, Vue Router, NGINX
- CI/CD: GitHub Actions (Terraform pipeline)

## Next Milestone
Finish Week 5 hardening/documentation and move to Week 6 observability implementation.
