# Phase 2: Migrasi Data - Execution

**Status**: IN PROGRESS  
**Date**: April 19, 2026

---

## 📋 Migrasi Data Step-by-Step

### Step 1: Backup Database

```sql
-- Backup existing data
CREATE TABLE projects_backup AS SELECT * FROM projects;
CREATE TABLE project_columns_backup AS SELECT * FROM project_columns;
CREATE TABLE project_records_backup AS SELECT * FROM project_records;
CREATE TABLE project_activity_logs_backup AS SELECT * FROM project_activity_logs;
CREATE TABLE project_members_backup AS SELECT * FROM project_members;

-- Verify backup
SELECT COUNT(*) as projects_count FROM projects_backup;
SELECT COUNT(*) as columns_count FROM project_columns_backup;
SELECT COUNT(*) as records_count FROM project_records_backup;
SELECT COUNT(*) as activities_count FROM project_activity_logs_backup;
SELECT COUNT(*) as members_count FROM project_members_backup;
```

### Step 2: Create Default Workflows

```sql
-- For each project, create a default workflow
INSERT INTO workflows (id, project_id, name, initial_status, created_at)
SELECT 
  gen_random_uuid(),
  p.id,
  'Default Workflow',
  'Backlog',
  NOW()
FROM projects p
WHERE NOT EXISTS (
  SELECT 1 FROM workflows w WHERE w.project_id = p.id
);

-- Verify workflows created
SELECT COUNT(*) as workflows_count FROM workflows;
```

### Step 3: Create Workflow Statuses from Columns

```sql
-- Create workflow statuses from project columns
INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
SELECT 
  gen_random_uuid(),
  w.id,
  pc.name,
  pc.position,
  NOW()
FROM project_columns pc
JOIN workflows w ON pc.project_id = w.project_id
WHERE NOT EXISTS (
  SELECT 1 FROM workflow_statuses ws 
  WHERE ws.workflow_id = w.id 
  AND ws.status_name = pc.name
);

-- Verify statuses created
SELECT COUNT(*) as statuses_count FROM workflow_statuses;
```

### Step 4: Create Default Issue Type Schemes

```sql
-- For each project, create a default issue type scheme
INSERT INTO issue_type_schemes (id, project_id, name, created_at)
SELECT 
  gen_random_uuid(),
  p.id,
  'Default Scheme',
  NOW()
FROM projects p
WHERE NOT EXISTS (
  SELECT 1 FROM issue_type_schemes iss WHERE iss.project_id = p.id
);

-- Verify schemes created
SELECT COUNT(*) as schemes_count FROM issue_type_schemes;
```

### Step 5: Add Issue Types to Schemes

```sql
-- Add all issue types to each scheme
INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
SELECT 
  gen_random_uuid(),
  iss.id,
  it.id,
  NOW()
FROM issue_type_schemes iss
CROSS JOIN issue_types it
WHERE NOT EXISTS (
  SELECT 1 FROM issue_type_scheme_items itsi
  WHERE itsi.scheme_id = iss.id
  AND itsi.issue_type_id = it.id
);

-- Verify scheme items created
SELECT COUNT(*) as scheme_items_count FROM issue_type_scheme_items;
```

### Step 6: Update project_records with Jira Fields

```sql
-- Add Jira fields to project_records if not exist
ALTER TABLE project_records
ADD COLUMN IF NOT EXISTS issue_type_id UUID REFERENCES issue_types(id),
ADD COLUMN IF NOT EXISTS status VARCHAR DEFAULT 'Backlog',
ADD COLUMN IF NOT EXISTS parent_record_id UUID REFERENCES project_records(id);

-- Set default issue type (Task) for all records
UPDATE project_records
SET issue_type_id = (SELECT id FROM issue_types WHERE name = 'Task')
WHERE issue_type_id IS NULL;

-- Set status from column name
UPDATE project_records pr
SET status = (
  SELECT ws.status_name
  FROM workflow_statuses ws
  JOIN workflows w ON ws.workflow_id = w.id
  JOIN project_columns pc ON w.project_id = pc.project_id
  WHERE pc.id = pr.column_id
  LIMIT 1
)
WHERE status = 'Backlog';

-- Verify updates
SELECT COUNT(*) as records_with_issue_type FROM project_records WHERE issue_type_id IS NOT NULL;
SELECT COUNT(*) as records_with_status FROM project_records WHERE status IS NOT NULL;
```

### Step 7: Verify Migration

```sql
-- Verify all data migrated correctly
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
SELECT 'Records with Status', COUNT(*) FROM project_records WHERE status IS NOT NULL;

-- Check for any records without issue type or status
SELECT COUNT(*) as records_without_issue_type FROM project_records WHERE issue_type_id IS NULL;
SELECT COUNT(*) as records_without_status FROM project_records WHERE status IS NULL;
```

---

## 🔄 Migration Execution Plan

### Pre-Migration Checklist
- [ ] Backup database
- [ ] Verify backup integrity
- [ ] Stop application
- [ ] Notify users about maintenance

### Migration Steps
- [ ] Step 1: Backup Database
- [ ] Step 2: Create Default Workflows
- [ ] Step 3: Create Workflow Statuses
- [ ] Step 4: Create Issue Type Schemes
- [ ] Step 5: Add Issue Types to Schemes
- [ ] Step 6: Update Records with Jira Fields
- [ ] Step 7: Verify Migration

### Post-Migration Checklist
- [ ] Verify all data migrated
- [ ] Check for any errors
- [ ] Start application
- [ ] Test all features
- [ ] Notify users

---

## 📊 Migration Statistics

### Before Migration
```
Projects: ?
Columns: ?
Records: ?
Activity Logs: ?
```

### After Migration
```
Projects: (same)
Workflows: (one per project)
Workflow Statuses: (one per column)
Issue Type Schemes: (one per project)
Issue Type Scheme Items: (5 per scheme)
Records with Issue Type: (all)
Records with Status: (all)
```

---

## ⚠️ Rollback Plan

If migration fails, rollback using backup:

```sql
-- Restore from backup
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS workflow_statuses;
DROP TABLE IF EXISTS workflow_transitions;
DROP TABLE IF EXISTS issue_type_schemes;
DROP TABLE IF EXISTS issue_type_scheme_items;

-- Restore original data
ALTER TABLE project_records
DROP COLUMN IF EXISTS issue_type_id,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS parent_record_id;

-- Verify rollback
SELECT COUNT(*) FROM projects;
SELECT COUNT(*) FROM project_columns;
SELECT COUNT(*) FROM project_records;
```

---

## 🎯 Next Steps

1. ✅ Phase 1 Complete - Analisis selesai
2. ⏳ Phase 2 - Migrasi Data (IN PROGRESS)
3. ⏳ Phase 3 - Update Backend
4. ⏳ Phase 4 - Update Frontend
5. ⏳ Phase 5 - Testing
6. ⏳ Phase 6 - Deployment

---

## 📝 Notes

- Migrasi data harus dilakukan di development terlebih dahulu
- Verify semua data termigrasi dengan benar
- Jika ada error, gunakan rollback plan
- Setelah sukses di development, jalankan di production

---

**Phase 2 Status**: ⏳ READY FOR EXECUTION

**Apakah Anda ingin saya jalankan migration scripts?** 🚀

Atau apakah Anda ingin menjalankannya sendiri di database?
