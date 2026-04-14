import api from './api'
import type { ApprovalDecision } from '../types'

export const approvalService = {
  decide: (ticketId: string, decision: ApprovalDecision, comment?: string) =>
    api.post('/approvals/decide', { ticket_id: ticketId, decision, comment }),
}
