# Jira-like Project Board API Reference

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer {token}
```

## Issue Types

### List Issue Types
Get all available issue types for a project.

```
GET /projects/{projectId}/issue-types
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Bug",
    "icon": "bug_report",
    "description": "A problem or defect",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Get Issue Type Scheme
Get the issue type scheme for a project.

```
GET /projects/{projectId}/issue-type-scheme
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440400",
  "project_id": "550e8400-e29b-41d4-a716-446655440100",
  "name": "Default Issue Type Scheme",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Create Issue Type Scheme
Create a new issue type scheme for a project.

```
POST /projects/{projectId}/issue-type-scheme
Content-Type: application/json

{
  "name": "Custom Scheme",
  "issue_type_ids": ["550e8400-e29b-41d4-a716-446655440001", "550e8400-e29b-41d4-a716-446655440002"]
}
```

## Custom Fields

### List Custom Fields
Get all custom fields for a project.

```
GET /projects/{projectId}/custom-fields
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440701",
    "project_id": "550e8400-e29b-41d4-a716-446655440100",
    "name": "Priority",
    "field_type": "dropdown",
    "is_required": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Create Custom Field
Create a new custom field for a project.

```
POST /projects/{projectId}/custom-fields
Content-Type: application/json

{
  "name": "Priority",
  "field_type": "dropdown",
  "is_required": true,
  "options": ["Low", "Medium", "High", "Critical"]
}
```

**Field Types:**
- `text` - Single-line text
- `textarea` - Multi-line text
- `dropdown` - Single-select dropdown
- `multiselect` - Multi-select dropdown
- `date` - Date picker
- `number` - Number input
- `checkbox` - Checkbox

### Update Custom Field
Update an existing custom field.

```
PATCH /projects/{projectId}/custom-fields/{fieldId}
Content-Type: application/json

{
  "name": "Priority",
  "options": ["Low", "Medium", "High", "Critical", "Blocker"]
}
```

### Delete Custom Field
Delete a custom field from a project.

```
DELETE /projects/{projectId}/custom-fields/{fieldId}
```

## Workflows

### Get Workflow
Get the workflow for a project.

```
GET /projects/{projectId}/workflow
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440200",
  "project_id": "550e8400-e29b-41d4-a716-446655440100",
  "name": "Default Workflow",
  "initial_status": "Backlog",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Create Workflow
Create a new workflow for a project.

```
POST /projects/{projectId}/workflow
Content-Type: application/json

{
  "name": "Custom Workflow",
  "initial_status": "Backlog",
  "statuses": ["Backlog", "To Do", "In Progress", "In Review", "Done"]
}
```

### Update Workflow
Update an existing workflow.

```
PATCH /projects/{projectId}/workflow
Content-Type: application/json

{
  "statuses": ["Backlog", "To Do", "In Progress", "In Review", "Testing", "Done"],
  "transitions": [
    {"from_status": "Backlog", "to_status": "To Do"},
    {"from_status": "To Do", "to_status": "In Progress"}
  ]
}
```

### List Workflow Statuses
Get all statuses for a workflow.

```
GET /workflows/{workflowId}/statuses
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440201",
    "workflow_id": "550e8400-e29b-41d4-a716-446655440200",
    "status_name": "Backlog",
    "status_order": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

## Sprints

### List Sprints
Get all sprints for a project.

```
GET /projects/{projectId}/sprints
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440500",
    "project_id": "550e8400-e29b-41d4-a716-446655440100",
    "name": "Sprint 1",
    "goal": "Complete initial setup",
    "start_date": "2024-01-01",
    "end_date": "2024-01-14",
    "status": "Active",
    "actual_start_date": "2024-01-01",
    "actual_end_date": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Get Active Sprint
Get the currently active sprint for a project.

```
GET /projects/{projectId}/sprints/active
```

### Create Sprint
Create a new sprint for a project.

```
POST /projects/{projectId}/sprints
Content-Type: application/json

{
  "name": "Sprint 2",
  "goal": "Implement user authentication",
  "start_date": "2024-01-15",
  "end_date": "2024-01-28"
}
```

### Update Sprint
Update an existing sprint (including status changes).

```
PATCH /projects/{projectId}/sprints/{sprintId}
Content-Type: application/json

{
  "name": "Sprint 2",
  "goal": "Implement user authentication",
  "status": "Active"
}
```

### Start Sprint
Start a sprint (change status from Planned to Active).

```
POST /sprints/{sprintId}/start
```

### Complete Sprint
Complete a sprint (change status from Active to Completed).

```
POST /sprints/{sprintId}/complete
```

**Response:**
```json
{
  "total_records": 10,
  "completed_records": 8,
  "completion_percentage": 80,
  "days_remaining": 0
}
```

### Get Sprint Records
Get all records assigned to a sprint.

```
GET /sprints/{sprintId}/records
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "project_id": "550e8400-e29b-41d4-a716-446655440100",
    "issue_type_id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "Implement login page",
    "description": "Create login page with email/password",
    "status": "In Progress",
    "assigned_to": "user-id",
    "due_date": "2024-01-14",
    "is_completed": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

## Backlog

### Get Backlog
Get all records not assigned to any sprint.

```
GET /projects/{projectId}/backlog
```

### Reorder Backlog
Reorder backlog records to change priority.

```
PATCH /projects/{projectId}/backlog/reorder
Content-Type: application/json

{
  "record_ids": ["id1", "id2", "id3"]
}
```

### Bulk Assign to Sprint
Assign multiple records to a sprint.

```
POST /projects/{projectId}/backlog/assign-sprint
Content-Type: application/json

{
  "sprint_id": "550e8400-e29b-41d4-a716-446655440500",
  "record_ids": ["id1", "id2", "id3"]
}
```

## Records

### Transition Record
Change the status of a record.

```
POST /records/{recordId}/transition
Content-Type: application/json

{
  "to_status_id": "550e8400-e29b-41d4-a716-446655440202"
}
```

### Assign Record to Sprint
Assign a record to a sprint or remove from sprint.

```
PATCH /projects/{projectId}/records/{recordId}/sprint
Content-Type: application/json

{
  "sprint_id": "550e8400-e29b-41d4-a716-446655440500"
}
```

To remove from sprint, set `sprint_id` to `null`.

## Comments

### Add Comment
Add a comment to a record.

```
POST /records/{recordId}/comments
Content-Type: application/json

{
  "text": "This is a comment with @mention @user-id"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "record_id": "550e8400-e29b-41d4-a716-446655440001",
  "author_id": "user-id",
  "author_name": "John Doe",
  "text": "This is a comment",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### List Comments
Get all comments for a record.

```
GET /records/{recordId}/comments
```

### Update Comment
Update an existing comment.

```
PATCH /comments/{commentId}
Content-Type: application/json

{
  "text": "Updated comment text"
}
```

### Delete Comment
Delete a comment.

```
DELETE /comments/{commentId}
```

## Attachments

### Upload Attachment
Upload a file to a record.

```
POST /records/{recordId}/attachments
Content-Type: multipart/form-data

file: <binary file data>
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "record_id": "550e8400-e29b-41d4-a716-446655440001",
  "file_name": "document.pdf",
  "file_size": 1024000,
  "file_type": "application/pdf",
  "file_path": "/uploads/document.pdf",
  "uploader_id": "user-id",
  "uploader_name": "John Doe",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### List Attachments
Get all attachments for a record.

```
GET /records/{recordId}/attachments
```

### Delete Attachment
Delete an attachment.

```
DELETE /records/{recordId}/attachments/{attachmentId}
```

## Labels

### List Labels
Get all labels for a project.

```
GET /projects/{projectId}/labels
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440601",
    "project_id": "550e8400-e29b-41d4-a716-446655440100",
    "name": "Frontend",
    "color": "#3b82f6",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Create Label
Create a new label for a project.

```
POST /projects/{projectId}/labels
Content-Type: application/json

{
  "name": "Frontend",
  "color": "#3b82f6"
}
```

### Update Label
Update an existing label.

```
PATCH /projects/{projectId}/labels/{labelId}
Content-Type: application/json

{
  "name": "Frontend",
  "color": "#3b82f6"
}
```

### Delete Label
Delete a label.

```
DELETE /projects/{projectId}/labels/{labelId}
```

### Add Label to Record
Add a label to a record.

```
POST /records/{recordId}/labels/{labelId}
```

### Remove Label from Record
Remove a label from a record.

```
DELETE /records/{recordId}/labels/{labelId}
```

## Bulk Operations

### Bulk Change Status
Change status for multiple records.

```
POST /projects/{projectId}/bulk/change-status
Content-Type: application/json

{
  "record_ids": ["id1", "id2", "id3"],
  "status_id": "550e8400-e29b-41d4-a716-446655440202"
}
```

### Bulk Assign
Assign multiple records to a user.

```
POST /projects/{projectId}/bulk/assign
Content-Type: application/json

{
  "record_ids": ["id1", "id2", "id3"],
  "assignee_id": "user-id"
}
```

### Bulk Add Label
Add a label to multiple records.

```
POST /projects/{projectId}/bulk/add-label
Content-Type: application/json

{
  "record_ids": ["id1", "id2", "id3"],
  "label_id": "550e8400-e29b-41d4-a716-446655440601"
}
```

### Bulk Delete
Delete multiple records.

```
POST /projects/{projectId}/bulk/delete
Content-Type: application/json

{
  "record_ids": ["id1", "id2", "id3"]
}
```

## Search and Filters

### Search Records
Search and filter records.

```
GET /projects/{projectId}/search?q=search+term&issue_type=bug&status=in-progress&assignee=user-id&label=frontend
```

**Query Parameters:**
- `q` - Search query (searches title and description)
- `issue_type` - Filter by issue type ID
- `status` - Filter by status
- `assignee` - Filter by assignee ID
- `label` - Filter by label ID
- `sprint` - Filter by sprint ID
- `due_date_from` - Filter by due date (from)
- `due_date_to` - Filter by due date (to)

### Save Filter
Save a search filter.

```
POST /projects/{projectId}/filters
Content-Type: application/json

{
  "name": "My Frontend Bugs",
  "filters": {
    "issue_type_id": "550e8400-e29b-41d4-a716-446655440001",
    "label_id": "550e8400-e29b-41d4-a716-446655440601"
  }
}
```

### List Saved Filters
Get all saved filters for a project.

```
GET /projects/{projectId}/filters
```

### Delete Saved Filter
Delete a saved filter.

```
DELETE /projects/{projectId}/filters/{filterId}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request",
  "message": "Field 'name' is required"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing authentication token"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden",
  "message": "You don't have permission to access this resource"
}
```

### 404 Not Found
```json
{
  "error": "Not Found",
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred"
}
```

## Rate Limiting

API requests are rate-limited to prevent abuse:
- 100 requests per minute per user
- 1000 requests per hour per user

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1609459200
```

## Pagination

List endpoints support pagination:

```
GET /projects/{projectId}/sprints?limit=10&offset=0
```

**Query Parameters:**
- `limit` - Number of items to return (default: 20, max: 100)
- `offset` - Number of items to skip (default: 0)

**Response:**
```json
{
  "data": [...],
  "total": 50,
  "limit": 10,
  "offset": 0
}
```

## Webhooks

Webhooks can be configured to receive notifications for events:

**Supported Events:**
- `record.created`
- `record.updated`
- `record.deleted`
- `record.status_changed`
- `comment.created`
- `comment.updated`
- `comment.deleted`
- `attachment.uploaded`
- `attachment.deleted`
- `sprint.started`
- `sprint.completed`

See webhook configuration in project settings.

---

**Last Updated**: 2024
**Version**: 1.0
