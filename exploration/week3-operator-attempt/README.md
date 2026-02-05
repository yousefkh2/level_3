# Week 3: Custom Operator Exploration

## Overview
Explored building a custom Kubernetes operator as a wrapper around CloudNativePG to provide a simplified, tier-based database provisioning API.

## Approach
Used Kubebuilder to scaffold a custom operator with:
- **ManagedDatabase CRD**: Custom resource with tier-based configuration (small/medium/large)
- **Controller Logic**: Reconciliation loop to manage underlying CNPG Cluster resources
- **Status Management**: Track provisioning state, service endpoints, and credentials

## Implementation Details
- **CRD Definition**: `api/v1alpha1/manageddatabase_types.go`
  - Spec: tier selection (enum validation)
  - Status: phase, message, service, secretName
- **Controller**: `internal/controller/manageddatabase_controller.go`
  - Basic reconciliation scaffolding
  - Tier-to-resource mapping logic outlined

## Decision to Pivot
After initial implementation, decided to use CloudNativePG directly instead:

### Reasons:
1. **Unnecessary Abstraction**: CNPG already provides excellent APIs and flexibility
2. **Maintenance Overhead**: Custom operator adds another layer to maintain and debug
3. **Limited Value**: Tier abstraction can be achieved with Helm values or simple templates
4. **Learning Focus**: Direct CNPG usage provides better learning of PostgreSQL clustering concepts

### What Was Learned:
- Kubebuilder project structure and scaffolding
- CRD design with validation markers
- Controller reconciliation patterns
- Importance of evaluating abstraction trade-offs

## Next Steps
Moving forward with direct CloudNativePG Cluster resources, potentially wrapped in Helm charts for tier-based deployments.

---
*Generated with: Kubebuilder v4.11.1*
