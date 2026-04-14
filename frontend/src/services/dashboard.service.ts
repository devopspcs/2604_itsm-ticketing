import api from './api'
import type { DashboardStats } from '../types'

export const dashboardService = {
  getStats: (params?: { date_from?: string; date_to?: string }) =>
    api.get<DashboardStats>('/dashboard', { params }),
}
