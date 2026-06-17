#!/bin/bash
# ===========================================
# Script Deploy ITSM ke Production (K3s/Minikube)
# ===========================================

set -e

echo "🚀 Deploy ITSM Ticketing ke Production Kubernetes"
echo "==================================================="

# Detect environment (K3s or Minikube)
if command -v k3s &> /dev/null; then
  echo "📋 Detected: K3s"
  RUNTIME="k3s"
elif command -v minikube &> /dev/null; then
  echo "📋 Detected: Minikube"
  RUNTIME="minikube"
else
  echo "❌ Tidak ditemukan K3s atau Minikube. Install salah satu dulu."
  exit 1
fi

# Build images
echo ""
echo "📦 Step 1: Build Docker images..."
docker build -t itsm-backend:latest ./backend
docker build -t itsm-frontend:latest ./frontend

if [ "$RUNTIME" = "k3s" ]; then
  echo "   Importing images ke K3s..."
  docker save itsm-backend:latest | sudo k3s ctr images import -
  docker save itsm-frontend:latest | sudo k3s ctr images import -
elif [ "$RUNTIME" = "minikube" ]; then
  echo "   Loading images ke Minikube..."
  minikube image load itsm-backend:latest
  minikube image load itsm-frontend:latest
fi
echo "✅ Images ready"

# Apply manifests
echo ""
echo "📦 Step 2: Apply Kubernetes manifests..."
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/production/secrets.yaml
kubectl apply -f k8s/production/configmap.yaml
kubectl apply -f k8s/postgres.yaml

echo "   Waiting for PostgreSQL..."
kubectl wait --for=condition=ready pod -l app=postgres -n ticketing --timeout=120s

kubectl apply -f k8s/backend.yaml
echo "   Waiting for Backend..."
kubectl wait --for=condition=ready pod -l app=backend -n ticketing --timeout=120s

kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/production/ingress.yaml

echo ""
echo "==================================================="
echo "🎉 Production deployment selesai!"
echo "==================================================="
echo ""
echo "📋 Status:"
kubectl get pods -n ticketing
echo ""
echo "📋 Services:"
kubectl get svc -n ticketing
echo ""
echo "📋 Ingress:"
kubectl get ingress -n ticketing
echo ""

if [ "$RUNTIME" = "minikube" ]; then
  echo "⚠️  Minikube: Jalankan 'minikube tunnel' di terminal lain untuk expose port 80/443"
fi
