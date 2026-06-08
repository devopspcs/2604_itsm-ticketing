# Kubernetes Deployment Guide (Minikube)

Panduan deploy aplikasi ITSM Ticketing di Kubernetes menggunakan Minikube untuk belajar.

## Prerequisites

1. **Docker** - [Install Docker](https://docs.docker.com/get-docker/)
2. **Minikube** - [Install Minikube](https://minikube.sigs.k8s.io/docs/start/)
3. **kubectl** - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)

### Install Minikube (Linux)

```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

### Install kubectl (Linux)

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install kubectl /usr/local/bin/kubectl
```

## Quick Start

```bash
# 1. Start minikube
minikube start --driver=docker --memory=4096 --cpus=2

# 2. Jalankan deploy script
chmod +x k8s/deploy.sh
./k8s/deploy.sh
```

## Manual Deployment (Step by Step)

Kalau mau belajar satu-satu:

### Step 1: Start Minikube

```bash
minikube start --driver=docker --memory=4096 --cpus=2
```

### Step 2: Enable Addons

```bash
minikube addons enable ingress
minikube addons enable storage-provisioner
```

### Step 3: Build Images di Minikube Docker

```bash
# Arahkan Docker CLI ke daemon minikube
eval $(minikube docker-env)

# Build images
docker build -t itsm-backend:latest ./backend
docker build -t itsm-frontend:latest ./frontend
```

> **Penting:** `eval $(minikube docker-env)` membuat Docker CLI kamu mengarah ke Docker daemon di dalam minikube. Jadi image yang di-build langsung tersedia di minikube tanpa perlu push ke registry.

### Step 4: Apply Manifests

```bash
# Buat namespace
kubectl apply -f k8s/namespace.yaml

# Buat secrets dan config
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml

# Deploy PostgreSQL (tunggu ready)
kubectl apply -f k8s/postgres.yaml
kubectl wait --for=condition=ready pod -l app=postgres -n ticketing --timeout=120s

# Deploy Backend (tunggu ready)
kubectl apply -f k8s/backend.yaml
kubectl wait --for=condition=ready pod -l app=backend -n ticketing --timeout=120s

# Deploy Frontend
kubectl apply -f k8s/frontend.yaml

# Setup Ingress
kubectl apply -f k8s/ingress.yaml
```

### Step 5: Setup Domain Lokal

```bash
# Tambah entry di /etc/hosts
echo "$(minikube ip) itsm.local" | sudo tee -a /etc/hosts
```

### Step 6: Akses Aplikasi

```bash
# Buka browser ke:
# http://itsm.local

# Jika ingress tidak jalan, gunakan tunnel:
minikube tunnel
```

## Perintah Berguna

```bash
# Lihat semua resources di namespace ticketing
kubectl get all -n ticketing

# Lihat pods
kubectl get pods -n ticketing

# Lihat logs backend
kubectl logs -f deployment/backend -n ticketing

# Lihat logs frontend  
kubectl logs -f deployment/frontend -n ticketing

# Masuk ke pod (debug)
kubectl exec -it deployment/backend -n ticketing -- sh

# Lihat events (untuk debugging)
kubectl get events -n ticketing --sort-by='.lastTimestamp'

# Scale deployment
kubectl scale deployment backend --replicas=3 -n ticketing

# Restart deployment
kubectl rollout restart deployment/backend -n ticketing

# Delete semua
kubectl delete namespace ticketing
```

## Arsitektur di Kubernetes

```
                    ┌──────────────────────────────────┐
                    │          Ingress (nginx)          │
                    │         host: itsm.local          │
                    └───────────┬───────────┬───────────┘
                                │           │
                    /api,/health,/sso       /
                                │           │
                    ┌───────────▼──┐  ┌─────▼──────────┐
                    │   Backend    │  │    Frontend     │
                    │  (Go:8080)   │  │  (Nginx:80)    │
                    │  replicas: 2 │  │  replicas: 2   │
                    └───────┬──────┘  └────────────────┘
                            │
                    ┌───────▼──────┐
                    │  PostgreSQL  │
                    │ (StatefulSet)│
                    │  replicas: 1 │
                    └──────────────┘
```

## Perbedaan dengan Docker Compose

| Aspek | Docker Compose | Kubernetes |
|-------|---------------|------------|
| Scaling | Manual (docker compose up --scale) | `kubectl scale` otomatis |
| Self-healing | Restart policy saja | Pod dihapus, langsung dibuat baru |
| Networking | Bridge network | Service discovery via DNS |
| Storage | Volume mount | PersistentVolumeClaim |
| Config | .env file | ConfigMap + Secret |
| Load Balancing | Tidak ada built-in | Service + Ingress |
| Rolling Update | Stop then start | Zero-downtime rolling update |

## Troubleshooting

### Pod stuck di "Pending"
```bash
kubectl describe pod <pod-name> -n ticketing
# Biasanya masalah resource atau PVC
```

### Pod CrashLoopBackOff
```bash
kubectl logs <pod-name> -n ticketing --previous
# Lihat error dari container sebelumnya
```

### Image pull error
```bash
# Pastikan sudah eval $(minikube docker-env) sebelum build
# Dan imagePullPolicy: IfNotPresent di manifest
```

### Ingress tidak bisa diakses
```bash
# Cek ingress controller running
kubectl get pods -n ingress-nginx

# Alternatif: gunakan minikube tunnel
minikube tunnel

# Atau port-forward langsung
kubectl port-forward service/frontend 3000:80 -n ticketing
kubectl port-forward service/backend 8080:8080 -n ticketing
```
