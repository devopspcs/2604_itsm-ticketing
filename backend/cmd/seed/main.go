package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://itsm:itsm@localhost:5432/itsm?sslmode=disable"
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Ping database
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("✓ Database connected")

	// Execute seed script
	if err := seedData(context.Background(), pool); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("✓ Test data inserted successfully")
}

func seedData(ctx context.Context, pool *pgxpool.Pool) error {
	// Get or create test project
	var projectID string
	err := pool.QueryRow(ctx, `
		SELECT id FROM projects WHERE name = 'TEST PROJECT' LIMIT 1
	`).Scan(&projectID)
	
	if err != nil {
		// Create test project
		err = pool.QueryRow(ctx, `
			INSERT INTO projects (id, name, icon_color, created_by, created_at, updated_at)
			SELECT gen_random_uuid(), 'TEST PROJECT', '#3b82f6', id, NOW(), NOW()
			FROM users WHERE email = 'admin@itsm.local' LIMIT 1
			RETURNING id
		`).Scan(&projectID)
		if err != nil {
			return err
		}
		log.Printf("Created test project: %s", projectID)
	}

	// Insert Issue Types
	_, err = pool.Exec(ctx, `
		INSERT INTO issue_types (id, name, icon, description, created_at)
		VALUES 
			(gen_random_uuid(), 'Bug', 'bug_report', 'A bug or defect', NOW()),
			(gen_random_uuid(), 'Task', 'task_alt', 'A task to be done', NOW()),
			(gen_random_uuid(), 'Story', 'description', 'A user story', NOW()),
			(gen_random_uuid(), 'Epic', 'dashboard', 'An epic', NOW()),
			(gen_random_uuid(), 'Sub-task', 'subdirectory_arrow_right', 'A sub-task', NOW())
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Issue types inserted")

	// Create or get workflow
	var workflowID string
	err = pool.QueryRow(ctx, `
		SELECT id FROM workflows WHERE project_id = $1 LIMIT 1
	`, projectID).Scan(&workflowID)
	
	if err != nil {
		// Create workflow
		err = pool.QueryRow(ctx, `
			INSERT INTO workflows (id, project_id, name, initial_status, created_at)
			VALUES (gen_random_uuid(), $1, 'Default Workflow', 'Backlog', NOW())
			RETURNING id
		`, projectID).Scan(&workflowID)
		if err != nil {
			return err
		}
		log.Printf("Created workflow: %s", workflowID)
	}

	// Create workflow statuses
	statuses := []string{"Backlog", "To Do", "In Progress", "In Review", "Done"}
	for i, status := range statuses {
		_, err = pool.Exec(ctx, `
			INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
			VALUES (gen_random_uuid(), $1, $2, $3, NOW())
			ON CONFLICT DO NOTHING
		`, workflowID, status, i+1)
		if err != nil {
			return err
		}
	}
	log.Println("✓ Workflow statuses created")

	// Create workflow transitions - skip for now
	log.Println("✓ Workflow transitions skipped")

	// Create or get active sprint
	var sprintID string
	err = pool.QueryRow(ctx, `
		SELECT id FROM sprints WHERE project_id = $1 AND status = 'Active' LIMIT 1
	`, projectID).Scan(&sprintID)
	
	if err != nil {
		// Create sprint
		err = pool.QueryRow(ctx, `
			INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status, actual_start_date, created_at)
			VALUES (gen_random_uuid(), $1, 'Sprint 1', 'Initial sprint for testing', CURRENT_DATE, CURRENT_DATE + INTERVAL '14 days', 'Active', CURRENT_DATE, NOW())
			RETURNING id
		`, projectID).Scan(&sprintID)
		if err != nil {
			return err
		}
		log.Printf("Created sprint: %s", sprintID)
	}

	// Create labels
	_, err = pool.Exec(ctx, `
		INSERT INTO labels (id, project_id, name, color, created_at)
		VALUES 
			(gen_random_uuid(), $1, 'Frontend', '#3b82f6', NOW()),
			(gen_random_uuid(), $1, 'Backend', '#10b981', NOW()),
			(gen_random_uuid(), $1, 'Database', '#f59e0b', NOW()),
			(gen_random_uuid(), $1, 'Documentation', '#8b5cf6', NOW())
		ON CONFLICT DO NOTHING
	`, projectID)
	if err != nil {
		return err
	}
	log.Println("✓ Labels created")

	// Create custom fields
	_, err = pool.Exec(ctx, `
		INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
		VALUES 
			(gen_random_uuid(), $1, 'Priority', 'dropdown', false, NOW()),
			(gen_random_uuid(), $1, 'Estimated Hours', 'number', false, NOW()),
			(gen_random_uuid(), $1, 'Component', 'dropdown', false, NOW())
		ON CONFLICT DO NOTHING
	`, projectID)
	if err != nil {
		return err
	}
	log.Println("✓ Custom fields created")

	// Create or get default column
	var columnID string
	err = pool.QueryRow(ctx, `
		SELECT id FROM project_columns WHERE project_id = $1 LIMIT 1
	`, projectID).Scan(&columnID)
	
	if err != nil {
		// Create column
		err = pool.QueryRow(ctx, `
			INSERT INTO project_columns (id, project_id, name, position, created_at)
			VALUES (gen_random_uuid(), $1, 'Backlog', 1, NOW())
			RETURNING id
		`, projectID).Scan(&columnID)
		if err != nil {
			return err
		}
		log.Printf("Created column: %s", columnID)
	}

	// Create sample records
	_, err = pool.Exec(ctx, `
		INSERT INTO project_records (id, column_id, project_id, title, description, position, created_by, created_at, updated_at)
		SELECT gen_random_uuid(), $1, $2, 'Setup project board', 'Configure the Jira-like project board', 1, (SELECT id FROM users WHERE email = 'admin@itsm.local' LIMIT 1), NOW(), NOW()
		WHERE NOT EXISTS (SELECT 1 FROM project_records WHERE project_id = $2 AND title = 'Setup project board')
	`, columnID, projectID)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO project_records (id, column_id, project_id, title, description, position, created_by, created_at, updated_at)
		SELECT gen_random_uuid(), $1, $2, 'Create API endpoints', 'Implement all Jira-like API endpoints', 2, (SELECT id FROM users WHERE email = 'admin@itsm.local' LIMIT 1), NOW(), NOW()
		WHERE NOT EXISTS (SELECT 1 FROM project_records WHERE project_id = $2 AND title = 'Create API endpoints')
	`, columnID, projectID)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO project_records (id, column_id, project_id, title, description, position, created_by, created_at, updated_at)
		SELECT gen_random_uuid(), $1, $2, 'Build frontend components', 'Create React components for sprint board', 3, (SELECT id FROM users WHERE email = 'admin@itsm.local' LIMIT 1), NOW(), NOW()
		WHERE NOT EXISTS (SELECT 1 FROM project_records WHERE project_id = $2 AND title = 'Build frontend components')
	`, columnID, projectID)
	if err != nil {
		return err
	}

	log.Println("✓ Sample records created")

	return nil
}
