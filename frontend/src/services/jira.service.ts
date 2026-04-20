import api from './api'
import type {
  IssueType,
  IssueTypeScheme,
  CustomField,
  Workflow,
  Sprint,
  Comment,
  Attachment,
  Label,
  SprintMetrics,
  CreateCustomFieldRequest,
  UpdateCustomFieldRequest,
  CreateWorkflowRequest,
  UpdateWorkflowRequest,
  CreateSprintRequest,
  AddCommentRequest,
  UpdateCommentRequest,
  CreateLabelRequest,
  BulkChangeStatusRequest,
  BulkAssignToRequest,
  BulkAddLabelRequest,
  BulkDeleteRequest,
  ReorderBacklogRequest,
  BulkAssignToSprintRequest,
  TransitionRecordRequest,
  SaveFilterRequest,
  SearchFilters,
  SavedFilter,
  JiraProjectRecord,
} from '../types/jira'

export const jiraService = {
  // Issue Types
  listIssueTypes: (projectId: string) =>
    api.get<IssueType[]>(`/projects/${projectId}/issue-types`),

  getIssueTypeScheme: (projectId: string) =>
    api.get<IssueTypeScheme>(`/projects/${projectId}/issue-type-scheme`),

  createIssueTypeScheme: (projectId: string, data: { name: string; issue_type_ids: string[] }) =>
    api.post<IssueTypeScheme>(`/projects/${projectId}/issue-type-scheme`, data),

  // Custom Fields
  createCustomField: (projectId: string, data: CreateCustomFieldRequest) =>
    api.post<CustomField>(`/projects/${projectId}/custom-fields`, data),

  listCustomFields: (projectId: string) =>
    api.get<CustomField[]>(`/projects/${projectId}/custom-fields`),

  updateCustomField: (projectId: string, fieldId: string, data: UpdateCustomFieldRequest) =>
    api.patch<CustomField>(`/projects/${projectId}/custom-fields/${fieldId}`, data),

  deleteCustomField: (projectId: string, fieldId: string) =>
    api.delete(`/projects/${projectId}/custom-fields/${fieldId}`),

  // Workflows
  getWorkflow: (projectId: string) =>
    api.get<Workflow>(`/projects/${projectId}/workflow`),

  createWorkflow: (projectId: string, data: CreateWorkflowRequest) =>
    api.post<Workflow>(`/projects/${projectId}/workflow`, data),

  updateWorkflow: (projectId: string, data: UpdateWorkflowRequest) =>
    api.patch(`/projects/${projectId}/workflow`, data),

  listWorkflowStatuses: (projectId: string, workflowId: string) =>
    api.get(`/projects/${projectId}/workflows/${workflowId}/statuses`),

  transitionRecord: (projectId: string, recordId: string, data: TransitionRecordRequest) =>
    api.post(`/projects/${projectId}/records/${recordId}/transition`, data),

  // Sprints
  createSprint: (projectId: string, data: CreateSprintRequest) =>
    api.post<Sprint>(`/projects/${projectId}/sprints`, data),

  listSprints: (projectId: string) =>
    api.get<Sprint[]>(`/projects/${projectId}/sprints`),

  getActiveSprint: (projectId: string) =>
    api.get<Sprint>(`/projects/${projectId}/sprints/active`),

  startSprint: (projectId: string, sprintId: string) =>
    api.post<Sprint>(`/projects/${projectId}/sprints/${sprintId}/start`),

  completeSprint: (projectId: string, sprintId: string) =>
    api.post<SprintMetrics>(`/projects/${projectId}/sprints/${sprintId}/complete`),

  getSprintRecords: (projectId: string, sprintId: string) =>
    api.get<JiraProjectRecord[]>(`/projects/${projectId}/sprints/${sprintId}/records`),

  // Backlog
  getBacklog: (projectId: string) =>
    api.get<JiraProjectRecord[]>(`/projects/${projectId}/backlog`),

  reorderBacklog: (projectId: string, data: ReorderBacklogRequest) =>
    api.patch(`/projects/${projectId}/backlog/reorder`, data),

  bulkAssignToSprint: (projectId: string, data: BulkAssignToSprintRequest) =>
    api.post(`/projects/${projectId}/backlog/assign-sprint`, data),

  // Comments
  addComment: (projectId: string, recordId: string, data: AddCommentRequest) =>
    api.post<Comment>(`/projects/${projectId}/records/${recordId}/comments`, data),

  listComments: (projectId: string, recordId: string) =>
    api.get<Comment[]>(`/projects/${projectId}/records/${recordId}/comments`),

  updateComment: (projectId: string, commentId: string, data: UpdateCommentRequest) =>
    api.patch<Comment>(`/projects/${projectId}/comments/${commentId}`, data),

  deleteComment: (projectId: string, commentId: string) =>
    api.delete(`/projects/${projectId}/comments/${commentId}`),

  // Attachments
  uploadAttachment: (projectId: string, recordId: string, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post<Attachment>(`/projects/${projectId}/records/${recordId}/attachments`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  listAttachments: (projectId: string, recordId: string) =>
    api.get<Attachment[]>(`/projects/${projectId}/records/${recordId}/attachments`),

  deleteAttachment: (projectId: string, attachmentId: string) =>
    api.delete(`/projects/${projectId}/attachments/${attachmentId}`),

  // Labels
  createLabel: (projectId: string, data: CreateLabelRequest) =>
    api.post<Label>(`/projects/${projectId}/labels`, data),

  listLabels: (projectId: string) =>
    api.get<Label[]>(`/projects/${projectId}/labels`),

  addLabelToRecord: (projectId: string, recordId: string, labelId: string) =>
    api.post(`/projects/${projectId}/records/${recordId}/labels/${labelId}`),

  getRecordLabels: (projectId: string, recordId: string) =>
    api.get<Label[]>(`/projects/${projectId}/records/${recordId}/labels`),

  removeLabelFromRecord: (projectId: string, recordId: string, labelId: string) =>
    api.delete(`/projects/${projectId}/records/${recordId}/labels/${labelId}`),

  deleteLabel: (projectId: string, labelId: string) =>
    api.delete(`/projects/${projectId}/labels/${labelId}`),

  // Bulk Operations
  bulkChangeStatus: (projectId: string, data: BulkChangeStatusRequest) =>
    api.post(`/projects/${projectId}/bulk/change-status`, data),

  bulkAssignTo: (projectId: string, data: BulkAssignToRequest) =>
    api.post(`/projects/${projectId}/bulk/assign`, data),

  bulkAddLabel: (projectId: string, data: BulkAddLabelRequest) =>
    api.post(`/projects/${projectId}/bulk/add-label`, data),

  bulkDelete: (projectId: string, data: BulkDeleteRequest) =>
    api.post(`/projects/${projectId}/bulk/delete`, data),

  // Search & Filters
  searchRecords: (projectId: string, query?: string, filters?: SearchFilters) =>
    api.get<JiraProjectRecord[]>(`/projects/${projectId}/search`, {
      params: { q: query, ...filters },
    }),

  saveFilter: (projectId: string, data: SaveFilterRequest) =>
    api.post(`/projects/${projectId}/filters`, data),

  listSavedFilters: (projectId: string) =>
    api.get<SavedFilter[]>(`/projects/${projectId}/filters`),
}
