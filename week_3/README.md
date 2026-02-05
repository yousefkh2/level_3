# Week 3: Platform-as-a-Service (PaaS) - Managed PostgreSQL

## Overview

Implementation of a production-ready managed PostgreSQL database service on SKE using the CloudNativePG operator.

## Status

✅ **Core Requirements Complete**
- SKE Cluster provisioned via Terraform
- CNPG Operator deployed
- Custom Resources working
- Full connectivity documented
- Complete lifecycle demonstration

⏳ **Bonus: GitOps** (In Progress)
- ArgoCD/Flux setup for automated deployment

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
