import { useState, useEffect } from 'react'
import { CommentSection } from './CommentSection'
import { AttachmentSection } from './AttachmentSection'
import { LabelManager } from './LabelManager'
import { formatDate, formatDateTime, isOverdue, getIssueTypeIcon } from '../../utils/jira.utils'
import { projectService } from '../../services/project.service'
import { jiraService } from '../../services/jira.service'
import type { Label, JiraProjectRecord, WorkflowStatus, CustomFieldValue } from '../../types/jira'
import type { ProjectRecord } from '../../types/project'
import type { Project } from '../../types/project'

interface RecordDetailModalProps {
  record: ProjectRecord | JiraProjectRecord
  project?: Project
  isOpen?: boolean
  onClose: () => void
  onUpdate?: () => void
  currentUserId?: string
  statuses?: WorkflowStatus[]
  onStatusChange?: (statusId: string) => Promise<void>
}

interface UserInfo { id: string; name: string; email: string }

export function RecordDetailModal({
  record,
  isOpen = true,
  onClose,
  onUpdate,
  currentUserId,
  statuses = [],
  onStatusChange,
}: RecordDetailModalProps) {
  const [selectedLabels, setSelectedLabels] = useState<Label[]>((record as any).labels || [])
  const [isEditingStatus, setIsEditingStatus] = useState(false)
  const [statusLoading, setStatusLoading] = useState(false)
  const [users, setUsers] = useState<UserInfo[]>([])
  const [editingAssignee, setEditingAssignee] = useState(false)
  const [editingDueDate, setEditingDueDate] = useState(false)
  const [dueDate, setDueDate] = useState((record as any).due_date?.split('T')[0] || '')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    projectService.listUsers()
      .then(res => setUsers(Array.isArray(res.data) ? res.data : []))
      .catch(() => setUsers([]))

    // Fetch record labels from API
    if (record?.id && record?.project_id) {
      jiraService.getRecordLabels(record.project_id, record.id)
        .then(res => setSelectedLabels(Array.isArray(res.data) ? res.data : []))
        .catch(() => setSelectedLabels([]))
    }
  }, [record?.id, record?.project_id])

  if (!isOpen || !record) return null

  const jiraRecord = record as JiraProjectRecord
  const overdue = jiraRecord.due_date && !jiraRecord.is_completed && isOverdue(jiraRecord.due_date)

  const handleStatusChange = async (statusId: string) => {
    if (!onStatusChange) return
    try {
      setStatusLoading(true)
      await onStatusChange(statusId)
      setIsEditingStatus(false)
    } catch (error) {
      console.error('Failed to change status:', error)
    } finally {
      setStatusLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-surface-container-lowest rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto shadow-lg"
        onClick={e => e.stopPropagation()}
      >
        {/* Header */}
        <div className="sticky top-0 bg-surface-container-low border-b border-outline-variant p-4 flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              {jiraRecord.issue_type && (
                <span className="text-lg">{getIssueTypeIcon(jiraRecord.issue_type.name)}</span>
              )}
              <h2 className="text-lg font-semibold text-on-surface">{record.title}</h2>
            </div>
            {record.description && (
              <p className="text-sm text-on-surface-variant line-clamp-2">{record.description}</p>
            )}
          </div>
          <button
            onClick={onClose}
            className="text-on-surface-variant hover:text-on-surface text-xl"
            aria-label="Close modal"
          >
            ✕
          </button>
        </div>

        <div className="p-4 space-y-6">
          {/* Main Fields Grid */}
          <div className="grid grid-cols-2 gap-4">
            {/* Status with Transition Options */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant block mb-1">Status</label>
              {isEditingStatus && statuses.length > 0 ? (
                <div className="space-y-2">
                  <select
                    value={jiraRecord.status || ''}
                    onChange={e => handleStatusChange(e.target.value)}
                    disabled={statusLoading}
                    className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
                  >
                    <option value="">Select status...</option>
                    {statuses.map(status => (
                      <option key={status.id} value={status.id}>{status.status_name}</option>
                    ))}
                  </select>
                  <button onClick={() => setIsEditingStatus(false)} className="text-[10px] text-on-surface-variant">Cancel</button>
                </div>
              ) : (
                <div
                  onClick={() => statuses.length > 0 && setIsEditingStatus(true)}
                  className={`text-sm text-on-surface p-2 rounded ${statuses.length > 0 ? 'cursor-pointer hover:bg-surface-container-high' : ''}`}
                >
                  {(() => {
                    const s = statuses.find(st => st.id === jiraRecord.status)
                    return s ? s.status_name : (statuses.find(st => st.status_name === jiraRecord.status)?.status_name || jiraRecord.status || 'No status')
                  })()}
                </div>
              )}
            </div>

            {/* Assignee */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant">Assignee</label>
              {editingAssignee ? (
                <div className="mt-1">
                  <select
                    defaultValue={jiraRecord.assigned_to || ''}
                    onChange={async (e) => {
                      const val = e.target.value
                      setSaving(true)
                      try {
                        const assignTo = val || null
                        await projectService.updateRecord(record.project_id, record.id, { assigned_to: assignTo } as any)
                        setEditingAssignee(false)
                      } catch (err) { console.error('Failed to update assignee', err) }
                      finally { setSaving(false) }
                    }}
                    className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
                  >
                    <option value="">Unassigned</option>
                    {users.map(u => (
                      <option key={u.id} value={u.id}>{u.name || u.email}</option>
                    ))}
                  </select>
                  <button onClick={() => setEditingAssignee(false)} className="text-[10px] text-on-surface-variant mt-1">Cancel</button>
                </div>
              ) : (
                <p className="text-sm text-on-surface mt-1 cursor-pointer hover:bg-surface-container-high p-1 rounded"
                  onClick={() => setEditingAssignee(true)}>
                  {jiraRecord.assigned_to ? (users.find(u => u.id === jiraRecord.assigned_to)?.name || users.find(u => u.id === jiraRecord.assigned_to)?.email || jiraRecord.assigned_to) : '— Klik untuk assign'}
                </p>
              )}
            </div>

            {/* Due Date */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant">Due Date</label>
              {editingDueDate ? (
                <div className="mt-1">
                  <input type="date" value={dueDate} onChange={e => setDueDate(e.target.value)}
                    className="w-full px-2 py-1 rounded border border-outline-variant text-sm" />
                  <div className="flex gap-2 mt-1">
                    <button onClick={async () => {
                      setSaving(true)
                      try {
                        const dueDateVal = dueDate ? new Date(dueDate).toISOString() : null
                        await projectService.updateRecord(record.project_id, record.id, { due_date: dueDateVal } as any)
                        setEditingDueDate(false)
                      } catch (err) { console.error('Failed to update due date', err) }
                      finally { setSaving(false) }
                    }} disabled={saving} className="text-[10px] text-primary font-semibold">Save</button>
                    <button onClick={() => setEditingDueDate(false)} className="text-[10px] text-on-surface-variant">Cancel</button>
                  </div>
                </div>
              ) : (
                <p className={`text-sm mt-1 cursor-pointer hover:bg-surface-container-high p-1 rounded ${overdue ? 'text-error' : 'text-on-surface'}`}
                  onClick={() => setEditingDueDate(true)}>
                  {jiraRecord.due_date ? formatDate(jiraRecord.due_date) : '— Klik untuk set'}
                  {overdue && <span className="ml-1 text-[10px]">(Overdue)</span>}
                </p>
              )}
            </div>

            {/* Created */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant">Created</label>
              <p className="text-sm text-on-surface mt-1">{formatDateTime(record.created_at)}</p>
            </div>
          </div>

          {/* Description */}
          {record.description && (
            <div className="border-t border-outline-variant pt-4">
              <label className="text-[10px] font-medium text-on-surface-variant block mb-2">
                Description
              </label>
              <p className="text-sm text-on-surface whitespace-pre-wrap">{record.description}</p>
            </div>
          )}

          {/* Custom Fields */}
          {jiraRecord.custom_fields && jiraRecord.custom_fields.length > 0 && (
            <div className="border-t border-outline-variant pt-4">
              <h3 className="text-sm font-semibold text-on-surface mb-3">Custom Fields</h3>
              <div className="grid grid-cols-2 gap-3">
                {jiraRecord.custom_fields.map(field => (
                  <div key={field.id}>
                    <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                      {field.field_id}
                    </label>
                    <p className="text-sm text-on-surface">{field.value || '-'}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Labels */}
          <div className="border-t border-outline-variant pt-4">
            <LabelManager
              recordId={record.id}
              projectId={record.project_id}
              selectedLabels={selectedLabels}
              onLabelsChange={setSelectedLabels}
            />
          </div>

          {/* Comments */}
          <div className="border-t border-outline-variant pt-4">
            <CommentSection
              recordId={record.id}
              projectId={record.project_id}
              currentUserId={currentUserId || ''}
              projectMembers={users.map(u => ({ id: u.id, name: u.name || u.email }))}
            />
          </div>

          {/* Attachments */}
          <div className="border-t border-outline-variant pt-4">
            <AttachmentSection recordId={record.id} projectId={record.project_id} currentUserId={currentUserId || ''} />
          </div>
        </div>
      </div>
    </div>
  )
}
