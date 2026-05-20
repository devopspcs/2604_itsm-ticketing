import api from './api'
import type { Application, AppWithAccess } from '../types'

export const appService = {
  getMyApps: () => api.get<AppWithAccess[]>('/me/apps'),
  listApps: () => api.get<Application[]>('/applications'),
  createApp: (data: { name: string; code: string; description?: string; icon?: string; color?: string }) =>
    api.post<Application>('/applications', data),
  updateApp: (id: string, data: { name?: string; description?: string; icon?: string; color?: string; is_active?: boolean }) =>
    api.patch<Application>(`/applications/${id}`, data),
  deleteApp: (id: string) => api.delete(`/applications/${id}`),
  getAppUsers: (appId: string) => api.get(`/applications/${appId}/users`),
  grantAccess: (appId: string, data: { user_id: string; role: string }) =>
    api.post(`/applications/${appId}/access`, { ...data, app_id: appId }),
  revokeAccess: (appId: string, userId: string) =>
    api.delete(`/applications/${appId}/access/${userId}`),
  bulkGrantAccess: (appId: string, data: { user_ids: string[]; role: string }) =>
    api.post(`/applications/${appId}/access/bulk`, { ...data, app_id: appId }),
}
