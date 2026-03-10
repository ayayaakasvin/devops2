#!/bin/bash
set -e  # stop if any command fails
set -o pipefail

echo "🚀 Applying database StatefulSet..."
kubectl apply -f database.yaml

echo "🛠 Applying database Service..."
kubectl apply -f postgres-service.yaml

echo "⏳ Waiting for PostgreSQL pod to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=60s

echo "⚡ Running migration Job..."
kubectl apply -f migration.yaml

# optional: wait for migration to finish
echo "⏳ Waiting for migration Job to complete..."
kubectl wait --for=condition=complete job/db-migration --timeout=60s

echo "🚀 Deploying backend Deployment..."
kubectl apply -f backend.yaml

echo "🛡 Deploying backend Service..."
kubectl apply -f backend-service.yaml

echo "📡 Applying NetworkPolicy..."
kubectl apply -f network-policy.yaml

echo "✅ All resources applied successfully!"