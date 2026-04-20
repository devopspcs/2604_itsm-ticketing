import { useState } from 'react'
import { DndContext, DragEndEvent, useDroppable } from '@dnd-kit/core'
import { useJiraBoard } from '../../hooks/useJiraBoard'
import { jiraService } from '../../services/jira.service'
import { RecordCard } from './RecordCard'
import { RecordDetailModal } from './RecordDetailModal'
import { getDaysRemaining, calculateSprintProgress } from '../../utils/jira.utils'
import type { WorkflowStatus } from '../../types/jira'

interface SprintBoardProps {
  projectId: string
  sprintId: string
  statuses: WorkflowStatus[]
}

interface StatusColumnProps {
  status: WorkflowStatus
  records: any[]
  onRecordClick: (record: any) => void
}

function StatusColumn({ status, records, onRecordClick }: StatusColumnProps) {
  const { setNodeRef } = useDroppable({
    id: status.id,
  })

  return (
    <div
      ref={setNodeRef}
      className="bg-surface-container-lowest rounded-lg border border-outline-variant/10 p-3 flex flex-col"
    >
      <h3 className="font-semibold text-on-surface mb-3 text-sm">{status.status_name}</h3>
      <div className="space-y-2 flex-1 min-h-96">
        {records.map((record) => (
          <div key={record.id} onClick={() => onRecordClick(record)}>
            <RecordCard record={record} onClick={() => onRecordClick(record)} />
          </div>
        ))}
      </div>
    </div>
  )
}

export function SprintBoard({ projectId, sprintId, statuses }: SprintBoardProps) {
  const { sprint, records, loading, error, handleDragEnd } = useJiraBoard(projectId, sprintId)
  const [selectedRecord, setSelectedRecord] = useState<any>(null)
  const [showDetailModal, setShowDetailModal] = useState(false)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-on-surface-variant">Memuat sprint board...</div>
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

  const getRecordsByStatus = (statusId: string) => {
    const status = statuses.find(s => s.id === statusId)
    return records.filter(r => 
      r.status === statusId || // match by UUID
      (status && r.status === status.status_name) // match by name
    )
  }

  const completedCount = records.filter(r => r.is_completed).length
  const totalCount = records.length
  const completionPercentage = calculateSprintProgress(completedCount, totalCount)
  const daysRemaining = sprint?.end_date ? getDaysRemaining(sprint.end_date) : 0

  return (
    <>
      <div className="space-y-4">
        {/* Sprint Header with Metrics */}
        {sprint && (
          <div className="bg-surface-container-low p-4 rounded-lg border border-outline-variant/10">
            <div className="flex items-center justify-between mb-3">
              <div>
                <h2 className="text-lg font-semibold text-on-surface">{sprint.name}</h2>
                {sprint.goal && <p className="text-sm text-on-surface-variant mt-1">{sprint.goal}</p>}
              </div>
              <div className="text-right">
                <p className="text-sm font-medium text-on-surface">{sprint.status}</p>
                {sprint.start_date && sprint.end_date && (
                  <p className="text-[10px] text-on-surface-variant">
                    {sprint.start_date} to {sprint.end_date}
                  </p>
                )}
              </div>
            </div>

            {/* Sprint Metrics */}
            <div className="grid grid-cols-4 gap-3">
              <div className="bg-surface-container-highest p-2 rounded">
                <p className="text-[10px] text-on-surface-variant">Total</p>
                <p className="text-lg font-semibold text-on-surface">{totalCount}</p>
              </div>
              <div className="bg-surface-container-highest p-2 rounded">
                <p className="text-[10px] text-on-surface-variant">Completed</p>
                <p className="text-lg font-semibold text-on-surface">{completedCount}</p>
              </div>
              <div className="bg-surface-container-highest p-2 rounded">
                <p className="text-[10px] text-on-surface-variant">Progress</p>
                <p className="text-lg font-semibold text-on-surface">{completionPercentage}%</p>
              </div>
              <div className="bg-surface-container-highest p-2 rounded">
                <p className="text-[10px] text-on-surface-variant">Days Left</p>
                <p className="text-lg font-semibold text-on-surface">{daysRemaining}</p>
              </div>
            </div>

            {/* Progress Bar */}
            <div className="mt-3 bg-surface-container-highest rounded-full h-2 overflow-hidden">
              <div
                className="bg-primary h-full transition-all"
                style={{ width: `${completionPercentage}%` }}
              />
            </div>
          </div>
        )}

        {/* Sprint Board Columns */}
        <DndContext onDragEnd={handleDragEnd}>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {statuses.map(status => (
              <StatusColumn
                key={status.id}
                status={status}
                records={getRecordsByStatus(status.id)}
                onRecordClick={(record: any) => {
                  setSelectedRecord(record)
                  setShowDetailModal(true)
                }}
              />
            ))}
          </div>
        </DndContext>
      </div>

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
          onStatusChange={async (statusId: string) => {
            try {
              await jiraService.transitionRecord(projectId, selectedRecord.id, { to_status_id: statusId })
              setSelectedRecord({ ...selectedRecord, status: statusId })
            } catch (err) {
              console.error('Failed to change status:', err)
            }
          }}
          onUpdate={() => {
            setShowDetailModal(false)
            setSelectedRecord(null)
          }}
        />
      )}
    </>
  )
}
