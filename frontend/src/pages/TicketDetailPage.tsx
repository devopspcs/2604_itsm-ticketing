import { useEffect, useState, useRef } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { ticketService } from '../services/ticket.service'
import { approvalService } from '../services/approval.service'
import api from '../services/api'
import type { Ticket, Approval, ActivityLog, User } from '../types'
import type { RootState } from '../store'

const actionIcon: Record<string, string> = {
  ticket_created: 'add_circle', status_changed: 'swap_horiz', assigned: 'person_add',
  reassigned: 'person_search', approval_requested: 'pending', approval_decided: 'verified',
  field_updated: 'edit_note',
}

// ── Modal wrapper ──────────────────────────────────────────────────────────
function Modal({ title, onClose, children }: { title: string; onClose: () => void; children: React.ReactNode }) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg mx-4 overflow-hidden">
        <div className="flex items-center justify-between px-6 py-4 border-b border-outline-variant/20">
          <h3 className="text-lg font-bold text-on-surface">{title}</h3>
          <button onClick={onClose} className="p-1 hover:bg-surface-container-high rounded-lg transition-colors">
            <span className="material-symbols-outlined text-on-surface-variant">close</span>
          </button>
        </div>
        <div className="px-6 py-5">{children}</div>
      </div>
    </div>
  )
}

// ── Attachment thumbnail with auth ─────────────────────────────────────────
function AttachmentThumb({ ticketId, attachmentId }: { ticketId: string; attachmentId: string }) {
  const [src, setSrc] = useState<string>('')
  useEffect(() => {
    api.get(`/tickets/${ticketId}/attachments/${attachmentId}`, { responseType: 'blob' })
      .then(res => {
        const url = URL.createObjectURL(res.data)
        setSrc(url)
      })
      .catch(() => {})
    return () => { if (src) URL.revokeObjectURL(src) }
  }, [ticketId, attachmentId])

  if (!src) return null
  return (
    <img src={src} alt="attachment" className="w-12 h-12 rounded-lg object-cover flex-shrink-0 border border-outline-variant/20" />
  )
}

// ── Note image with auth ───────────────────────────────────────────────────
function NoteImage({ ticketId, noteId, imageName }: { ticketId: string; noteId: string; imageName?: string }) {
  const [src, setSrc] = useState<string>('')
  useEffect(() => {
    api.get(`/tickets/${ticketId}/notes/${noteId}/image`, { responseType: 'blob' })
      .then(res => setSrc(URL.createObjectURL(res.data)))
      .catch(() => {})
    return () => { if (src) URL.revokeObjectURL(src) }
  }, [ticketId, noteId])

  if (!src) return null
  return (
    <div className="mt-3">
      <img
        src={src}
        alt={imageName ?? 'Note image'}
        className="max-w-full max-h-64 rounded-xl border border-outline-variant/20 cursor-pointer hover:opacity-90 transition-opacity"
        onClick={() => window.open(src, '_blank')}
      />
      {imageName && <p className="text-[10px] text-on-surface-variant mt-1">{imageName}</p>}
    </div>
  )
}

export function TicketDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [ticket, setTicket] = useState<Ticket | null>(null)
  const [approvals, setApprovals] = useState<Approval[]>([])
  const [activities, setActivities] = useState<ActivityLog[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [attachments, setAttachments] = useState<{id:string;filename:string;file_size:number;mime_type:string;created_at:string;uploaded_by?:string}[]>([])
  const [notes, setNotes] = useState<{id:string;content:string;has_image:boolean;image_name?:string;created_at:string;author_id?:string}[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [uploading, setUploading] = useState(false)
  const role = useSelector((s: RootState) => s.auth.role) ?? 'user'
  const canApprove = role === 'admin' || role === 'approver'
  const canAssign = role === 'admin' || role === 'approver'

  // Modal states
  const [modal, setModal] = useState<'edit' | 'assign' | 'note' | 'approve' | null>(null)

  // Form states
  const [editTitle, setEditTitle] = useState('')
  const [editDesc, setEditDesc] = useState('')
  const [editPriority, setEditPriority] = useState('')
  const [editCategory, setEditCategory] = useState('')
  const [assigneeId, setAssigneeId] = useState('')
  const [noteText, setNoteText] = useState('')
  const [approvalComment, setApprovalComment] = useState('')
  const [saving, setSaving] = useState(false)
  const [actionMsg, setActionMsg] = useState('')

  const load = async () => {
    if (!id) return
    try {
      const [t, a, act, att, notesRes] = await Promise.all([
        ticketService.get(id),
        ticketService.getApprovals(id),
        ticketService.getActivities(id),
        api.get(`/tickets/${id}/attachments`),
        api.get(`/tickets/${id}/notes`),
      ])
      setTicket(t.data)
      setApprovals(a.data ?? [])
      setActivities(act.data ?? [])
      setAttachments(att.data ?? [])
      setNotes(notesRes.data ?? [])
      setEditTitle(t.data.title)
      setEditDesc(t.data.description)
      setEditPriority(t.data.priority)
      setEditCategory(t.data.category ?? '')
    } catch {
      setError('Failed to load ticket')
    } finally {
      setLoading(false)
    }
  }

  const loadUsers = async () => {
    try {
      const res = await api.get<User[]>('/users/list')
      setUsers(res.data ?? [])
    } catch {}
  }

  const userName = (userId?: string) => {
    if (!userId) return 'Unknown'
    const u = users.find(u => u.id === userId)
    return u ? u.full_name : userId.slice(0, 8) + '...'
  }

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file || !id) return
    setUploading(true)
    try {
      const formData = new FormData()
      formData.append('file', file)
      await api.post(`/tickets/${id}/attachments`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      await load()
      showSuccess('File uploaded successfully')
    } catch {
      showSuccess('Failed to upload file')
    } finally {
      setUploading(false)
      if (fileInputRef.current) fileInputRef.current.value = ''
    }
  }

  const handleDeleteAttachment = async (attachmentId: string) => {
    if (!id || !window.confirm('Delete this attachment?')) return
    try {
      await api.delete(`/tickets/${id}/attachments/${attachmentId}`)
      await load()
      showSuccess('Attachment deleted')
    } catch {
      showSuccess('Failed to delete attachment')
    }
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  useEffect(() => { load(); loadUsers() }, [id])

  const closeModal = () => { setModal(null); setActionMsg('') }

  const showSuccess = (msg: string) => {
    setActionMsg(msg)
    setTimeout(() => setActionMsg(''), 3000)
  }

  // ── Actions ──────────────────────────────────────────────────────────────

  const handleEdit = async () => {
    if (!id) return
    setSaving(true)
    try {
      await ticketService.update(id, {
        title: editTitle,
        description: editDesc,
        priority: editPriority as Ticket['priority'],
        category: editCategory,
      })
      await load()
      closeModal()
      showSuccess('Ticket updated successfully')
    } catch { setActionMsg('Failed to update ticket') }
    finally { setSaving(false) }
  }

  const handleAssign = async () => {
    if (!id || !assigneeId) return
    setSaving(true)
    try {
      await ticketService.assign(id, assigneeId)
      await load()
      closeModal()
      showSuccess('Ticket assigned successfully')
    } catch { setActionMsg('Failed to assign ticket') }
    finally { setSaving(false) }
  }

  const handleSubmitApproval = async () => {
    if (!id) return
    setSaving(true)
    try {
      await ticketService.submit(id)
      await load()
      closeModal()
      showSuccess('Submitted for approval')
    } catch { setActionMsg('Failed to submit for approval') }
    finally { setSaving(false) }
  }

  const handleDecide = async (decision: 'approved' | 'rejected') => {
    if (!id) return
    setSaving(true)
    try {
      await approvalService.decide(id, decision, approvalComment)
      await load()
      closeModal()
      showSuccess(`Ticket ${decision}`)
    } catch { setActionMsg('Failed to record decision') }
    finally { setSaving(false) }
  }

  const handleArchive = async () => {
    if (!id || !ticket) return
    if (!window.confirm('Archive this ticket? It will be marked as Done.')) return
    setSaving(true)
    try {
      await ticketService.update(id, { status: 'done' } as Partial<Ticket>)
      await load()
      showSuccess('Ticket archived')
    } catch { showSuccess('Failed to archive ticket') }
    finally { setSaving(false) }
  }

  const openAssign = () => {
    loadUsers()
    setModal('assign')
  }

  if (loading) return (
    <div className="flex items-center justify-center py-20">
      <div className="flex items-center gap-3 text-on-surface-variant">
        <span className="material-symbols-outlined animate-spin">refresh</span>
        Loading ticket...
      </div>
    </div>
  )

  if (error || !ticket) return (
    <div className="max-w-7xl mx-auto p-8">
      <div className="bg-error-container/30 text-error px-6 py-4 rounded-xl">{error || 'Ticket not found'}</div>
    </div>
  )

  return (
    <div className="max-w-7xl mx-auto px-8 pb-12 pt-8">

      {/* Success/Error toast */}
      {actionMsg && (
        <div className={`fixed top-20 right-6 z-50 px-5 py-3 rounded-xl shadow-lg font-semibold text-sm flex items-center gap-2 ${
          actionMsg.startsWith('Failed') ? 'bg-error-container text-error' : 'bg-emerald-100 text-emerald-800'
        }`}>
          <span className="material-symbols-outlined text-sm">{actionMsg.startsWith('Failed') ? 'error' : 'check_circle'}</span>
          {actionMsg}
        </div>
      )}

      {/* Header */}
      <header className="mb-10 flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <span className="text-xs font-semibold tracking-wider text-primary py-1 px-2 bg-primary-fixed rounded">
              {ticket.type === 'incident' ? 'INC' : ticket.type === 'change_request' ? 'CHG' : 'REQ'}-{ticket.id.slice(0, 6).toUpperCase()}
            </span>
            <span className={`flex items-center gap-1 text-xs font-bold px-2 py-1 rounded ${
              ticket.priority === 'critical' ? 'text-error bg-error-container' :
              ticket.priority === 'high' ? 'text-amber-700 bg-amber-100' :
              'text-on-surface-variant bg-surface-container-high'
            }`}>
              <span className="material-symbols-outlined text-sm" style={{ fontVariationSettings: "'FILL' 1" }}>priority_high</span>
              {ticket.priority.toUpperCase()}
            </span>
            <span className={`text-xs font-bold px-2 py-1 rounded capitalize ${
              ticket.status === 'open' ? 'bg-red-100 text-red-700' :
              ticket.status === 'in_progress' ? 'bg-amber-100 text-amber-700' :
              ticket.status === 'pending_approval' ? 'bg-purple-100 text-purple-700' :
              ticket.status === 'approved' ? 'bg-emerald-100 text-emerald-700' :
              ticket.status === 'rejected' ? 'bg-red-100 text-red-700' :
              'bg-surface-container-high text-on-surface-variant'
            }`}>{ticket.status.replace(/_/g, ' ')}</span>
          </div>
          <h1 className="text-3xl font-extrabold text-on-surface leading-tight">{ticket.title}</h1>
        </div>
        <div className="flex gap-3">
          <Link to="/tickets" className="bg-surface-container-high text-on-secondary-container px-5 py-2.5 rounded-xl font-semibold flex items-center gap-2 hover:bg-slate-200 transition-all">
            <span className="material-symbols-outlined">arrow_back</span>
            Back
          </Link>
        </div>
      </header>

      {/* Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">

        {/* Main Content */}
        <div className="lg:col-span-8 space-y-6">

          {/* Description */}
          <section className="bg-surface-container-lowest rounded-xl p-8 shadow-sm">
            <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
              <span className="material-symbols-outlined text-primary">description</span>
              Description
            </h3>
            <p className="text-on-surface-variant leading-relaxed">{ticket.description}</p>
            <div className="mt-6 p-4 bg-surface-container-low rounded-xl border-l-4 border-primary grid grid-cols-2 gap-3 text-sm">
              <div><span className="text-on-surface-variant">Type:</span> <span className="font-medium capitalize ml-1">{ticket.type.replace(/_/g, ' ')}</span></div>
              <div><span className="text-on-surface-variant">Category:</span> <span className="font-medium ml-1">{ticket.category || '-'}</span></div>
              <div><span className="text-on-surface-variant">Priority:</span> <span className="font-medium capitalize ml-1">{ticket.priority}</span></div>
              <div><span className="text-on-surface-variant">Created by:</span> <span className="font-medium ml-1">{userName(ticket.created_by)}</span></div>
              <div><span className="text-on-surface-variant">Assigned to:</span> <span className="font-medium ml-1">{ticket.assigned_to ? userName(ticket.assigned_to) : 'Unassigned'}</span></div>
            </div>
          </section>

          {/* Approval Workflow */}
          <section className="bg-surface-container-lowest rounded-xl p-8 shadow-sm">
            <h3 className="text-lg font-bold mb-6 flex items-center gap-2">
              <span className="material-symbols-outlined text-primary">schema</span>
              Approval Workflow
            </h3>

            {ticket.status === 'pending_approval' && canApprove && (
              <div className="mb-6 p-5 bg-primary-fixed rounded-xl border-l-4 border-primary">
                <div className="flex justify-between items-start mb-3">
                  <span className="text-xs font-bold text-primary uppercase">Awaiting Decision</span>
                  <span className="material-symbols-outlined text-primary animate-pulse">pending</span>
                </div>
                <textarea
                  value={approvalComment}
                  onChange={(e) => setApprovalComment(e.target.value)}
                  placeholder="Add a comment (optional)..."
                  rows={2}
                  className="w-full bg-white/60 border-none rounded-xl px-3 py-2 text-sm mb-3 outline-none resize-none"
                />
                <div className="flex gap-2">
                  <button onClick={() => handleDecide('approved')} disabled={saving}
                    className="bg-primary text-white text-xs px-4 py-2 rounded-lg font-bold hover:bg-primary-container transition-colors disabled:opacity-70 flex items-center gap-1">
                    <span className="material-symbols-outlined text-sm">check_circle</span> Approve
                  </button>
                  <button onClick={() => handleDecide('rejected')} disabled={saving}
                    className="bg-white text-error text-xs px-4 py-2 rounded-lg font-bold hover:bg-error-container transition-colors disabled:opacity-70 flex items-center gap-1">
                    <span className="material-symbols-outlined text-sm">cancel</span> Reject
                  </button>
                </div>
              </div>
            )}

            {approvals.length === 0 ? (
              <p className="text-sm text-on-surface-variant text-center py-6">No approval history yet.</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {approvals.map((a, i) => (
                  <div key={a.id} className={`p-5 rounded-xl border-l-4 ${
                    a.decision === 'approved' ? 'bg-surface-container-low border-green-500' :
                    a.decision === 'rejected' ? 'bg-error-container/20 border-error' :
                    'bg-surface-container-low border-outline-variant opacity-60'
                  }`}>
                    <div className="flex justify-between items-start mb-2">
                      <span className="text-xs font-bold text-on-surface-variant uppercase">Step {i + 1} · Level {a.level}</span>
                      <span className="material-symbols-outlined text-sm" style={{ fontVariationSettings: "'FILL' 1" }}>
                        {a.decision === 'approved' ? 'check_circle' : a.decision === 'rejected' ? 'cancel' : 'hourglass_empty'}
                      </span>
                    </div>
                    <p className="text-sm font-bold text-on-surface capitalize">{a.decision ?? 'Pending'}</p>
                    <p className="text-xs text-on-surface-variant">by {userName(a.approver_id)}</p>
                    {a.comment && <p className="mt-2 text-xs italic text-on-surface-variant">"{a.comment}"</p>}
                    {a.decided_at && <p className="text-xs text-on-surface-variant mt-1">{new Date(a.decided_at).toLocaleString()}</p>}
                  </div>
                ))}
              </div>
            )}
          </section>

          {/* Activity Log */}
          <section className="bg-surface-container-lowest rounded-xl p-8 shadow-sm">
            <h3 className="text-lg font-bold mb-6 flex items-center gap-2">
              <span className="material-symbols-outlined text-primary">history</span>
              Activity Log
            </h3>
            {activities.length === 0 ? (
              <p className="text-sm text-on-surface-variant text-center py-6">No activity yet.</p>
            ) : (
              <div className="space-y-4">
                {activities.map((a) => (
                  <div key={a.id} className="flex gap-4">
                    <div className="flex-shrink-0 w-8 h-8 rounded-xl bg-primary-fixed flex items-center justify-center text-primary">
                      <span className="material-symbols-outlined text-sm">{actionIcon[a.action] ?? 'info'}</span>
                    </div>
                    <div className="flex-grow">
                      <div className="flex justify-between items-start">
                        <div>
                          <p className="text-sm font-semibold text-on-surface capitalize">{a.action.replace(/_/g, ' ')}</p>
                          <p className="text-[10px] text-on-surface-variant">by {userName(a.actor_id)}</p>
                        </div>
                        <span className="text-[10px] text-on-surface-variant">{new Date(a.created_at).toLocaleString()}</span>
                      </div>
                      {a.new_value && <p className="text-xs text-on-surface-variant mt-0.5">→ {a.new_value}</p>}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>

          {/* Attachments */}
          <section className="bg-surface-container-lowest rounded-xl p-8 shadow-sm">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-lg font-bold flex items-center gap-2">
                <span className="material-symbols-outlined text-primary">attachment</span>
                Attachments
              </h3>
              <label className="px-4 py-2 bg-primary text-white rounded-xl text-xs font-bold cursor-pointer hover:opacity-90 transition-opacity flex items-center gap-1">
                <span className="material-symbols-outlined text-sm">upload</span>
                Upload File
                <input
                  type="file"
                  className="hidden"
                  onChange={handleUpload}
                  accept="image/*,.pdf,.doc,.docx,.txt,.csv,.xlsx"
                />
              </label>
            </div>
            {uploading && (
              <div className="flex items-center gap-2 text-sm text-on-surface-variant mb-4">
                <span className="material-symbols-outlined animate-spin text-sm">refresh</span>
                Uploading...
              </div>
            )}
            {attachments.length === 0 ? (
              <p className="text-sm text-on-surface-variant text-center py-6">No attachments yet.</p>
            ) : (
              <div className="space-y-3">
                {attachments.map((a) => (
                  <div key={a.id} className="flex items-center justify-between p-3 bg-surface-container-low rounded-xl group">
                    <div className="flex items-center gap-3 min-w-0">
                      <span className="material-symbols-outlined text-primary flex-shrink-0">
                        {a.mime_type?.startsWith('image/') ? 'image' : 'description'}
                      </span>
                      <div className="min-w-0">
                        <button
                          onClick={async () => {
                            try {
                              const res = await api.get(`/tickets/${id}/attachments/${a.id}`, { responseType: 'blob' })
                              const url = window.URL.createObjectURL(res.data)
                              if (a.mime_type?.startsWith('image/') || a.mime_type === 'application/pdf') {
                                window.open(url, '_blank')
                              } else {
                                const link = document.createElement('a')
                                link.href = url
                                link.download = a.filename
                                link.click()
                                window.URL.revokeObjectURL(url)
                              }
                            } catch { showSuccess('Failed to download file') }
                          }}
                          className="text-sm font-medium text-primary hover:underline truncate block text-left"
                        >
                          {a.filename}
                        </button>
                        <p className="text-[10px] text-on-surface-variant">
                          {(a.file_size / 1024).toFixed(1)} KB · uploaded by {userName(a.uploaded_by)} · {new Date(a.created_at).toLocaleString()}
                        </p>
                      </div>
                    </div>
                    {a.mime_type?.startsWith('image/') && (
                      <AttachmentThumb ticketId={id!} attachmentId={a.id} />
                    )}
                  </div>
                ))}
              </div>
            )}
          </section>

          {/* Internal Notes */}
          <section className="bg-surface-container-lowest rounded-xl p-8 shadow-sm">
            <h3 className="text-lg font-bold mb-6 flex items-center gap-2">
              <span className="material-symbols-outlined text-primary">sticky_note_2</span>
              Internal Notes
            </h3>
            {notes.length === 0 ? (
              <p className="text-sm text-on-surface-variant text-center py-6">No notes yet. Click "Add Internal Note" to add one.</p>
            ) : (
              <div className="space-y-4">
                {notes.map((n) => (
                  <div key={n.id} className="p-4 bg-surface-container-low rounded-xl border-l-4 border-secondary">
                    <p className="text-sm text-on-surface leading-relaxed">{n.content}</p>
                    {n.has_image && (
                      <NoteImage ticketId={id!} noteId={n.id} imageName={n.image_name} />
                    )}
                    <p className="text-[10px] text-on-surface-variant mt-2">by {userName(n.author_id)} · {new Date(n.created_at).toLocaleString()}</p>
                  </div>
                ))}
              </div>
            )}
          </section>
        </div>

        {/* Sidebar */}
        <aside className="lg:col-span-4 space-y-6">
          {/* Metadata */}
          <section className="bg-surface-container-low rounded-xl p-6">
            <h3 className="text-sm font-bold text-on-surface-variant uppercase tracking-widest mb-4">Ticket Metadata</h3>
            <div className="space-y-3">
              {[
                { label: 'Status', value: ticket.status.replace(/_/g, ' ') },
                { label: 'Priority', value: ticket.priority },
                { label: 'Type', value: ticket.type.replace(/_/g, ' ') },
                { label: 'Category', value: ticket.category || '-' },
                { label: 'Created', value: new Date(ticket.created_at).toLocaleDateString() },
                { label: 'Updated', value: new Date(ticket.updated_at).toLocaleDateString() },
              ].map(({ label, value }) => (
                <div key={label} className="flex justify-between items-center">
                  <span className="text-xs text-on-surface-variant">{label}</span>
                  <span className="text-xs font-bold text-on-surface bg-white px-2 py-1 rounded shadow-sm capitalize">{value}</span>
                </div>
              ))}
            </div>
          </section>

          {/* Quick Actions */}
          <section className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
            <h3 className="text-sm font-bold text-on-surface-variant uppercase tracking-widest mb-4">Quick Actions</h3>
            <div className="space-y-1">
              <button
                onClick={() => setModal('edit')}
                className="w-full text-left p-3 rounded-xl hover:bg-surface-container-low transition-colors flex items-center gap-3"
              >
                <span className="material-symbols-outlined text-primary">edit</span>
                <span className="text-sm font-medium">Edit Ticket Detail</span>
              </button>

              <button
                onClick={() => setModal('note')}
                className="w-full text-left p-3 rounded-xl hover:bg-surface-container-low transition-colors flex items-center gap-3"
              >
                <span className="material-symbols-outlined text-primary">add_comment</span>
                <span className="text-sm font-medium">Add Internal Note</span>
              </button>

              <button
                onClick={openAssign}
                className={`w-full text-left p-3 rounded-xl hover:bg-surface-container-low transition-colors flex items-center gap-3 ${!canAssign ? 'hidden' : ''}`}
              >
                <span className="material-symbols-outlined text-primary">person_add</span>
                <span className="text-sm font-medium">Assign Ticket</span>
              </button>

              {(ticket.status === 'open' || ticket.status === 'in_progress') ? (
                <button
                  onClick={handleSubmitApproval}
                  disabled={saving}
                  className="w-full text-left p-3 rounded-xl hover:bg-surface-container-low transition-colors flex items-center gap-3 disabled:opacity-60"
                >
                  <span className="material-symbols-outlined text-primary">send</span>
                  <span className="text-sm font-medium">Request Approval</span>
                </button>
              ) : null}

              {ticket.status !== 'done' && (
                <button
                  onClick={handleArchive}
                  disabled={saving}
                  className="w-full text-left p-3 rounded-xl hover:bg-error-container/20 text-error transition-colors flex items-center gap-3 mt-2 disabled:opacity-60"
                >
                  <span className="material-symbols-outlined">archive</span>
                  <span className="text-sm font-medium">Archive Ticket</span>
                </button>
              )}
            </div>
          </section>
        </aside>
      </div>

      {/* ── MODALS ─────────────────────────────────────────────────────────── */}

      {/* Edit Modal */}
      {modal === 'edit' && (
        <Modal title="Edit Ticket" onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Title</label>
              <input value={editTitle} onChange={(e) => setEditTitle(e.target.value)}
                className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-primary" />
            </div>
            <div>
              <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Description</label>
              <textarea value={editDesc} onChange={(e) => setEditDesc(e.target.value)} rows={4}
                className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-primary resize-none" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Priority</label>
                <select value={editPriority} onChange={(e) => setEditPriority(e.target.value)}
                  className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none appearance-none">
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="critical">Critical</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Category</label>
                <input value={editCategory} onChange={(e) => setEditCategory(e.target.value)}
                  className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-primary" />
              </div>
            </div>
            {actionMsg && <p className="text-sm text-error">{actionMsg}</p>}
            <div className="flex justify-end gap-3 pt-2">
              <button onClick={closeModal} className="px-5 py-2.5 rounded-xl font-semibold text-on-surface-variant hover:bg-surface-container-high transition-colors">Cancel</button>
              <button onClick={handleEdit} disabled={saving}
                className="px-6 py-2.5 bg-primary text-white rounded-xl font-bold hover:opacity-90 transition-all disabled:opacity-70">
                {saving ? 'Saving...' : 'Save Changes'}
              </button>
            </div>
          </div>
        </Modal>
      )}

      {/* Assign Modal */}
      {modal === 'assign' && (
        <Modal title="Assign Ticket" onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Select Assignee</label>
              {users.length === 0 ? (
                <p className="text-sm text-on-surface-variant">Loading users...</p>
              ) : (
                <select value={assigneeId} onChange={(e) => setAssigneeId(e.target.value)}
                  className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none appearance-none">
                  <option value="">-- Select user --</option>
                  {users.filter(u => u.is_active).map(u => (
                    <option key={u.id} value={u.id}>{u.full_name} ({u.role})</option>
                  ))}
                </select>
              )}
            </div>
            {actionMsg && <p className="text-sm text-error">{actionMsg}</p>}
            <div className="flex justify-end gap-3 pt-2">
              <button onClick={closeModal} className="px-5 py-2.5 rounded-xl font-semibold text-on-surface-variant hover:bg-surface-container-high transition-colors">Cancel</button>
              <button onClick={handleAssign} disabled={saving || !assigneeId}
                className="px-6 py-2.5 bg-primary text-white rounded-xl font-bold hover:opacity-90 transition-all disabled:opacity-70">
                {saving ? 'Assigning...' : 'Assign'}
              </button>
            </div>
          </div>
        </Modal>
      )}

      {/* Add Note Modal */}
      {modal === 'note' && (
        <Modal title="Add Internal Note" onClose={closeModal}>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Note</label>
              <textarea value={noteText} onChange={(e) => setNoteText(e.target.value)} rows={4}
                placeholder="Write your internal note here..."
                className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-primary resize-none" />
            </div>
            <div>
              <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Attach Image (optional)</label>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                className="w-full text-sm text-on-surface-variant file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-semibold file:bg-primary-fixed file:text-primary hover:file:bg-primary-fixed-dim"
              />
              <p className="text-xs text-on-surface-variant mt-1">PNG, JPG, GIF up to 10MB — great for evidence screenshots</p>
            </div>
            <div className="flex justify-end gap-3 pt-2">
              <button onClick={closeModal} className="px-5 py-2.5 rounded-xl font-semibold text-on-surface-variant hover:bg-surface-container-high transition-colors">Cancel</button>
              <button
                onClick={async () => {
                  if (!noteText.trim() || !id) return
                  setSaving(true)
                  try {
                    const formData = new FormData()
                    formData.append('content', noteText)
                    const file = fileInputRef.current?.files?.[0]
                    if (file) formData.append('image', file)
                    await api.post(`/tickets/${id}/notes`, formData, {
                      headers: { 'Content-Type': 'multipart/form-data' },
                    })
                    setNoteText('')
                    if (fileInputRef.current) fileInputRef.current.value = ''
                    closeModal()
                    showSuccess('Note added')
                    await load()
                  } catch { showSuccess('Failed to add note') }
                  finally { setSaving(false) }
                }}
                disabled={saving || !noteText.trim()}
                className="px-6 py-2.5 bg-primary text-white rounded-xl font-bold hover:opacity-90 transition-all disabled:opacity-70">
                {saving ? 'Adding...' : 'Add Note'}
              </button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  )
}
