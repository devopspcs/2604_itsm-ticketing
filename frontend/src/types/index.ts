export type Role = 'user' | 'approver' | 'admin'
export type TicketType = 'change_request' | 'incident' | 'helpdesk_request'
export type Priority = 'low' | 'medium' | 'high' | 'critical'
export type TicketStatus = 'open' | 'in_progress' | 'pending_approval' | 'approved' | 'rejected' | 'done'
export type ApprovalDecision = 'approved' | 'rejected'

export type Position = 'division_manager' | 'manager' | 'leader' | 'staff'

export interface Department {
  id: string
  name: string
  code: string
  created_at: string
  updated_at: string
}

export interface Division {
  id: string
  department_id: string
  name: string
  code: string
  created_at: string
  updated_at: string
}

export interface Team {
  id: string
  division_id: string
  name: string
  created_at: string
  updated_at: string
}

export interface User {
  id: string
  full_name: string
  email: string
  role: Role
  is_active: boolean
  created_at: string
  updated_at: string
  department_id?: string
  division_id?: string
  team_id?: string
  position?: Position
}

export interface Ticket {
  id: string
  title: string
  description: string
  type: TicketType
  category: string
  priority: Priority
  status: TicketStatus
  created_by: string
  assigned_to?: string
  created_at: string
  updated_at: string
}

export interface Approval {
  id: string
  ticket_id: string
  approver_id: string
  level: number
  decision?: ApprovalDecision
  comment: string
  decided_at?: string
}

export interface ActivityLog {
  id: string
  ticket_id: string
  actor_id: string
  action: string
  old_value?: string
  new_value?: string
  created_at: string
}

export interface Notification {
  id: string
  user_id: string
  ticket_id: string
  message: string
  is_read: boolean
  created_at: string
}

export interface DashboardStats {
  total_tickets: number
  by_status: Record<TicketStatus, number>
  by_type: Record<TicketType, number>
  by_priority: Record<Priority, number>
  recent_tickets: Ticket[]
}

export interface PaginatedTickets {
  tickets: Ticket[]
  total: number
  page: number
  page_size: number
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}
