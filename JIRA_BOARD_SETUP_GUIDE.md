# Jira-like Project Board Setup Guide

## Overview

This guide provides step-by-step instructions to set up and verify the Jira-like project board features, including database migrations, test data insertion, API integration verification, and frontend component fixes.

## Prerequisites

- PostgreSQL database running
- Backend Go server running
- Frontend React development server running
- Docker (optional, for containerized setup)
- Migration tool (golang-migrate or similar)

## Part 1: Database Setup

### 1.1 Run Database Migrations

The database schema for Jira-like features has already been created in migration `000009_jira_features.up.sql`. To verify it's applied:

```bash
# Using golang-migrate
migrate -path backend/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

# Or check migration status
migrate -path backend/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" version
```

### 1.2 Insert Test Data

Run the test data migration to populate predefined issue types, test project, workflow, sprint, labels, and custom fields:

```bash
# Using golang-migrate
migrate -path backend/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

# Or manually run the SQL script
psql -U $DB_USER -d $DB_NAME -h $DB_HOST -f backend/migrations/000011_seed_jira_test_data.up.sql
```

### 1.3 Verify Test Data

Connect to PostgreSQL and verify the data was inserted:

```sql
-- Check issue types
SELECT * FROM issue_types;
-- Expected: 5 rows (Bug, Task, Story, Epic, Sub-task)

-- Check test project
SELECT * FROM projects WHERE name = 'Test Project - Jira Board';
-- Expected: 1 row

-- Check workflow
SELECT * FROM workflows WHERE name = 'Default Workflow';
-- Expected: 1 row

-- Check workflow statuses
SELECT * FROM workflow_statuses WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Default Workflow');
-- Expected: 5 rows (Backlog, To Do, In Progress, In Review, Done)

-- Check active sprint
SELECT * FROM sprints WHERE name = 'Sprint 1 - Test' AND status = 'Active';
-- Expected: 1 row

-- Check labels
SELECT * FROM labels WHERE project_id = (SELECT id FROM projects WHERE name = 'Test Project - Jira Board');
-- Expected: 4 rows (Frontend, Backend, Database, Documentation)

-- Check custom fields
SELECT * FROM custom_fields WHERE project_id = (SELECT id FROM projects WHERE name = 'Test Project - Jira Board');
-- Expected: 3 rows (Priority, Estimated Hours, Component)
```

## Part 2: Backend API Verification

### 2.1 Check API Endpoints

Verify that all Jira-like API endpoints are implemented and accessible:

```bash
# Get issue types for project
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/issue-types \
  -H "Authorization: Bearer {token}"

# Get workflow
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/workflow \
  -H "Authorization: Bearer {token}"

# Get active sprint
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/sprints/active \
  -H "Authorization: Bearer {token}"

# Get sprint records
curl -X GET http://localhost:8080/api/v1/sprints/{sprintId}/records \
  -H "Authorization: Bearer {token}"

# Get backlog
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/backlog \
  -H "Authorization: Bearer {token}"

# Get custom fields
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/custom-fields \
  -H "Authorization: Bearer {token}"

# Get labels
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/labels \
  -H "Authorization: Bearer {token}"
```

### 2.2 Verify Error Handling

Test error scenarios:

```bash
# Missing authentication
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/issue-types
# Expected: 401 Unauthorized

# Invalid project ID
curl -X GET http://localhost:8080/api/v1/projects/invalid-id/issue-types \
  -H "Authorization: Bearer {token}"
# Expected: 404 Not Found or 400 Bad Request

# Unauthorized user (not project member)
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/issue-types \
  -H "Authorization: Bearer {otherUserToken}"
# Expected: 403 Forbidden
```

## Part 3: Frontend API Integration

### 3.1 Verify jira.service.ts

The frontend service file `frontend/src/services/jira.service.ts` contains all API calls. Verify the endpoints match your backend:

```typescript
// Key endpoints to verify:
- listIssueTypes(projectId)
- getWorkflow(projectId)
- getActiveSprint(projectId)
- getSprintRecords(sprintId)
- getBacklog(projectId)
- listCustomFields(projectId)
- listLabels(projectId)
- addComment(recordId, data)
- uploadAttachment(recordId, file)
```

### 3.2 Check API Base URL

Verify the API base URL in `frontend/src/services/api.ts`:

```typescript
// Should match your backend URL
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1'
```

Update `.env` file if needed:

```bash
VITE_API_URL=http://localhost:8080/api/v1
```

## Part 4: Frontend Component Fixes

### 4.1 ProjectBoardPage Updates

The `ProjectBoardPage.tsx` has been updated to:
- Fetch workflow and statuses from API
- Display loading state while fetching
- Show error messages if API calls fail
- Display sprint board when active sprint exists
- Show default board layout when no active sprint

**Key changes:**
```typescript
// Now fetches from API instead of using hardcoded statuses
const sprintRes = await jiraService.getActiveSprint(projectId)
const workflowRes = await jiraService.getWorkflow(projectId)
const statusesRes = await jiraService.listWorkflowStatuses(workflowRes.data.id)
```

### 4.2 SprintBoard Component

The `SprintBoard.tsx` component:
- Uses `useJiraBoard` hook to fetch sprint data
- Displays records organized by status columns
- Implements drag-and-drop between columns
- Shows sprint metrics (total, completed, progress, days remaining)
- Handles record detail modal

**Verify drag-and-drop:**
1. Open project board
2. Drag a record from one column to another
3. Record should move and status should update

### 4.3 RecordCard Component

The `RecordCard.tsx` component displays:
- Issue type icon and label
- Record title
- Assignee avatar
- Labels as colored badges
- Due date with overdue indicator
- Attachment and comment count

## Part 5: Testing the Setup

### 5.1 Manual Testing Checklist

- [ ] Database migrations applied successfully
- [ ] Test data inserted (issue types, project, sprint, labels, custom fields)
- [ ] Backend API endpoints responding correctly
- [ ] Frontend can fetch project data
- [ ] ProjectBoardPage displays sprint board
- [ ] Drag-and-drop works between columns
- [ ] Record detail modal opens when clicking a card
- [ ] Comments can be added to records
- [ ] Attachments can be uploaded
- [ ] Labels can be added/removed from records
- [ ] Custom fields display correctly
- [ ] Backlog view shows unassigned records
- [ ] Sprint metrics display correctly

### 5.2 Automated Testing

Run the test suite:

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### 5.3 Integration Testing

Test the full flow:

1. **Create a record:**
   - Navigate to project board
   - Create a new record with issue type, title, description
   - Verify record appears in "To Do" column

2. **Transition record:**
   - Drag record from "To Do" to "In Progress"
   - Verify status updates in database
   - Verify record moves in UI

3. **Add comment:**
   - Click on record to open detail modal
   - Add a comment with @mention
   - Verify comment appears
   - Verify mention notification sent

4. **Upload attachment:**
   - Click on record to open detail modal
   - Upload a file
   - Verify file appears in attachments section
   - Verify file can be downloaded

5. **Add label:**
   - Click on record to open detail modal
   - Add a label
   - Verify label appears as colored badge on card

6. **Manage sprint:**
   - Create a new sprint
   - Assign records to sprint
   - Start sprint
   - Verify sprint board displays
   - Complete sprint
   - Verify incomplete records moved to backlog

## Part 6: Troubleshooting

### Issue: Project board shows "No Active Sprint"

**Solution:**
1. Check if sprint exists: `SELECT * FROM sprints WHERE project_id = '{projectId}' AND status = 'Active'`
2. If no sprint, create one via API or database
3. Verify sprint has valid start_date and end_date

### Issue: API returns 401 Unauthorized

**Solution:**
1. Verify authentication token is valid
2. Check token expiration
3. Verify Authorization header format: `Authorization: Bearer {token}`

### Issue: Drag-and-drop not working

**Solution:**
1. Verify @dnd-kit/core is installed: `npm list @dnd-kit/core`
2. Check browser console for errors
3. Verify SprintBoard component is rendering correctly
4. Check that statuses are being fetched from API

### Issue: Custom fields not displaying

**Solution:**
1. Verify custom fields exist in database: `SELECT * FROM custom_fields WHERE project_id = '{projectId}'`
2. Verify custom field values are created for records
3. Check RecordDetailModal component is rendering custom fields
4. Verify field types are supported (text, dropdown, date, number, checkbox)

### Issue: Comments not saving

**Solution:**
1. Verify comment endpoint is implemented in backend
2. Check API response for errors
3. Verify user has permission to add comments
4. Check database for comment records

### Issue: Attachments not uploading

**Solution:**
1. Verify file upload endpoint is implemented
2. Check file size limit (should be 50MB)
3. Verify file type is supported
4. Check file storage path is writable
5. Verify multipart/form-data header is set

## Part 7: Performance Optimization

### 7.1 Database Indexes

Verify all indexes are created:

```sql
-- Check indexes
SELECT * FROM pg_indexes WHERE tablename IN (
  'issue_types', 'workflows', 'sprints', 'sprint_records',
  'comments', 'attachments', 'labels', 'custom_fields'
);
```

### 7.2 Query Optimization

For large projects, consider:
- Pagination for record lists
- Lazy loading for custom fields
- Caching for workflow and issue type schemes
- Full-text search indexes for record titles

### 7.3 Frontend Performance

- Use React.memo for RecordCard component
- Implement virtual scrolling for large record lists
- Lazy load record detail modal
- Cache API responses with React Query

## Part 8: Deployment

### 8.1 Docker Deployment

Build and run with Docker:

```bash
# Build backend
docker build -f backend/Dockerfile -t itsm-backend:latest backend/

# Build frontend
docker build -f frontend/Dockerfile -t itsm-frontend:latest frontend/

# Run with docker-compose
docker-compose up -d
```

### 8.2 Environment Variables

Set required environment variables:

```bash
# Backend
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=itsm_db
JWT_SECRET=your-secret-key

# Frontend
VITE_API_URL=http://backend:8080/api/v1
```

### 8.3 Database Migrations in Production

Run migrations before deploying:

```bash
# In production environment
migrate -path backend/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require" up
```

## Part 9: Monitoring and Logging

### 9.1 Backend Logging

Enable debug logging:

```bash
LOG_LEVEL=debug
```

### 9.2 Frontend Error Tracking

Implement error tracking:

```typescript
// In frontend error boundary
console.error('Error:', error)
// Send to error tracking service (Sentry, etc.)
```

### 9.3 Database Monitoring

Monitor database performance:

```sql
-- Check slow queries
SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) 
FROM pg_tables 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Part 10: Maintenance

### 10.1 Regular Backups

```bash
# Backup database
pg_dump -U postgres -d itsm_db > backup.sql

# Restore database
psql -U postgres -d itsm_db < backup.sql
```

### 10.2 Cleanup Old Data

```sql
-- Archive completed sprints older than 90 days
DELETE FROM sprints WHERE status = 'Completed' AND actual_end_date < NOW() - INTERVAL '90 days';

-- Archive old comments
DELETE FROM comments WHERE created_at < NOW() - INTERVAL '1 year';
```

### 10.3 Update Dependencies

```bash
# Backend
go get -u ./...

# Frontend
npm update
```

## Support and Resources

- **Backend Documentation**: See `backend/README.md`
- **Frontend Documentation**: See `frontend/README.md`
- **API Documentation**: See `.kiro/specs/jira-like-project-board/design.md`
- **Requirements**: See `.kiro/specs/jira-like-project-board/requirements.md`

## Next Steps

1. Complete all setup steps above
2. Run manual testing checklist
3. Deploy to staging environment
4. Run integration tests
5. Deploy to production
6. Monitor for issues
7. Gather user feedback
8. Iterate and improve

---

**Last Updated**: 2024
**Version**: 1.0
