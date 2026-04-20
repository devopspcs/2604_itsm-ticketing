# Jira-like Project Board Troubleshooting Guide

## Common Issues and Solutions

### 1. Project Board Shows "No Active Sprint"

**Symptoms:**
- Project board page displays "No Active Sprint" message
- Sprint board is not visible
- Cannot see any records on the board

**Root Causes:**
1. No sprint exists for the project
2. Sprint exists but status is not "Active"
3. Sprint has invalid dates
4. API call to get active sprint is failing

**Solutions:**

**Step 1: Check if sprint exists**
```sql
SELECT * FROM sprints WHERE project_id = '{projectId}';
```

If no sprints exist, create one:
```sql
INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status, actual_start_date, created_at)
VALUES (
  gen_random_uuid(),
  '{projectId}',
  'Sprint 1',
  'Initial sprint',
  CURRENT_DATE,
  CURRENT_DATE + INTERVAL '14 days',
  'Active',
  CURRENT_DATE,
  NOW()
);
```

**Step 2: Check sprint status**
```sql
SELECT id, name, status, start_date, end_date FROM sprints WHERE project_id = '{projectId}';
```

If status is not "Active", update it:
```sql
UPDATE sprints SET status = 'Active', actual_start_date = CURRENT_DATE WHERE id = '{sprintId}';
```

**Step 3: Verify API endpoint**
```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/sprints/active \
  -H "Authorization: Bearer {token}"
```

Expected response should include sprint data.

**Step 4: Check browser console**
- Open browser DevTools (F12)
- Check Console tab for JavaScript errors
- Check Network tab to see if API call succeeded

---

### 2. API Returns 401 Unauthorized

**Symptoms:**
- API calls fail with 401 status
- "Unauthorized" error message in browser console
- Cannot fetch project data

**Root Causes:**
1. Missing or invalid authentication token
2. Token has expired
3. Token format is incorrect
4. Backend JWT secret doesn't match frontend

**Solutions:**

**Step 1: Verify token is being sent**
```bash
# Check if Authorization header is present
curl -v http://localhost:8080/api/v1/projects/{projectId}/issue-types \
  -H "Authorization: Bearer {token}"
```

Look for `Authorization: Bearer` in the request headers.

**Step 2: Check token validity**
```bash
# Decode JWT token (use jwt.io or similar)
# Verify:
# - Token is not expired (check 'exp' claim)
# - Token has correct format (3 parts separated by dots)
# - Token is for correct user
```

**Step 3: Verify backend JWT configuration**
```bash
# Check backend environment variables
echo $JWT_SECRET
echo $JWT_EXPIRATION
```

**Step 4: Get new token**
```bash
# Login to get new token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}'
```

**Step 5: Update frontend token**
- Clear browser localStorage: `localStorage.clear()`
- Login again to get new token
- Retry API call

---

### 3. Drag-and-Drop Not Working

**Symptoms:**
- Cannot drag records between columns
- Records don't move when dragged
- No visual feedback when dragging

**Root Causes:**
1. @dnd-kit library not installed
2. DndContext not properly configured
3. Droppable zones not set up correctly
4. CSS conflicts with drag-and-drop

**Solutions:**

**Step 1: Verify @dnd-kit is installed**
```bash
cd frontend
npm list @dnd-kit/core @dnd-kit/utilities
```

If not installed:
```bash
npm install @dnd-kit/core @dnd-kit/utilities @dnd-kit/sortable
```

**Step 2: Check SprintBoard component**
```typescript
// Verify DndContext is wrapping the board
<DndContext onDragEnd={handleDragEnd}>
  {/* Board columns */}
</DndContext>

// Verify useDroppable is used in columns
const { setNodeRef } = useDroppable({ id: status.id })
```

**Step 3: Check browser console**
- Look for errors related to @dnd-kit
- Check if handleDragEnd is being called
- Verify status IDs are unique

**Step 4: Test with simple example**
```typescript
// Add console.log to verify drag events
const handleDragEnd = (event: DragEndEvent) => {
  console.log('Drag ended:', event)
  // ... rest of logic
}
```

**Step 5: Clear browser cache**
```bash
# Hard refresh browser
Ctrl+Shift+R (Windows/Linux)
Cmd+Shift+R (Mac)
```

---

### 4. Custom Fields Not Displaying

**Symptoms:**
- Custom fields section is empty
- Custom field values not showing on records
- Cannot add/edit custom fields

**Root Causes:**
1. Custom fields not created in database
2. Custom field values not associated with records
3. Frontend not fetching custom fields
4. RecordDetailModal not rendering custom fields

**Solutions:**

**Step 1: Check if custom fields exist**
```sql
SELECT * FROM custom_fields WHERE project_id = '{projectId}';
```

If empty, create custom fields:
```sql
INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
VALUES (
  gen_random_uuid(),
  '{projectId}',
  'Priority',
  'dropdown',
  true,
  NOW()
);
```

**Step 2: Check custom field values**
```sql
SELECT * FROM custom_field_values WHERE record_id = '{recordId}';
```

If empty, create values for records:
```sql
INSERT INTO custom_field_values (id, record_id, field_id, value, created_at, updated_at)
SELECT
  gen_random_uuid(),
  pr.id,
  cf.id,
  'Medium',
  NOW(),
  NOW()
FROM project_records pr
CROSS JOIN custom_fields cf
WHERE pr.project_id = '{projectId}' AND cf.project_id = '{projectId}'
ON CONFLICT DO NOTHING;
```

**Step 3: Verify API endpoint**
```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/custom-fields \
  -H "Authorization: Bearer {token}"
```

Should return list of custom fields.

**Step 4: Check frontend component**
- Verify RecordDetailModal is calling `jiraService.listCustomFields()`
- Check if custom fields are being rendered in the modal
- Look for errors in browser console

**Step 5: Verify field types**
```typescript
// Supported field types
const fieldTypes = ['text', 'textarea', 'dropdown', 'multiselect', 'date', 'number', 'checkbox']

// Check if field type is supported
if (!fieldTypes.includes(field.field_type)) {
  console.error('Unsupported field type:', field.field_type)
}
```

---

### 5. Comments Not Saving

**Symptoms:**
- Cannot add comments to records
- Comments disappear after refresh
- Error message when trying to save comment

**Root Causes:**
1. Comment endpoint not implemented in backend
2. User doesn't have permission to add comments
3. Record ID is invalid
4. API request is malformed

**Solutions:**

**Step 1: Verify comment endpoint exists**
```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/comments \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"text": "Test comment"}'
```

Should return 200 OK with comment data.

**Step 2: Check if record exists**
```sql
SELECT * FROM project_records WHERE id = '{recordId}';
```

If not found, use correct record ID.

**Step 3: Verify user permissions**
```sql
SELECT * FROM project_members WHERE project_id = '{projectId}' AND user_id = '{userId}';
```

User must be a project member to add comments.

**Step 4: Check comment in database**
```sql
SELECT * FROM comments WHERE record_id = '{recordId}' ORDER BY created_at DESC;
```

**Step 5: Check backend logs**
```bash
# Look for errors in backend logs
tail -f backend.log | grep -i comment
```

**Step 6: Verify @mention parsing**
```sql
-- Check if mentions were parsed
SELECT * FROM comment_mentions WHERE comment_id = '{commentId}';
```

---

### 6. Attachments Not Uploading

**Symptoms:**
- File upload fails silently
- Error message: "File too large" or "Invalid file type"
- Attachment doesn't appear after upload

**Root Causes:**
1. File size exceeds limit (50MB)
2. File type not supported
3. Upload endpoint not implemented
4. File storage path not writable
5. Multipart form data not properly formatted

**Solutions:**

**Step 1: Check file size**
```bash
# Maximum file size is 50MB
ls -lh file.pdf
# If larger than 50MB, compress or split the file
```

**Step 2: Check file type**
```bash
# Supported file types
# Images: jpg, png, gif
# Documents: pdf, doc, docx, xls, xlsx
# Archives: zip, rar
# Text: txt, md, csv

# Check file type
file file.pdf
```

**Step 3: Verify upload endpoint**
```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/attachments \
  -H "Authorization: Bearer {token}" \
  -F "file=@file.pdf"
```

Should return 200 OK with attachment data.

**Step 4: Check file storage path**
```bash
# Verify upload directory exists and is writable
ls -la /path/to/uploads/
chmod 755 /path/to/uploads/
```

**Step 5: Check attachment in database**
```sql
SELECT * FROM attachments WHERE record_id = '{recordId}' ORDER BY created_at DESC;
```

**Step 6: Verify file was saved**
```bash
# Check if file exists in storage
ls -la /path/to/uploads/
```

---

### 7. Labels Not Appearing

**Symptoms:**
- Cannot add labels to records
- Labels section is empty
- Label colors not displaying

**Root Causes:**
1. Labels not created in database
2. Labels not associated with records
3. Frontend not fetching labels
4. RecordDetailModal not rendering labels

**Solutions:**

**Step 1: Check if labels exist**
```sql
SELECT * FROM labels WHERE project_id = '{projectId}';
```

If empty, create labels:
```sql
INSERT INTO labels (id, project_id, name, color, created_at)
VALUES (
  gen_random_uuid(),
  '{projectId}',
  'Frontend',
  '#3b82f6',
  NOW()
);
```

**Step 2: Check record labels**
```sql
SELECT * FROM record_labels WHERE record_id = '{recordId}';
```

**Step 3: Verify API endpoint**
```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/labels \
  -H "Authorization: Bearer {token}"
```

Should return list of labels.

**Step 4: Add label to record**
```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/labels/{labelId} \
  -H "Authorization: Bearer {token}"
```

**Step 5: Check frontend component**
- Verify LabelManager component is rendering
- Check if labels are being fetched
- Verify label colors are valid hex codes

---

### 8. Workflow Transitions Not Working

**Symptoms:**
- Cannot transition record to certain statuses
- Error: "Invalid status transition"
- Record status doesn't change

**Root Causes:**
1. Workflow transitions not defined
2. Invalid status ID
3. Transition validation rules not met
4. User doesn't have permission

**Solutions:**

**Step 1: Check workflow transitions**
```sql
SELECT wt.*, ws1.status_name as from_status, ws2.status_name as to_status
FROM workflow_transitions wt
JOIN workflow_statuses ws1 ON wt.from_status_id = ws1.id
JOIN workflow_statuses ws2 ON wt.to_status_id = ws2.id
WHERE wt.workflow_id = '{workflowId}';
```

**Step 2: Create missing transitions**
```sql
INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, created_at)
SELECT
  gen_random_uuid(),
  w.id,
  ws1.id,
  ws2.id,
  NOW()
FROM workflows w
JOIN workflow_statuses ws1 ON w.id = ws1.workflow_id AND ws1.status_name = 'To Do'
JOIN workflow_statuses ws2 ON w.id = ws2.workflow_id AND ws2.status_name = 'In Progress'
WHERE w.id = '{workflowId}'
ON CONFLICT DO NOTHING;
```

**Step 3: Check validation rules**
```sql
SELECT * FROM workflow_transitions WHERE id = '{transitionId}';
```

If validation_rule is set, verify it's being enforced.

**Step 4: Verify transition endpoint**
```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/transition \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"to_status_id": "{statusId}"}'
```

---

### 9. Sprint Metrics Not Calculating

**Symptoms:**
- Sprint metrics show 0 or incorrect values
- Progress bar not updating
- Days remaining shows wrong number

**Root Causes:**
1. Records not assigned to sprint
2. Completion status not updated
3. Sprint dates are invalid
4. Metrics calculation logic has bug

**Solutions:**

**Step 1: Check sprint records**
```sql
SELECT COUNT(*) as total, SUM(CASE WHEN is_completed THEN 1 ELSE 0 END) as completed
FROM project_records pr
JOIN sprint_records sr ON pr.id = sr.record_id
WHERE sr.sprint_id = '{sprintId}';
```

**Step 2: Assign records to sprint**
```sql
INSERT INTO sprint_records (id, sprint_id, record_id, priority, created_at)
SELECT
  gen_random_uuid(),
  '{sprintId}',
  pr.id,
  ROW_NUMBER() OVER (ORDER BY pr.created_at),
  NOW()
FROM project_records pr
WHERE pr.project_id = '{projectId}' AND pr.status != 'Done'
ON CONFLICT DO NOTHING;
```

**Step 3: Update record completion**
```sql
UPDATE project_records SET is_completed = true WHERE status = 'Done';
```

**Step 4: Check sprint dates**
```sql
SELECT id, name, start_date, end_date, status FROM sprints WHERE id = '{sprintId}';
```

Dates should be valid and end_date should be in the future for active sprints.

**Step 5: Verify metrics calculation**
```typescript
// Check calculation logic
const completedCount = records.filter(r => r.is_completed).length
const totalCount = records.length
const completionPercentage = (completedCount / totalCount) * 100
const daysRemaining = Math.ceil((new Date(sprint.end_date) - new Date()) / (1000 * 60 * 60 * 24))
```

---

### 10. Database Connection Errors

**Symptoms:**
- "Connection refused" error
- "Database connection failed"
- Backend cannot start

**Root Causes:**
1. PostgreSQL not running
2. Wrong database credentials
3. Database doesn't exist
4. Firewall blocking connection

**Solutions:**

**Step 1: Check if PostgreSQL is running**
```bash
# Linux/Mac
sudo systemctl status postgresql

# Or check if port 5432 is listening
netstat -an | grep 5432
```

**Step 2: Verify database credentials**
```bash
# Test connection
psql -h localhost -U postgres -d itsm_db

# If fails, check .env file
cat .env | grep DB_
```

**Step 3: Create database if missing**
```bash
psql -U postgres -c "CREATE DATABASE itsm_db;"
```

**Step 4: Check firewall**
```bash
# Allow PostgreSQL port
sudo ufw allow 5432/tcp
```

**Step 5: Check backend logs**
```bash
# Look for connection errors
tail -f backend.log | grep -i "connection\|database"
```

---

### 11. Frontend Build Errors

**Symptoms:**
- `npm run build` fails
- TypeScript compilation errors
- Module not found errors

**Root Causes:**
1. Missing dependencies
2. TypeScript errors
3. Incorrect import paths
4. Node version mismatch

**Solutions:**

**Step 1: Install dependencies**
```bash
cd frontend
npm install
```

**Step 2: Check TypeScript errors**
```bash
npm run type-check
```

**Step 3: Fix import paths**
```bash
# Check if all imports are correct
grep -r "from.*jira" src/
```

**Step 4: Check Node version**
```bash
node --version
# Should be 16+ for this project
```

**Step 5: Clear cache and rebuild**
```bash
rm -rf node_modules package-lock.json
npm install
npm run build
```

---

### 12. Backend Build Errors

**Symptoms:**
- `go build` fails
- Compilation errors
- Missing packages

**Root Causes:**
1. Go version mismatch
2. Missing dependencies
3. Syntax errors
4. Import path issues

**Solutions:**

**Step 1: Check Go version**
```bash
go version
# Should be 1.21+
```

**Step 2: Download dependencies**
```bash
cd backend
go mod download
go mod tidy
```

**Step 3: Check for syntax errors**
```bash
go build ./...
```

**Step 4: Run tests**
```bash
go test ./...
```

---

## Performance Issues

### Slow API Responses

**Solutions:**
1. Add database indexes
2. Implement pagination
3. Cache frequently accessed data
4. Optimize queries

### High Memory Usage

**Solutions:**
1. Implement connection pooling
2. Limit concurrent requests
3. Clear old data regularly
4. Monitor memory usage

### Slow Frontend

**Solutions:**
1. Implement virtual scrolling
2. Lazy load components
3. Optimize images
4. Use React.memo for expensive components

---

## Getting Help

If you encounter issues not covered here:

1. **Check logs:**
   - Backend: `tail -f backend.log`
   - Frontend: Browser DevTools Console
   - Database: PostgreSQL logs

2. **Enable debug mode:**
   - Backend: `LOG_LEVEL=debug`
   - Frontend: `VITE_DEBUG=true`

3. **Check documentation:**
   - API Reference: `JIRA_API_REFERENCE.md`
   - Setup Guide: `JIRA_BOARD_SETUP_GUIDE.md`
   - Requirements: `.kiro/specs/jira-like-project-board/requirements.md`

4. **Run verification script:**
   ```bash
   ./verify-jira-setup.sh
   ```

---

**Last Updated**: 2024
**Version**: 1.0
