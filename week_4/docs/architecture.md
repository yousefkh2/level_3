# PaaS Database Platform - Architecture & Flow Diagrams

## 1. System Architecture (Week 4)

```mermaid
graph TB
    subgraph "External Access"
        User[User/Client]
    end
    
    subgraph "SKE Cluster"
        subgraph "Control Plane Namespace (paas-control-plane)"
            Service[Service<br/>paas-api:80]
            Pod[API Deployment<br/>Port 8080]
            SA[ServiceAccount<br/>paas-api-sa]
            Role[ClusterRole<br/>CNPG Manager]
        end
        
        subgraph "Kubernetes API Server"
            K8sAPI[API Server<br/>RBAC Enforcement]
        end
        
        subgraph "Database Namespace (default)"
            CR[Cluster CR<br/>postgresql.cnpg.io]
            Operator[CNPG Operator<br/>Watches CRs]
            DB1[(PostgreSQL Pod 1<br/>Primary)]
            DB2[(PostgreSQL Pod 2<br/>Replica)]
            DB3[(PostgreSQL Pod 3<br/>Replica)]
        end
    end
    
    User -->|HTTP POST /databases| Service
    Service --> Pod
    Pod -->|Uses| SA
    SA -->|Has Permissions| Role
    Pod -->|Create Cluster CR| K8sAPI
    K8sAPI -->|Creates| CR
    Operator -->|Watches| CR
    Operator -->|Spawns| DB1
    Operator -->|Spawns| DB2
    Operator -->|Spawns| DB3
    
    style Pod fill:#e1f5ff
    style SA fill:#fff4e1
    style Role fill:#fff4e1
    style Operator fill:#e8f5e9
    style CR fill:#f3e5f5
```

## 2. Database Creation Flow (Sequence Diagram)

```mermaid
sequenceDiagram
    actor User
    participant API as PaaS API<br/>(Pod)
    participant K8s as Kubernetes<br/>API Server
    participant CNPG as CNPG<br/>Operator
    participant DB as PostgreSQL<br/>Pods
    
    User->>+API: POST /databases<br/>{"name":"my-db","instances":3,"storage":"1Gi"}
    
    Note over API: 1. Parse & Validate Request
    API->>API: Bind JSON to struct
    
    Note over API: 2. Call Service Layer
    API->>API: DBService.CreateDatabase()
    
    Note over API: 3. Create Cluster CR
    API->>+K8s: Create Cluster CR<br/>using ServiceAccount credentials
    
    Note over K8s: 4. RBAC Check
    K8s->>K8s: Verify paas-api-sa has<br/>permission to create Clusters
    
    K8s-->>-API: CR Created Successfully
    
    API-->>-User: 201 Created<br/>{"name":"my-db","status":"creating",...}
    
    Note over CNPG: 5. Operator Detects Change
    K8s--)CNPG: Watch Event:<br/>New Cluster CR
    
    Note over CNPG: 6. Reconciliation Loop
    CNPG->>CNPG: Read Cluster spec<br/>(instances=3, storage=1Gi)
    
    Note over CNPG: 7. Create Resources
    CNPG->>+DB: Create StatefulSet<br/>(3 replicas)
    CNPG->>DB: Create Services<br/>(-rw, -ro, -r)
    CNPG->>DB: Create PVCs<br/>(1Gi each)
    
    Note over DB: 8. PostgreSQL Initialization
    DB->>DB: Initialize primary
    DB->>DB: Configure replication
    DB-->>-CNPG: Pods Ready
    
    Note over CNPG: 9. Update Status
    CNPG->>K8s: Update Cluster CR status<br/>"Cluster in healthy state"
    
    Note over User,DB: User can now query<br/>GET /databases/my-db<br/>to retrieve connection info
```

## 3. Create Flow (Simplified Flowchart)

```mermaid
flowchart TD
    Start([User sends POST /databases]) --> Parse[Parse JSON Request Body]
    Parse --> Validate{Valid Request?}
    Validate -->|No| Error400[Return 400 Bad Request]
    Validate -->|Yes| CreateCR[Create Cluster CR via K8s Client]
    
    CreateCR --> RBACCheck{ServiceAccount<br/>has permission?}
    RBACCheck -->|No| Error401[Return 401 Unauthorized]
    RBACCheck -->|Yes| K8sCreate[Kubernetes creates CR]
    
    K8sCreate --> Exists{Cluster<br/>already exists?}
    Exists -->|Yes| Error409[Return 409 Conflict]
    Exists -->|No| Return201[Return 201 Created<br/>with cluster metadata]
    
    Return201 --> Background[CNPG Operator processes CR in background]
    Background --> SpawnPods[Operator spawns PostgreSQL pods]
    SpawnPods --> End([Database becomes healthy])
    
    Error400 --> EndError([End])
    Error401 --> EndError
    Error409 --> EndError
    
    style Start fill:#e3f2fd
    style End fill:#c8e6c9
    style EndError fill:#ffcdd2
    style CreateCR fill:#fff9c4
    style SpawnPods fill:#c8e6c9
```

## 4. Key Endpoints and Their Flow

### POST /databases
1. User → Service → API Pod
2. API validates request body
3. API creates Cluster CR in Kubernetes
4. Returns 201 with metadata
5. CNPG Operator asynchronously provisions pods

### GET /databases
1. User → Service → API Pod
2. API lists all Cluster CRs
3. Returns array of database summaries

### GET /databases/{name}
1. User → Service → API Pod
2. API fetches specific Cluster CR
3. Generates connection info
4. Returns database details + connection string

### DELETE /databases/{name}
1. User → Service → API Pod
2. API deletes Cluster CR
3. CNPG Operator tears down pods
4. Returns 204 No Content
