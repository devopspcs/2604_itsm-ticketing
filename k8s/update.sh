#!/bin/bash
# ===========================================
# Script Update ITSM Ticketing di Minikube
# Jalankan setelah ada perubahan code
# ===========================================
# Usage:
#   ./k8s/update.sh          → update backend + frontend
#   ./k8s/update.sh backend  → update backend saja
#   ./k8s/update.sh frontend → update frontend saja

set -e

COMPONENT=${1:-all}

echo "🔄 Update ITSM Ticketing di Minikube"
echo "======================================"

# Arahkan Docker ke minikube
eval $(minikube docker-env)

if [ "$COMPONENT" = "all" ] || [ "$COMPONENT" = "backend" ]; then
  echo ""
  echo "📦 Rebuilding backend..."
  docker build -t itsm-backend:latest ./backend
  echo "♻️  Restarting backend pods..."
  kubectl rollout restart deployment/backend -n ticketing
  kubectl rollout status deployment/backend -n ticketing
  echo "✅ Backend updated"
fi

if [ "$COMPONENT" = "all" ] || [ "$COMPONENT" = "frontend" ]; then
  echo ""
  echo "📦 Rebuilding frontend..."
  docker build -t itsm-frontend:latest ./frontend
  echo "♻️  Restarting frontend pods..."
  kubectl rollout restart deployment/frontend -n ticketing
  kubectl rollout status deployment/frontend -n ticketing
  echo "✅ Frontend updated"
fi

echo ""
echo "======================================"
echo "🎉 Update selesai!"
echo ""
echo "📋 Status pods:"
kubectl get pods -n ticketing
echo ""
echo "📍 Akses: http://itsm.local"
