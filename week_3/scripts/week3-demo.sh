#!/bin/bash
set -e

INSTANCE_NAME="lifecycle-demo-db"
NAMESPACE="default"
INITIAL_INSTANCES=3
SCALED_INSTANCES=5

echo "=== 1. Creating PostgreSQL Instance ==="
kubectl apply -f - <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: ${INSTANCE_NAME}
  namespace: ${NAMESPACE}
spec:
  instances: ${INITIAL_INSTANCES}
  storage:
    size: 1Gi
EOF

echo "=== 2. Waiting for cluster to be ready ==="
sleep 5
kubectl wait --for=jsonpath='{.status.phase}'=Cluster\ in\ healthy\ state cluster/${INSTANCE_NAME} -n ${NAMESPACE} --timeout=300s
kubectl get pods -l cnpg.io/cluster=${INSTANCE_NAME} -n ${NAMESPACE}

echo "=== 3. Extracting credentials ==="
export DB_USER=$(kubectl get secret ${INSTANCE_NAME}-app -n ${NAMESPACE} -o jsonpath='{.data.username}' | base64 -d)
export DB_PASSWORD=$(kubectl get secret ${INSTANCE_NAME}-app -n ${NAMESPACE} -o jsonpath='{.data.password}' | base64 -d)
export DB_NAME=$(kubectl get secret ${INSTANCE_NAME}-app -n ${NAMESPACE} -o jsonpath='{.data.dbname}' | base64 -d)
echo "Username: ${DB_USER}"
echo "Database: ${DB_NAME}"

echo "=== 4. Testing write to primary ==="
kubectl run psql-client-write --rm -i --restart=Never --image=postgres:16 -- \
  psql "postgresql://${DB_USER}:${DB_PASSWORD}@${INSTANCE_NAME}-rw.${NAMESPACE}.svc:5432/${DB_NAME}" <<EOF
CREATE TABLE IF NOT EXISTS test_data (id SERIAL PRIMARY KEY, value TEXT, created_at TIMESTAMP DEFAULT NOW());
INSERT INTO test_data (value) VALUES ('test-entry-1'), ('test-entry-2');
SELECT * FROM test_data;
EOF

echo "=== 5. Testing read from replica ==="
kubectl run psql-client-read --rm -i --restart=Never --image=postgres:16 -- \
  psql "postgresql://${DB_USER}:${DB_PASSWORD}@${INSTANCE_NAME}-ro.${NAMESPACE}.svc:5432/${DB_NAME}" \
  -c "SELECT * FROM test_data;"

echo "=== 6. Testing self-healing ==="
echo "Pods before deletion:"
kubectl get pods -l cnpg.io/cluster=${INSTANCE_NAME} -n ${NAMESPACE}
kubectl delete pod ${INSTANCE_NAME}-2 -n ${NAMESPACE}
echo "Waiting for self-healing..."
kubectl wait --for=jsonpath="{.status.readyInstances}"=${INITIAL_INSTANCES} cluster/${INSTANCE_NAME} -n ${NAMESPACE} --timeout=120s
echo "Pods after self-healing:"
kubectl get pods -l cnpg.io/cluster=${INSTANCE_NAME} -n ${NAMESPACE}

echo "=== 7. Scaling up ==="
kubectl patch cluster ${INSTANCE_NAME} -n ${NAMESPACE} --type='merge' -p '{"spec":{"instances":${SCALED_INSTANCES}}}'
echo "Waiting for scaling to complete..."
DESIRED_INSTANCES=$(kubectl get cluster ${INSTANCE_NAME} -n ${NAMESPACE} -o jsonpath='{.spec.instances}')
kubectl wait --for=jsonpath="{.status.readyInstances}"=${DESIRED_INSTANCES} cluster/${INSTANCE_NAME} -n ${NAMESPACE} --timeout=180s
kubectl get pods -l cnpg.io/cluster=${INSTANCE_NAME} -n ${NAMESPACE}

echo "=== 8. Cleanup ==="
kubectl delete cluster ${INSTANCE_NAME} -n ${NAMESPACE}
kubectl delete pvc -l cnpg.io/cluster=${INSTANCE_NAME} -n ${NAMESPACE}