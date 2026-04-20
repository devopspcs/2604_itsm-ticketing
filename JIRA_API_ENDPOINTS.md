# Jira-like Project Board API Endpoints

## Base URL
```
/api/v1/projects/{id}
```

All endpoints require authentication via JWT token in Authorization header.

## Issue Type Endpoints

### List Issue Types
```
GET /issue-types
Response: []*IssueType
```

### Get Issue Type Scheme
```
GET /issue-type-scheme
Response: *IssueTypeScheme
```

### Create Issue Type Scheme
```
POST /issue-type-scheme
Request: CreateIssueTypeSchemeRequest
Response: *IssueTypeScheme
```

## Custom Field Endpoints

### Create Custom Field
```
POST /custom-fields
Request: CreateCustomFieldRequest
Response: *CustomField
```

### List Custom Fields
```
GET /custom-fields
Response: []*CustomField
```

### Update Custom Field
```
PUT /custom-fields/{fieldId}
Request: UpdateCustomFieldRequest
Response: *CustomField
```

### Delete Custom Field
```
DELETE /custom-fields/{fieldId}
Response: 204 No Content
```

## Workflow Endpoints

### Get Workflow
```
GET /workflow
Response: *Workflow
```

### Create Workflow
```
POST /workflow
Request: CreateWorkflowRequest
Response: *Workflow
```

### Update Workflow
```
PUT /workflow/{workflowId}
Request: UpdateWorkflowRequest
Response: 204 No Content
```

### Transition Record
```
POST /records/{recordId}/transition
Request: TransitionRecordRequest
Response: 204 No Content
```

## Sprint Endpoints

### Create Sprint
```
POST /sprints
Request: CreateSprintRequest
Response: *Sprint
```

### List Sprints
```
GET /sprints
Response: []*Sprint
```

### Get Active Sprint
```
GET /sprints/active
Response: *Sprint
```

### Start Sprint
```
POST /sprints/{sprintId}/start
Response: *Sprint
```

### Complete Sprint
```
POST /sprints/{sprintId}/complete
Response: *SprintMetrics
```

### Get Sprint Records
```
GET /sprints/{sprintId}/records
Response: []*ProjectRecord
```

## Backlog Endpoints

### Get Backlog
```
GET /backlog
Response: []*ProjectRecord
```

### Reorder Backlog
```
PUT /backlog/reorder
Request: ReorderBacklogRequest
Response: 204 No Content
```

### Bulk Assign to Sprint
```
POST /backlog/assign-sprint
Request: BulkAssignToSprintRequest
Response: 204 No Content
```

## Comment Endpoints

### Add Comment
```
POST /records/{recordId}/comments
Request: AddCommentRequest
Response: *Comment
```

### List Comments
```
GET /records/{recordId}/comments
Response: []*Comment
```

### Update Comment
```
PUT /comments/{commentId}
Request: UpdateCommentRequest
Response: *Comment
```

### Delete Comment
```
DELETE /comments/{commentId}
Response: 204 No Content
```

## Attachment Endpoints

### Upload Attachment
```
POST /records/{recordId}/attachments
Content-Type: multipart/form-data
Form Field: file
Response: *Attachment
```

### List Attachments
```
GET /records/{recordId}/attachments
Response: []*Attachment
```

### Delete Attachment
```
DELETE /attachments/{attachmentId}
Response: 204 No Content
```

## Label Endpoints

### Create Label
```
POST /labels
Request: CreateLabelRequest
Response: *Label
```

### List Labels
```
GET /labels
Response: []*Label
```

### Add Label to Record
```
POST /records/{recordId}/labels/{labelId}
Response: 204 No Content
```

### Remove Label from Record
```
DELETE /records/{recordId}/labels/{labelId}
Response: 204 No Content
```

### Delete Label
```
DELETE /labels/{labelId}
Response: 204 No Content
```

## Bulk Operation Endpoints

### Bulk Change Status
```
POST /bulk/change-status
Request: BulkChangeStatusRequest
Response: 204 No Content
```

### Bulk Assign To
```
POST /bulk/assign
Request: BulkAssignToRequest
Response: 204 No Content
```

### Bulk Add Label
```
POST /bulk/add-label
Request: BulkAddLabelRequest
Response: 204 No Content
```

### Bulk Delete
```
POST /bulk/delete
Request: BulkDeleteRequest
Response: 204 No Content
```

## Search Endpoints

### Search Records
```
GET /search?q=query
Query Parameters:
  - q: search query
Response: []*ProjectRecord
```

### Save Filter
```
POST /filters
Request: SaveFilterRequest
Response: 201 Created
```

### List Saved Filters
```
GET /filters
Response: []*SavedFilter
```

## Request/Response Types

### CreateIssueTypeSchemeRequest
```json
{
  "name": "string",
  "issue_type_ids": ["uuid", "uuid"]
}
```

### CreateCustomFieldRequest
```json
{
  "name": "string",
  "field_type": "text|textarea|dropdown|multiselect|date|number|checkbox",
  "is_required": boolean,
  "options": ["option1", "option2"]
}
```

### CreateWorkflowRequest
```json
{
  "name": "string",
  "initial_status": "string",
  "statuses": ["status1", "status2"]
}
```

### CreateSprintRequest
```json
{
  "name": "string",
  "goal": "string",
  "start_date": "2026-04-19T00:00:00Z",
  "end_date": "2026-04-26T00:00:00Z"
}
```

### AddCommentRequest
```json
{
  "text": "string"
}
```

### CreateLabelRequest
```json
{
  "name": "string",
  "color": "#FF0000"
}
```

### BulkChangeStatusRequest
```json
{
  "record_ids": ["uuid", "uuid"],
  "status_id": "uuid"
}
```

### BulkAssignToSprintRequest
```json
{
  "sprint_id": "uuid",
  "record_ids": ["uuid", "uuid"]
}
```

### SaveFilterRequest
```json
{
  "name": "string",
  "filters": {
    "issue_type": "uuid",
    "status": "string",
    "assignee": "uuid",
    "label": "uuid",
    "sprint": "uuid",
    "due_date_from": "2026-04-19T00:00:00Z",
    "due_date_to": "2026-04-26T00:00:00Z",
    "custom_fields": {"field_id": "value"}
  }
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "validation_error",
  "message": "Invalid request parameters"
}
```

### 401 Unauthorized
```json
{
  "error": "unauthorized",
  "message": "Missing or invalid authentication token"
}
```

### 403 Forbidden
```json
{
  "error": "forbidden",
  "message": "You don't have permission to perform this action"
}
```

### 404 Not Found
```json
{
  "error": "not_found",
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal_error",
  "message": "An unexpected error occurred"
}
```

## Authentication

All endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

The token is obtained by logging in:

```
POST /api/v1/auth/login
Request: {
  "email": "user@example.com",
  "password": "password"
}
Response: {
  "access_token": "token",
  "refresh_token": "token",
  "expires_in": 3600
}
```

## Rate Limiting

All endpoints are rate-limited to prevent abuse:
- 100 requests per minute per user
- 1000 requests per minute per IP

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1713607200
```
