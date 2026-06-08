#!/bin/bash
# ===========================================
# Script Deploy ITSM Ticketing ke Minikube
# ===========================================

set -e

echo "🚀 Deploy ITSM Ticketing ke Kubernetes (Minikube)"
echo "=================================================="

# 1. Pastikan minikube running
echo ""
echo "📋 Step 1: Cek Minikube status..."
if ! minikube status | grep -q "Running"; then
  echo "❌ Minikube belum jalan. Jalankan dulu:"
  echo "   minikube start --driver=docker --memory=4096 --cpus=2"
  exit 1
fi
echo "✅ Minikube running"

# 2. Enable addons yang diperlukan
echo ""
echo "📋 Step 2: Enable addons..."
minikube addons enable ingress
minikube addons enable storage-provisioner
echo "✅ Addons enabled"

# 3. Gunakan Docker daemon minikube (supaya image lokal bisa dipakai)
echo ""
echo "📋 Step 3: Setup Docker environment minikube..."
eval $(minikube docker-env)
echo "✅ Docker env pointing ke minikube"

# 4. Build images
echo ""
echo "📋 Step 4: Build Docker images..."
echo "   Building backend..."
docker build -t itsm-backend:latest ./backend
echo "   Building frontend..."
docker build -t itsm-frontend:latest ./frontend
echo "✅ Images built"

# 5. Apply Kubernetes manifests (urutan penting!)
echo ""
echo "📋 Step 5: Apply Kubernetes manifests..."
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/postgres.yaml
echo "   Waiting for PostgreSQL ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n ticketing --timeout=120s
kubectl apply -f k8s/backend.yaml
echo "   Waiting for Backend ready..."
kubectl wait --for=condition=ready pod -l app=backend -n ticketing --timeout=120s
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/ingress.yaml
echo "✅ All manifests applied"

# 6. Tambah host entry
echo ""
echo "📋 Step 6: Setup /etc/hosts..."
MINIKUBE_IP=$(minikube ip)
if ! grep -q "itsm.local" /etc/hosts; then
  echo "   Menambah entry di /etc/hosts (perlu sudo)..."
  echo "$MINIKUBE_IP itsm.local" | sudo tee -a /etc/hosts
else
  echo "   Entry itsm.local sudah ada di /etc/hosts"
fi
echo "✅ Host configured"

# 7. Summary
echo ""
echo "=================================================="
echo "🎉 Deployment selesai!"
echo "=================================================="
echo ""
echo "📍 Akses aplikasi: http://itsm.local"
echo "📍 Minikube IP: $MINIKUBE_IP"
echo ""
echo "📋 Cek status pods:"
echo "   kubectl get pods -n ticketing"
echo ""
echo "📋 Lihat logs:"
echo "   kubectl logs -f deployment/backend -n ticketing"
echo "   kubectl logs -f deployment/frontend -n ticketing"
echo ""
echo "📋 Jika ingress belum jalan, bisa pakai tunnel:"
echo "   minikube tunnel"
echo ""
