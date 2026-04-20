import { useState } from 'react'
import { useBacklog } from '../../hooks/useBacklog'
import { RecordCard } from './RecordCard'
import { RecordDetailModal } from './RecordDetailModal'
import { BulkOperationsBar } from './BulkOperationsBar'
import type { Sprint, Label, WorkflowStatus } from '../../types/jira'

interface BacklogViewProps {
  projectId: string
  sprints: Sprint[]
  labels?: Label[]
  statuses?: WorkflowStatus[]
}

export function BacklogView({ projectId, sprints, labels = [], statuses = [] }: BacklogViewProps) {
  const {
    records,
    loading,
    error,
    selectedRecords,
    toggleRecordSelection,
    assignToSprint,
    clearSelection,
  } = useBacklog(projectId)

  const [selectedRecord, setSelectedRecord] = useState<any>(null)
  const [showDetailModal, setShowDetailModal] = useState(false)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-on-surface-variant">Memuat backlog...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-error">{error}</div>
      </div>
    )
  }

  const handleAssignToSprint = async (sprintId: string) => {
    await assignToSprint(sprintId, Array.from(selectedRecords))
    clearSelection()
  }

  const handleSelectAll = () => {
    if (selectedRecords.size === records.length) {
      clearSelection()
    } else {
      records.forEach(r => {
        if (!selectedRecords.has(r.id)) {
          toggleRecordSelection(r.id)
        }
      })
    }
  }

  return (
    <>
      <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
        {/* Backlog Records */}
        <div className="lg:col-span-3">
          <div className="bg-surface-container-lowest rounded-lg border border-outline-variant/10 p-4">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-on-surface">Backlog ({records.length})</h2>
              {records.length > 0 && (
                <button
                  onClick={handleSelectAll}
                  className="text-sm text-primary hover:underline"
                >
                  {selectedRecords.size === records.length ? 'Deselect All' : 'Select All'}
                </button>
              )}
            </div>

            {records.length === 0 ? (
              <div className="text-center py-8">
                <p className="text-on-surface-variant">No backlog items</p>
              </div>
            ) : (
              <div className="space-y-2">
                {records.map((record, index) => (
                  <div
                    key={record.id}
                    className="flex items-center gap-3 p-2 rounded hover:bg-surface-container-high transition-colors"
                  >
                    <input
                      type="checkbox"
                      checked={selectedRecords.has(record.id)}
                      onChange={() => toggleRecordSelection(record.id)}
                      className="w-4 h-4 cursor-pointer"
                      aria-label={`Select ${record.title}`}
                    />
                    <span className="text-[10px] text-on-surface-variant font-medium w-6 text-right">
                      {index + 1}
                    </span>
                    <div
                      className="flex-1 cursor-pointer"
                      onClick={() => {
                        setSelectedRecord(record)
                        setShowDetailModal(true)
                      }}
                    >
                      <div className="flex items-center gap-2">
                        <RecordCard record={record} onClick={() => {}} />
                        {record.status && (
                          <span className="px-2 py-0.5 rounded text-[10px] font-semibold bg-primary/10 text-primary whitespace-nowrap">
                            {statuses.find(s => s.id === record.status)?.status_name || 
                             statuses.find(s => s.status_name === record.status)?.status_name || 
                             ''}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Sprints Sidebar */}
        <div className="lg:col-span-1">
          <div className="bg-surface-container-lowest rounded-lg border border-outline-variant/10 p-4 sticky top-4">
            <h3 className="font-semibold text-on-surface mb-3">Sprints</h3>
            <div className="space-y-2">
              {sprints.length === 0 ? (
                <p className="text-sm text-on-surface-variant">No sprints available</p>
              ) : (
                sprints
                  .filter(s => s.status !== 'Completed')
                  .map(sprint => (
                    <button
                      key={sprint.id}
                      onClick={() => handleAssignToSprint(sprint.id)}
                      disabled={selectedRecords.size === 0}
                      className="w-full text-left p-2 rounded bg-surface-container-high hover:bg-surface-container-highest disabled:opacity-50 disabled:cursor-not-allowed text-sm text-on-surface transition-colors"
                      title={`Assign ${selectedRecords.size} record(s) to ${sprint.name}`}
                    >
                      <div className="flex items-center justify-between">
                        <span>{sprint.name}</span>
                        {sprint.status === 'Active' && (
                          <span className="text-primary text-lg">●</span>
                        )}
                      </div>
                      {sprint.goal && (
                        <p className="text-[10px] text-on-surface-variant mt-1 line-clamp-1">
                          {sprint.goal}
                        </p>
                      )}
                    </button>
                  ))
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Bulk Operations Bar */}
      {selectedRecords.size > 0 && (
        <BulkOperationsBar
          projectId={projectId}
          selectedRecordIds={Array.from(selectedRecords)}
          statuses={statuses}
          labels={labels}
          sprints={sprints}
          onOperationComplete={() => {
            clearSelection()
            window.location.reload()
          }}
        />
      )}

      {/* Record Detail Modal */}
      {selectedRecord && (
        <RecordDetailModal
          record={selectedRecord}
          isOpen={showDetailModal}
          onClose={() => {
            setShowDetailModal(false)
            setSelectedRecord(null)
          }}
          currentUserId=""
          statuses={statuses}
        />
      )}
    </>
  )
}
