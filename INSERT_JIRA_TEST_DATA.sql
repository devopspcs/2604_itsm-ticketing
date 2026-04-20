-- ============================================
-- Jira-like Project Board - Test Data Insertion
-- ============================================

-- Get the test project ID (or create one if needed)
DO $$
DECLARE
    v_project_id UUID;
    v_workflow_id UUID;
    v_backlog_status_id UUID;
    v_todo_status_id UUID;
    v_inprogress_status_id UUID;
    v_inreview_status_id UUID;
    v_done_status_id UUID;
    v_sprint_id UUID;
BEGIN
    -- Get or create test project
    SELECT id INTO v_project_id FROM projects WHERE name = 'TEST PROJECT' LIMIT 1;
    
    IF v_project_id IS NULL THEN
        INSERT INTO projects (id, name, icon_color, created_by, created_at, updated_at)
        SELECT gen_random_uuid(), 'TEST PROJECT', '#3b82f6', id, NOW(), NOW()
        FROM users WHERE email = 'admin@itsm.local' LIMIT 1
        RETURNING id INTO v_project_id;
    END IF;

    -- ============================================
    -- 1. Insert Issue Types (if not exist)
    -- ============================================
    INSERT INTO issue_types (id, name, icon, description, created_at)
    VALUES 
        (gen_random_uuid(), 'Bug', 'bug_report', 'A bug or defect', NOW()),
        (gen_random_uuid(), 'Task', 'task_alt', 'A task to be done', NOW()),
        (gen_random_uuid(), 'Story', 'description', 'A user story', NOW()),
        (gen_random_uuid(), 'Epic', 'dashboard', 'An epic', NOW()),
        (gen_random_uuid(), 'Sub-task', 'subdirectory_arrow_right', 'A sub-task', NOW())
    ON CONFLICT (name) DO NOTHING;

    -- ============================================
    -- 2. Create Workflow for Project
    -- ============================================
    SELECT id INTO v_workflow_id FROM workflows WHERE project_id = v_project_id LIMIT 1;
    
    IF v_workflow_id IS NULL THEN
        INSERT INTO workflows (id, project_id, name, initial_status, created_at)
        VALUES (gen_random_uuid(), v_project_id, 'Default Workflow', 'Backlog', NOW())
        RETURNING id INTO v_workflow_id;
    END IF;

    -- ============================================
    -- 3. Create Workflow Statuses
    -- ============================================
    -- Backlog
    INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
    SELECT gen_random_uuid(), v_workflow_id, 'Backlog', 1, NOW()
    WHERE NOT EXISTS (SELECT 1 FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'Backlog')
    RETURNING id INTO v_backlog_status_id;
    
    IF v_backlog_status_id IS NULL THEN
        SELECT id INTO v_backlog_status_id FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'Backlog';
    END IF;

    -- To Do
    INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
    SELECT gen_random_uuid(), v_workflow_id, 'To Do', 2, NOW()
    WHERE NOT EXISTS (SELECT 1 FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'To Do')
    RETURNING id INTO v_todo_status_id;
    
    IF v_todo_status_id IS NULL THEN
        SELECT id INTO v_todo_status_id FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'To Do';
    END IF;

    -- In Progress
    INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
    SELECT gen_random_uuid(), v_workflow_id, 'In Progress', 3, NOW()
    WHERE NOT EXISTS (SELECT 1 FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'In Progress')
    RETURNING id INTO v_inprogress_status_id;
    
    IF v_inprogress_status_id IS NULL THEN
        SELECT id INTO v_inprogress_status_id FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'In Progress';
    END IF;

    -- In Review
    INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
    SELECT gen_random_uuid(), v_workflow_id, 'In Review', 4, NOW()
    WHERE NOT EXISTS (SELECT 1 FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'In Review')
    RETURNING id INTO v_inreview_status_id;
    
    IF v_inreview_status_id IS NULL THEN
        SELECT id INTO v_inreview_status_id FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'In Review';
    END IF;

    -- Done
    INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
    SELECT gen_random_uuid(), v_workflow_id, 'Done', 5, NOW()
    WHERE NOT EXISTS (SELECT 1 FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'Done')
    RETURNING id INTO v_done_status_id;
    
    IF v_done_status_id IS NULL THEN
        SELECT id INTO v_done_status_id FROM workflow_statuses WHERE workflow_id = v_workflow_id AND status_name = 'Done';
    END IF;

    -- ============================================
    -- 4. Create Workflow Transitions
    -- ============================================
    INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
    SELECT gen_random_uuid(), v_workflow_id, v_backlog_status_id, v_todo_status_id, NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM workflow_transitions 
        WHERE workflow_id = v_workflow_id 
        AND from_status_id = v_backlog_status_id 
        AND to_status_id = v_todo_status_id
    );

    INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
    SELECT gen_random_uuid(), v_workflow_id, v_todo_status_id, v_inprogress_status_id, NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM workflow_transitions 
        WHERE workflow_id = v_workflow_id 
        AND from_status_id = v_todo_status_id 
        AND to_status_id = v_inprogress_status_id
    );

    INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
    SELECT gen_random_uuid(), v_workflow_id, v_inprogress_status_id, v_inreview_status_id, NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM workflow_transitions 
        WHERE workflow_id = v_workflow_id 
        AND from_status_id = v_inprogress_status_id 
        AND to_status_id = v_inreview_status_id
    );

    INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
    SELECT gen_random_uuid(), v_workflow_id, v_inreview_status_id, v_done_status_id, NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM workflow_transitions 
        WHERE workflow_id = v_workflow_id 
        AND from_status_id = v_inreview_status_id 
        AND to_status_id = v_done_status_id
    );

    -- ============================================
    -- 5. Create Sprint
    -- ============================================
    SELECT id INTO v_sprint_id FROM sprints WHERE project_id = v_project_id AND status = 'Active' LIMIT 1;
    
    IF v_sprint_id IS NULL THEN
        INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status, actual_start_date, created_at)
        VALUES (
            gen_random_uuid(),
            v_project_id,
            'Sprint 1',
            'Initial sprint for testing',
            CURRENT_DATE,
            CURRENT_DATE + INTERVAL '14 days',
            'Active',
            CURRENT_DATE,
            NOW()
        )
        RETURNING id INTO v_sprint_id;
    END IF;

    -- ============================================
    -- 6. Create Labels
    -- ============================================
    INSERT INTO labels (id, project_id, name, color, created_at)
    VALUES 
        (gen_random_uuid(), v_project_id, 'Frontend', '#3b82f6', NOW()),
        (gen_random_uuid(), v_project_id, 'Backend', '#10b981', NOW()),
        (gen_random_uuid(), v_project_id, 'Database', '#f59e0b', NOW()),
        (gen_random_uuid(), v_project_id, 'Documentation', '#8b5cf6', NOW())
    ON CONFLICT DO NOTHING;

    -- ============================================
    -- 7. Create Custom Fields
    -- ============================================
    INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
    VALUES 
        (gen_random_uuid(), v_project_id, 'Priority', 'dropdown', false, NOW()),
        (gen_random_uuid(), v_project_id, 'Estimated Hours', 'number', false, NOW()),
        (gen_random_uuid(), v_project_id, 'Component', 'dropdown', false, NOW())
    ON CONFLICT DO NOTHING;

    -- ============================================
    -- 8. Create Sample Records
    -- ============================================
    INSERT INTO project_records (id, project_id, title, description, status, created_by, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_project_id,
        'Setup project board',
        'Configure the Jira-like project board',
        'To Do',
        (SELECT id FROM users WHERE email = 'admin@itsm.local' LIMIT 1),
        NOW(),
        NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM project_records 
        WHERE project_id = v_project_id 
        AND title = 'Setup project board'
    );

    INSERT INTO project_records (id, project_id, title, description, status, created_by, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_project_id,
        'Create API endpoints',
        'Implement all Jira-like API endpoints',
        'In Progress',
        (SELECT id FROM users WHERE email = 'admin@itsm.local' LIMIT 1),
        NOW(),
        NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM project_records 
        WHERE project_id = v_project_id 
        AND title = 'Create API endpoints'
    );

    INSERT INTO project_records (id, project_id, title, description, status, created_by, created_at, updated_at)
    SELECT 
        gen_random_uuid(),
        v_project_id,
        'Build frontend components',
        'Create React components for sprint board',
        'Backlog',
        (SELECT id FROM users WHERE email = 'admin@itsm.local' LIMIT 1),
        NOW(),
        NOW()
    WHERE NOT EXISTS (
        SELECT 1 FROM project_records 
        WHERE project_id = v_project_id 
        AND title = 'Build frontend components'
    );

    RAISE NOTICE 'Test data inserted successfully for project: %', v_project_id;
END $$;

-- Verify data was inserted
SELECT 'Issue Types' as entity, COUNT(*) as count FROM issue_types
UNION ALL
SELECT 'Workflows', COUNT(*) FROM workflows
UNION ALL
SELECT 'Workflow Statuses', COUNT(*) FROM workflow_statuses
UNION ALL
SELECT 'Sprints', COUNT(*) FROM sprints
UNION ALL
SELECT 'Labels', COUNT(*) FROM labels
UNION ALL
SELECT 'Custom Fields', COUNT(*) FROM custom_fields
UNION ALL
SELECT 'Project Records', COUNT(*) FROM project_records;
