import { useEffect, useState } from 'react'
import api from '../../services/api'
import { orgService } from '../../services/org.service'
import type { User, Role, Department, Division, Team, Position } from '../../types'

const POSITION_LABELS: Record<string, string> = {
  manager: 'Manager',
  leader: 'Leader',
  staff: 'Staff',
}

interface StaffDetailPanelProps {
  user: User
  users: User[]
  onClose: () => void
  onUpdated: () => void
}

function StaffDetailPanel({ user, users, onClose, onUpdated }: StaffDetailPanelProps) {
  const [editing, setEditing] = useState<string | null>(null)
  const [divisions, setDivisions] = useState<Division[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [teams, setTeams] = useState<Team[]>([])
  const [form, setForm] = useState({
    role: user.role,
    division_id: user.division_id ?? '',
    department_id: user.department_id ?? '',
    team_id: user.team_id ?? '',
    position: user.position ?? '',
    reports_to: user.reports_to ?? '',
  })
  const [saving, setSaving] = useState(false)
  const [msg, setMsg] = useState('')

  useEffect(() => {
    orgService.listDivisions().then(r => setDivisions(r.data ?? []))
  }, [])

  useEffect(() => {
    if (form.division_id) {
      orgService.listDepartments(form.division_id).then(r => setDepartments(r.data ?? []))
    } else setDepartments([])
  }, [form.division_id])

  useEffect(() => {
    if (form.department_id) {
      orgService.listTeams(form.department_id).then(r => setTeams(r.data ?? []))
    } else setTeams([])
  }, [form.department_id])

  const reportingManager = users.find(u => u.id === user.reports_to)
  const division = divisions.find(d => d.id === (form.division_id || user.division_id))
  const department = departments.find(d => d.id === (form.department_id || user.department_id))
  const team = teams.find(t => t.id === (form.team_id || user.team_id))

  const saveField = async (field: string) => {
    setSaving(true)
    try {
      if (field === 'role') {
        await api.patch(`/users/${user.id}/role`, { role: form.role })
      } else {
        const orgPayload: Record<string, string | null> = {
          division_id: form.division_id || null,
          department_id: form.department_id || null,
          team_id: form.team_id || null,
          position: form.position || null,
          reports_to: form.reports_to || null,
        }
        await api.patch(`/users/${user.id}/org`, orgPayload)
      }
      setMsg('Saved')
      setTimeout(() => setMsg(''), 2000)
      setEditing(null)
      onUpdated()
    } catch {
      setMsg('Failed to save')
      setTimeout(() => setMsg(''), 2000)
    } finally {
      setSaving(false)
    }
  }

  const initial = user.full_name.charAt(0).toUpperCase()

  return (
    <div className="fixed inset-y-0 right-0 w-full max-w-md bg-surface-container-lowest shadow-2xl z-50 flex flex-col overflow-hidden border-l border-outline-variant/20">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-4 border-b border-outline-variant/20 bg-surface-container-low">
        <h2 className="text-sm font-bold text-on-surface uppercase tracking-wider">Staff Details</h2>
        <button onClick={onClose} className="p-1.5 hover:bg-surface-container-high rounded-lg transition-colors">
          <span className="material-symbols-outlined text-on-surface-variant">close</span>
        </button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* User Info */}
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 rounded-full bg-gradient-to-br from-primary to-primary-container flex items-center justify-center text-white text-xl font-bold">
            {initial}
          </div>
          <div>
            <h3 className="text-lg font-bold text-on-surface">{user.full_name}</h3>
            <p className="text-sm text-on-surface-variant">{user.position ? POSITION_LABELS[user.position] || user.position : 'Employee'}</p>
            <div className="flex items-center gap-2 mt-1 text-xs text-on-surface-variant">
              <span className="material-symbols-outlined text-[14px]">mail</span>
              {user.email}
            </div>
            <div className="flex items-center gap-2 mt-0.5 text-xs text-on-surface-variant">
              <span className="material-symbols-outlined text-[14px]">event</span>
              Joined {new Date(user.created_at).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })}
            </div>
          </div>
        </div>

        {/* Status */}
        <div className="flex items-center gap-2">
          <span className={`text-xs font-bold px-3 py-1 rounded-full ${user.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'}`}>
            {user.is_active ? '● Active' : '● Inactive'}
          </span>
          <span className={`text-xs font-bold px-3 py-1 rounded-full capitalize ${
            user.role === 'admin' ? 'bg-purple-100 text-purple-700' :
            user.role === 'approver' ? 'bg-amber-100 text-amber-700' :
            user.role === 'agent' ? 'bg-blue-100 text-blue-700' :
            'bg-slate-100 text-slate-600'
          }`}>{user.role}</span>
        </div>

        {msg && (
          <div className="px-3 py-2 rounded-lg bg-emerald-100 text-emerald-700 text-xs font-bold">
            {msg}
          </div>
        )}

        {/* Organization Section */}
        <div className="space-y-1">
          <h4 className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-3">Organization</h4>

          {/* Role */}
          <DetailRow
            label="Role"
            value={<span className="capitalize">{form.role}</span>}
            editing={editing === 'role'}
            onEdit={() => setEditing('role')}
            onCancel={() => setEditing(null)}
            onSave={() => saveField('role')}
            saving={saving}
            editContent={
              <select value={form.role} onChange={e => setForm({ ...form, role: e.target.value as Role })}
                className="w-full text-sm border border-outline-variant rounded-lg px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30">
                <option value="user">User</option>
                <option value="agent">Agent</option>
                <option value="approver">Approver</option>
                <option value="admin">Admin</option>
              </select>
            }
          />

          {/* Position */}
          <DetailRow
            label="Position"
            value={form.position ? POSITION_LABELS[form.position] || form.position : '—'}
            editing={editing === 'position'}
            onEdit={() => setEditing('position')}
            onCancel={() => setEditing(null)}
            onSave={() => saveField('position')}
            saving={saving}
            editContent={
              <select value={form.position} onChange={e => setForm({ ...form, position: e.target.value })}
                className="w-full text-sm border border-outline-variant rounded-lg px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30">
                <option value="">— None —</option>
                <option value="manager">Manager</option>
                <option value="leader">Leader</option>
                <option value="staff">Staff</option>
              </select>
            }
          />

          {/* Division */}
          <DetailRow
            label="Division"
            value={division?.name || '—'}
            editing={editing === 'division'}
            onEdit={() => setEditing('division')}
            onCancel={() => setEditing(null)}
            onSave={() => saveField('division')}
            saving={saving}
            editContent={
              <select value={form.division_id} onChange={e => setForm({ ...form, division_id: e.target.value, department_id: '', team_id: '' })}
                className="w-full text-sm border border-outline-variant rounded-lg px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30">
                <option value="">— None —</option>
                {divisions.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
              </select>
            }
          />

          {/* Department */}
          <DetailRow
            label="Department"
            value={department?.name || '—'}
            editing={editing === 'department'}
            onEdit={() => setEditing('department')}
            onCancel={() => setEditing(null)}
            onSave={() => saveField('department')}
            saving={saving}
            editContent={
              <select value={form.department_id} onChange={e => setForm({ ...form, department_id: e.target.value, team_id: '' })}
                className="w-full text-sm border border-outline-variant rounded-lg px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30">
                <option value="">— None —</option>
                {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
              </select>
            }
          />

          {/* Team */}
          <DetailRow
            label="Team"
            value={team?.name || '—'}
            editing={editing === 'team'}
            onEdit={() => setEditing('team')}
            onCancel={() => setEditing(null)}
            onSave={() => saveField('team')}
            saving={saving}
            editContent={
              <select value={form.team_id} onChange={e => setForm({ ...form, team_id: e.target.value })}
                className="w-full text-sm border border-outline-variant rounded-lg px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30">
                <option value="">— None —</option>
                {teams.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
              </select>
            }
          />

          {/* Reporting Manager */}
          <DetailRow
            label="Reporting Manager"
            value={
              reportingManager ? (
                <div className="flex items-center gap-2">
                  <div className="w-6 h-6 rounded-full bg-primary-fixed flex items-center justify-center text-primary text-[10px] font-bold">
                    {reportingManager.full_name.charAt(0)}
                  </div>
                  <div>
                    <span className="text-sm text-on-surface">{reportingManager.full_name}</span>
                    <p className="text-[10px] text-on-surface-variant">{reportingManager.email}</p>
                  </div>
                </div>
              ) : '—'
            }
            editing={editing === 'reports_to'}
            onEdit={() => setEditing('reports_to')}
            onCancel={() => setEditing(null)}
            onSave={() => saveField('reports_to')}
            saving={saving}
            editContent={
              <select value={form.reports_to} onChange={e => setForm({ ...form, reports_to: e.target.value })}
                className="w-full text-sm border border-outline-variant rounded-lg px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30">
                <option value="">— None —</option>
                {users.filter(u => u.id !== user.id && u.is_active).map(u => (
                  <option key={u.id} value={u.id}>{u.full_name}</option>
                ))}
              </select>
            }
          />
        </div>
      </div>
    </div>
  )
}

function DetailRow({ label, value, editing, onEdit, onCancel, onSave, saving, editContent }: {
  label: string
  value: React.ReactNode
  editing: boolean
  onEdit: () => void
  onCancel: () => void
  onSave: () => void
  saving: boolean
  editContent: React.ReactNode
}) {
  return (
    <div className="py-3 border-b border-outline-variant/10">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider">{label}</span>
        {!editing && (
          <button onClick={onEdit} className="text-[11px] font-bold text-primary hover:underline">Edit</button>
        )}
      </div>
      {editing ? (
        <div className="space-y-2">
          {editContent}
          <div className="flex gap-2">
            <button onClick={onSave} disabled={saving}
              className="text-xs font-bold text-white bg-primary px-3 py-1.5 rounded-lg hover:opacity-90 disabled:opacity-50">
              {saving ? 'Saving...' : 'Save'}
            </button>
            <button onClick={onCancel} className="text-xs font-bold text-on-surface-variant hover:text-on-surface px-3 py-1.5">
              Cancel
            </button>
          </div>
        </div>
      ) : (
        <div className="text-sm text-on-surface">{value}</div>
      )}
    </div>
  )
}

export function StaffTable() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive'>('active')
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [sortField, setSortField] = useState<'name' | 'position' | 'role' | 'status'>('name')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('asc')

  const fetchUsers = () => {
    api.get<User[]>('/users').then(r => setUsers(r.data ?? [])).finally(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [])

  const handleSort = (field: 'name' | 'position' | 'role' | 'status') => {
    if (sortField === field) {
      setSortDir(sortDir === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortDir('asc')
    }
  }

  const sortIcon = (field: string) => {
    if (sortField !== field) return '↕'
    return sortDir === 'asc' ? '↑' : '↓'
  }

  const filteredUsers = users
    .filter(u => {
      if (statusFilter === 'active') return u.is_active
      if (statusFilter === 'inactive') return !u.is_active
      return true
    })
    .filter(u => {
      if (!searchQuery) return true
      const q = searchQuery.toLowerCase()
      return u.full_name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q)
    })
    .sort((a, b) => {
      let cmp = 0
      switch (sortField) {
        case 'name':
          cmp = a.full_name.localeCompare(b.full_name)
          break
        case 'position':
          cmp = (a.position || 'zzz').localeCompare(b.position || 'zzz')
          break
        case 'role':
          cmp = a.role.localeCompare(b.role)
          break
        case 'status':
          cmp = (a.is_active === b.is_active) ? 0 : a.is_active ? -1 : 1
          break
      }
      return sortDir === 'asc' ? cmp : -cmp
    })

  const getReportingManager = (userId?: string) => {
    if (!userId) return null
    return users.find(u => u.id === userId)
  }

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin w-8 h-8 border-4 border-primary border-t-transparent rounded-full" />
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1 text-sm">
            <span className="text-on-surface-variant">View:</span>
            <select value={statusFilter} onChange={e => setStatusFilter(e.target.value as any)}
              className="text-sm font-bold text-on-surface bg-transparent border-none outline-none cursor-pointer">
              <option value="active">Active ({users.filter(u => u.is_active).length})</option>
              <option value="inactive">Inactive ({users.filter(u => !u.is_active).length})</option>
              <option value="all">All ({users.length})</option>
            </select>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative">
            <span className="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-[16px] text-on-surface-variant">search</span>
            <input
              type="text"
              placeholder="Search staff..."
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              className="pl-8 pr-3 py-2 text-sm border border-outline-variant/30 rounded-lg bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30 w-56"
            />
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-surface-container-low/50">
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant w-8"></th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant cursor-pointer hover:text-on-surface select-none" onClick={() => handleSort('name')}>Person {sortIcon('name')}</th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant cursor-pointer hover:text-on-surface select-none" onClick={() => handleSort('position')}>Position {sortIcon('position')}</th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Reporting Manager</th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant cursor-pointer hover:text-on-surface select-none" onClick={() => handleSort('status')}>Work Status {sortIcon('status')}</th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant cursor-pointer hover:text-on-surface select-none" onClick={() => handleSort('role')}>Role {sortIcon('role')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-outline-variant/10">
            {filteredUsers.map(u => {
              const mgr = getReportingManager(u.reports_to)
              const initial = u.full_name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()
              return (
                <tr key={u.id} className="hover:bg-surface-container-low/30 transition-colors cursor-pointer"
                  onClick={() => setSelectedUser(u)}>
                  <td className="px-4 py-3">
                    <div className="w-8 h-8 rounded-full bg-surface-container-high flex items-center justify-center text-[10px] font-bold text-on-surface-variant">
                      {initial}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div>
                      <p className="text-sm font-semibold text-on-surface hover:text-primary transition-colors">{u.full_name}</p>
                      <p className="text-[11px] text-on-surface-variant">{u.email}</p>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-xs text-on-surface-variant">
                      {u.position ? POSITION_LABELS[u.position] || u.position : '—'}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    {mgr ? (
                      <div className="flex items-center gap-2">
                        <div className="w-6 h-6 rounded-full bg-primary-fixed flex items-center justify-center text-[9px] font-bold text-primary">
                          {mgr.full_name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()}
                        </div>
                        <div>
                          <p className="text-xs font-medium text-on-surface">{mgr.full_name}</p>
                          <p className="text-[10px] text-on-surface-variant">{mgr.email}</p>
                        </div>
                      </div>
                    ) : (
                      <span className="text-xs text-on-surface-variant">—</span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`text-[10px] font-bold px-2.5 py-1 rounded-full ${u.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'}`}>
                      {u.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`text-[10px] font-bold px-2.5 py-1 rounded-full capitalize ${
                      u.role === 'admin' ? 'bg-purple-100 text-purple-700' :
                      u.role === 'approver' ? 'bg-amber-100 text-amber-700' :
                      u.role === 'agent' ? 'bg-blue-100 text-blue-700' :
                      'bg-slate-100 text-slate-600'
                    }`}>{u.role}</span>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {filteredUsers.length === 0 && (
          <div className="text-center py-12 text-on-surface-variant text-sm">No staff found</div>
        )}
      </div>

      {/* Detail Panel */}
      {selectedUser && (
        <>
          <div className="fixed inset-0 bg-black/30 z-40" onClick={() => setSelectedUser(null)} />
          <StaffDetailPanel
            user={selectedUser}
            users={users}
            onClose={() => setSelectedUser(null)}
            onUpdated={() => { fetchUsers(); }}
          />
        </>
      )}
    </div>
  )
}
