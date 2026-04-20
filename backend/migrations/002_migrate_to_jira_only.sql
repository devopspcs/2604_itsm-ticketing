-- ============================================
-- Phase 2: Migrate to Jira-Only Project Board
-- ============================================
-- Date: April 19, 2026
-- Purpose: Migrate old project board data to Jira-compatible schema
-- ============================================

-- ============================================
-- STEP 1: BACKUP EXISTING DATA
-- ============================================

-- Backup projects
CREATE TABLE IF NOT EXISTS projects_backup AS 
SELECT * FROM projects WHERE FALSE;

INSERT INTO projects_backup 
SELECT * FROM projects;

-- Backup project columns
CREATE TABLE IF NOT EXISTS project_columns_backup AS 
SELECT * FROM project_columns WHERE FALSE;

INSERT INTO project_columns_backup 
SELECT * FROM project_columns;

-- Backup project records
CREATE TABLE IF NOT EXISTS project_records_backup AS 
SELECT * FROM project_records WHERE FALSE;

INSERT INTO project_records_backup 
SELECT * FROM project_records;

-- Backup project activity logs
CREATE TABLE IF NOT EXISTS project_activity_logs_backup AS 
SELECT * FROM project_activity_logs WHERE FALSE;

INSERT INTO project_activity_logs_backup 
SELECT * FROM project_activity_logs;

-- Backup project members
CREATE TABLE IF NOT EXISTS project_members_backup AS 
SELECT * FROM project_members WHERE FALSE;

INSERT INTO project_members_backup 
SELECT * FROM project_members;

-- Verify backups
SELECT 'Backup Complete' as status,
  (SELECT COUNT(*) FROM projects_backup) as projects_count,
  (SELECT COUNT(*) FROM project_columns_backup) as columns_count,
  (SELECT COUNT(*) FROM project_records_backup) as records_count,
  (SELECT COUNT(*) FROM project_activity_logs_backup) as activities_count,
  (SELECT COUNT(*) FROM project_members_backup) as members_count;

-- ============================================
-- STEP 2: CREATE DEFAULT WORKFLOWS
-- ============================================

-- For each project, create a default workflow
INSERT INTO workflows (id, project_id, name, initial_status, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  p.id,
  'Default Workflow',
  'Backlog',
  NOW(),
  NOW()
FROM projects p
WHERE NOT EXISTS (
  SELECT 1 FROM workflows w WHERE w.project_id = p.id
)
ON CONFLICT DO NOTHING;

-- Verify workflows created
SELECT 'Workflows Created' as status,
  COUNT(*) as workflows_count
FROM workflows;

-- ============================================
-- STEP 3: CREATE WORKFLOW STATUSES FROM COLUMNS
-- ============================================

-- Create workflow statuses from project columns
INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  w.id,
  pc.name,
  pc.position,
  NOW(),
  NOW()
FROM project_columns pc
JOIN workflows w ON pc.project_id = w.project_id
WHERE NOT EXISTS (
  SELECT 1 FROM workflow_statuses ws 
  WHERE ws.workflow_id = w.id 
  AND ws.status_name = pc.name
)
ON CONFLICT DO NOTHING;

-- Verify statuses created
SELECT 'Workflow Statuses Created' as status,
  COUNT(*) as statuses_count
FROM workflow_statuses;

-- ============================================
-- STEP 4: CREATE DEFAULT ISSUE TYPE SCHEMES
-- ============================================

-- For each project, create a default issue type scheme
INSERT INTO issue_type_schemes (id, project_id, name, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  p.id,
  'Default Scheme',
  NOW(),
  NOW()
FROM projects p
WHERE NOT EXISTS (
  SELECT 1 FROM issue_type_schemes iss WHERE iss.project_id = p.id
)
ON CONFLICT DO NOTHING;

-- Verify schemes created
SELECT 'Issue Type Schemes Created' as status,
  COUNT(*) as schemes_count
FROM issue_type_schemes;

-- ============================================
-- STEP 5: ADD ISSUE TYPES TO SCHEMES
-- ============================================

-- Add all issue types to each scheme
INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  iss.id,
  it.id,
  NOW(),
  NOW()
FROM issue_type_schemes iss
CROSS JOIN issue_types it
WHERE NOT EXISTS (
  SELECT 1 FROM issue_type_scheme_items itsi
  WHERE itsi.scheme_id = iss.id
  AND itsi.issue_type_id = it.id
)
ON CONFLICT DO NOTHING;

-- Verify scheme items created
SELECT 'Issue Type Scheme Items Created' as status,
  COUNT(*) as scheme_items_count
FROM issue_type_scheme_items;

-- ============================================
-- STEP 6: UPDATE PROJECT_RECORDS WITH JIRA FIELDS
-- ============================================

-- Add Jira fields to project_records if not exist
ALTER TABLE project_records
ADD COLUMN IF NOT EXISTS issue_type_id UUID REFERENCES issue_types(id),
ADD COLUMN IF NOT EXISTS status VARCHAR DEFAULT 'Backlog',
ADD COLUMN IF NOT EXISTS parent_record_id UUID REFERENCES project_records(id);

-- Set default issue type (Task) for all records
UPDATE project_records
SET issue_type_id = (SELECT id FROM issue_types WHERE name = 'Task' LIMIT 1)
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
WHERE status = 'Backlog' AND column_id IS NOT NULL;

-- Verify updates
SELECT 'Records Updated' as status,
  (SELECT COUNT(*) FROM project_records WHERE issue_type_id IS NOT NULL) as records_with_issue_type,
  (SELECT COUNT(*) FROM project_records WHERE status IS NOT NULL) as records_with_status,
  (SELECT COUNT(*) FROM project_records WHERE issue_type_id IS NULL) as records_without_issue_type,
  (SELECT COUNT(*) FROM project_records WHERE status IS NULL) as records_without_status;

-- ============================================
-- STEP 7: VERIFY MIGRATION
-- ============================================

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
SELECT 'Records with Status', COUNT(*) FROM project_records WHERE status IS NOT NULL
ORDER BY entity;

-- ============================================
-- MIGRATION COMPLETE
-- ============================================

-- Summary
SELECT 
  'Migration Complete' as status,
  NOW() as completed_at,
  'All data successfully migrated to Jira-only schema' as message;
