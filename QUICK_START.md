# 🚀 ITSM Platform - Quick Start Guide

## ⚡ 30-Second Setup

```bash
# 1. Clone repository
git clone <repository-url>
cd itsm-platform

# 2. Start services
sudo docker compose up -d

# 3. Access application
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Database: localhost:5432
```

---

## 🌐 Access Points

| Service | URL | Status |
|---------|-----|--------|
| Frontend | http://localhost:3000 | ✅ Running |
| Backend API | http://localhost:8080 | ✅ Running |
| Database | localhost:5432 | ✅ Running |
| Health Check | http://localhost:8080/health | ✅ OK |

---

## 🔐 Default Credentials

**Database:**
- Username: `itsm`
- Password: `itsm`
- Database: `itsm`

**Application:**
- Create admin user on first login
- Or use SSO (Keycloak) if configured

---

## 📋 Common Commands

### View Status
```bash
sudo docker compose ps
```

### View Logs
```bash
sudo docker compose logs -f
sudo docker compose logs -f backend
sudo docker compose logs -f frontend
```

### Restart Services
```bash
sudo docker compose restart
sudo docker compose restart backend
```

### Stop Services
```bash
sudo docker compose stop
sudo docker compose down
```

### Database Access
```bash
sudo docker compose exec postgres psql -U itsm -d itsm
```

---

## 🧪 Test API

### Health Check
```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

---

## ⚙️ Configuration

### Environment Variables
Edit `.env` file to configure:
- Database credentials
- JWT secrets
- Email SMTP settings
- Keycloak SSO settings
- Webhook configuration

### Important for Production
Change these in `.env`:
- `JWT_SECRET` → Strong random string
- `JWT_REFRESH_SECRET` → Strong random string
- `WEBHOOK_SECRET` → Strong random string
- `POSTGRES_PASSWORD` → Strong password

---

## 📊 Service Status

### PostgreSQL
- **Status**: ✅ Healthy
- **Port**: 5432
- **Volume**: postgres_data (persistent)

### Backend
- **Status**: ✅ Running
- **Port**: 8080
- **Health**: ✅ Passing

### Frontend
- **Status**: ✅ Running
- **Port**: 3000
- **Health**: ✅ Passing

---

## 🆘 Troubleshooting

### Services won't start
```bash
sudo docker compose build --no-cache
sudo docker compose up -d
```

### Check logs
```bash
sudo docker compose logs
```

### Restart specific service
```bash
sudo docker compose restart backend
```

### Remove everything and start fresh
```bash
sudo docker compose down -v
sudo docker compose up -d
```

---

## 📚 More Information

- **Full Deployment Guide**: See `DEPLOYMENT.md`
- **Deployment Success Report**: See `DEPLOYMENT_SUCCESS.md`
- **Architecture**: See `README.md`

---

## ✅ Deployment Status

- ✅ Docker images built
- ✅ All containers running
- ✅ Health checks passing
- ✅ Database connected
- ✅ API responding
- ✅ Frontend accessible

**Ready to use!** 🎉
