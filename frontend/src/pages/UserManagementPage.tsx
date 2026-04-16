import { useEffect, useState } from 'react'
import api from '../services/api'
import { orgService } from '../services/org.service'
import type { User, Role, Department, Division, Team, Position } from '../types'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import { ErrorMessage } from '../components/common/ErrorMessage'

const POSITION_LABELS: Record<Position, string> = {
  division_manager: 'Division Manager',
  manager: 'Manager',
  leader: 'Leader',
  staff: 'Staff',
}

export function UserManagementPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [actionMsg, setActionMsg] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [form, setForm] = useState({
    full_name: '', email: '', password: '', role: 'user' as Role,
    department_id: '', division_id: '', team_id: '', position: '' as string,
  })

  const [departments, setDepartments] = useState<Department[]>([])
  const [divisions, setDivisions] = useState<Division[]>([])
  const [teams, setTeams] = useState<Team[]>([])

  const showMsg = (msg: string) => { setActionMsg(msg); setTimeout(() => setActionMsg(''), 3000) }

  const fetchUsers = () => {
    api.get<User[]>('/users').then((r) => setUsers(r.data ?? [])).catch(() => setError('Failed to load users')).finally(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [])

  useEffect(() => {
    if (showForm || editingUser) orgService.listDepartments().then(r => setDepartments(r.data ?? []))
  }, [showForm, editingUser])

  useEffect(() => {
    if (form.department_id) {
      orgService.listDivisions(form.department_id).then(r => setDivisions(r.data ?? []))
    } else { setDivisions([]) }
  }, [form.department_id])

  useEffect(() => {
    if (form.division_id) {
      orgService.listTeams(form.division_id).then(r => setTeams(r.data ?? []))
    } else { setTeams([]) }
  }, [form.division_id])

  const resetForm = () => {
    setForm({ full_name: '', email: '', password: '', role: 'user', department_id: '', division_id: '', team_id: '', position: '' })
    setEditingUser(null)
    setShowForm(false)
  }

  const openCreate = () => {
    resetForm()
    setShowForm(true)
  }

  const openEdit = (u: User) => {
    setEditingUser(u)
    setForm({
      full_name: u.full_name,
      email: u.email,
      password: '',
      role: u.role,
      department_id: u.department_id ?? '',
      division_id: u.division_id ?? '',
      team_id: u.team_id ?? '',
      position: u.position ?? '',
    })
    setShowForm(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    try {
      if (editingUser) {
        // Update role
        await api.patch(`/users/${editingUser.id}/role`, { role: form.role })
        // Update org
        const orgPayload: Record<string, string | null> = {}
        orgPayload.department_id = form.department_id || null
        orgPayload.division_id = form.division_id || null
        orgPayload.team_id = form.team_id || null
        orgPayload.position = form.position || null
        await api.patch(`/users/${editingUser.id}/org`, orgPayload)
        showMsg('User updated')
      } else {
        const payload: Record<string, string> = {
          full_name: form.full_name, email: form.email, password: form.password, role: form.role,
        }
        if (form.department_id) payload.department_id = form.department_id
        if (form.division_id) payload.division_id = form.division_id
        if (form.team_id) payload.team_id = form.team_id
        if (form.position) payload.position = form.position
        await api.post('/users', payload)
        showMsg('User created')
      }
      resetForm()
      fetchUsers()
    } catch { setError(editingUser ? 'Failed to update user' : 'Failed to create user') }
  }

  const handleDeactivate = async (id: string) => {
    try { await api.patch(`/users/${id}/deactivate`); fetchUsers(); showMsg('User deactivated') }
    catch { setError('Failed to deactivate user') }
  }

  const handleActivate = async (id: string) => {
    try { await api.patch(`/users/${id}/activate`); fetchUsers(); showMsg('User activated') }
    catch { setError('Failed to activate user') }
  }

  const inputStyle = 'w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-primary'

  if (loading) return <LoadingSpinner />

  return (
    <div className="max-w-7xl mx-auto p-6">
      {actionMsg && (
        <div className="fixed top-20 right-6 z-50 px-5 py-3 rounded-xl shadow-lg font-semibold text-sm bg-emerald-100 text-emerald-800 flex items-center gap-2">
          <span className="material-symbols-outlined text-sm">check_circle</span>
          {actionMsg}
        </div>
      )}

      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-extrabold tracking-tight text-on-surface mb-1">User Management</h1>
          <p className="text-on-surface-variant text-sm">Manage users, roles, positions, and organizational assignments.</p>
        </div>
        <button onClick={openCreate}
          className="bg-gradient-to-r from-primary to-primary-container text-white px-6 py-3 rounded-xl font-bold flex items-center gap-2 shadow-lg shadow-primary/20 hover:opacity-90 transition-all">
          <span className="material-symbols-outlined">person_add</span>
          New User
        </button>
      </div>

      {error && <div className="mb-4"><ErrorMessage message={error} /></div>}

      {/* Create/Edit Form */}
      {showForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg mx-4 overflow-hidden">
            <div className="flex items-center justify-between px-6 py-4 border-b border-outline-variant/20">
              <h3 className="text-lg font-bold text-on-surface">{editingUser ? 'Edit User' : 'Create New User'}</h3>
              <button onClick={resetForm} className="p-1 hover:bg-surface-container-high rounded-lg transition-colors">
                <span className="material-symbols-outlined text-on-surface-variant">close</span>
              </button>
            </div>
            <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4 max-h-[70vh] overflow-y-auto">
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Full Name</label>
                <input value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })}
                  required={!editingUser} disabled={!!editingUser} placeholder="Full Name" className={inputStyle} />
              </div>
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Email</label>
                <input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })}
                  required={!editingUser} disabled={!!editingUser} placeholder="Email" className={inputStyle} />
              </div>
              {!editingUser && (
                <div>
                  <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Password</label>
                  <input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })}
                    required placeholder="Min 8 characters" className={inputStyle} />
                </div>
              )}
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Role</label>
                <select value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value as Role })} className={inputStyle + ' appearance-none'}>
                  <option value="user">User</option>
                  <option value="approver">Approver</option>
                  <option value="admin">Admin</option>
                </select>
              </div>

              <div className="pt-2 border-t border-outline-variant/20">
                <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-3">Organization Assignment</p>
              </div>

              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Department</label>
                <select value={form.department_id} onChange={(e) => setForm({ ...form, department_id: e.target.value, division_id: '', team_id: '' })} className={inputStyle + ' appearance-none'}>
                  <option value="">-- None --</option>
                  {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Division</label>
                <select value={form.division_id} onChange={(e) => setForm({ ...form, division_id: e.target.value, team_id: '' })}
                  disabled={!form.department_id} className={inputStyle + ' appearance-none'}>
                  <option value="">-- None --</option>
                  {divisions.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Team</label>
                <select value={form.team_id} onChange={(e) => setForm({ ...form, team_id: e.target.value })}
                  disabled={!form.division_id} className={inputStyle + ' appearance-none'}>
                  <option value="">-- None --</option>
                  {teams.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Position</label>
                <select value={form.position} onChange={(e) => setForm({ ...form, position: e.target.value })} className={inputStyle + ' appearance-none'}>
                  <option value="">-- None --</option>
                  {(Object.keys(POSITION_LABELS) as Position[]).map(p => <option key={p} value={p}>{POSITION_LABELS[p]}</option>)}
                </select>
              </div>

              <div className="flex justify-end gap-3 pt-4">
                <button type="button" onClick={resetForm}
                  className="px-5 py-2.5 rounded-xl font-semibold text-on-surface-variant hover:bg-surface-container-high transition-colors">Cancel</button>
                <button type="submit"
                  className="px-6 py-2.5 bg-primary text-white rounded-xl font-bold hover:opacity-90 transition-all">
                  {editingUser ? 'Save Changes' : 'Create User'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* User Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 overflow-x-auto">
        <table className="w-full min-w-[900px] text-left border-collapse">
          <thead>
            <tr className="bg-surface-container-low/50">
              {['Name', 'Email', 'Role', 'Position', 'Status', 'Actions'].map(h => (
                <th key={h} className="px-6 py-4 text-[11px] font-black uppercase tracking-widest text-on-surface-variant">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-surface-container-low">
            {users.map((u) => (
              <tr key={u.id} className="hover:bg-surface-container-low/30 transition-colors">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-full bg-primary-fixed flex items-center justify-center text-primary font-bold text-sm">
                      {u.full_name.charAt(0).toUpperCase()}
                    </div>
                    <span className="text-sm font-semibold text-on-surface">{u.full_name}</span>
                  </div>
                </td>
                <td className="px-6 py-4 text-sm text-on-surface-variant">{u.email}</td>
                <td className="px-6 py-4">
                  <span className={`text-xs px-2 py-1 rounded font-bold capitalize ${
                    u.role === 'admin' ? 'bg-primary-fixed text-primary' :
                    u.role === 'approver' ? 'bg-amber-100 text-amber-700' :
                    'bg-surface-container-high text-on-surface-variant'
                  }`}>{u.role}</span>
                </td>
                <td className="px-6 py-4 text-sm text-on-surface-variant">
                  {u.position ? POSITION_LABELS[u.position] : '—'}
                </td>
                <td className="px-6 py-4">
                  <span className={`text-xs font-bold px-2 py-1 rounded ${u.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'}`}>
                    {u.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex gap-2">
                    <button onClick={() => openEdit(u)}
                      className="text-xs font-bold text-primary hover:bg-primary-fixed px-3 py-1.5 rounded-xl transition-colors">
                      Edit
                    </button>
                    {u.is_active ? (
                      <button onClick={() => handleDeactivate(u.id)}
                        className="text-xs font-bold text-error hover:bg-error-container px-3 py-1.5 rounded-xl transition-colors">
                        Deactivate
                      </button>
                    ) : (
                      <button onClick={() => handleActivate(u.id)}
                        className="text-xs font-bold text-emerald-700 hover:bg-emerald-100 px-3 py-1.5 rounded-xl transition-colors">
                        Activate
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
