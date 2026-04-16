import { useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import api from '../services/api'
import type { User, Department, Division, Team, Position } from '../types'
import type { RootState } from '../store'
import { orgService } from '../services/org.service'

const POSITION_LABELS: Record<Position, string> = {
  division_manager: 'Division Manager', manager: 'Manager', leader: 'Leader', staff: 'Staff',
}

export function ProfilePage() {
  const userId = useSelector((s: RootState) => s.auth.userId)
  const role = useSelector((s: RootState) => s.auth.role)
  const [user, setUser] = useState<User | null>(null)
  const [deptName, setDeptName] = useState('')
  const [divName, setDivName] = useState('')
  const [teamName, setTeamName] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!userId) return
    // Get user list and find current user
    api.get<User[]>('/users/list')
      .then(res => {
        const me = (res.data ?? []).find(u => u.id === userId)
        if (me) setUser(me)
      })
      .finally(() => setLoading(false))
  }, [userId])

  useEffect(() => {
    if (!user) return
    if (user.department_id) {
      orgService.listDepartments().then(r => {
        const d = (r.data ?? []).find(d => d.id === user.department_id)
        if (d) setDeptName(d.name)
      })
    }
    if (user.division_id) {
      orgService.listDivisions().then(r => {
        const d = (r.data ?? []).find(d => d.id === user.division_id)
        if (d) setDivName(d.name)
      })
    }
    if (user.team_id) {
      orgService.listTeams().then(r => {
        const t = (r.data ?? []).find(t => t.id === user.team_id)
        if (t) setTeamName(t.name)
      })
    }
  }, [user])

  if (loading) return (
    <div className="flex items-center justify-center py-20">
      <span className="material-symbols-outlined animate-spin text-primary">refresh</span>
    </div>
  )

  if (!user) return (
    <div className="max-w-2xl mx-auto p-8 text-center">
      <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 block mb-3">person_off</span>
      <p className="text-on-surface-variant">Profile not found</p>
    </div>
  )

  return (
    <div className="max-w-4xl mx-auto p-8">
      <h1 className="text-3xl font-extrabold tracking-tight text-on-surface mb-8">My Profile</h1>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Profile Card */}
        <div className="lg:col-span-4">
          <div className="bg-surface-container-lowest rounded-xl p-8 shadow-sm text-center">
            <div className="w-24 h-24 rounded-full bg-gradient-to-br from-primary to-primary-container flex items-center justify-center text-white text-3xl font-black mx-auto mb-4">
              {user.full_name.charAt(0).toUpperCase()}
            </div>
            <h2 className="text-xl font-bold text-on-surface">{user.full_name}</h2>
            <p className="text-on-surface-variant text-sm mt-1">{user.email}</p>
            <div className="mt-4 flex justify-center gap-2">
              <span className={`text-xs px-3 py-1 rounded-full font-bold capitalize ${
                user.role === 'admin' ? 'bg-primary-fixed text-primary' :
                user.role === 'approver' ? 'bg-amber-100 text-amber-700' :
                'bg-surface-container-high text-on-surface-variant'
              }`}>{user.role}</span>
              {user.position && (
                <span className="text-xs px-3 py-1 rounded-full font-bold bg-tertiary-fixed text-on-tertiary-fixed">
                  {POSITION_LABELS[user.position]}
                </span>
              )}
            </div>
            <div className="mt-4 flex items-center justify-center gap-2">
              <span className={`w-2 h-2 rounded-full ${user.is_active ? 'bg-emerald-500' : 'bg-red-500'}`} />
              <span className="text-xs text-on-surface-variant">{user.is_active ? 'Active' : 'Inactive'}</span>
            </div>
          </div>
        </div>

        {/* Details */}
        <div className="lg:col-span-8 space-y-6">
          {/* Account Info */}
          <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
            <h3 className="text-sm font-bold text-on-surface-variant uppercase tracking-widest mb-4">Account Information</h3>
            <div className="grid grid-cols-2 gap-4">
              {[
                { label: 'Full Name', value: user.full_name },
                { label: 'Email', value: user.email },
                { label: 'Role', value: user.role },
                { label: 'Status', value: user.is_active ? 'Active' : 'Inactive' },
                { label: 'Member Since', value: new Date(user.created_at).toLocaleDateString() },
                { label: 'Last Updated', value: new Date(user.updated_at).toLocaleDateString() },
              ].map(({ label, value }) => (
                <div key={label}>
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">{label}</p>
                  <p className="text-sm font-semibold text-on-surface capitalize">{value}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Organization */}
          <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
            <h3 className="text-sm font-bold text-on-surface-variant uppercase tracking-widest mb-4">Organization</h3>
            <div className="grid grid-cols-2 gap-4">
              {[
                { label: 'Department', value: deptName || '—' },
                { label: 'Division', value: divName || '—' },
                { label: 'Team', value: teamName || '—' },
                { label: 'Position', value: user.position ? POSITION_LABELS[user.position] : '—' },
              ].map(({ label, value }) => (
                <div key={label}>
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">{label}</p>
                  <p className="text-sm font-semibold text-on-surface">{value}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
