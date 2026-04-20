# ITSM Platform - Docker Deployment Guide

## Prerequisites

- Docker (version 20.10+)
- Docker Compose (version 1.29+)
- Git
- At least 2GB free disk space

## Quick Start

### 1. Clone Repository
```bash
git clone <repository-url>
cd itsm-platform
```

### 2. Configure Environment
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your configuration
nano .env
```

**Important environment variables to update:**
- `JWT_SECRET` - Change to a strong random string
- `JWT_REFRESH_SECRET` - Change to a strong random string
- `WEBHOOK_SECRET` - Change to a strong random string
- `APP_EMAIL_SMTP_PASS` - Your email password
- `KEYCLOAK_CLIENT_SECRET` - Your Keycloak client secret

### 3. Deploy with Docker

#### Option A: Using Deployment Script (Recommended)
```bash
# Make script executable
chmod +x deploy/docker-deploy.sh

# Run deployment
./deploy/docker-deploy.sh
```

#### Option B: Manual Docker Compose
```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check status
docker-compose ps
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5432

### 5. Default Credentials

**Database:**
- Username: `itsm`
- Password: `itsm`
- Database: `itsm`

**Application:**
- Login with credentials created during first setup
- Or use SSO (Keycloak) if configured

## Services

### PostgreSQL
- **Port**: 5432
- **Container**: itsm-postgres
- **Volume**: postgres_data (persistent)
- **Health Check**: Enabled

### Backend (Go)
- **Port**: 8080
- **Container**: itsm-backend
- **Health Check**: Enabled
- **Auto-restart**: Yes

### Frontend (React + Nginx)
- **Port**: 3000
- **Container**: itsm-frontend
- **Health Check**: Enabled
- **Auto-restart**: Yes

## Common Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### Restart Services
```bash
# All services
docker-compose restart

# Specific service
docker-compose restart backend
docker-compose restart frontend
```

### Stop Services
```bash
# Stop without removing volumes
docker-compose stop

# Stop and remove containers (keep volumes)
docker-compose down

# Stop and remove everything including volumes
docker-compose down -v
```

### Execute Commands in Container
```bash
# Backend shell
docker-compose exec backend sh

# Database shell
docker-compose exec postgres psql -U itsm -d itsm

# Frontend shell
docker-compose exec frontend sh
```

### View Database
```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U itsm -d itsm

# List tables
\dt

# Exit
\q
```

## Troubleshooting

### Services Won't Start
```bash
# Check logs
docker-compose logs

# Rebuild images
docker-compose build --no-cache

# Start again
docker-compose up -d
```

### Database Connection Error
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Frontend Shows Blank Page
```bash
# Check frontend logs
docker-compose logs frontend

# Rebuild frontend
docker-compose build --no-cache frontend

# Restart frontend
docker-compose restart frontend
```

### Port Already in Use
```bash
# Find process using port
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :5432  # Database

# Kill process
kill -9 <PID>

# Or change port in docker-compose.yml
```

## Production Deployment

### Security Considerations

1. **Change All Secrets**
   - Update `JWT_SECRET`
   - Update `JWT_REFRESH_SECRET`
   - Update `WEBHOOK_SECRET`
   - Use strong random strings (32+ characters)

2. **Database Security**
   - Change default password from `itsm` to strong password
   - Use environment variables for credentials
   - Enable SSL for database connections

3. **Email Configuration**
   - Use production SMTP server
   - Enable TLS/SSL
   - Use app-specific passwords

4. **Keycloak Configuration**
   - Update `KEYCLOAK_CLIENT_SECRET`
   - Configure proper redirect URIs
   - Enable HTTPS

### Performance Optimization

1. **Database**
   - Increase `shared_buffers` for PostgreSQL
   - Configure connection pooling
   - Enable query caching

2. **Backend**
   - Increase worker processes
   - Configure rate limiting appropriately
   - Enable caching layer (Redis)

3. **Frontend**
   - Enable gzip compression (already configured)
   - Configure CDN for static assets
   - Enable browser caching

### Monitoring

```bash
# Monitor resource usage
docker stats

# Monitor logs in real-time
docker-compose logs -f --tail=100

# Check container health
docker-compose ps
```

### Backup

```bash
# Backup database
docker-compose exec postgres pg_dump -U itsm itsm > backup.sql

# Backup volumes
docker run --rm -v itsm-platform_postgres_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/postgres_backup.tar.gz -C /data .

# Restore database
docker-compose exec -T postgres psql -U itsm itsm < backup.sql
```

## Scaling

### Horizontal Scaling

For production, consider:

1. **Load Balancer** (Nginx, HAProxy)
   - Route traffic to multiple backend instances
   - Session affinity for WebSocket connections

2. **Multiple Backend Instances**
   ```bash
   # Scale backend to 3 instances
   docker-compose up -d --scale backend=3
   ```

3. **Database Replication**
   - Primary-replica setup
   - Read replicas for scaling reads

4. **Caching Layer**
   - Redis for session storage
   - Cache frequently accessed data

## Maintenance

### Regular Tasks

1. **Daily**
   - Monitor logs for errors
   - Check disk space
   - Verify backups

2. **Weekly**
   - Review performance metrics
   - Check for security updates
   - Test backup restoration

3. **Monthly**
   - Update Docker images
   - Review and optimize queries
   - Audit access logs

### Updates

```bash
# Pull latest code
git pull origin main

# Rebuild images
docker-compose build --no-cache

# Restart services
docker-compose up -d
```

## Support

For issues or questions:
1. Check logs: `docker-compose logs`
2. Review troubleshooting section above
3. Check application documentation
4. Contact support team

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Application Documentation](./README.md)
