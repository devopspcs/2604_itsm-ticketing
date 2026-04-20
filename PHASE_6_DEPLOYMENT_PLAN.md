# Phase 6: Deployment - Jira-Only Project Board

**Status**: READY FOR EXECUTION  
**Date**: April 19, 2026

---

## 📋 Deployment Plan

Comprehensive deployment plan untuk meluncurkan Jira-only project board ke production.

---

## 🎯 Deployment Strategy

### Deployment Type: Blue-Green Deployment

**Advantages:**
- Zero downtime
- Easy rollback if issues occur
- Can test in production environment
- Smooth transition for users

**Process:**
1. Deploy new version (Green) alongside current version (Blue)
2. Test Green environment
3. Switch traffic to Green
4. Keep Blue as fallback
5. Monitor Green environment
6. Remove Blue after stable

---

## 📋 Pre-Deployment Checklist

### 1. Code Review
- [ ] All code changes reviewed
- [ ] No security vulnerabilities
- [ ] No performance issues
- [ ] All tests passing

### 2. Database Backup
- [ ] Full database backup created
- [ ] Backup verified
- [ ] Backup stored securely
- [ ] Rollback plan documented

### 3. Environment Preparation
- [ ] Production environment ready
- [ ] All dependencies installed
- [ ] Environment variables configured
- [ ] SSL certificates valid

### 4. Documentation
- [ ] Deployment guide prepared
- [ ] Rollback guide prepared
- [ ] User documentation updated
- [ ] Support team trained

### 5. Communication
- [ ] Stakeholders notified
- [ ] Users notified about maintenance window
- [ ] Support team on standby
- [ ] Monitoring team ready

---

## 🚀 Deployment Steps

### Step 1: Database Migration

**Objective:** Migrate data from old project board to Jira-only schema

**Commands:**
```bash
# Connect to production database
psql -h prod-db.example.com -U itsm -d itsm

# Run migration script
\i backend/migrations/002_migrate_to_jira_only.sql

# Verify migration
SELECT COUNT(*) FROM workflows;
SELECT COUNT(*) FROM workflow_statuses;
SELECT COUNT(*) FROM issue_type_schemes;
```

**Verification:**
```sql
-- Check migration status
SELECT 
  'Projects' as entity,
  COUNT(*) as count
FROM projects
UNION ALL
SELECT 'Workflows', COUNT(*) FROM workflows
UNION ALL
SELECT 'Workflow Statuses', COUNT(*) FROM workflow_statuses
UNION ALL
SELECT 'Issue Type Schemes', COUNT(*) FROM issue_type_schemes
UNION ALL
SELECT 'Records with Issue Type', COUNT(*) FROM project_records WHERE issue_type_id IS NOT NULL
UNION ALL
SELECT 'Records with Status', COUNT(*) FROM project_records WHERE status IS NOT NULL;
```

**Rollback:**
```bash
# If migration fails, rollback using backup
psql -h prod-db.example.com -U itsm -d itsm

# Drop new tables
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS workflow_statuses;
DROP TABLE IF EXISTS workflow_transitions;
DROP TABLE IF EXISTS issue_type_schemes;
DROP TABLE IF EXISTS issue_type_scheme_items;

# Restore original data
ALTER TABLE project_records
DROP COLUMN IF EXISTS issue_type_id,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS parent_record_id;
```

---

### Step 2: Backend Deployment

**Objective:** Deploy updated backend to production

**Commands:**
```bash
# Build backend
cd backend
go build -o itsm-backend

# Create backup of current backend
cp /opt/itsm/backend /opt/itsm/backend.backup

# Deploy new backend
cp itsm-backend /opt/itsm/backend

# Restart backend service
systemctl restart itsm-backend

# Verify backend is running
curl http://localhost:8080/health
```

**Verification:**
```bash
# Check backend logs
tail -f /var/log/itsm/backend.log

# Check API endpoints
curl -X GET http://localhost:8080/api/sprints
curl -X GET http://localhost:8080/api/workflows
curl -X GET http://localhost:8080/api/records
```

**Rollback:**
```bash
# Restore previous backend
cp /opt/itsm/backend.backup /opt/itsm/backend

# Restart backend service
systemctl restart itsm-backend

# Verify backend is running
curl http://localhost:8080/health
```

---

### Step 3: Frontend Deployment

**Objective:** Deploy updated frontend to production

**Commands:**
```bash
# Build frontend
cd frontend
npm run build

# Create backup of current frontend
cp -r /opt/itsm/frontend/dist /opt/itsm/frontend/dist.backup

# Deploy new frontend
cp -r dist/* /opt/itsm/frontend/dist/

# Restart frontend service (if using Node.js)
systemctl restart itsm-frontend

# Or if using Nginx, reload configuration
nginx -s reload
```

**Verification:**
```bash
# Check frontend is accessible
curl http://localhost:3000

# Check frontend logs
tail -f /var/log/itsm/frontend.log

# Check Nginx logs
tail -f /var/log/nginx/access.log
```

**Rollback:**
```bash
# Restore previous frontend
rm -rf /opt/itsm/frontend/dist
cp -r /opt/itsm/frontend/dist.backup /opt/itsm/frontend/dist

# Restart frontend service
systemctl restart itsm-frontend

# Or reload Nginx
nginx -s reload
```

---

### Step 4: Docker Deployment (Alternative)

**Objective:** Deploy using Docker containers

**Commands:**
```bash
# Build Docker images
docker build -t itsm-backend:latest backend/
docker build -t itsm-frontend:latest frontend/

# Tag images for production
docker tag itsm-backend:latest itsm-backend:prod-v2.0
docker tag itsm-frontend:latest itsm-frontend:prod-v2.0

# Push to registry
docker push itsm-backend:prod-v2.0
docker push itsm-frontend:prod-v2.0

# Update docker-compose.prod.yml
# Change image tags to prod-v2.0

# Deploy using docker-compose
docker-compose -f docker-compose.prod.yml up -d

# Verify containers are running
docker-compose -f docker-compose.prod.yml ps
```

**Verification:**
```bash
# Check container logs
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f frontend

# Check API endpoints
curl http://localhost:8080/health
curl http://localhost:3000
```

**Rollback:**
```bash
# Revert to previous version
docker-compose -f docker-compose.prod.yml down

# Update docker-compose.prod.yml to use previous version
# Change image tags back to prod-v1.0

# Deploy previous version
docker-compose -f docker-compose.prod.yml up -d

# Verify containers are running
docker-compose -f docker-compose.prod.yml ps
```

---

### Step 5: Health Check

**Objective:** Verify that deployment is successful

**Checks:**
```bash
# 1. Backend health check
curl -X GET http://localhost:8080/health

# 2. Frontend health check
curl -X GET http://localhost:3000

# 3. Database connection check
curl -X GET http://localhost:8080/api/projects

# 4. API endpoints check
curl -X GET http://localhost:8080/api/sprints
curl -X GET http://localhost:8080/api/workflows
curl -X GET http://localhost:8080/api/records

# 5. Frontend page load check
curl -X GET http://localhost:3000/projects
```

**Expected Results:**
```
✅ Backend health: 200 OK
✅ Frontend health: 200 OK
✅ Database connection: 200 OK
✅ API endpoints: 200 OK
✅ Frontend pages: 200 OK
```

---

### Step 6: Smoke Testing

**Objective:** Quick verification that main features work

**Test Cases:**
```bash
# 1. Create a new sprint
curl -X POST http://localhost:8080/api/sprints \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "...",
    "name": "Sprint 1",
    "goal": "Test sprint"
  }'

# 2. Get sprints
curl -X GET http://localhost:8080/api/sprints

# 3. Get workflows
curl -X GET http://localhost:8080/api/workflows

# 4. Get records
curl -X GET http://localhost:8080/api/records

# 5. Create a new record
curl -X POST http://localhost:8080/api/records \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "...",
    "title": "Test record",
    "status": "To Do"
  }'
```

---

### Step 7: Monitoring

**Objective:** Monitor production environment after deployment

**Metrics to Monitor:**
- [ ] CPU usage
- [ ] Memory usage
- [ ] Disk usage
- [ ] Network traffic
- [ ] API response time
- [ ] Error rate
- [ ] User activity

**Monitoring Tools:**
```bash
# Using Prometheus
curl http://localhost:9090/api/v1/query?query=up

# Using Grafana
# Access at http://localhost:3000/grafana

# Using ELK Stack
# Access at http://localhost:5601

# Using New Relic
# Access at https://newrelic.com
```

**Alert Thresholds:**
- CPU usage > 80% → Alert
- Memory usage > 85% → Alert
- Disk usage > 90% → Alert
- API response time > 1000ms → Alert
- Error rate > 1% → Alert

---

## 📊 Deployment Timeline

| Step | Duration | Status |
|------|----------|--------|
| Pre-Deployment Checklist | 30 min | ⏳ |
| Database Migration | 15 min | ⏳ |
| Backend Deployment | 10 min | ⏳ |
| Frontend Deployment | 10 min | ⏳ |
| Health Check | 5 min | ⏳ |
| Smoke Testing | 10 min | ⏳ |
| Monitoring Setup | 10 min | ⏳ |
| **Total** | **90 min** | ⏳ |

---

## 🔄 Rollback Plan

### Scenario 1: Database Migration Failed

**Steps:**
1. Stop backend service
2. Restore database from backup
3. Restore backend to previous version
4. Restart backend service
5. Verify system is working

**Commands:**
```bash
# Stop backend
systemctl stop itsm-backend

# Restore database
psql -h prod-db.example.com -U itsm -d itsm < backup.sql

# Restore backend
cp /opt/itsm/backend.backup /opt/itsm/backend

# Restart backend
systemctl start itsm-backend

# Verify
curl http://localhost:8080/health
```

### Scenario 2: Backend Deployment Failed

**Steps:**
1. Restore backend to previous version
2. Restart backend service
3. Verify system is working

**Commands:**
```bash
# Restore backend
cp /opt/itsm/backend.backup /opt/itsm/backend

# Restart backend
systemctl restart itsm-backend

# Verify
curl http://localhost:8080/health
```

### Scenario 3: Frontend Deployment Failed

**Steps:**
1. Restore frontend to previous version
2. Restart frontend service
3. Verify system is working

**Commands:**
```bash
# Restore frontend
rm -rf /opt/itsm/frontend/dist
cp -r /opt/itsm/frontend/dist.backup /opt/itsm/frontend/dist

# Restart frontend
systemctl restart itsm-frontend

# Verify
curl http://localhost:3000
```

### Scenario 4: Critical Issues in Production

**Steps:**
1. Immediately rollback all changes
2. Restore database from backup
3. Restore backend to previous version
4. Restore frontend to previous version
5. Notify stakeholders
6. Investigate issues
7. Plan re-deployment

**Commands:**
```bash
# Rollback database
psql -h prod-db.example.com -U itsm -d itsm < backup.sql

# Rollback backend
cp /opt/itsm/backend.backup /opt/itsm/backend
systemctl restart itsm-backend

# Rollback frontend
rm -rf /opt/itsm/frontend/dist
cp -r /opt/itsm/frontend/dist.backup /opt/itsm/frontend/dist
systemctl restart itsm-frontend

# Verify
curl http://localhost:8080/health
curl http://localhost:3000
```

---

## 📝 Deployment Checklist

### Pre-Deployment
- [ ] Code review completed
- [ ] All tests passing
- [ ] Database backup created
- [ ] Environment variables configured
- [ ] SSL certificates valid
- [ ] Stakeholders notified
- [ ] Support team ready

### During Deployment
- [ ] Database migration executed
- [ ] Backend deployed
- [ ] Frontend deployed
- [ ] Health checks passed
- [ ] Smoke tests passed
- [ ] Monitoring setup

### Post-Deployment
- [ ] Monitor system for 24 hours
- [ ] Check error logs
- [ ] Verify user feedback
- [ ] Document any issues
- [ ] Update documentation
- [ ] Notify stakeholders

---

## 📊 Deployment Report Template

```markdown
# Deployment Report - Phase 6

**Date**: [Date]
**Deployed By**: [Name]
**Status**: [SUCCESS/FAILED]

## Summary
- Deployment Start: [Time]
- Deployment End: [Time]
- Total Duration: [Duration]
- Rollback Required: [YES/NO]

## Pre-Deployment
- [ ] Code review: PASS/FAIL
- [ ] Tests passing: PASS/FAIL
- [ ] Database backup: PASS/FAIL
- [ ] Environment ready: PASS/FAIL

## Deployment Steps
- [ ] Database migration: PASS/FAIL
- [ ] Backend deployment: PASS/FAIL
- [ ] Frontend deployment: PASS/FAIL
- [ ] Health checks: PASS/FAIL
- [ ] Smoke tests: PASS/FAIL

## Post-Deployment
- [ ] Monitoring setup: PASS/FAIL
- [ ] Error logs checked: PASS/FAIL
- [ ] User feedback: [Feedback]
- [ ] Issues found: [Issues]

## Issues Found
1. [Issue 1]
2. [Issue 2]
3. [Issue 3]

## Resolutions
1. [Resolution 1]
2. [Resolution 2]
3. [Resolution 3]

## Sign-off
- Deployed By: [Name]
- Approved By: [Name]
- Date: [Date]
- Status: [APPROVED/REJECTED]
```

---

## 🎯 Post-Deployment Tasks

### 1. User Communication
- [ ] Send deployment notification to users
- [ ] Provide user guide for new features
- [ ] Set up support channel for issues
- [ ] Schedule training sessions

### 2. Documentation Update
- [ ] Update user documentation
- [ ] Update API documentation
- [ ] Update deployment guide
- [ ] Update troubleshooting guide

### 3. Monitoring
- [ ] Monitor system for 24 hours
- [ ] Check error logs daily
- [ ] Monitor performance metrics
- [ ] Collect user feedback

### 4. Optimization
- [ ] Analyze performance data
- [ ] Identify bottlenecks
- [ ] Optimize queries
- [ ] Optimize frontend

---

## ✅ Deployment Success Criteria

- ✅ Database migration successful
- ✅ Backend deployed without errors
- ✅ Frontend deployed without errors
- ✅ All health checks passed
- ✅ All smoke tests passed
- ✅ No critical errors in logs
- ✅ System stable for 24 hours
- ✅ Users can access all features

---

## 🎯 Next Steps

After Phase 6 (Deployment) is complete:
1. Monitor production environment
2. Collect user feedback
3. Fix any issues found
4. Plan next features
5. Schedule next deployment

---

## 📞 Support Contacts

**During Deployment:**
- Backend Team: [Contact]
- Frontend Team: [Contact]
- DevOps Team: [Contact]
- Database Team: [Contact]

**After Deployment:**
- Support Team: [Contact]
- Product Team: [Contact]
- Management: [Contact]

---

**Phase 6 Status**: ⏳ READY FOR EXECUTION

**Apakah Anda ingin saya jalankan deployment?** 🚀

