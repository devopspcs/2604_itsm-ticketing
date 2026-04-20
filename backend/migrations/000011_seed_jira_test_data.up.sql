-- Seed predefined issue types
INSERT INTO issue_types (id, name, icon, description, created_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'Bug', 'bug_report', 'A problem or defect', NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'Task', 'task', 'A unit of work', NOW()),
  ('550e8400-e29b-41d4-a716-446655440003', 'Story', 'description', 'A user story', NOW()),
  ('550e8400-e29b-41d4-a716-446655440004', 'Epic', 'flag', 'A large body of work', NOW()),
  ('550e8400-e29b-41d4-a716-446655440005', 'Sub-task', 'subdirectory_arrow_right', 'A subtask of a story', NOW())
ON CONFLICT DO NOTHING;

-- Create test project (if not exists)
INSERT INTO projects (id, name, icon_color, created_by, created_at, updated_at)
SELECT 
  '550e8400-e29b-41d4-a716-446655440100',
  'Test Project - Jira Board',
  '#3b82f6',
  u.id,
  NOW(),
  NOW()
FROM users u
WHERE u.email = 'admin@example.com'
LIMIT 1
ON CONFLICT DO NOTHING;

-- Create default workflow for test project
INSERT INTO workflows (id, project_id, name, initial_status, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440200',
  p.id,
  'Default Workflow',
  'Backlog',
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

-- Create workflow statuses
INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440201',
  w.id,
  'Backlog',
  0,
  NOW()
FROM workflows w
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440202',
  w.id,
  'To Do',
  1,
  NOW()
FROM workflows w
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440203',
  w.id,
  'In Progress',
  2,
  NOW()
FROM workflows w
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440204',
  w.id,
  'In Review',
  3,
  NOW()
FROM workflows w
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440205',
  w.id,
  'Done',
  4,
  NOW()
FROM workflows w
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

-- Create workflow transitions
INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440301',
  w.id,
  ws1.id,
  ws2.id,
  NOW()
FROM workflows w
JOIN workflow_statuses ws1 ON w.id = ws1.workflow_id AND ws1.status_name = 'Backlog'
JOIN workflow_statuses ws2 ON w.id = ws2.workflow_id AND ws2.status_name = 'To Do'
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440302',
  w.id,
  ws1.id,
  ws2.id,
  NOW()
FROM workflows w
JOIN workflow_statuses ws1 ON w.id = ws1.workflow_id AND ws1.status_name = 'To Do'
JOIN workflow_statuses ws2 ON w.id = ws2.workflow_id AND ws2.status_name = 'In Progress'
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440303',
  w.id,
  ws1.id,
  ws2.id,
  NOW()
FROM workflows w
JOIN workflow_statuses ws1 ON w.id = ws1.workflow_id AND ws1.status_name = 'In Progress'
JOIN workflow_statuses ws2 ON w.id = ws2.workflow_id AND ws2.status_name = 'In Review'
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440304',
  w.id,
  ws1.id,
  ws2.id,
  NOW()
FROM workflows w
JOIN workflow_statuses ws1 ON w.id = ws1.workflow_id AND ws1.status_name = 'In Review'
JOIN workflow_statuses ws2 ON w.id = ws2.workflow_id AND ws2.status_name = 'Done'
WHERE w.name = 'Default Workflow'
ON CONFLICT DO NOTHING;

-- Create issue type scheme
INSERT INTO issue_type_schemes (id, project_id, name, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440400',
  p.id,
  'Default Issue Type Scheme',
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

-- Add issue types to scheme
INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440401',
  s.id,
  it.id,
  NOW()
FROM issue_type_schemes s
JOIN issue_types it ON it.name = 'Bug'
WHERE s.name = 'Default Issue Type Scheme'
ON CONFLICT DO NOTHING;

INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440402',
  s.id,
  it.id,
  NOW()
FROM issue_type_schemes s
JOIN issue_types it ON it.name = 'Task'
WHERE s.name = 'Default Issue Type Scheme'
ON CONFLICT DO NOTHING;

INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440403',
  s.id,
  it.id,
  NOW()
FROM issue_type_schemes s
JOIN issue_types it ON it.name = 'Story'
WHERE s.name = 'Default Issue Type Scheme'
ON CONFLICT DO NOTHING;

INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440404',
  s.id,
  it.id,
  NOW()
FROM issue_type_schemes s
JOIN issue_types it ON it.name = 'Epic'
WHERE s.name = 'Default Issue Type Scheme'
ON CONFLICT DO NOTHING;

INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440405',
  s.id,
  it.id,
  NOW()
FROM issue_type_schemes s
JOIN issue_types it ON it.name = 'Sub-task'
WHERE s.name = 'Default Issue Type Scheme'
ON CONFLICT DO NOTHING;

-- Create test sprint (Active status)
INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status, actual_start_date, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440500',
  p.id,
  'Sprint 1 - Test',
  'Complete initial setup and testing',
  CURRENT_DATE,
  CURRENT_DATE + INTERVAL '14 days',
  'Active',
  CURRENT_DATE,
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

-- Create test labels
INSERT INTO labels (id, project_id, name, color, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440601',
  p.id,
  'Frontend',
  '#3b82f6',
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

INSERT INTO labels (id, project_id, name, color, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440602',
  p.id,
  'Backend',
  '#ef4444',
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

INSERT INTO labels (id, project_id, name, color, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440603',
  p.id,
  'Database',
  '#8b5cf6',
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

INSERT INTO labels (id, project_id, name, color, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440604',
  p.id,
  'Documentation',
  '#10b981',
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

-- Create custom fields
INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440701',
  p.id,
  'Priority',
  'dropdown',
  true,
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440702',
  p.id,
  'Estimated Hours',
  'number',
  false,
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440703',
  p.id,
  'Component',
  'dropdown',
  false,
  NOW()
FROM projects p
WHERE p.name = 'Test Project - Jira Board'
ON CONFLICT DO NOTHING;

-- Add options to Priority field
INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440801',
  cf.id,
  'Low',
  1,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Priority' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440802',
  cf.id,
  'Medium',
  2,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Priority' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440803',
  cf.id,
  'High',
  3,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Priority' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440804',
  cf.id,
  'Critical',
  4,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Priority' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

-- Add options to Component field
INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440901',
  cf.id,
  'UI',
  1,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Component' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440902',
  cf.id,
  'API',
  2,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Component' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440903',
  cf.id,
  'Database',
  3,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Component' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;

INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
SELECT
  '550e8400-e29b-41d4-a716-446655440904',
  cf.id,
  'Infrastructure',
  4,
  NOW()
FROM custom_fields cf
WHERE cf.name = 'Component' AND cf.field_type = 'dropdown'
ON CONFLICT DO NOTHING;
