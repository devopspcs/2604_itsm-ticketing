import { useEffect, useState } from 'react'
import { projectService } from '../../services/project.service'
import api from '../../services/api'
import type { ProjectRecord, ProjectDetail, ProjectActivityLog } from '../../types/project'
import type { User } from '../../types'

interface RecordDetailModalProps {
  record: ProjectRecord
  project: ProjectDetail
  onClose: () => void
  onUpdate: () => void
}

export function RecordDetailModal({ record, project, onClose, onUpdate }: RecordDetailModalProps) {
  const [title, setTitle] = useState(record.title)
  const [description, setDescription] = useState(record.description)
  const [assignedTo, setAssignedTo] = useState(record.assigned_to ?? '')
  const [dueDate, setDueDate] = useState(record.due_date ?? '')
  const [editingTitle, setEditingTitle] = useState(false)
  const [editingDesc, setEditingDesc] = useState(false)
  const [saving, setSaving] = useState(false)
  const [users, setUsers] = useState<User[]>([])
  const [activities, setActivities] = useState<ProjectActivityLog[]>([])

  const column = project.columns.find(c => c.id === record.column_id)

  useEffect(() => {
    api.get<User[]>('/users/list').then(res => setUsers(res.data ?? [])).catch(() => {})
    projectService.getActivities(project.id).then(res => setActivities(res.data ?? [])).catch(() => {})
  }, [project.id])

  const save = async (data: Partial<ProjectRecord>) => {
    setSaving(true)
    try {
      await projectService.updateRecord(project.id, record.id, data)
    } catch { /* ignore */ }
    setSaving(false)
  }

  const handleTitleBlur = () => {
    setEditingTitle(false)
    if (title.trim() && title !== record.title) save({ title: title.trim() })
  }

  const handleDescBlur = () => {
    setEditingDesc(false)
    if (description !== record.description) save({ description })
  }

  const handleAssigneeChange = (val: string) => {
    setAssignedTo(val)
    save({ assigned_to: val || undefined } as Partial<ProjectRecord>)
  }

  const handleDueDateChange = (val: string) => {
    setDueDate(val)
    save({ due_date: val || undefined } as Partial<ProjectRecord>)
  }

  const handleComplete = async () => {
    try {
      await projectService.completeRecord(project.id, record.id)
      onUpdate()
    } catch { /* ignore */ }
  }

  const recordActivities = activities.filter(a => a.record_id === record.id)

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-20 bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl mx-4 max-h-[80vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
        {/* Header */}
        <div className="px-6 pt-5 pb-3 border-b border-outline-variant/10">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-xs text-on-surface-variant">
              <span>{project.name}</span>
              <span className="material-symbols-outlined text-[12px]">chevron_right</span>
              <span>{column?.name ?? 'Unknown'}</span>
            </div>
            <button onClick={onClose} className="p-1 hover:bg-surface-container-high rounded-lg transition-colors">
              <span className="material-symbols-outlined text-on-surface-variant">close</span>
            </button>
          </div>
        </div>

        <div className="p-6 flex flex-col gap-5">
          {/* Title */}
          {editingTitle ? (
            <input
              value={title}
              onChange={e => setTitle(e.target.value)}
              onBlur={handleTitleBlur}
              onKeyDown={e => { if (e.key === 'Enter') handleTitleBlur() }}
              className="text-xl font-bold text-on-surface bg-transparent border-b-2 border-primary outline-none w-full"
              autoFocus
            />
          ) : (
            <h2
              onClick={() => setEditingTitle(true)}
              className="text-xl font-bold text-on-surface cursor-pointer hover:text-primary transition-colors"
            >
              {title}
            </h2>
          )}

          {/* Fields */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1 block">Due Date</label>
              <input
                type="date"
                value={dueDate}
                onChange={e => handleDueDateChange(e.target.value)}
                className="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20"
                placeholder="Tidak ada due date"
              />
              {!dueDate && <p className="text-xs text-on-surface-variant mt-1">Tidak ada due date</p>}
            </div>
            <div>
              <label className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1 block">Assignee</label>
              <select
                value={assignedTo}
                onChange={e => handleAssigneeChange(e.target.value)}
                className="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20"
              >
                <option value="">Belum ditugaskan</option>
                {users.map(u => (
                  <option key={u.id} value={u.id}>{u.full_name}</option>
                ))}
              </select>
            </div>
          </div>

          {/* Description */}
          <div>
            <label className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1 block">Deskripsi</label>
            {editingDesc ? (
              <textarea
                value={description}
                onChange={e => setDescription(e.target.value)}
                onBlur={handleDescBlur}
                rows={4}
                className="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20 resize-none"
                autoFocus
              />
            ) : description ? (
              <p
                onClick={() => setEditingDesc(true)}
                className="text-sm text-on-surface cursor-pointer hover:bg-surface-container-low p-2 rounded-lg transition-colors whitespace-pre-wrap"
              >
                {description}
              </p>
            ) : (
              <button
                onClick={() => setEditingDesc(true)}
                className="text-sm font-bold text-primary hover:underline"
              >
                TAMBAH DESKRIPSI
              </button>
            )}
          </div>

          {/* Complete button */}
          {!record.is_completed && (
            <button
              onClick={handleComplete}
              className="flex items-center gap-2 px-4 py-2 text-sm font-bold text-on-primary bg-primary rounded-xl hover:opacity-90 transition-opacity self-start"
            >
              <span className="material-symbols-outlined text-[18px]">check_circle</span>
              TANDAI SELESAI
            </button>
          )}
          {record.is_completed && (
            <div className="flex items-center gap-2 text-sm text-on-surface-variant">
              <span className="material-symbols-outlined text-[18px] text-primary">check_circle</span>
              Selesai pada {record.completed_at ? new Date(record.completed_at).toLocaleString('id-ID') : '-'}
            </div>
          )}

          {/* Activity Feed */}
          <div>
            <h3 className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-3">Aktivitas</h3>
            {recordActivities.length === 0 ? (
              <p className="text-xs text-on-surface-variant/60">Belum ada aktivitas</p>
            ) : (
              <div className="flex flex-col gap-2">
                {recordActivities.map(a => (
                  <div key={a.id} className="flex items-start gap-2 text-xs text-on-surface-variant">
                    <span className="material-symbols-outlined text-[14px] mt-0.5">history</span>
                    <div>
                      <span className="font-medium">{a.action}</span>
                      {a.detail && <span className="ml-1">{a.detail}</span>}
                      <span className="ml-2 text-on-surface-variant/50">
                        {new Date(a.created_at).toLocaleString('id-ID')}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Last updated */}
          <p className="text-[10px] text-on-surface-variant/50 text-right">
            Terakhir diubah: {new Date(record.updated_at).toLocaleString('id-ID')}
          </p>
        </div>
      </div>
    </div>
  )
}
