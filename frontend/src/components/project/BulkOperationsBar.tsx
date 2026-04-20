import { useState } from 'react'
import { jiraService } from '../../services/jira.service'
import type { WorkflowStatus, Label, Sprint } from '../../types/jira'

interface BulkOperationsBarProps {
  projectId: string
  selectedRecordIds: string[]
  statuses: WorkflowStatus[]
  labels: Label[]
  sprints?: Sprint[]
  onOperationComplete: () => void
}

export function BulkOperationsBar({
  projectId,
  selectedRecordIds,
  statuses,
  labels,
  sprints = [],
  onOperationComplete,
}: BulkOperationsBarProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (selectedRecordIds.length === 0) {
    return null
  }

  const handleBulkChangeStatus = async (statusId: string) => {
    if (!statusId) return
    try {
      setLoading(true)
      setError(null)
      await jiraService.bulkChangeStatus(projectId, {
        record_ids: selectedRecordIds,
        status_id: statusId,
      })
      onOperationComplete()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal mengubah status'
      setError(errorMessage)
      setTimeout(() => setError(null), 4000)
    } finally {
      setLoading(false)
    }
  }

  const handleBulkAssignTo = async (assigneeId: string) => {
    if (!assigneeId) return
    try {
      setLoading(true)
      setError(null)
      await jiraService.bulkAssignTo(projectId, {
        record_ids: selectedRecordIds,
        assignee_id: assigneeId,
      })
      onOperationComplete()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menugaskan'
      setError(errorMessage)
      setTimeout(() => setError(null), 4000)
    } finally {
      setLoading(false)
    }
  }

  const handleBulkAddLabel = async (labelId: string) => {
    if (!labelId) return
    try {
      setLoading(true)
      setError(null)
      await jiraService.bulkAddLabel(projectId, {
        record_ids: selectedRecordIds,
        label_id: labelId,
      })
      onOperationComplete()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menambah label'
      setError(errorMessage)
      setTimeout(() => setError(null), 4000)
    } finally {
      setLoading(false)
    }
  }

  const handleBulkAssignToSprint = async (sprintId: string) => {
    if (!sprintId) return
    try {
      setLoading(true)
      setError(null)
      await jiraService.bulkAssignToSprint(projectId, {
        sprint_id: sprintId,
        record_ids: selectedRecordIds,
      })
      onOperationComplete()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menugaskan ke sprint'
      setError(errorMessage)
      setTimeout(() => setError(null), 4000)
    } finally {
      setLoading(false)
    }
  }

  const handleBulkDelete = async () => {
    if (!confirm(`Hapus ${selectedRecordIds.length} record? Tindakan ini tidak dapat dibatalkan.`)) {
      return
    }
    try {
      setLoading(true)
      setError(null)
      await jiraService.bulkDelete(projectId, {
        record_ids: selectedRecordIds,
      })
      onOperationComplete()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menghapus records'
      setError(errorMessage)
      setTimeout(() => setError(null), 4000)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed bottom-4 left-4 right-4 bg-surface-container-low border border-outline-variant rounded-lg p-4 shadow-lg z-40">
      <div className="space-y-3">
        {error && (
          <div className="text-error text-sm p-2 bg-error/10 rounded">{error}</div>
        )}

        <div className="flex items-center justify-between gap-4">
          <span className="text-sm font-medium text-on-surface">
            {selectedRecordIds.length} record dipilih
          </span>

          <div className="flex items-center gap-2 flex-wrap">
            {/* Change Status */}
            <select
              onChange={e => handleBulkChangeStatus(e.target.value)}
              disabled={loading || statuses.length === 0}
              className="px-3 py-1 rounded border border-outline-variant text-sm bg-surface-container-highest text-on-surface disabled:opacity-50 disabled:cursor-not-allowed"
              aria-label="Change status for selected records"
            >
              <option value="">Ubah Status</option>
              {statuses.map(status => (
                <option key={status.id} value={status.id}>
                  {status.status_name}
                </option>
              ))}
            </select>

            {/* Add Label */}
            <select
              onChange={e => handleBulkAddLabel(e.target.value)}
              disabled={loading || labels.length === 0}
              className="px-3 py-1 rounded border border-outline-variant text-sm bg-surface-container-highest text-on-surface disabled:opacity-50 disabled:cursor-not-allowed"
              aria-label="Add label to selected records"
            >
              <option value="">Tambah Label</option>
              {labels.map(label => (
                <option key={label.id} value={label.id}>
                  {label.name}
                </option>
              ))}
            </select>

            {/* Assign to Sprint */}
            {sprints.length > 0 && (
              <select
                onChange={e => handleBulkAssignToSprint(e.target.value)}
                disabled={loading}
                className="px-3 py-1 rounded border border-outline-variant text-sm bg-surface-container-highest text-on-surface disabled:opacity-50 disabled:cursor-not-allowed"
                aria-label="Assign to sprint"
              >
                <option value="">Pindah ke Sprint</option>
                {sprints
                  .filter(s => s.status !== 'Completed')
                  .map(sprint => (
                    <option key={sprint.id} value={sprint.id}>
                      {sprint.name}
                    </option>
                  ))}
              </select>
            )}

            {/* Delete */}
            <button
              onClick={handleBulkDelete}
              disabled={loading}
              className="px-3 py-1 bg-error text-on-error rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-90"
              aria-label="Delete selected records"
            >
              Hapus
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
