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
