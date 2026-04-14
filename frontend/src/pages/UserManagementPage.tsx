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
  const [showForm, setShowForm] = useState(false)
  const [newUser, setNewUser] = useState({
    full_name: '', email: '', password: '', role: 'user' as Role,
    department_id: '', division_id: '', team_id: '', position: '' as string,
  })

  // Org data for dropdowns
  const [departments, setDepartments] = useState<Department[]>([])
  const [divisions, setDivisions] = useState<Division[]>([])
  const [teams, setTeams] = useState<Team[]>([])

  const fetchUsers = () => {
    api.get<User[]>('/users').then((r) => setUsers(r.data)).catch(() => setError('Failed to load users')).finally(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [])

  // Load departments when form opens
  useEffect(() => {
    if (showForm) orgService.listDepartments().then(r => setDepartments(r.data ?? []))
  }, [showForm])

  // Cascade: load divisions when department changes
  useEffect(() => {
    if (newUser.department_id) {
      orgService.listDivisions(newUser.department_id).then(r => setDivisions(r.data ?? []))
    } else {
      setDivisions([])
    }
    setNewUser(u => ({ ...u, division_id: '', team_id: '' }))
  }, [newUser.department_id])

  // Cascade: load teams when division changes
  useEffect(() => {
    if (newUser.division_id) {
      orgService.listTeams(newUser.division_id).then(r => setTeams(r.data ?? []))
    } else {
      setTeams([])
    }
    setNewUser(u => ({ ...u, team_id: '' }))
  }, [newUser.division_id])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const payload: Record<string, string> = {
        full_name: newUser.full_name, email: newUser.email, password: newUser.password, role: newUser.role,
      }
      if (newUser.department_id) payload.department_id = newUser.department_id
      if (newUser.division_id) payload.division_id = newUser.division_id
      if (newUser.team_id) payload.team_id = newUser.team_id
      if (newUser.position) payload.position = newUser.position
      await api.post('/users', payload)
      setShowForm(false)
      setNewUser({ full_name: '', email: '', password: '', role: 'user', department_id: '', division_id: '', team_id: '', position: '' })
      fetchUsers()
    } catch { setError('Failed to create user') }
  }

  const handleDeactivate = async (id: string) => {
    try { await api.patch(`/users/${id}/deactivate`); fetchUsers() }
    catch { setError('Failed to deactivate user') }
  }

  const handleActivate = async (id: string) => {
    try { await api.patch(`/users/${id}/activate`); fetchUsers() }
    catch { setError('Failed to activate user') }
  }

  const inputStyle = { width: '100%', padding: '8px 12px', border: '1px solid #cbd5e1', borderRadius: 4, boxSizing: 'border-box' as const }

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h2 style={{ fontSize: 22, fontWeight: 700 }}>User Management</h2>
        <button onClick={() => setShowForm(!showForm)} style={{ background: '#1e40af', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: 4, cursor: 'pointer' }}>
          + New User
        </button>
      </div>

      {error && <ErrorMessage message={error} />}

      {showForm && (
        <form onSubmit={handleCreate} style={{ background: '#fff', padding: '1.5rem', borderRadius: 8, marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
          <input placeholder="Full Name" value={newUser.full_name} onChange={(e) => setNewUser({ ...newUser, full_name: e.target.value })} required style={inputStyle} />
          <input placeholder="Email" type="email" value={newUser.email} onChange={(e) => setNewUser({ ...newUser, email: e.target.value })} required style={inputStyle} />
          <input placeholder="Password" type="password" value={newUser.password} onChange={(e) => setNewUser({ ...newUser, password: e.target.value })} required style={inputStyle} />
          <select value={newUser.role} onChange={(e) => setNewUser({ ...newUser, role: e.target.value as Role })} style={inputStyle}>
            <option value="user">User</option>
            <option value="approver">Approver</option>
            <option value="admin">Admin</option>
          </select>

          {/* Org assignment fields */}
          <select value={newUser.department_id} onChange={(e) => setNewUser({ ...newUser, department_id: e.target.value })} style={inputStyle}>
            <option value="">-- Department (optional) --</option>
            {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
          </select>
          <select value={newUser.division_id} onChange={(e) => setNewUser({ ...newUser, division_id: e.target.value })} style={inputStyle} disabled={!newUser.department_id}>
            <option value="">-- Division (optional) --</option>
            {divisions.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
          </select>
          <select value={newUser.team_id} onChange={(e) => setNewUser({ ...newUser, team_id: e.target.value })} style={inputStyle} disabled={!newUser.division_id}>
            <option value="">-- Team (optional) --</option>
            {teams.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
          </select>
          <select value={newUser.position} onChange={(e) => setNewUser({ ...newUser, position: e.target.value })} style={inputStyle}>
            <option value="">-- Position (optional) --</option>
            {(Object.keys(POSITION_LABELS) as Position[]).map(p => <option key={p} value={p}>{POSITION_LABELS[p]}</option>)}
          </select>

          <div style={{ display: 'flex', gap: 8 }}>
            <button type="submit" style={{ background: '#1e40af', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: 4, cursor: 'pointer' }}>Create</button>
            <button type="button" onClick={() => setShowForm(false)} style={{ background: '#e2e8f0', border: 'none', padding: '8px 16px', borderRadius: 4, cursor: 'pointer' }}>Cancel</button>
          </div>
        </form>
      )}

      <table style={{ width: '100%', borderCollapse: 'collapse', background: '#fff', borderRadius: 8 }}>
        <thead><tr style={{ background: '#f1f5f9' }}>
          <th style={{ padding: '10px 16px', textAlign: 'left' }}>Name</th>
          <th style={{ padding: '10px 16px', textAlign: 'left' }}>Email</th>
          <th style={{ padding: '10px 16px', textAlign: 'left' }}>Role</th>
          <th style={{ padding: '10px 16px', textAlign: 'left' }}>Position</th>
          <th style={{ padding: '10px 16px', textAlign: 'left' }}>Status</th>
          <th style={{ padding: '10px 16px', textAlign: 'left' }}>Actions</th>
        </tr></thead>
        <tbody>{users.map((u) => (
          <tr key={u.id} style={{ borderTop: '1px solid #e2e8f0' }}>
            <td style={{ padding: '10px 16px' }}>{u.full_name}</td>
            <td style={{ padding: '10px 16px' }}>{u.email}</td>
            <td style={{ padding: '10px 16px' }}>{u.role}</td>
            <td style={{ padding: '10px 16px' }}>{u.position ? POSITION_LABELS[u.position] : '—'}</td>
            <td style={{ padding: '10px 16px' }}>
              <span style={{ color: u.is_active ? '#10b981' : '#ef4444' }}>{u.is_active ? 'Active' : 'Inactive'}</span>
            </td>
            <td style={{ padding: '10px 16px' }}>
              {u.is_active ? (
                <button onClick={() => handleDeactivate(u.id)} style={{ background: '#ef4444', color: '#fff', border: 'none', padding: '4px 10px', borderRadius: 4, cursor: 'pointer', fontSize: 12 }}>
                  Deactivate
                </button>
              ) : (
                <button onClick={() => handleActivate(u.id)} style={{ background: '#10b981', color: '#fff', border: 'none', padding: '4px 10px', borderRadius: 4, cursor: 'pointer', fontSize: 12 }}>
                  Activate
                </button>
              )}
            </td>
          </tr>
        ))}</tbody>
      </table>
    </div>
  )
}
