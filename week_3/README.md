# Week 3: Platform-as-a-Service (PaaS) - Managed PostgreSQL

## Overview

Implementation of a production-ready managed PostgreSQL database service on SKE using the CloudNativePG operator.

## Goal

Design and implement a Platform-as-a-Service offering on top of Kubernetes.

## Requirements & Status

### Core Requirements

- [x] **SKE Cluster Creation**: Using the STACKIT Terraform Provider to provision an SKE (STACKIT Kubernetes Engine) Cluster
- [x] **PaaS Product Implementation** (Managed Database): Design and technical implementation of a simple PaaS service
  - [x] **Operator deployment**: Provisioning of an Operator (CloudNativePG)
  - [x] **Product Component Management**: Utilization of Custom Kubernetes Resources (CRs) for the provisioning and management of product components
  - [x] **Connectivity**: Documentation and demonstration of connecting to and using the PaaS product
    - See [Lifecycle Demo Script](scripts/week3-demo.sh) for complete demonstration
- [x] **Understanding Kubernetes Concepts**: Deepening knowledge of Custom Resource Definitions (CRDs) and the functioning of Operators (Reconciler Pattern)
  - See [documentation](docs/WEEK3-DOCUMENTATION.md) for detailed explanations

### Bonus

- [ ] **Automating the Deployment**: Introduction of a GitOps approach and CI/CD integration for automated provisioning of the SKE and the PaaS service

## Implementation Notes

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

The [demo script](scripts/week3-demo.sh) uses `kubectl wait` with specific conditions instead of arbitrary sleep times:

```bash
# Wait for exact desired state, not arbitrary time
kubectl wait --for=jsonpath='{.status.readyInstances}'=3 cluster/${INSTANCE_NAME}
```

**Benefits:**
- Waits for actual desired state, not guessed duration
- Faster when resources are ready quickly
- More reliable when resources take longer
- Self-documenting - shows what we're waiting for

## Structure

```
week_3/
├── README.md                    # This file
├── infrastructure/              # SKE cluster (Terraform)
│   ├── main.tf
│   ├── terraform.tfvars
│   └── terraform.tfstate
├── manifests/                   # Kubernetes manifests (GitOps watches here)
│   └── databases/               # Database cluster definitions
│       └── first-db.yaml
├── docs/                        # Documentation
│   ├── WEEK3-README.md          # Main overview
│   ├── WEEK3-DOCUMENTATION.md   # Complete guide
│   ├── WEEK3-QUICKREF.md        # Command reference
│   ├── WEEK3-NEXT-STEPS.md      # What to do next
│   └── week3-architecture.svg   # Architecture diagram
└── scripts/                     # Helper scripts
    └── week3-demo.sh            # Lifecycle demonstration
```

## Quick Start

### Run the Demo
```bash
cd scripts
./week3-demo.sh
```

### Create a Database
```bash
kubectl apply -f manifests/databases/first-db.yaml
```

## Documentation

- **[Main Documentation](docs/WEEK3-DOCUMENTATION.md)** - Complete technical guide
- **[Quick Reference](docs/WEEK3-QUICKREF.md)** - Common commands
- **[Next Steps](docs/WEEK3-NEXT-STEPS.md)** - What to do next

## Scripts

- **[Lifecycle Demo](scripts/week3-demo.sh)** - Full CREATE → USE → DESTROY demo

## Infrastructure

The SKE cluster is managed in [`infrastructure/`](infrastructure/) using Terraform.

## Next Week

Week 4 will build a REST API on top of this PaaS offering for programmatic database provisioning.
