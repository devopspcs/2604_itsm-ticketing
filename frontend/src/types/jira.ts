// Issue Types
export interface IssueType {
  id: string
  name: string
  icon: string
  description: string
  created_at: string
}

export interface IssueTypeScheme {
  id: string
  project_id: string
  name: string
  created_at: string
}

export interface IssueTypeSchemeItem {
  id: string
  scheme_id: string
  issue_type_id: string
  created_at: string
}

// Custom Fields
export interface CustomField {
  id: string
  project_id: string
  name: string
  field_type: 'text' | 'textarea' | 'dropdown' | 'multiselect' | 'date' | 'number' | 'checkbox'
  is_required: boolean
  created_at: string
}

export interface CustomFieldOption {
  id: string
  field_id: string
  option_value: string
  option_order: number
  created_at: string
}

export interface CustomFieldValue {
  id: string
  record_id: string
  field_id: string
  value: string
  created_at: string
  updated_at: string
}

// Workflows
export interface Workflow {
  id: string
  project_id: string
  name: string
  initial_status: string
  created_at: string
}

export interface WorkflowStatus {
  id: string
  workflow_id: string
  status_name: string
  status_order: number
  created_at: string
}

export interface WorkflowTransition {
  id: string
  workflow_id: string
  from_status_id: string
  to_status_id: string
  validation_rule?: string
  created_at: string
}

// Sprints
export interface Sprint {
  id: string
  project_id: string
  name: string
  goal: string
  start_date?: string
  end_date?: string
  status: 'Planned' | 'Active' | 'Completed'
  actual_start_date?: string
  actual_end_date?: string
  created_at: string
}

export interface SprintRecord {
  id: string
  sprint_id: string
  record_id: string
  priority: number
  created_at: string
}

export interface SprintMetrics {
  total_records: number
  completed_records: number
  completion_percentage: number
  days_remaining: number
}

// Comments
export interface Comment {
  id: string
  record_id: string
  author_id: string
  author_name: string
  author_avatar?: string
  text: string
  created_at: string
  updated_at: string
  edited_at?: string
}

export interface CommentMention {
  id: string
  comment_id: string
  mentioned_user_id: string
  created_at: string
}

// Attachments
export interface Attachment {
  id: string
  record_id: string
  file_name: string
  file_size: number
  file_type: string
  file_path: string
  uploader_id: string
  uploader_name: string
  created_at: string
}

// Labels
export interface Label {
  id: string
  project_id: string
  name: string
  color: string
  created_at: string
}

export interface RecordLabel {
  id: string
  record_id: string
  label_id: string
  created_at: string
}

// Search & Filters
export interface SearchFilters {
  issue_type_id?: string
  status?: string
  assignee_id?: string
  label_id?: string
  sprint_id?: string
  custom_fields?: Record<string, string>
  due_date_from?: string
  due_date_to?: string
}

export interface SavedFilter {
  id: string
  project_id: string
  name: string
  filters: SearchFilters
  created_at: string
}

// Extended ProjectRecord with Jira features
export interface JiraProjectRecord {
  id: string
  column_id: string
  project_id: string
  issue_type_id: string
  issue_type?: IssueType
  title: string
  description: string
  status: string
  assigned_to?: string
  assignees: string[]
  due_date?: string
  parent_record_id?: string
  position: number
  is_completed: boolean
  completed_at?: string
  created_by: string
  created_at: string
  updated_at: string
  custom_fields?: CustomFieldValue[]
  labels?: Label[]
  comments_count?: number
  attachments_count?: number
}

// Request/Response DTOs
export interface CreateCustomFieldRequest {
  name: string
  field_type: string
  is_required: boolean
  options?: string[]
}

export interface UpdateCustomFieldRequest {
  name?: string
  options?: string[]
}

export interface CreateWorkflowRequest {
  name: string
  initial_status: string
  statuses: string[]
}

export interface UpdateWorkflowRequest {
  statuses?: string[]
  transitions?: Array<{
    from_status: string
    to_status: string
  }>
}

export interface CreateSprintRequest {
  name: string
  goal?: string
  start_date?: string
  end_date: string
}

export interface AddCommentRequest {
  text: string
}

export interface UpdateCommentRequest {
  text: string
}

export interface CreateLabelRequest {
  name: string
  color: string
}

export interface BulkChangeStatusRequest {
  record_ids: string[]
  status_id: string
}

export interface BulkAssignToRequest {
  record_ids: string[]
  assignee_id: string
}

export interface BulkAddLabelRequest {
  record_ids: string[]
  label_id: string
}

export interface BulkDeleteRequest {
  record_ids: string[]
}

export interface ReorderBacklogRequest {
  record_ids: string[]
}

export interface BulkAssignToSprintRequest {
  sprint_id: string
  record_ids: string[]
}

export interface TransitionRecordRequest {
  to_status_id: string
}

export interface SaveFilterRequest {
  name: string
  filters: SearchFilters
}

export interface SearchRecordsRequest {
  query?: string
  filters?: SearchFilters
}
