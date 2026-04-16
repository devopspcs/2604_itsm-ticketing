import api from './api'
import type { Project, ProjectDetail, ProjectColumn, ProjectRecord, ProjectActivityLog, ProjectHomeData } from '../types/project'

export const projectService = {
  // Projects
  list: () => api.get<Project[]>('/projects'),
  get: (id: string) => api.get<ProjectDetail>(`/projects/${id}`),
  create: (data: { name: string; icon_color: string }) => api.post<Project>('/projects', data),
  update: (id: string, data: Partial<Project>) => api.patch<Project>(`/projects/${id}`, data),
  delete: (id: string) => api.delete(`/projects/${id}`),

  // Columns
  createColumn: (projectId: string, data: { name: string }) =>
    api.post<ProjectColumn>(`/projects/${projectId}/columns`, data),
  updateColumn: (projectId: string, columnId: string, data: { name: string }) =>
    api.patch<ProjectColumn>(`/projects/${projectId}/columns/${columnId}`, data),
  deleteColumn: (projectId: string, columnId: string) =>
    api.delete(`/projects/${projectId}/columns/${columnId}`),
  reorderColumns: (projectId: string, columnIds: string[]) =>
    api.patch(`/projects/${projectId}/columns/reorder`, { column_ids: columnIds }),

  // Records
  createRecord: (projectId: string, data: Partial<ProjectRecord>) =>
    api.post<ProjectRecord>(`/projects/${projectId}/records`, data),
  getRecord: (projectId: string, recordId: string) =>
    api.get<ProjectRecord>(`/projects/${projectId}/records/${recordId}`),
  updateRecord: (projectId: string, recordId: string, data: Partial<ProjectRecord>) =>
    api.patch<ProjectRecord>(`/projects/${projectId}/records/${recordId}`, data),
  deleteRecord: (projectId: string, recordId: string) =>
    api.delete(`/projects/${projectId}/records/${recordId}`),
  moveRecord: (projectId: string, recordId: string, data: { target_column_id: string; position: number }) =>
    api.patch(`/projects/${projectId}/records/${recordId}/move`, data),
  completeRecord: (projectId: string, recordId: string) =>
    api.patch(`/projects/${projectId}/records/${recordId}/complete`),

  // Views
  getHome: () => api.get<ProjectHomeData>('/projects/home'),
  getCalendar: (month: number, year: number) =>
    api.get<ProjectRecord[]>('/projects/calendar', { params: { month, year } }),
  getActivities: (projectId: string) =>
    api.get<ProjectActivityLog[]>(`/projects/${projectId}/activities`),

  // Members
  listMembers: (projectId: string) =>
    api.get<{ project_id: string; user_id: string; role: string; created_at: string }[]>(`/projects/${projectId}/members`),
  inviteMember: (projectId: string, userId: string) =>
    api.post(`/projects/${projectId}/members`, { user_id: userId }),
  removeMember: (projectId: string, userId: string) =>
    api.delete(`/projects/${projectId}/members/${userId}`),

  // Comments
  addComment: (projectId: string, recordId: string, text: string) =>
    api.post(`/projects/${projectId}/records/${recordId}/comments`, { text }),
}
