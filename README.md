# Building a Cloud-Native Platform: From IaaS to PaaS

A 6-week build of a managed PostgreSQL Platform-as-a-Service on Kubernetes, starting from OpenStack infrastructure and moving up to API, UI, and GitOps operations.

## Current Status (February 2026)

| Week | Topic | Status |
|------|-------|--------|
| 1 | IaaS with OpenStack (DevStack) | Complete |
| 2 | Kubernetes foundation on OpenStack (Terraform) | Complete |
| 3 | Managed PostgreSQL PaaS (CloudNativePG) | Complete |
| 4 | REST API control plane (Go + Echo + OpenAPI) | Complete |
| 5 | Web UI + Ingress + TLS integration | In Progress (near completion) |
| 6 | Observability (Prometheus/Grafana/Loki) | Not started |

## What Is Implemented Today

### Platform Capabilities
- Managed PostgreSQL provisioning on SKE via CloudNativePG CRs.
- REST API for create/list/get/update/delete database instances.
- JWT-based API auth (`/auth/login`) for protected database operations.
- Vue 3 web UI for login and database lifecycle actions.
- NGINX ingress routing:
  - `/api` -> PaaS API service
  - `/` -> PaaS UI service
- GitOps automation with Argo CD app-of-apps and sync waves.

### GitOps Sync Order
The repo currently defines child apps in `gitops/apps/` with sync waves:
1. `cnpg-operator` (wave 1)
2. `databases` (wave 2)
3. `paas-api` (wave 3)
4. `paas-ui` (wave 4)

Root app: `gitops/root-app.yaml`

## Repository Structure

```text
cloud/
├── README.md
├── gitops/
│   ├── root-app.yaml
│   └── apps/
│       ├── cnpg-operator.yaml
│       ├── databases.yaml
│       ├── paas-api.yaml
│       └── paas-ui.yaml
├── week_1_2/                  # OpenStack + Kubernetes foundation
├── week_3/                    # Managed PostgreSQL product on SKE
│   ├── infrastructure/        # Terraform for SKE + base setup
│   ├── manifests/             # Database CR manifests
│   ├── docs/
│   └── scripts/
├── week_4/                    # API control plane
│   ├── api/                   # Go API code (Echo, JWT, handlers, tests)
│   ├── manifests/             # Kustomize, deployment, RBAC, HPA
│   └── docs/                  # OpenAPI + architecture
├── week_5/                    # UI + ingress layer
│   ├── ui/                    # Vue 3 + Vite frontend
│   ├── manifests/             # UI deployment + ingress (kustomize)
│   └── cluster-issuer.yaml    # cert-manager ClusterIssuer manifest
└── week_6/                    # Observability (planned)
```

## Week Highlights

### Week 1-2: Infrastructure Foundation
- DevStack/OpenStack setup.
- Terraform-based VM/network provisioning for Kubernetes groundwork.
- Main reference: `week_1_2/README.md`

### Week 3: Managed Database Product
- CloudNativePG operator + PostgreSQL cluster provisioning via CRs.
- Lifecycle demonstration and product documentation.
- Main reference: `week_3/README.md`

### Week 4: API Control Plane
- Go/Echo API with health + auth + CRUD endpoints:
  - `GET /health`
  - `POST /auth/login`
  - `POST /databases`
  - `GET /databases`
  - `GET /databases/:name`
  - `PATCH /databases/:name`
  - `DELETE /databases/:name`
- RBAC + service account model for managing CNPG `Cluster` resources.
- Containerized deployment and Kustomize overlays in `week_4/manifests/`.
- Main reference: `week_4/README.md`

### Week 5: UI, Ingress, TLS
- Vue UI for login and database operations.
- API integration through `VITE_API_URL` (production uses `/api`).
- Ingress routes UI and API behind one host (`paas.null.stackit.run`).
- TLS setup prepared with cert-manager `ClusterIssuer` manifest.
- Main references: `week_5/ui/`, `week_5/manifests/`

### Week 6: Observability (Next)
- Planned: Prometheus, Grafana, Loki, and alerting/audit visibility.

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
