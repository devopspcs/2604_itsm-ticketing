# Phase 5: Testing - Jira-Only Project Board

**Status**: READY FOR EXECUTION  
**Date**: April 19, 2026

---

## 📋 Testing Plan

Comprehensive testing untuk memastikan migrasi ke Jira-only project board berhasil dan semua fitur berfungsi dengan baik.

---

## 🧪 Testing Categories

### 1. Data Migration Testing

**Objective:** Verifikasi bahwa semua data termigrasi dengan benar

**Test Cases:**

#### 1.1 Backup Verification
- [ ] Backup tables exist (projects_backup, project_columns_backup, etc.)
- [ ] Backup data count matches original data
- [ ] Backup data integrity verified

**SQL Query:**
```sql
SELECT 
  'Projects' as table_name,
  (SELECT COUNT(*) FROM projects) as current_count,
  (SELECT COUNT(*) FROM projects_backup) as backup_count
UNION ALL
SELECT 'Columns', COUNT(*) FROM project_columns, (SELECT COUNT(*) FROM project_columns_backup)
UNION ALL
SELECT 'Records', COUNT(*) FROM project_records, (SELECT COUNT(*) FROM project_records_backup);
```

#### 1.2 Workflow Creation
- [ ] Default workflows created for all projects
- [ ] Workflow count = Project count
- [ ] All workflows have initial_status = 'Backlog'

**SQL Query:**
```sql
SELECT 
  p.id as project_id,
  p.name as project_name,
  w.id as workflow_id,
  w.name as workflow_name,
  w.initial_status
FROM projects p
LEFT JOIN workflows w ON p.id = w.project_id
ORDER BY p.id;
```

#### 1.3 Workflow Statuses
- [ ] Workflow statuses created from project columns
- [ ] Status count = Column count
- [ ] Status names match column names
- [ ] Status order preserved

**SQL Query:**
```sql
SELECT 
  w.id as workflow_id,
  w.name as workflow_name,
  ws.status_name,
  ws.status_order,
  COUNT(*) as status_count
FROM workflows w
LEFT JOIN workflow_statuses ws ON w.id = ws.workflow_id
GROUP BY w.id, w.name, ws.status_name, ws.status_order
ORDER BY w.id, ws.status_order;
```

#### 1.4 Issue Type Schemes
- [ ] Issue type schemes created for all projects
- [ ] Scheme count = Project count
- [ ] All issue types added to schemes

**SQL Query:**
```sql
SELECT 
  p.id as project_id,
  p.name as project_name,
  iss.id as scheme_id,
  iss.name as scheme_name,
  COUNT(itsi.id) as issue_type_count
FROM projects p
LEFT JOIN issue_type_schemes iss ON p.id = iss.project_id
LEFT JOIN issue_type_scheme_items itsi ON iss.id = itsi.scheme_id
GROUP BY p.id, p.name, iss.id, iss.name
ORDER BY p.id;
```

#### 1.5 Records Migration
- [ ] All records have issue_type_id (Task)
- [ ] All records have status (from column name)
- [ ] No records without issue_type_id
- [ ] No records without status

**SQL Query:**
```sql
SELECT 
  COUNT(*) as total_records,
  COUNT(CASE WHEN issue_type_id IS NOT NULL THEN 1 END) as with_issue_type,
  COUNT(CASE WHEN status IS NOT NULL THEN 1 END) as with_status,
  COUNT(CASE WHEN issue_type_id IS NULL THEN 1 END) as without_issue_type,
  COUNT(CASE WHEN status IS NULL THEN 1 END) as without_status
FROM project_records;
```

---

### 2. Frontend Testing

**Objective:** Verifikasi bahwa UI berfungsi dengan baik setelah migrasi

**Test Cases:**

#### 2.1 Page Navigation
- [ ] ProjectBoardPage loads without errors
- [ ] Sidebar navigation items visible
- [ ] All navigation links clickable
- [ ] No console errors

**Steps:**
1. Navigate to `/projects/{id}`
2. Verify page loads
3. Check browser console for errors
4. Click each sidebar link
5. Verify navigation works

#### 2.2 Sprint Board Display
- [ ] Sprint board displays correctly
- [ ] Columns visible (Backlog, To Do, In Progress, Done)
- [ ] Cards display in correct columns
- [ ] Sprint metrics visible (Total, Completed, Progress, Days Left)

**Steps:**
1. Navigate to `/projects/{id}`
2. Verify sprint board displays
3. Check columns are visible
4. Verify cards are in correct columns
5. Check sprint metrics

#### 2.3 Drag and Drop
- [ ] Cards can be dragged
- [ ] Cards can be dropped in different columns
- [ ] Status updates after drop
- [ ] No errors during drag-drop

**Steps:**
1. Navigate to sprint board
2. Drag a card from one column to another
3. Verify card moves
4. Verify status updates in database
5. Refresh page and verify status persisted

#### 2.4 Record Detail Modal
- [ ] Modal opens when card clicked
- [ ] Record details display correctly
- [ ] Comments section visible
- [ ] Attachments section visible
- [ ] Labels section visible
- [ ] Modal closes properly

**Steps:**
1. Click on a card
2. Verify modal opens
3. Check all sections visible
4. Close modal
5. Verify modal closes

#### 2.5 Sidebar Navigation
- [ ] Backlog link works
- [ ] Board link works
- [ ] Sprint link works
- [ ] Reports link works
- [ ] Releases link works
- [ ] Components link works
- [ ] Issues link works
- [ ] Repository link works
- [ ] Settings link works

**Steps:**
1. Click each sidebar link
2. Verify page loads
3. Verify correct page displays
4. Check no errors

---

### 3. Backend API Testing

**Objective:** Verifikasi bahwa API endpoints berfungsi dengan baik

**Test Cases:**

#### 3.1 Sprint Endpoints
- [ ] GET /api/sprints - List sprints
- [ ] GET /api/sprints/{id} - Get sprint details
- [ ] POST /api/sprints - Create sprint
- [ ] PUT /api/sprints/{id} - Update sprint
- [ ] DELETE /api/sprints/{id} - Delete sprint

**Test Commands:**
```bash
# List sprints
curl -X GET http://localhost:8080/api/sprints

# Get sprint details
curl -X GET http://localhost:8080/api/sprints/{id}

# Create sprint
curl -X POST http://localhost:8080/api/sprints \
  -H "Content-Type: application/json" \
  -d '{"project_id":"...","name":"Sprint 1","goal":"..."}'

# Update sprint
curl -X PUT http://localhost:8080/api/sprints/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Sprint 1 Updated"}'

# Delete sprint
curl -X DELETE http://localhost:8080/api/sprints/{id}
```

#### 3.2 Workflow Endpoints
- [ ] GET /api/workflows - List workflows
- [ ] GET /api/workflows/{id} - Get workflow details
- [ ] GET /api/workflows/{id}/statuses - Get workflow statuses

**Test Commands:**
```bash
# List workflows
curl -X GET http://localhost:8080/api/workflows

# Get workflow details
curl -X GET http://localhost:8080/api/workflows/{id}

# Get workflow statuses
curl -X GET http://localhost:8080/api/workflows/{id}/statuses
```

#### 3.3 Issue Type Endpoints
- [ ] GET /api/issue-types - List issue types
- [ ] GET /api/issue-types/{id} - Get issue type details

**Test Commands:**
```bash
# List issue types
curl -X GET http://localhost:8080/api/issue-types

# Get issue type details
curl -X GET http://localhost:8080/api/issue-types/{id}
```

#### 3.4 Record Endpoints
- [ ] GET /api/records - List records
- [ ] GET /api/records/{id} - Get record details
- [ ] POST /api/records - Create record
- [ ] PUT /api/records/{id} - Update record
- [ ] DELETE /api/records/{id} - Delete record

**Test Commands:**
```bash
# List records
curl -X GET http://localhost:8080/api/records

# Get record details
curl -X GET http://localhost:8080/api/records/{id}

# Create record
curl -X POST http://localhost:8080/api/records \
  -H "Content-Type: application/json" \
  -d '{"project_id":"...","title":"...","status":"..."}'

# Update record
curl -X PUT http://localhost:8080/api/records/{id} \
  -H "Content-Type: application/json" \
  -d '{"status":"In Progress"}'

# Delete record
curl -X DELETE http://localhost:8080/api/records/{id}
```

---

### 4. Jira Features Testing

**Objective:** Verifikasi bahwa semua Jira features berfungsi dengan baik

**Test Cases:**

#### 4.1 Issue Types
- [ ] Task issue type available
- [ ] Bug issue type available
- [ ] Story issue type available
- [ ] Epic issue type available
- [ ] Sub-task issue type available
- [ ] Can create records with different issue types

#### 4.2 Custom Fields
- [ ] Custom fields visible in record detail
- [ ] Can add custom field values
- [ ] Can update custom field values
- [ ] Custom field values persist

#### 4.3 Workflows
- [ ] Workflow statuses visible in board
- [ ] Can move records between statuses
- [ ] Status updates in database
- [ ] Workflow transitions work correctly

#### 4.4 Sprints
- [ ] Can create sprints
- [ ] Can start sprints
- [ ] Can end sprints
- [ ] Sprint board displays correctly
- [ ] Sprint metrics calculate correctly

#### 4.5 Backlog
- [ ] Backlog page displays correctly
- [ ] Can drag records from backlog to sprint
- [ ] Can drag records between sprints
- [ ] Backlog updates correctly

#### 4.6 Comments
- [ ] Can add comments to records
- [ ] Comments display correctly
- [ ] Can mention users in comments
- [ ] Comments persist

#### 4.7 Attachments
- [ ] Can upload attachments
- [ ] Attachments display correctly
- [ ] Can download attachments
- [ ] Attachments persist

#### 4.8 Labels
- [ ] Can add labels to records
- [ ] Labels display correctly
- [ ] Can filter by labels
- [ ] Labels persist

---

### 5. Performance Testing

**Objective:** Verifikasi bahwa aplikasi memiliki performa yang baik

**Test Cases:**

#### 5.1 Page Load Time
- [ ] ProjectBoardPage loads in < 2 seconds
- [ ] Sprint board renders in < 1 second
- [ ] Backlog page loads in < 2 seconds

**Measurement:**
```javascript
// In browser console
performance.measure('page-load', 'navigationStart', 'loadEventEnd');
console.log(performance.getEntriesByName('page-load')[0].duration);
```

#### 5.2 API Response Time
- [ ] GET /api/sprints responds in < 500ms
- [ ] GET /api/records responds in < 500ms
- [ ] POST /api/records responds in < 1000ms

**Measurement:**
```bash
# Using curl with timing
curl -w "Time: %{time_total}s\n" http://localhost:8080/api/sprints
```

#### 5.3 Database Query Performance
- [ ] Workflow status queries < 100ms
- [ ] Record queries < 200ms
- [ ] Sprint queries < 150ms

---

### 6. Backward Compatibility Testing

**Objective:** Verifikasi bahwa data lama masih bisa diakses

**Test Cases:**

#### 6.1 Old Project Board Data
- [ ] Old project columns still accessible
- [ ] Old project records still accessible
- [ ] Old activity logs still accessible
- [ ] Old project members still accessible

#### 6.2 Data Integrity
- [ ] No data loss during migration
- [ ] All records have correct issue_type_id
- [ ] All records have correct status
- [ ] All relationships intact

---

### 7. Error Handling Testing

**Objective:** Verifikasi bahwa error handling berfungsi dengan baik

**Test Cases:**

#### 7.1 Invalid Data
- [ ] Cannot create record without title
- [ ] Cannot create record without project_id
- [ ] Cannot create record with invalid status
- [ ] Error messages display correctly

#### 7.2 Not Found Errors
- [ ] 404 error when accessing non-existent record
- [ ] 404 error when accessing non-existent sprint
- [ ] 404 error when accessing non-existent project

#### 7.3 Permission Errors
- [ ] Cannot access records from other projects
- [ ] Cannot modify records without permission
- [ ] Cannot delete records without permission

---

## 📊 Testing Checklist

### Data Migration
- [ ] Backup verification
- [ ] Workflow creation
- [ ] Workflow statuses
- [ ] Issue type schemes
- [ ] Records migration

### Frontend
- [ ] Page navigation
- [ ] Sprint board display
- [ ] Drag and drop
- [ ] Record detail modal
- [ ] Sidebar navigation

### Backend API
- [ ] Sprint endpoints
- [ ] Workflow endpoints
- [ ] Issue type endpoints
- [ ] Record endpoints

### Jira Features
- [ ] Issue types
- [ ] Custom fields
- [ ] Workflows
- [ ] Sprints
- [ ] Backlog
- [ ] Comments
- [ ] Attachments
- [ ] Labels

### Performance
- [ ] Page load time
- [ ] API response time
- [ ] Database query performance

### Backward Compatibility
- [ ] Old project board data
- [ ] Data integrity

### Error Handling
- [ ] Invalid data
- [ ] Not found errors
- [ ] Permission errors

---

## 🎯 Testing Execution

### Step 1: Setup Test Environment
```bash
# Start development server
npm run dev

# Start backend server
go run main.go

# Start database
docker-compose up -d postgres
```

### Step 2: Run Data Migration Tests
```bash
# Execute migration script
psql -U itsm -d itsm -f backend/migrations/002_migrate_to_jira_only.sql

# Verify migration
psql -U itsm -d itsm -c "SELECT * FROM workflows;"
```

### Step 3: Run Frontend Tests
```bash
# Open browser
# Navigate to http://localhost:3000/projects/{id}
# Test all features manually
```

### Step 4: Run API Tests
```bash
# Test endpoints using curl or Postman
curl -X GET http://localhost:8080/api/sprints
```

### Step 5: Run Performance Tests
```bash
# Measure page load time
# Measure API response time
# Measure database query performance
```

---

## ✅ Testing Success Criteria

- ✅ All data migrated correctly
- ✅ No data loss
- ✅ All Jira features working
- ✅ Frontend loads without errors
- ✅ API endpoints responding correctly
- ✅ Performance acceptable (< 2s page load)
- ✅ Backward compatibility maintained
- ✅ Error handling working correctly

---

## 📝 Testing Report Template

```markdown
# Testing Report - Phase 5

**Date**: [Date]
**Tester**: [Name]
**Status**: [PASS/FAIL]

## Summary
- Total Tests: [Number]
- Passed: [Number]
- Failed: [Number]
- Skipped: [Number]

## Data Migration
- [ ] Backup verification: PASS/FAIL
- [ ] Workflow creation: PASS/FAIL
- [ ] Workflow statuses: PASS/FAIL
- [ ] Issue type schemes: PASS/FAIL
- [ ] Records migration: PASS/FAIL

## Frontend
- [ ] Page navigation: PASS/FAIL
- [ ] Sprint board display: PASS/FAIL
- [ ] Drag and drop: PASS/FAIL
- [ ] Record detail modal: PASS/FAIL
- [ ] Sidebar navigation: PASS/FAIL

## Backend API
- [ ] Sprint endpoints: PASS/FAIL
- [ ] Workflow endpoints: PASS/FAIL
- [ ] Issue type endpoints: PASS/FAIL
- [ ] Record endpoints: PASS/FAIL

## Jira Features
- [ ] Issue types: PASS/FAIL
- [ ] Custom fields: PASS/FAIL
- [ ] Workflows: PASS/FAIL
- [ ] Sprints: PASS/FAIL
- [ ] Backlog: PASS/FAIL
- [ ] Comments: PASS/FAIL
- [ ] Attachments: PASS/FAIL
- [ ] Labels: PASS/FAIL

## Performance
- [ ] Page load time: PASS/FAIL
- [ ] API response time: PASS/FAIL
- [ ] Database query performance: PASS/FAIL

## Issues Found
1. [Issue 1]
2. [Issue 2]
3. [Issue 3]

## Recommendations
1. [Recommendation 1]
2. [Recommendation 2]
3. [Recommendation 3]

## Sign-off
- Tester: [Name]
- Date: [Date]
- Status: [APPROVED/REJECTED]
```

---

## 🎯 Next Steps

After Phase 5 (Testing) is complete:
1. Review testing report
2. Fix any issues found
3. Re-test if needed
4. Proceed to Phase 6 (Deployment)

---

**Phase 5 Status**: ⏳ READY FOR EXECUTION

**Apakah Anda ingin saya jalankan testing?** 🚀

