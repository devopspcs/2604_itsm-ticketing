import type { Ticket, TicketStatus } from './index'

export const KANBAN_COLUMNS: { status: TicketStatus; title: string }[] = [
  { status: 'open', title: 'Open' },
  { status: 'in_progress', title: 'In Progress' },
  { status: 'pending_approval', title: 'Pending Approval' },
  { status: 'approved', title: 'Approved' },
  { status: 'done', title: 'Done' },
]

export type ColumnState = Record<TicketStatus, Ticket[]>
export type ColumnLoadingState = Record<TicketStatus, boolean>
