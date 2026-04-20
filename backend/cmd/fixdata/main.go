package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://itsm:itsm@localhost:5432/itsm?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	projectID := "f3705d12-213e-4895-ae1a-90883b63735b"
	workflowID := "14767610-e62a-4ea6-9c96-bb3dd1251997"

	// Step 1: Remove duplicate workflow statuses, keep only one per status_name
	log.Println("=== Cleaning duplicate workflow statuses ===")
	_, err = pool.Exec(ctx, `
		DELETE FROM workflow_statuses
		WHERE id NOT IN (
			SELECT DISTINCT ON (status_name) id
			FROM workflow_statuses
			WHERE workflow_id = $1
			ORDER BY status_name, created_at ASC
		)
		AND workflow_id = $1
	`, workflowID)
	if err != nil {
		log.Fatalf("Failed to clean statuses: %v", err)
	}

	// Fix status_order to be sequential 0-4
	statuses := []string{"Backlog", "To Do", "In Progress", "In Review", "Done"}
	for i, s := range statuses {
		_, err = pool.Exec(ctx, `
			UPDATE workflow_statuses SET status_order = $1
			WHERE workflow_id = $2 AND status_name = $3
		`, i, workflowID, s)
		if err != nil {
			log.Fatalf("Failed to update status order: %v", err)
		}
	}
	log.Println("✓ Cleaned duplicate statuses and fixed order")

	// Step 2: Get the status IDs
	type statusInfo struct {
		ID   string
		Name string
	}
	rows, err := pool.Query(ctx, `
		SELECT id, status_name FROM workflow_statuses
		WHERE workflow_id = $1 ORDER BY status_order
	`, workflowID)
	if err != nil {
		log.Fatalf("Failed to query statuses: %v", err)
	}
	var statusList []statusInfo
	for rows.Next() {
		var s statusInfo
		rows.Scan(&s.ID, &s.Name)
		statusList = append(statusList, s)
		log.Printf("  Status: %s = %s", s.Name, s.ID)
	}
	rows.Close()

	if len(statusList) < 5 {
		log.Fatalf("Expected 5 statuses, got %d", len(statusList))
	}

	// Step 3: Get sprint ID
	var sprintID string
	err = pool.QueryRow(ctx, `
		SELECT id FROM sprints WHERE project_id = $1 AND status = 'Active' LIMIT 1
	`, projectID).Scan(&sprintID)
	if err != nil {
		log.Fatalf("No active sprint found: %v", err)
	}
	log.Printf("✓ Active sprint: %s", sprintID)

	// Step 4: Get issue type IDs
	type issueTypeInfo struct {
		ID   string
		Name string
	}
	rows2, err := pool.Query(ctx, `SELECT id, name FROM issue_types ORDER BY name`)
	if err != nil {
		log.Fatalf("Failed to query issue types: %v", err)
	}
	issueTypes := make(map[string]string)
	for rows2.Next() {
		var it issueTypeInfo
		rows2.Scan(&it.ID, &it.Name)
		issueTypes[it.Name] = it.ID
	}
	rows2.Close()

	// Step 5: Check if sprint_id column exists on project_records
	var hasSprintCol bool
	err = pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'project_records' AND column_name = 'sprint_id'
		)
	`).Scan(&hasSprintCol)
	if err != nil {
		log.Fatalf("Failed to check column: %v", err)
	}

	if !hasSprintCol {
		log.Println("Adding sprint_id column to project_records...")
		_, err = pool.Exec(ctx, `
			ALTER TABLE project_records ADD COLUMN IF NOT EXISTS sprint_id UUID REFERENCES sprints(id)
		`)
		if err != nil {
			log.Fatalf("Failed to add sprint_id column: %v", err)
		}
		log.Println("✓ Added sprint_id column")
	}

	// Step 6: Update existing records with sprint_id, status, and issue_type_id
	type recordInfo struct {
		ID    string
		Title string
	}
	rows3, err := pool.Query(ctx, `
		SELECT id, title FROM project_records WHERE project_id = $1
	`, projectID)
	if err != nil {
		log.Fatalf("Failed to query records: %v", err)
	}
	var records []recordInfo
	for rows3.Next() {
		var r recordInfo
		rows3.Scan(&r.ID, &r.Title)
		records = append(records, r)
	}
	rows3.Close()

	// Assign records to sprint with different statuses
	for i, rec := range records {
		statusIdx := i % len(statusList)
		issueTypeName := "Task"
		if i == 0 {
			issueTypeName = "Bug"
		} else if i == 2 {
			issueTypeName = "Story"
		}

		issueTypeID := issueTypes[issueTypeName]
		statusID := statusList[statusIdx].ID

		_, err = pool.Exec(ctx, `
			UPDATE project_records
			SET sprint_id = $1, status = $2, issue_type_id = $3
			WHERE id = $4
		`, sprintID, statusID, issueTypeID, rec.ID)
		if err != nil {
			log.Printf("Warning: Failed to update record %s: %v", rec.ID, err)
		} else {
			log.Printf("✓ Updated record '%s' → sprint=%s, status=%s (%s), type=%s",
				rec.Title, sprintID, statusList[statusIdx].Name, statusID, issueTypeName)
		}
	}

	// Step 7: Verify sprint records
	var count int
	pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM project_records WHERE sprint_id = $1
	`, sprintID).Scan(&count)
	fmt.Printf("\n✓ Total records in sprint: %d\n", count)

	// Step 8: Verify statuses
	rows4, err := pool.Query(ctx, `
		SELECT id, status_name, status_order FROM workflow_statuses
		WHERE workflow_id = $1 ORDER BY status_order
	`, workflowID)
	if err == nil {
		fmt.Println("\n=== Final Workflow Statuses ===")
		for rows4.Next() {
			var id, name string
			var order int
			rows4.Scan(&id, &name, &order)
			fmt.Printf("  [%d] %s = %s\n", order, name, id)
		}
		rows4.Close()
	}

	log.Println("\n✓ Data fix complete!")
}
