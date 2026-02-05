# Connecting to Your Managed PostgreSQL Instance

## Prerequisites
- `kubectl` configured with access to the cluster
- The name of your database instance (e.g., `first-db`)

## Step 1: Extract Connection Credentials

### Username
```bash
kubectl get secret <instance-name>-app -o jsonpath='{.data.username}' | base64 -d
```

### Password
```bash
kubectl get secret <instance-name>-app -o jsonpath='{.data.password}' | base64 -d
```

### Database Name
```bash
kubectl get secret <instance-name>-app -o jsonpath='{.data.dbname}' | base64 -d
```

### Host
```bash
kubectl get service <instance-name>-rw -o jsonpath='{.metadata.name}'
```

### Port
```bash
kubectl get service <instance-name>-rw -o jsonpath='{.spec.ports[0].port}'
```

## Step 2: Build Connection String

Format:
```
postgresql://username:password@host:port/dbname?sslmode=require
```

## Step 3: Connect

Using `psql`:
```bash
psql "postgresql://username:password@host:port/dbname?sslmode=require"
```
