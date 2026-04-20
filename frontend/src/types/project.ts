export interface Project {
  id: string
  name: string
  icon_color: string
  created_by: string
  created_at: string
  updated_at: string
}

export interface ProjectColumn {
  id: string
  project_id: string
  name: string
  position: number
  created_at: string
}

export interface ProjectRecord {
  id: string
  column_id: string
  project_id: string
  title: string
  description: string
  assigned_to?: string
  assignees: string[]
  due_date?: string
  position: number
  is_completed: boolean
  completed_at?: string
  created_by: string
  created_at: string
  updated_at: string
  // Jira-like features (optional for backward compatibility)
  issue_type_id?: string
  status?: string
  parent_record_id?: string
}

export interface ProjectActivityLog {
  id: string
  project_id: string
  record_id?: string
  actor_id: string
  action: string
  detail: string
  created_at: string
}

export interface ProjectDetail extends Project {
  columns: (ProjectColumn & { records: ProjectRecord[] })[]
}

export interface ProjectHomeData {
  overdue_count: number
  recent_activities: ProjectActivityLog[]
}

// --- Reports ---
export interface ReportsSummary {
  total_records: number
  completed_count: number
  open_count: number
  by_status: Record<string, number>
}

export interface VelocityDataPoint {
  sprint_name: string
  total_records: number
  completed_count: number
}

export interface BurndownData {
  sprint_name: string
  start_date: string
  end_date: string
  total_count: number
  done_count: number
  has_active: boolean
}

// --- Releases ---
export interface Release {
  id: string
  project_id: string
  name: string
  version: string
  description: string
  start_date: string
  release_date: string
  status: string
  created_by: string
  created_at: string
  updated_at: string
}

export interface ReleaseWithProgress extends Release {
  total_records: number
  completed_count: number
  progress_percent: number
}

export interface CreateReleaseRequest {
  name: string
  version: string
  description?: string
  start_date?: string
  release_date?: string
  status?: string
}

// --- Components ---
export interface Component {
  id: string
  project_id: string
  name: string
  description: string
  lead_user_id: string | null
  created_at: string
  updated_at: string
}

export interface ComponentWithCount extends Component {
  record_count: number
}

export interface CreateComponentRequest {
  name: string
  description?: string
  lead_user_id?: string
}

// --- Issues ---
export interface PaginatedRecords {
  records: ProjectRecord[]
  total: number
  page: number
  page_size: number
}

export interface IssuesFilterParams {
  search?: string
  status_id?: string
  assignee_id?: string
  issue_type?: string
  label_id?: string
  page?: number
  page_size?: number
}

// --- Activity Log ---
export interface ActivityLogEntry {
  id: string
  project_id: string
  record_id: string | null
  actor_id: string
  action: string
  detail: string
  created_at: string
}

export interface PaginatedActivityLogs {
  logs: ActivityLogEntry[]
  total: number
  page: number
  page_size: number
}

export interface ActivityLogFilterParams {
  action_type?: string
  page?: number
  page_size?: number
}
