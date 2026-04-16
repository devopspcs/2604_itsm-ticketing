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

type Tab = 'everything' | 'comments' | 'activity'

export function RecordDetailModal({ record, project, onClose, onUpdate }: RecordDetailModalProps) {
  const [title, setTitle] = useState(record.title)
  const [description, setDescription] = useState(record.description)
  const [assignees, setAssignees] = useState<string[]>(record.assignees ?? [])
  const [dueDate, setDueDate] = useState(record.due_date ?? '')
  const [editingTitle, setEditingTitle] = useState(false)
  const [editingDesc, setEditingDesc] = useState(false)
  const [users, setUsers] = useState<User[]>([])
  const [members, setMembers] = useState<string[]>([])
  const [activities, setActivities] = useState<ProjectActivityLog[]>([])
  const [activeTab, setActiveTab] = useState<Tab>('everything')
  const [comment, setComment] = useState('')

  const column = project.columns.find(c => c.id === record.column_id)

  useEffect(() => {
    // Load all users for name resolution
    api.get<User[]>('/users/list').then(res => setUsers(res.data ?? [])).catch(() => {})
    // Load project members to restrict assignee options
    projectService.listMembers(project.id).then(res => {
      setMembers((res.data ?? []).map(m => m.user_id))
    }).catch(() => {})
    projectService.getActivities(project.id).then(res => setActivities(res.data ?? [])).catch(() => {})
  }, [project.id])

  const save = async (data: Partial<ProjectRecord>) => {
    try { await projectService.updateRecord(project.id, record.id, data) } catch {}
  }

  const handleTitleBlur = () => {
    setEditingTitle(false)
    if (title.trim() && title !== record.title) save({ title: title.trim() })
  }

  const handleDescBlur = () => {
    setEditingDesc(false)
    if (description !== record.description) save({ description })
  }

  const toggleAssignee = (userId: string) => {
    const next = assignees.includes(userId)
      ? assignees.filter(id => id !== userId)
      : [...assignees, userId]
    setAssignees(next)
    save({ assignees: next } as any)
  }

  const handleDueDateChange = (val: string) => {
    setDueDate(val)
    save({ due_date: val || undefined } as Partial<ProjectRecord>)
  }

  const handleComplete = async () => {
    try { await projectService.completeRecord(project.id, record.id); onUpdate() } catch {}
  }

  const recordActivities = activities.filter(a => a.record_id === record.id)
  const lastEdit = record.updated_at ? new Date(record.updated_at) : null
  const timeSinceEdit = lastEdit ? formatTimeAgo(lastEdit) : ''

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-10 bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-4xl mx-4 max-h-[85vh] flex flex-col overflow-hidden" onClick={e => e.stopPropagation()}>

        {/* Top bar */}
        <div className="flex items-center justify-between px-6 py-3 border-b border-outline-variant/10">
          <div className="flex items-center gap-3">
            {!record.is_completed ? (
              <button onClick={handleComplete}
                className="px-4 py-1.5 text-xs font-bold border-2 border-primary text-primary rounded-lg hover:bg-primary hover:text-on-primary transition-colors uppercase tracking-wider">
                Mark Complete
              </button>
            ) : (
              <span className="px-4 py-1.5 text-xs font-bold border-2 border-emerald-500 text-emerald-600 rounded-lg uppercase tracking-wider flex items-center gap-1">
                <span className="material-symbols-outlined text-[14px]" style={{ fontVariationSettings: "'FILL' 1" }}>check_circle</span>
                Completed
              </span>
            )}
          </div>
          <div className="flex items-center gap-3 text-on-surface-variant">
            {timeSinceEdit && <span className="text-xs">Last edit was {timeSinceEdit}</span>}
            <button className="p-1 hover:bg-surface-container-high rounded-lg"><span className="material-symbols-outlined text-[18px]">attach_file</span></button>
            <button className="p-1 hover:bg-surface-container-high rounded-lg"><span className="material-symbols-outlined text-[18px]">push_pin</span></button>
            <button className="p-1 hover:bg-surface-container-high rounded-lg"><span className="material-symbols-outlined text-[18px]">settings</span></button>
            <button className="p-1 hover:bg-surface-container-high rounded-lg"><span className="material-symbols-outlined text-[18px]">more_horiz</span></button>
            <button onClick={onClose} className="p-1 hover:bg-surface-container-high rounded-lg"><span className="material-symbols-outlined">close</span></button>
          </div>
        </div>

        {/* Two-panel layout */}
        <div className="flex flex-1 overflow-hidden">
          {/* Left panel - Details */}
          <div className="flex-1 overflow-y-auto p-6 border-r border-outline-variant/10">
            {/* Title */}
            {editingTitle ? (
              <input value={title} onChange={e => setTitle(e.target.value)} onBlur={handleTitleBlur}
                onKeyDown={e => { if (e.key === 'Enter') handleTitleBlur() }}
                className="text-xl font-bold text-on-surface bg-transparent border-b-2 border-primary outline-none w-full mb-4" autoFocus />
            ) : (
              <h2 onClick={() => setEditingTitle(true)}
                className="text-xl font-bold text-on-surface cursor-pointer hover:text-primary transition-colors mb-4">
                {title}
              </h2>
            )}

            {/* Breadcrumb */}
            <div className="flex items-center gap-2 text-xs text-on-surface-variant mb-6">
              <span>In list</span>
              <span className="font-bold text-primary">{column?.name ?? 'Unknown'}</span>
              <span className="material-symbols-outlined text-[10px]">chevron_right</span>
              <span>In project</span>
              <span className="font-bold text-primary">{project.name}</span>
            </div>

            {/* Fields */}
            <div className="flex flex-col gap-4 mb-6">
              <div className="flex items-start gap-3">
                <span className="material-symbols-outlined text-on-surface-variant text-[18px] mt-0.5">calendar_today</span>
                <div className="flex-1">
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1">Due date</p>
                  <input type="date" value={dueDate} onChange={e => handleDueDateChange(e.target.value)}
                    className="px-3 py-1.5 bg-surface-container-lowest border border-outline-variant/20 rounded-lg text-sm outline-none focus:ring-2 focus:ring-primary/20" />
                  {!dueDate && <p className="text-xs text-on-surface-variant mt-1">No due date</p>}
                </div>
              </div>

              <div className="flex items-start gap-3">
                <span className="material-symbols-outlined text-on-surface-variant text-[18px] mt-0.5">group</span>
                <div className="flex-1">
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1">Assignees</p>
                  <div className="flex flex-wrap gap-2">
                    {assignees.length === 0 && <span className="text-xs text-on-surface-variant">Unassigned</span>}
                    {assignees.map(uid => {
                      const u = users.find(u => u.id === uid)
                      return u ? (
                        <span key={uid} className="inline-flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary text-xs font-medium rounded-lg">
                          {u.full_name}
                          <button onClick={() => toggleAssignee(uid)} className="hover:text-error"><span className="material-symbols-outlined text-[12px]">close</span></button>
                        </span>
                      ) : null
                    })}
                  </div>
                  <div className="mt-2 flex flex-wrap gap-1">
                    {users.filter(u => !assignees.includes(u.id) && members.includes(u.id)).map(u => (
                      <button key={u.id} onClick={() => toggleAssignee(u.id)}
                        className="text-[10px] px-2 py-1 bg-surface-container-high text-on-surface-variant rounded-lg hover:bg-primary/10 hover:text-primary transition-colors">
                        + {u.full_name}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <span className="material-symbols-outlined text-on-surface-variant text-[18px] mt-0.5">sell</span>
                <div>
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1">Tags</p>
                  <span className="text-xs text-on-surface-variant">No tags</span>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <span className="material-symbols-outlined text-on-surface-variant text-[18px] mt-0.5">share</span>
                <div>
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-1">Dependencies</p>
                  <span className="text-xs text-on-surface-variant">No dependencies</span>
                </div>
              </div>
            </div>

            {/* Description */}
            <div className="mb-4">
              {editingDesc ? (
                <div>
                  <textarea value={description} onChange={e => setDescription(e.target.value)} onBlur={handleDescBlur}
                    rows={4} className="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20 resize-none" autoFocus />
                </div>
              ) : description ? (
                <p onClick={() => setEditingDesc(true)} className="text-sm text-on-surface cursor-pointer hover:bg-surface-container-low p-2 rounded-lg transition-colors whitespace-pre-wrap">{description}</p>
              ) : (
                <button onClick={() => setEditingDesc(true)} className="flex items-center gap-2 text-sm font-bold text-primary hover:underline">
                  <span className="material-symbols-outlined text-[16px]">edit_note</span> ADD DESCRIPTION
                </button>
              )}
            </div>

            <button className="flex items-center gap-2 text-sm font-bold text-primary hover:underline">
              <span className="material-symbols-outlined text-[16px]">checklist</span> ADD CHECKLIST
            </button>
          </div>

          {/* Right panel - Activity/Comments */}
          <div className="w-[380px] flex flex-col bg-surface-container-lowest/50">
            {/* Tabs */}
            <div className="flex border-b border-outline-variant/10">
              {(['everything', 'comments', 'activity'] as Tab[]).map(tab => (
                <button key={tab} onClick={() => setActiveTab(tab)}
                  className={`flex-1 py-3 text-[10px] font-bold uppercase tracking-widest transition-colors ${
                    activeTab === tab ? 'text-primary border-b-2 border-primary' : 'text-on-surface-variant hover:text-on-surface'
                  }`}>
                  {tab === 'everything' ? 'Everything' : tab === 'comments' ? 'Comments Only' : 'Activity Only'}
                </button>
              ))}
            </div>

            {/* Feed */}
            <div className="flex-1 overflow-y-auto p-4">
              {(() => {
                const filtered = activeTab === 'comments'
                  ? recordActivities.filter(a => a.action === 'comment')
                  : activeTab === 'activity'
                  ? recordActivities.filter(a => a.action !== 'comment')
                  : recordActivities
                return filtered.length === 0 ? (
                <p className="text-xs text-on-surface-variant/60 text-center py-8">
                  {activeTab === 'comments' ? 'No comments yet' : activeTab === 'activity' ? 'No activity yet' : 'No activity yet'}
                </p>
              ) : (
                <div className="flex flex-col gap-3">
                  {filtered.map(a => {
                    const actor = users.find(u => u.id === a.actor_id)
                    return (
                      <div key={a.id} className="flex items-start gap-3">
                        <div className="w-7 h-7 rounded-full bg-primary/10 flex items-center justify-center shrink-0 mt-0.5">
                          <span className="text-[10px] font-bold text-primary">{actor?.full_name?.charAt(0) ?? '?'}</span>
                        </div>
                        <div className="text-xs text-on-surface-variant">
                          <span className="font-semibold text-on-surface">{actor?.full_name ?? 'Unknown'}</span>
                          <span className="ml-1">{a.detail}</span>
                          <span className="block text-on-surface-variant/50 mt-0.5">
                            {formatTimeAgo(new Date(a.created_at))}
                          </span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )
              })()}
            </div>

            {/* Comment input */}
            <div className="p-4 border-t border-outline-variant/10">
              <input value={comment} onChange={e => setComment(e.target.value)}
                placeholder="Enter a comment..."
                className="w-full px-3 py-2.5 bg-white border border-outline-variant/20 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20"
                onKeyDown={async e => {
                  if (e.key === 'Enter' && comment.trim()) {
                    const text = comment.trim()
                    setComment('')
                    try {
                      await projectService.addComment(project.id, record.id, text)
                      // Refresh activities
                      const res = await projectService.getActivities(project.id)
                      setActivities(res.data ?? [])
                    } catch {}
                  }
                }} />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function formatTimeAgo(date: Date): string {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins} minutes ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} hours ago`
  const days = Math.floor(hours / 24)
  return `${days} days ago`
}
