import { useState } from 'react'
import { useDroppable } from '@dnd-kit/core'
import type { ProjectColumn, ProjectRecord } from '../../types/project'
import { ProjectRecordCard } from './ProjectRecordCard'

interface ProjectBoardColumnProps {
  column: ProjectColumn
  records: ProjectRecord[]
  onAddRecord: (columnId: string, title: string) => void
  onClickRecord: (record: ProjectRecord) => void
  onEditColumn: (columnId: string, name: string) => void
  onDeleteColumn: (columnId: string) => void
}

export function ProjectBoardColumn({ column, records, onAddRecord, onClickRecord, onEditColumn, onDeleteColumn }: ProjectBoardColumnProps) {
  const { setNodeRef, isOver } = useDroppable({ id: column.id })
  const [showAdd, setShowAdd] = useState(false)
  const [newTitle, setNewTitle] = useState('')
  const [showMenu, setShowMenu] = useState(false)
  const [editing, setEditing] = useState(false)
  const [editName, setEditName] = useState(column.name)

  const handleAdd = () => {
    if (!newTitle.trim()) return
    onAddRecord(column.id, newTitle.trim())
    setNewTitle('')
    setShowAdd(false)
  }

  const handleRename = () => {
    if (editName.trim() && editName.trim() !== column.name) {
      onEditColumn(column.id, editName.trim())
    }
    setEditing(false)
  }

  return (
    <div className="w-[300px] min-w-[280px] flex flex-col gap-3 shrink-0">
      {/* Column Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 flex-1 min-w-0">
          <div className="w-1 h-5 bg-primary rounded-full" />
          {editing ? (
            <input
              value={editName}
              onChange={e => setEditName(e.target.value)}
              onBlur={handleRename}
              onKeyDown={e => { if (e.key === 'Enter') handleRename(); if (e.key === 'Escape') setEditing(false) }}
              className="text-sm font-bold text-on-surface bg-transparent border-b border-primary outline-none flex-1"
              autoFocus
            />
          ) : (
            <h2 className="text-sm font-bold tracking-tight text-on-surface truncate">{column.name}</h2>
          )}
          <span className="text-[10px] font-medium text-on-surface-variant bg-surface-container-high px-1.5 py-0.5 rounded-full">
            {records.length}
          </span>
        </div>
        <div className="flex items-center gap-0.5">
          <button
            onClick={() => setShowAdd(true)}
            className="w-6 h-6 flex items-center justify-center rounded-md hover:bg-surface-container-high transition-colors"
          >
            <span className="material-symbols-outlined text-on-surface-variant text-[16px]">add</span>
          </button>
          <div className="relative">
            <button
              onClick={() => setShowMenu(!showMenu)}
              className="w-6 h-6 flex items-center justify-center rounded-md hover:bg-surface-container-high transition-colors"
            >
              <span className="material-symbols-outlined text-on-surface-variant text-[16px]">more_horiz</span>
            </button>
            {showMenu && (
              <div className="absolute right-0 top-full mt-1 w-40 bg-white rounded-xl shadow-xl border border-outline-variant/15 overflow-hidden z-50">
                <button
                  onClick={() => { setEditing(true); setShowMenu(false) }}
                  className="flex items-center gap-2 w-full px-3 py-2 text-sm text-on-surface hover:bg-surface-container-low transition-colors"
                >
                  <span className="material-symbols-outlined text-[16px]">edit</span>
                  Ubah nama
                </button>
                <button
                  onClick={() => { onDeleteColumn(column.id); setShowMenu(false) }}
                  className="flex items-center gap-2 w-full px-3 py-2 text-sm text-error hover:bg-error-container/20 transition-colors"
                >
                  <span className="material-symbols-outlined text-[16px]">delete</span>
                  Hapus kolom
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Droppable Area */}
      <div
        ref={setNodeRef}
        className={`flex-1 bg-surface-container-low rounded-xl p-3 flex flex-col gap-3 overflow-y-auto transition-all min-h-[200px] ${
          isOver ? 'ring-2 ring-primary/20 bg-surface-container-low/80' : ''
        }`}
      >
        {showAdd && (
          <div className="bg-surface-container-lowest p-3 rounded-xl border border-primary/20">
            <input
              value={newTitle}
              onChange={e => setNewTitle(e.target.value)}
              onKeyDown={e => { if (e.key === 'Enter') handleAdd(); if (e.key === 'Escape') setShowAdd(false) }}
              placeholder="Judul record..."
              className="w-full text-sm bg-transparent outline-none text-on-surface"
              autoFocus
            />
            <div className="flex gap-2 mt-2">
              <button onClick={handleAdd} className="px-3 py-1 text-xs font-bold text-on-primary bg-primary rounded-lg">
                Tambah
              </button>
              <button onClick={() => setShowAdd(false)} className="px-3 py-1 text-xs font-medium text-on-surface-variant">
                Batal
              </button>
            </div>
          </div>
        )}
        {records.length === 0 && !showAdd ? (
          <div className="flex-1 flex items-center justify-center">
            <p className="text-xs text-on-surface-variant/50 font-medium">Tidak ada record</p>
          </div>
        ) : (
          records.map(record => (
            <ProjectRecordCard key={record.id} record={record} onClick={() => onClickRecord(record)} />
          ))
        )}
      </div>
    </div>
  )
}
