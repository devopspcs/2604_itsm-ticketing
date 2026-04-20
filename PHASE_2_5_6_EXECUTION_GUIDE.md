# Phase 2, 5, 6 - Execution Guide

**Status**: READY FOR EXECUTION  
**Date**: April 19, 2026  
**Environment**: Production/Development with Database Access

---

## 📋 Pre-Execution Checklist

Sebelum menjalankan Phase 2, 5, dan 6, pastikan:

- [ ] Database backup sudah dibuat
- [ ] Stakeholders sudah diberitahu
- [ ] Support team siap standby
- [ ] Maintenance window sudah dijadwalkan
- [ ] Rollback plan sudah disiapkan
- [ ] Environment variables sudah dikonfigurasi

---

## 🚀 PHASE 2: DATABASE MIGRATION

### Objective
Migrasi data dari old project board schema ke Jira-compatible schema

### Duration
5-10 minutes

### Prerequisites
- PostgreSQL client (psql) installed
- Database credentials ready
- Backup created

### Execution Steps

#### Step 1: Backup Database (CRITICAL!)

**Using pg_dump:**
```bash
# Create backup
pg_dump -h localhost -U itsm -d itsm > backup_$(date +%Y%m%d_%H%M%S).sql

# Verify backup
ls -lh backup_*.sql
```

**Using Docker:**
```bash
# If using Docker Compose
docker-compose exec postgres pg_dump -U itsm itsm > backup_$(date +%Y%m%d_%H%M%S).sql
```

**Verify backup size:**
```bash
# Backup should be at least 1MB
ls -lh backup_*.sql
```

#### Step 2: Run Migration Script

**Option A: Using psql directly**
```bash
# Connect to database and run migration
psql -h localhost -U itsm -d itsm -f backend/migrations/002_migrate_to_jira_only.sql
```

**Option B: Using Docker Compose**
```bash
# Copy migration script to container
docker cp backend/migrations/002_migrate_to_jira_only.sql itsm-postgres:/tmp/

# Run migration in container
docker-compose exec postgres psql -U itsm -d itsm -f /tmp/002_migrate_to_jira_only.sql
```

**Option C: Using Docker directly**
```bash
# Run migration
docker exec itsm-postgres psql -U itsm -d itsm -f /migrations/002_migrate_to_jira_only.sql
```

#### Step 3: Verify Migration

**Check workflows created:**
```bash
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) as workflows FROM workflows;"
```

**Expected output:**
```
 workflows
-----------
    [number of projects]
(1 row)
```

**Check workflow statuses:**
```bash
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) as statuses FROM workflow_statuses;"
```

**Check issue type schemes:**
```bash
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) as schemes FROM issue_type_schemes;"
```

**Check records updated:**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  COUNT(*) as total_records,
  COUNT(CASE WHEN issue_type_id IS NOT NULL THEN 1 END) as with_issue_type,
  COUNT(CASE WHEN status IS NOT NULL THEN 1 END) as with_status
FROM project_records;
"
```

**Expected output:**
```
 total_records | with_issue_type | with_status
---------------+-----------------+-------------
    [number]   |    [number]     |  [number]
(1 row)
```

#### Step 4: Full Verification Query

```bash
psql -h localhost -U itsm -d itsm << EOF
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
SELECT 'Issue Type Scheme Items', COUNT(*) FROM issue_type_scheme_items
UNION ALL
SELECT 'Records with Issue Type', COUNT(*) FROM project_records WHERE issue_type_id IS NOT NULL
UNION ALL
SELECT 'Records with Status', COUNT(*) FROM project_records WHERE status IS NOT NULL
ORDER BY entity;
EOF
```

### Troubleshooting

#### Error: "relation does not exist"
```bash
# Check if tables exist
psql -h localhost -U itsm -d itsm -c "\dt"

# If tables don't exist, run backend migrations first
go run cmd/migrate/main.go up
```

#### Error: "permission denied"
```bash
# Check user permissions
psql -h localhost -U itsm -d itsm -c "SELECT current_user;"

# Grant permissions if needed
psql -h localhost -U postgres -d itsm -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO itsm;"
```

#### Error: "connection refused"
```bash
# Check if database is running
psql -h localhost -U itsm -d itsm -c "SELECT 1;"

# If not running, start database
docker-compose up -d postgres
```

### Rollback (If Needed)

**If migration fails, rollback immediately:**

```bash
# Restore from backup
psql -h localhost -U itsm -d itsm < backup_YYYYMMDD_HHMMSS.sql

# Verify rollback
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) FROM workflows;"
# Should return 0 if rollback successful
```

---

## 🧪 PHASE 5: TESTING

### Objective
Test semua fitur untuk memastikan migrasi berhasil

### Duration
2-4 hours

### Test Categories

#### 1. Data Migration Tests (7 tests)

**Test 1.1: Backup Verification**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  (SELECT COUNT(*) FROM projects) as projects,
  (SELECT COUNT(*) FROM projects_backup) as projects_backup,
  (SELECT COUNT(*) FROM project_records) as records,
  (SELECT COUNT(*) FROM project_records_backup) as records_backup;
"
```

**Test 1.2: Workflow Creation**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  COUNT(DISTINCT p.id) as projects,
  COUNT(DISTINCT w.id) as workflows,
  CASE WHEN COUNT(DISTINCT p.id) = COUNT(DISTINCT w.id) THEN 'PASS' ELSE 'FAIL' END as status
FROM projects p
LEFT JOIN workflows w ON p.id = w.project_id;
"
```

**Test 1.3: Workflow Statuses**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  COUNT(DISTINCT pc.id) as columns,
  COUNT(DISTINCT ws.id) as statuses,
  CASE WHEN COUNT(DISTINCT pc.id) = COUNT(DISTINCT ws.id) THEN 'PASS' ELSE 'FAIL' END as status
FROM project_columns pc
LEFT JOIN workflow_statuses ws ON pc.name = ws.status_name;
"
```

**Test 1.4: Issue Type Schemes**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  COUNT(DISTINCT p.id) as projects,
  COUNT(DISTINCT iss.id) as schemes,
  CASE WHEN COUNT(DISTINCT p.id) = COUNT(DISTINCT iss.id) THEN 'PASS' ELSE 'FAIL' END as status
FROM projects p
LEFT JOIN issue_type_schemes iss ON p.id = iss.project_id;
"
```

**Test 1.5: Records Migration**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  COUNT(*) as total_records,
  COUNT(CASE WHEN issue_type_id IS NOT NULL THEN 1 END) as with_issue_type,
  COUNT(CASE WHEN status IS NOT NULL THEN 1 END) as with_status,
  CASE WHEN COUNT(*) = COUNT(CASE WHEN issue_type_id IS NOT NULL THEN 1 END) THEN 'PASS' ELSE 'FAIL' END as status
FROM project_records;
"
```

**Test 1.6: Data Integrity**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  COUNT(*) as total_records,
  COUNT(CASE WHEN project_id IS NOT NULL THEN 1 END) as with_project,
  COUNT(CASE WHEN created_by IS NOT NULL THEN 1 END) as with_creator,
  CASE WHEN COUNT(*) = COUNT(CASE WHEN project_id IS NOT NULL THEN 1 END) THEN 'PASS' ELSE 'FAIL' END as status
FROM project_records;
"
```

**Test 1.7: No Data Loss**
```bash
psql -h localhost -U itsm -d itsm -c "
SELECT 
  (SELECT COUNT(*) FROM project_records_backup) as original_count,
  (SELECT COUNT(*) FROM project_records) as current_count,
  CASE WHEN (SELECT COUNT(*) FROM project_records_backup) = (SELECT COUNT(*) FROM project_records) THEN 'PASS' ELSE 'FAIL' END as status;
"
```

#### 2. Frontend Tests (5 tests)

**Test 2.1: Page Navigation**
```bash
# Open browser and navigate to:
# http://localhost:3000/projects/{project_id}
# Check: Page loads without errors
# Check: Sidebar visible
# Check: No console errors
```

**Test 2.2: Sprint Board Display**
```bash
# Check: Sprint board displays
# Check: Columns visible (Backlog, To Do, In Progress, Done)
# Check: Cards display in columns
# Check: Sprint metrics visible
```

**Test 2.3: Drag and Drop**
```bash
# Drag a card from one column to another
# Check: Card moves
# Check: Status updates
# Refresh page
# Check: Status persisted
```

**Test 2.4: Record Detail Modal**
```bash
# Click on a card
# Check: Modal opens
# Check: Record details display
# Check: Comments section visible
# Check: Attachments section visible
# Close modal
# Check: Modal closes
```

**Test 2.5: Sidebar Navigation**
```bash
# Click each sidebar link:
# - Backlog
# - Board
# - Sprint
# - Reports
# - Releases
# - Components
# - Issues
# - Repository
# - Settings
# Check: Each page loads correctly
```

#### 3. Backend API Tests (4 tests)

**Test 3.1: Sprint Endpoints**
```bash
# List sprints
curl -X GET http://localhost:8080/api/sprints

# Get sprint details
curl -X GET http://localhost:8080/api/sprints/{id}

# Create sprint
curl -X POST http://localhost:8080/api/sprints \
  -H "Content-Type: application/json" \
  -d '{"project_id":"...","name":"Test Sprint"}'
```

**Test 3.2: Workflow Endpoints**
```bash
# List workflows
curl -X GET http://localhost:8080/api/workflows

# Get workflow details
curl -X GET http://localhost:8080/api/workflows/{id}

# Get workflow statuses
curl -X GET http://localhost:8080/api/workflows/{id}/statuses
```

**Test 3.3: Issue Type Endpoints**
```bash
# List issue types
curl -X GET http://localhost:8080/api/issue-types

# Get issue type details
curl -X GET http://localhost:8080/api/issue-types/{id}
```

**Test 3.4: Record Endpoints**
```bash
# List records
curl -X GET http://localhost:8080/api/records

# Get record details
curl -X GET http://localhost:8080/api/records/{id}

# Create record
curl -X POST http://localhost:8080/api/records \
  -H "Content-Type: application/json" \
  -d '{"project_id":"...","title":"Test","status":"To Do"}'

# Update record
curl -X PUT http://localhost:8080/api/records/{id} \
  -H "Content-Type: application/json" \
  -d '{"status":"In Progress"}'
```

#### 4. Jira Features Tests (8 tests)

**Test 4.1: Issue Types**
- [ ] Can create records with Task type
- [ ] Can create records with Bug type
- [ ] Can create records with Story type
- [ ] Can create records with Epic type

**Test 4.2: Custom Fields**
- [ ] Custom fields visible in record detail
- [ ] Can add custom field values
- [ ] Can update custom field values

**Test 4.3: Workflows**
- [ ] Workflow statuses visible in board
- [ ] Can move records between statuses
- [ ] Status updates in database

**Test 4.4: Sprints**
- [ ] Can create sprints
- [ ] Can start sprints
- [ ] Can end sprints
- [ ] Sprint board displays correctly

**Test 4.5: Backlog**
- [ ] Backlog page displays
- [ ] Can drag records to sprint
- [ ] Can drag records between sprints

**Test 4.6: Comments**
- [ ] Can add comments
- [ ] Comments display correctly
- [ ] Can mention users

**Test 4.7: Attachments**
- [ ] Can upload attachments
- [ ] Attachments display
- [ ] Can download attachments

**Test 4.8: Labels**
- [ ] Can add labels
- [ ] Labels display
- [ ] Can filter by labels

#### 5. Performance Tests (3 tests)

**Test 5.1: Page Load Time**
```javascript
// In browser console
performance.measure('page-load', 'navigationStart', 'loadEventEnd');
console.log(performance.getEntriesByName('page-load')[0].duration);
// Should be < 2000ms
```

**Test 5.2: API Response Time**
```bash
# Measure API response time
time curl -X GET http://localhost:8080/api/sprints
# Should be < 500ms
```

**Test 5.3: Database Query Performance**
```bash
# Check slow queries
psql -h localhost -U itsm -d itsm -c "
SELECT query, calls, mean_time 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;
"
# Should be < 200ms
```

#### 6. Backward Compatibility Tests (2 tests)

**Test 6.1: Old Data Accessible**
```bash
# Check old project columns still exist
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) FROM project_columns;"

# Check old project records still exist
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) FROM project_records;"
```

**Test 6.2: Data Integrity**
```bash
# Check no data loss
psql -h localhost -U itsm -d itsm -c "
SELECT 
  (SELECT COUNT(*) FROM project_records_backup) as original,
  (SELECT COUNT(*) FROM project_records) as current,
  CASE WHEN (SELECT COUNT(*) FROM project_records_backup) = (SELECT COUNT(*) FROM project_records) THEN 'OK' ELSE 'LOST' END as status;
"
```

#### 7. Error Handling Tests (3 tests)

**Test 7.1: Invalid Data**
```bash
# Try to create record without title
curl -X POST http://localhost:8080/api/records \
  -H "Content-Type: application/json" \
  -d '{"project_id":"..."}'
# Should return 400 error
```

**Test 7.2: Not Found Errors**
```bash
# Try to access non-existent record
curl -X GET http://localhost:8080/api/records/invalid-id
# Should return 404 error
```

**Test 7.3: Permission Errors**
```bash
# Try to access records from other projects
curl -X GET http://localhost:8080/api/records?project_id=other-project
# Should return 403 error or empty list
```

### Testing Report

Create report using this template:

```markdown
# Testing Report - Phase 5

**Date**: [Date]
**Tester**: [Name]
**Status**: [PASS/FAIL]

## Summary
- Total Tests: 32
- Passed: [Number]
- Failed: [Number]
- Skipped: [Number]

## Data Migration Tests
- [ ] Backup verification: PASS/FAIL
- [ ] Workflow creation: PASS/FAIL
- [ ] Workflow statuses: PASS/FAIL
- [ ] Issue type schemes: PASS/FAIL
- [ ] Records migration: PASS/FAIL
- [ ] Data integrity: PASS/FAIL
- [ ] No data loss: PASS/FAIL

## Frontend Tests
- [ ] Page navigation: PASS/FAIL
- [ ] Sprint board display: PASS/FAIL
- [ ] Drag and drop: PASS/FAIL
- [ ] Record detail modal: PASS/FAIL
- [ ] Sidebar navigation: PASS/FAIL

## Backend API Tests
- [ ] Sprint endpoints: PASS/FAIL
- [ ] Workflow endpoints: PASS/FAIL
- [ ] Issue type endpoints: PASS/FAIL
- [ ] Record endpoints: PASS/FAIL

## Jira Features Tests
- [ ] Issue types: PASS/FAIL
- [ ] Custom fields: PASS/FAIL
- [ ] Workflows: PASS/FAIL
- [ ] Sprints: PASS/FAIL
- [ ] Backlog: PASS/FAIL
- [ ] Comments: PASS/FAIL
- [ ] Attachments: PASS/FAIL
- [ ] Labels: PASS/FAIL

## Performance Tests
- [ ] Page load time: PASS/FAIL
- [ ] API response time: PASS/FAIL
- [ ] Database query performance: PASS/FAIL

## Backward Compatibility Tests
- [ ] Old data accessible: PASS/FAIL
- [ ] Data integrity: PASS/FAIL

## Error Handling Tests
- [ ] Invalid data: PASS/FAIL
- [ ] Not found errors: PASS/FAIL
- [ ] Permission errors: PASS/FAIL

## Issues Found
1. [Issue 1]
2. [Issue 2]

## Sign-off
- Tester: [Name]
- Date: [Date]
- Status: [APPROVED/REJECTED]
```

---

## 🚀 PHASE 6: DEPLOYMENT

### Objective
Deploy Jira-only project board ke production

### Duration
90 minutes

### Prerequisites
- Phase 2 (Database Migration) completed
- Phase 5 (Testing) passed
- Database backup created
- Backend built
- Frontend built

### Deployment Steps

#### Step 1: Pre-Deployment Checklist

```bash
# 1. Verify database migration
psql -h localhost -U itsm -d itsm -c "SELECT COUNT(*) FROM workflows;"

# 2. Verify backend build
ls -la backend/itsm-backend

# 3. Verify frontend build
ls -la frontend/dist/

# 4. Create backup
pg_dump -h localhost -U itsm -d itsm > backup_pre_deployment_$(date +%Y%m%d_%H%M%S).sql

# 5. Verify backup
ls -lh backup_pre_deployment_*.sql
```

#### Step 2: Backend Deployment

**Using Docker Compose:**
```bash
# Build backend image
docker-compose build backend

# Stop current backend
docker-compose stop backend

# Start new backend
docker-compose up -d backend

# Verify backend is running
docker-compose logs backend | tail -20

# Check health
curl http://localhost:8080/health
```

**Using systemd:**
```bash
# Build backend
cd backend
go build -o itsm-backend

# Create backup of current backend
cp /opt/itsm/backend /opt/itsm/backend.backup

# Deploy new backend
cp itsm-backend /opt/itsm/backend

# Restart service
systemctl restart itsm-backend

# Verify
curl http://localhost:8080/health
```

#### Step 3: Frontend Deployment

**Using Docker Compose:**
```bash
# Build frontend image
docker-compose build frontend

# Stop current frontend
docker-compose stop frontend

# Start new frontend
docker-compose up -d frontend

# Verify frontend is running
docker-compose logs frontend | tail -20

# Check health
curl http://localhost:3000
```

**Using Nginx:**
```bash
# Build frontend
cd frontend
npm run build

# Create backup
cp -r /opt/itsm/frontend/dist /opt/itsm/frontend/dist.backup

# Deploy new frontend
cp -r dist/* /opt/itsm/frontend/dist/

# Reload Nginx
nginx -s reload

# Verify
curl http://localhost:3000
```

#### Step 4: Health Checks

```bash
# 1. Backend health
curl -X GET http://localhost:8080/health
# Expected: 200 OK

# 2. Frontend health
curl -X GET http://localhost:3000
# Expected: 200 OK

# 3. Database connection
curl -X GET http://localhost:8080/api/projects
# Expected: 200 OK

# 4. API endpoints
curl -X GET http://localhost:8080/api/sprints
curl -X GET http://localhost:8080/api/workflows
curl -X GET http://localhost:8080/api/records
# Expected: 200 OK
```

#### Step 5: Smoke Testing

```bash
# 1. Create sprint
curl -X POST http://localhost:8080/api/sprints \
  -H "Content-Type: application/json" \
  -d '{"project_id":"test","name":"Sprint 1","goal":"Test"}'

# 2. Get sprints
curl -X GET http://localhost:8080/api/sprints

# 3. Get workflows
curl -X GET http://localhost:8080/api/workflows

# 4. Get records
curl -X GET http://localhost:8080/api/records
```

#### Step 6: Monitoring

```bash
# Monitor backend logs
docker-compose logs -f backend

# Monitor frontend logs
docker-compose logs -f frontend

# Monitor database
psql -h localhost -U itsm -d itsm -c "SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"

# Monitor system resources
top
```

### Rollback Plan

**If deployment fails:**

```bash
# 1. Stop services
docker-compose stop backend frontend

# 2. Restore database
psql -h localhost -U itsm -d itsm < backup_pre_deployment_YYYYMMDD_HHMMSS.sql

# 3. Restore backend
cp /opt/itsm/backend.backup /opt/itsm/backend

# 4. Restore frontend
rm -rf /opt/itsm/frontend/dist
cp -r /opt/itsm/frontend/dist.backup /opt/itsm/frontend/dist

# 5. Start services
docker-compose up -d backend frontend

# 6. Verify
curl http://localhost:8080/health
curl http://localhost:3000
```

---

## 📊 Execution Timeline

| Step | Duration | Status |
|------|----------|--------|
| Phase 2: Database Migration | 5-10 min | ⏳ |
| Phase 5: Testing | 2-4 hours | ⏳ |
| Phase 6: Deployment | 90 min | ⏳ |
| **TOTAL** | **3-5 hours** | ⏳ |

---

## ✅ Success Criteria

### Phase 2: Database Migration
- ✅ All data migrated
- ✅ No data loss
- ✅ All workflows created
- ✅ All records updated

### Phase 5: Testing
- ✅ All 32 tests passed
- ✅ No critical issues
- ✅ Performance acceptable
- ✅ Backward compatible

### Phase 6: Deployment
- ✅ Database migrated
- ✅ Backend deployed
- ✅ Frontend deployed
- ✅ All health checks passed
- ✅ System stable 24h

---

## 🎯 Next Steps

1. Execute Phase 2 (Database Migration)
2. Execute Phase 5 (Testing)
3. Execute Phase 6 (Deployment)
4. Monitor production
5. Collect user feedback

---

**Ready to execute!** 🚀

