import api from './api'
import type { Department, Division, Team } from '../types'

export const orgService = {
  listDivisions: () => api.get<Division[]>('/divisions'),
  createDivision: (data: { name: string; code: string }) => api.post<Division>('/divisions', data),
  updateDivision: (id: string, data: { name: string; code: string }) => api.patch<Division>(`/divisions/${id}`, data),
  deleteDivision: (id: string) => api.delete(`/divisions/${id}`),

  listDepartments: (divisionId?: string) => api.get<Department[]>('/departments', { params: divisionId ? { division_id: divisionId } : {} }),
  createDepartment: (data: { division_id: string; name: string; code: string }) => api.post<Department>('/departments', data),
  updateDepartment: (id: string, data: { division_id: string; name: string; code: string }) => api.patch<Department>(`/departments/${id}`, data),
  deleteDepartment: (id: string) => api.delete(`/departments/${id}`),

  listTeams: (departmentId?: string) => api.get<Team[]>('/teams', { params: departmentId ? { department_id: departmentId } : {} }),
  createTeam: (data: { department_id: string; name: string; email?: string }) => api.post<Team>('/teams', data),
  updateTeam: (id: string, data: { department_id: string; name: string; email?: string }) => api.patch<Team>(`/teams/${id}`, data),
  deleteTeam: (id: string) => api.delete(`/teams/${id}`),

  updateUserOrg: (userId: string, data: { department_id?: string; division_id?: string; team_id?: string; position?: string }) =>
    api.patch(`/users/${userId}/org`, data),
}
