import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { DndContext, DragOverlay, PointerSensor, useSensor, useSensors } from '@dnd-kit/core'
import { useProjectBoard } from '../hooks/useProjectBoard'
import { ProjectBoardColumn } from '../components/project/ProjectBoardColumn'
import { ProjectRecordCard } from '../components/project/ProjectRecordCard'
import { ProjectFilterBar } from '../components/project/ProjectFilterBar'
import { RecordDetailModal } from '../components/project/RecordDetailModal'
import { MemberManagement } from '../components/project/MemberManagement'
import type { ProjectRecord } from '../types/project'

export function ProjectBoardPage() {
  const { id } = useParams()
  const board = useProjectBoard(id)
  const [selectedRecord, setSelectedRecord] = useState<ProjectRecord | null>(null)
  const [showAddColumn, setShowAddColumn] = useState(false)
  const [showMembers, setShowMembers] = useState(false)
  const [newColumnName, setNewColumnName] = useState('')

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } })
  )

  const handleAddColumn = () => {
    if (!newColumnName.trim()) return
    board.addColumn(newColumnName.trim())
    setNewColumnName('')
    setShowAddColumn(false)
  }

  if (board.loading) {
    return (
      <div className="flex flex-col h-[calc(100vh-64px)]">
        <div className="px-8 pt-8 pb-4">
          <div className="h-7 w-48 bg-surface-container-high rounded-lg animate-pulse mb-2" />
          <div className="h-4 w-64 bg-surface-container-high rounded-lg animate-pulse" />
        </div>
        <div className="flex gap-6 px-8 pb-8 flex-1">
          {[1, 2, 3].map(i => (
            <div key={i} className="w-[300px] shrink-0">
              <div className="h-5 w-24 bg-surface-container-high rounded animate-pulse mb-3" />
              <div className="bg-surface-container-low rounded-xl p-3 flex flex-col gap-3 min-h-[200px]">
                {[1, 2].map(j => (
                  <div key={j} className="bg-surface-container-lowest p-3 rounded-xl animate-pulse">
                    <div className="h-4 bg-surface-container-high rounded w-full mb-2" />
                    <div className="h-3 bg-surface-container-high rounded w-2/3" />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (!board.project) {
    return (
      <div className="flex items-center justify-center h-[calc(100vh-64px)]">
        <p className="text-on-surface-variant">Project tidak ditemukan</p>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-[calc(100vh-64px)]">
      <div className="px-8 pt-8 pb-4">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">{board.project.name}</h1>
          <button onClick={() => setShowMembers(true)}
            className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-surface-variant bg-surface-container-high rounded-xl hover:bg-surface-container-highest transition-colors">
            <span className="material-symbols-outlined text-[18px]">group</span>
            Members
          </button>
        </div>
        <div className="mt-3">
          <ProjectFilterBar filters={board.filters} onChange={board.setFilters} />
        </div>
      </div>

      {board.error && (
        <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
          <span className="material-symbols-outlined text-lg">error</span>
          <span>{board.error}</span>
        </div>
      )}
      {board.dragError && (
        <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
          <span className="material-symbols-outlined text-lg">warning</span>
          <span>{board.dragError}</span>
        </div>
      )}

      <DndContext sensors={sensors} onDragStart={board.handleDragStart} onDragEnd={board.handleDragEnd}>
        <div className="flex-1 overflow-x-auto overflow-y-hidden px-8 pb-8">
          <div className="flex gap-6 h-full min-w-max">
            {board.project.columns.map(col => (
              <ProjectBoardColumn
                key={col.id}
                column={col}
                records={board.getFilteredRecords(col.records)}
                onAddRecord={board.addRecord}
                onClickRecord={setSelectedRecord}
                onEditColumn={board.editColumn}
                onDeleteColumn={board.deleteColumn}
              />
            ))}

            {/* Add Column */}
            <div className="w-[300px] min-w-[280px] shrink-0">
              {showAddColumn ? (
                <div className="bg-surface-container-low rounded-xl p-3">
                  <input
                    value={newColumnName}
                    onChange={e => setNewColumnName(e.target.value)}
                    onKeyDown={e => { if (e.key === 'Enter') handleAddColumn(); if (e.key === 'Escape') setShowAddColumn(false) }}
                    placeholder="Nama kolom..."
                    className="w-full text-sm bg-transparent outline-none text-on-surface mb-2"
                    autoFocus
                  />
                  <div className="flex gap-2">
                    <button onClick={handleAddColumn} className="px-3 py-1 text-xs font-bold text-on-primary bg-primary rounded-lg">
                      Tambah
                    </button>
                    <button onClick={() => setShowAddColumn(false)} className="px-3 py-1 text-xs font-medium text-on-surface-variant">
                      Batal
                    </button>
                  </div>
                </div>
              ) : (
                <button
                  onClick={() => setShowAddColumn(true)}
                  className="flex items-center gap-2 px-4 py-2.5 text-sm font-medium text-on-surface-variant hover:bg-surface-container-low rounded-xl transition-colors w-full"
                >
                  <span className="material-symbols-outlined text-[18px]">add</span>
                  Tambah kolom
                </button>
              )}
            </div>
          </div>
        </div>
        <DragOverlay>
          {board.activeRecord ? <ProjectRecordCard record={board.activeRecord} onClick={() => {}} /> : null}
        </DragOverlay>
      </DndContext>

      {selectedRecord && board.project && (
        <RecordDetailModal
          record={selectedRecord}
          project={board.project}
          onClose={() => setSelectedRecord(null)}
          onUpdate={() => { setSelectedRecord(null); board.refresh() }}
        />
      )}

      {showMembers && id && (
        <MemberManagement projectId={id} onClose={() => setShowMembers(false)} />
      )}
    </div>
  )
}
