# Building a Cloud-Native Platform: From IaaS to PaaS

**A complete 6-week journey building a production-grade Platform-as-a-Service on Kubernetes**

## Project Overview

This project demonstrates building a complete cloud platform stack end-to-end: starting from raw infrastructure with OpenStack, moving through Kubernetes, and ending with a production-grade PaaS product featuring APIs, UI, automation, and observability.

**Final Outcome**: A fully functional managed database PaaS on Kubernetes (SKE), complete with REST API, Web UI, GitOps automation, and full observability.

## Project Structure

```
cloud/
├── README.md               # This file
├── week_1_2/              # Week 1-2: OpenStack + K8s
├── week_3/                # Week 3: PaaS Product - Managed Database
│   ├── infrastructure/    # SKE cluster (Terraform)
│   ├── manifests/         # K8s resources (GitOps watches)
│   ├── docs/              # Documentation
│   └── scripts/           # Demo & validation scripts
├── week_4/                # Week 4: REST API (Coming soon)
├── week_5/                # Week 5: Web UI & Ingress (Coming soon)
└── week_6/                # Week 6: Observability (Coming soon)
```

## Progress

| Week | Topic | Status |
|------|-------|--------|
| 1 | IaaS with OpenStack | ✅ Complete |
| 2 | K8s on OpenStack (Terraform) | ✅ Complete |
| 3 | PaaS Product (Managed DB) | ✅ Core Complete, ⏳ Bonus |
| 4 | REST API | 📋 Planned |
| 5 | Web UI & Ingress | 📋 Planned |
| 6 | Observability | 📋 Planned |

## Week Summaries

### Week 1-2: Infrastructure Foundation
- OpenStack deployment with DevStack
- Kubernetes cluster provisioning with Terraform
- Location: [`week_1_2/`](week_1_2/)

### Week 3: PaaS Product - Managed PostgreSQL ✅
**What was built:**
- SKE cluster on STACKIT Cloud
- CloudNativePG operator deployment
- Managed database provisioning via Custom Resources
- Complete documentation and lifecycle demonstration

**Key deliverables:**
- Working managed database service
- Self-healing & auto-scaling capabilities
- Full lifecycle demo script
- Comprehensive documentation

📁 [Week 3 Details](week_3/README.md)

### Week 4: REST API
Build a RESTful API for programmatic database provisioning with OpenAPI specs and unit tests.

📁 [Week 4 Details](week_4/) (Coming soon)

### Week 5: Web UI & Ingress
User-facing web interface with secure authentication and SSL/TLS.

📁 [Week 5 Details](week_5/) (Coming soon)

### Week 6: Observability
Prometheus, Grafana, Loki for monitoring and audit logging.

📁 [Week 6 Details](week_6/) (Coming soon)

## Technology Stack

- **Infrastructure**: OpenStack, Terraform, STACKIT SKE
- **Kubernetes**: K8s 1.28+, Custom Resources, Operators
- **Database**: PostgreSQL via CloudNativePG
- **GitOps**: ArgoCD/Flux (Week 3 Bonus)
- **API**: Go/Python/Node.js (Week 4)
- **UI**: Vue.js/React (Week 5)
- **Observability**: Prometheus, Grafana, Loki (Week 6)

## Learning Outcomes

By completing this project, you will understand:
- Infrastructure-as-Code with Terraform
- Kubernetes operators and Custom Resources
- PaaS architecture and design patterns
- RESTful API development
- GitOps workflows
- Production observability and monitoring

## Architecture Evolution

The architecture grows each week:
- **Week 1-2**: OpenStack → VMs → Kubernetes
- **Week 3**: + SKE → CNPG Operator → PostgreSQL Clusters
- **Week 4**: + REST API → Programmatic provisioning
- **Week 5**: + Web UI → Ingress → User-facing platform
- **Week 6**: + Prometheus/Grafana → Full observability

## Contributing

This is a learning project following the STACKIT Cloud-Native Platform track.

## License

Apache License 2.0 - See LICENSE file for details.

---

**Current Focus**: Week 3 - Completing GitOps bonus and preparing for Week 4
