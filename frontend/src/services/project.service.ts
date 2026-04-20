import api from './api'
import type {
  Project, ProjectDetail, ProjectColumn, ProjectRecord, ProjectActivityLog, ProjectHomeData,
  ReportsSummary, VelocityDataPoint, BurndownData,
  ReleaseWithProgress, Release, CreateReleaseRequest,
  ComponentWithCount, Component, CreateComponentRequest,
  PaginatedRecords, IssuesFilterParams,
  PaginatedActivityLogs, ActivityLogFilterParams,
} from '../types/project'

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

  // Reports
  getReportsSummary: (projectId: string) =>
    api.get<ReportsSummary>(`/projects/${projectId}/reports/summary`),
  getReportsVelocity: (projectId: string) =>
    api.get<VelocityDataPoint[]>(`/projects/${projectId}/reports/velocity`),
  getReportsBurndown: (projectId: string) =>
    api.get<BurndownData>(`/projects/${projectId}/reports/burndown`),

  // Releases
  listReleases: (projectId: string) =>
    api.get<ReleaseWithProgress[]>(`/projects/${projectId}/releases`),
  createRelease: (projectId: string, data: CreateReleaseRequest) =>
    api.post<Release>(`/projects/${projectId}/releases`, data),
  updateRelease: (projectId: string, releaseId: string, data: Partial<Release>) =>
    api.patch<Release>(`/projects/${projectId}/releases/${releaseId}`, data),
  deleteRelease: (projectId: string, releaseId: string) =>
    api.delete(`/projects/${projectId}/releases/${releaseId}`),

  // Components
  listComponents: (projectId: string) =>
    api.get<ComponentWithCount[]>(`/projects/${projectId}/components`),
  createComponent: (projectId: string, data: CreateComponentRequest) =>
    api.post<Component>(`/projects/${projectId}/components`, data),
  updateComponent: (projectId: string, componentId: string, data: Partial<Component>) =>
    api.patch<Component>(`/projects/${projectId}/components/${componentId}`, data),
  deleteComponent: (projectId: string, componentId: string) =>
    api.delete(`/projects/${projectId}/components/${componentId}`),

  // Issues
  listIssues: (projectId: string, params: IssuesFilterParams) =>
    api.get<PaginatedRecords>(`/projects/${projectId}/issues`, { params }),

  // Activity Log
  listActivityLog: (projectId: string, params: ActivityLogFilterParams) =>
    api.get<PaginatedActivityLogs>(`/projects/${projectId}/activity-log`, { params }),

  // Users
  listUsers: () =>
    api.get<{ id: string; name: string; email: string }[]>('/users/list'),
}
