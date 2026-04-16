import { useEffect, useState } from 'react'
import { projectService } from '../../services/project.service'
import api from '../../services/api'
import type { User } from '../../types'

interface Member {
  project_id: string
  user_id: string
  role: string
  created_at: string
}

interface MemberManagementProps {
  projectId: string
  onClose: () => void
}

export function MemberManagement({ projectId, onClose }: MemberManagementProps) {
  const [members, setMembers] = useState<Member[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  const fetchMembers = () => {
    projectService.listMembers(projectId)
      .then(res => setMembers(res.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    fetchMembers()
    api.get<User[]>('/users/list').then(res => setUsers(res.data ?? [])).catch(() => {})
  }, [projectId])

  const memberIds = new Set(members.map(m => m.user_id))
  const nonMembers = users.filter(u => !memberIds.has(u.id) && u.full_name.toLowerCase().includes(search.toLowerCase()))

  const invite = async (userId: string) => {
    try {
      await projectService.inviteMember(projectId, userId)
      fetchMembers()
    } catch {}
  }

  const remove = async (userId: string) => {
    try {
      await projectService.removeMember(projectId, userId)
      fetchMembers()
    } catch {}
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg mx-4 max-h-[80vh] flex flex-col overflow-hidden" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between px-6 py-4 border-b border-outline-variant/10">
          <h2 className="text-lg font-bold text-on-surface font-headline">Project Members</h2>
          <button onClick={onClose} className="p-1 hover:bg-surface-container-high rounded-lg">
            <span className="material-symbols-outlined text-on-surface-variant">close</span>
          </button>
        </div>

        <div className="p-4 border-b border-outline-variant/10">
          <input value={search} onChange={e => setSearch(e.target.value)}
            placeholder="Search users to invite..."
            className="w-full px-4 py-2 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20" />
        </div>

        <div className="flex-1 overflow-y-auto">
          {/* Current members */}
          <div className="p-4">
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-3">
              Members ({members.length})
            </p>
            {loading ? (
              <div className="animate-pulse space-y-2">
                {[1,2].map(i => <div key={i} className="h-10 bg-surface-container-high rounded-lg" />)}
              </div>
            ) : (
              <div className="flex flex-col gap-2">
                {members.map(m => {
                  const u = users.find(u => u.id === m.user_id)
                  return (
                    <div key={m.user_id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-surface-container-low">
                      <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-primary-container flex items-center justify-center text-white text-xs font-bold">
                        {u?.full_name?.charAt(0) ?? '?'}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-on-surface truncate">{u?.full_name ?? m.user_id}</p>
                        <p className="text-[10px] text-on-surface-variant">{m.role}</p>
                      </div>
                      {m.role !== 'owner' && (
                        <button onClick={() => remove(m.user_id)}
                          className="text-xs text-error hover:underline font-medium">Remove</button>
                      )}
                    </div>
                  )
                })}
              </div>
            )}
          </div>

          {/* Invite section */}
          {nonMembers.length > 0 && (
            <div className="p-4 border-t border-outline-variant/10">
              <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-3">
                Invite Users
              </p>
              <div className="flex flex-col gap-1">
                {nonMembers.slice(0, 10).map(u => (
                  <button key={u.id} onClick={() => invite(u.id)}
                    className="flex items-center gap-3 p-2 rounded-lg hover:bg-primary/5 transition-colors text-left w-full">
                    <div className="w-8 h-8 rounded-full bg-surface-container-high flex items-center justify-center text-on-surface-variant text-xs font-bold">
                      {u.full_name.charAt(0)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-on-surface truncate">{u.full_name}</p>
                      <p className="text-[10px] text-on-surface-variant">{u.email}</p>
                    </div>
                    <span className="text-xs font-bold text-primary">+ Invite</span>
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
