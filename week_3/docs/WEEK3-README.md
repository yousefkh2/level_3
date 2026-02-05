# Week 3: Platform-as-a-Service (PaaS) - Complete ✅

## What We Built

A **production-ready managed PostgreSQL database service** running on SKE (STACKIT Kubernetes Engine), demonstrating real PaaS functionality.

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| **[week3-demo.sh](../scripts/week3-demo.sh)** | Lifecycle demonstration script (CREATE → USE → DESTROY) |
| **[connecting-to-database.md](./connecting-to-database.md)** | Connection guide |

## 🎯 Week 3 Requirements - Status

### Core Requirements

- [x] **SKE Cluster Creation**: Using the STACKIT Terraform Provider to provision an SKE (STACKIT Kubernetes Engine) Cluster
- [x] **PaaS Product Implementation** (Managed Database): Design and technical implementation of a simple PaaS service
  - [x] **Operator deployment**: Provisioning of an Operator (CloudNativePG)
  - [x] **Product Component Management**: Utilization of Custom Kubernetes Resources (CRs) for the provisioning and management of product components
  - [x] **Connectivity**: Documentation and demonstration of connecting to and using the PaaS product
    - See [Lifecycle Demo Script](../scripts/week3-demo.sh) for complete demonstration
- [x] **Understanding Kubernetes Concepts**: Deepening knowledge of Custom Resource Definitions (CRDs) and the functioning of Operators (Reconciler Pattern)

### Bonus

- [ ] **Automating the Deployment**: Introduction of a GitOps approach and CI/CD integration for automated provisioning of the SKE and the PaaS service

## 💡 Implementation Notes

### Why connect through services, not pods?

The demo connects to databases via Kubernetes services (`*-rw`, `*-ro`) instead of directly to pods:

```bash
# Through service (correct)
psql "postgresql://${USER}:${PASS}@${CLUSTER}-rw.default.svc:5432/${DB}"

# NOT directly to pod
# kubectl exec ${CLUSTER}-1 -- psql ...
```

**Why this matters:**
- **Production pattern**: Applications connect via services, not pod IPs
- **Service discovery**: DNS-based, stable endpoints
- **Load balancing**: Services route to healthy pods automatically
- **Read-write splitting**: `-rw` for writes, `-ro` for reads
- **Failover**: If primary fails, service routes to new primary
- **Real-world simulation**: Demonstrates how actual applications would connect

### Why kubectl wait instead of sleep?

The [demo script](../scripts/week3-demo.sh) uses `kubectl wait` with specific conditions instead of arbitrary sleep times:

```bash
# Wait for exact desired state, not arbitrary time
kubectl wait --for=jsonpath='{.status.readyInstances}'=3 cluster/${INSTANCE_NAME}
```

**Benefits:**
- Waits for actual desired state, not guessed duration
- Faster when resources are ready quickly
- More reliable when resources take longer
- Self-documenting - shows what we're waiting for

## 🚀 Quick Start

### 1. Run the Full Demo
```bash
cd ../scripts
./week3-demo.sh
```
This interactive demo shows:
- Creating a database cluster
- Extracting credentials
- Connecting and using the database
- Testing self-healing and scaling
- Destroying resources

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│          STACKIT Cloud (OpenStack)              │
│                                                  │
│  ┌───────────────────────────────────────────┐  │
│  │     SKE Managed Kubernetes Cluster        │  │
│  │                                            │  │
│  │  ┌──────────────────────────────────────┐ │  │
│  │  │  CNPG Operator (Control Plane)       │ │  │
│  │  │  - Watches for Cluster resources     │ │  │
│  │  │  - Reconciles desired vs actual      │ │  │
│  │  │  - Manages lifecycle operations      │ │  │
│  │  └──────────────────────────────────────┘ │  │
│  │                    ↓                       │  │
│  │  ┌──────────────────────────────────────┐ │  │
│  │  │  PostgreSQL Clusters (Data Plane)    │ │  │
│  │  │                                       │ │  │
│  │  │  ┌────────┐  ┌────────┐  ┌────────┐ │ │  │
│  │  │  │ Pod 1  │  │ Pod 2  │  │ Pod 3  │ │ │  │
│  │  │  │Primary │  │Replica │  │Replica │ │ │  │
│  │  │  └────────┘  └────────┘  └────────┘ │ │  │
│  │  │       ↓          ↓           ↓       │ │  │
│  │  │  ┌────────┐  ┌────────┐  ┌────────┐ │ │  │
│  │  │  │  PVC   │  │  PVC   │  │  PVC   │ │ │  │
│  │  │  └────────┘  └────────┘  └────────┘ │ │  │
│  │  │                                       │ │  │
│  │  │  Services:                            │ │  │
│  │  │  - cluster-rw (primary)               │ │  │
│  │  │  - cluster-ro (replicas)              │ │  │
│  │  │  - cluster-r  (all)                   │ │  │
│  │  │                                       │ │  │
│  │  │  Secrets:                             │ │  │
│  │  │  - cluster-app (credentials)          │ │  │
│  │  └──────────────────────────────────────┘ │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘

User → kubectl apply → Cluster CR → Operator → PostgreSQL Pods
```

## 📁 File Structure

```
week_3/
├── infrastructure/
│   └── main.tf                  # SKE cluster provisioning
├── manifests/
│   └── databases/
│       └── first-db.yaml        # Example cluster definition
├── docs/
│   ├── WEEK3-README.md          # This file ⭐
│   └── connecting-to-database.md
└── scripts/
    └── week3-demo.sh            # Interactive demo script ⭐
```

## 🔗 Additional Resources

- [CloudNativePG Documentation](https://cloudnative-pg.io/)
- [Kubernetes Operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)

