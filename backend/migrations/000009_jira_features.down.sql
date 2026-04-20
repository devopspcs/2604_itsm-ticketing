-- Drop indexes on extended project_records columns
DROP INDEX IF EXISTS idx_project_records_parent_record_id;
DROP INDEX IF EXISTS idx_project_records_status;
DROP INDEX IF EXISTS idx_project_records_issue_type_id;

-- Remove columns from project_records table
ALTER TABLE project_records DROP COLUMN IF EXISTS parent_record_id;
ALTER TABLE project_records DROP COLUMN IF EXISTS status;
ALTER TABLE project_records DROP COLUMN IF EXISTS issue_type_id;

-- Drop record_labels table
DROP TABLE IF EXISTS record_labels;

-- Drop labels table
DROP TABLE IF EXISTS labels;

-- Drop attachments table
DROP TABLE IF EXISTS attachments;

-- Drop comment_mentions table
DROP TABLE IF EXISTS comment_mentions;

-- Drop comments table
DROP TABLE IF EXISTS comments;

-- Drop sprint_records table
DROP TABLE IF EXISTS sprint_records;

-- Drop sprints table
DROP TABLE IF EXISTS sprints;

-- Drop workflow_transitions table
DROP TABLE IF EXISTS workflow_transitions;

-- Drop workflow_statuses table
DROP TABLE IF EXISTS workflow_statuses;

-- Drop workflows table
DROP TABLE IF EXISTS workflows;

-- Drop custom_field_values table
DROP TABLE IF EXISTS custom_field_values;

-- Drop custom_field_options table
DROP TABLE IF EXISTS custom_field_options;

-- Drop custom_fields table
DROP TABLE IF EXISTS custom_fields;

-- Drop issue_type_scheme_items table
DROP TABLE IF EXISTS issue_type_scheme_items;

-- Drop issue_type_schemes table
DROP TABLE IF EXISTS issue_type_schemes;

-- Drop issue_types table
DROP TABLE IF EXISTS issue_types;
