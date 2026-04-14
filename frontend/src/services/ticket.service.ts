import api from './api'
import type { Ticket, PaginatedTickets, Approval, ActivityLog } from '../types'

export const ticketService = {
  list: (params?: Record<string, string | number>) =>
    api.get<PaginatedTickets>('/tickets', { params }),

  get: (id: string) => api.get<Ticket>(`/tickets/${id}`),

  create: (data: Partial<Ticket>) => api.post<Ticket>('/tickets', data),

  update: (id: string, data: Partial<Ticket>) =>
    api.patch<Ticket>(`/tickets/${id}`, data),

  submit: (id: string) => api.post(`/tickets/${id}/submit`),

  assign: (id: string, assigneeId: string) =>
    api.post(`/tickets/${id}/assign`, { assignee_id: assigneeId }),

  getApprovals: (id: string) =>
    api.get<Approval[]>(`/tickets/${id}/approvals`),

  getActivities: (id: string) =>
    api.get<ActivityLog[]>(`/tickets/${id}/activities`),
}
