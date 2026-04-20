# Docker Deployment Guide - April 19, 2026

## Status: ✅ Ready for Deployment

Semua code sudah compile dan siap untuk di-deploy ke production. Docker images dapat dibangun dan dijalankan di server production.

## Prerequisites

Pastikan server production memiliki:
- Docker 20.10+
- Docker Compose 2.0+
- Git
- 4GB RAM minimum
- 20GB disk space

## Deployment Steps

### 1. Clone Repository
```bash
git clone <repository-url>
cd ticketing-system
```

### 2. Setup Environment Variables

**Development:**
```bash
cp .env.example .env
# Edit .env dengan konfigurasi lokal
```

**Production:**
```bash
cp .env.example .env.prod
# Edit .env.prod dengan konfigurasi production
```

### 3. Build Docker Images

```bash
# Build semua images
docker compose build

# Atau build individual images
docker build -t itsm-backend:latest ./backend
docker build -t itsm-frontend:latest ./frontend
```

### 4. Start Services

**Development:**
```bash
docker compose up -d
```

**Production:**
```bash
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

### 5. Verify Deployment

```bash
# Check container status
docker compose ps

# Check backend health
curl http://localhost:8080/health

# Check frontend
curl http://localhost:3000

# View logs
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

## Services

### PostgreSQL (Port 5432)
- Database: itsm
- User: itsm
- Password: itsm (change in production!)
- Volume: postgres_data

### Backend (Port 8080)
- Go API server
- Health endpoint: /health
- API endpoints: /api/v1/*

### Frontend (Port 3000)
- React web application
- Nginx reverse proxy
- Health endpoint: /health
- API proxy: /api/* → backend:8080

## Docker Compose Services

```yaml
services:
  postgres:
    - PostgreSQL 16 Alpine
    - Port: 5432
    - Health check: pg_isready
    
  backend:
    - Go application
    - Port: 8080
    - Health check: wget /health
    - Depends on: postgres
    
  frontend:
    - React + Nginx
    - Port: 3000
    - Health check: wget /health
    - Depends on: backend
```

## Environment Variables

### Backend
```
DATABASE_URL=postgres://itsm:itsm@postgres:5432/itsm?sslmode=disable
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret
WEBHOOK_SECRET=your-webhook-secret
PORT=8080
BASE_URL=http://localhost:3000
KEYCLOAK_URL=https://your-keycloak-server
KEYCLOAK_REALM=your-realm
KEYCLOAK_CLIENT_ID=your-client-id
KEYCLOAK_CLIENT_SECRET=your-client-secret
```

### Frontend
```
VITE_API_URL=http://localhost:8080/api/v1
```

## Troubleshooting

### Container won't start
```bash
# Check logs
docker compose logs backend
docker compose logs frontend

# Rebuild without cache
docker compose build --no-cache
```

### Database connection error
```bash
# Check PostgreSQL is running
docker compose ps postgres

# Check database exists
docker compose exec postgres psql -U itsm -d itsm -c "\dt"
```

### Port already in use
```bash
# Change ports in docker-compose.yml
# Or kill existing process
lsof -i :8080
kill -9 <PID>
```

### Memory issues
```bash
# Increase Docker memory limit
# Edit Docker daemon config or use:
docker run --memory=2g ...
```

## Backup & Recovery

### Backup Database
```bash
docker compose exec postgres pg_dump -U itsm itsm > backup.sql
```

### Restore Database
```bash
docker compose exec -T postgres psql -U itsm itsm < backup.sql
```

### Backup Volumes
```bash
docker run --rm -v postgres_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/postgres_backup.tar.gz -C /data .
```

## Monitoring

### View Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

### Check Resource Usage
```bash
docker stats
```

### Check Container Health
```bash
docker compose ps
```

## Scaling

### Horizontal Scaling (Multiple Backend Instances)
```bash
# Scale backend to 3 instances
docker compose up -d --scale backend=3
```

### Load Balancing
Use Nginx or HAProxy in front of multiple backend instances.

## Security Considerations

1. **Change Default Passwords**
   - PostgreSQL password
   - JWT secrets
   - Webhook secrets

2. **Use HTTPS**
   - Configure SSL/TLS certificates
   - Use Apache/Nginx reverse proxy

3. **Environment Variables**
   - Never commit .env files
   - Use secrets management (Docker Secrets, Vault)

4. **Database**
   - Regular backups
   - Restricted access
   - Encrypted connections

5. **API Security**
   - Rate limiting (configured)
   - CORS configuration
   - Input validation

## Production Checklist

- [ ] All environment variables configured
- [ ] Database backups scheduled
- [ ] SSL/TLS certificates installed
- [ ] Reverse proxy configured (Apache/Nginx)
- [ ] Monitoring and logging setup
- [ ] Backup and recovery tested
- [ ] Security audit completed
- [ ] Performance testing done
- [ ] User acceptance testing passed
- [ ] Deployment documentation reviewed

## Rollback Procedure

```bash
# Stop current deployment
docker compose down

# Restore previous version
git checkout <previous-commit>

# Rebuild and restart
docker compose build
docker compose up -d

# Restore database if needed
docker compose exec -T postgres psql -U itsm itsm < backup.sql
```

## Performance Optimization

### Frontend
- Gzip compression enabled
- Static asset caching (1 year)
- Code splitting
- Lazy loading

### Backend
- Connection pooling
- Query optimization
- Caching layer
- Rate limiting

### Database
- Indexes on frequently queried columns
- Query optimization
- Connection pooling
- Regular maintenance

## Maintenance

### Regular Tasks
- Monitor disk space
- Check logs for errors
- Update dependencies
- Backup database
- Review performance metrics

### Monthly
- Security updates
- Dependency updates
- Performance review
- Backup verification

### Quarterly
- Full security audit
- Performance optimization
- Capacity planning
- Disaster recovery drill

---

**Deployment Guide Version:** 1.0
**Last Updated:** April 19, 2026
**Status:** Ready for Production
