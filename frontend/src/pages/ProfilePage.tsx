import { useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import api from '../services/api'
import type { User, Position } from '../types'
import type { RootState } from '../store'
import { orgService } from '../services/org.service'

const POSITION_LABELS: Record<Position, string> = {
  division_manager: 'Division Manager', manager: 'Manager', leader: 'Leader', staff: 'Staff',
}

export function ProfilePage() {
  const userId = useSelector((s: RootState) => s.auth.userId)
  const [user, setUser] = useState<User | null>(null)
  const [deptName, setDeptName] = useState('')
  const [divName, setDivName] = useState('')
  const [teamName, setTeamName] = useState('')
  const [loading, setLoading] = useState(true)
  const [ticketCount, setTicketCount] = useState(0)

  useEffect(() => {
    if (!userId) return
    api.get<User[]>('/users/list')
      .then(res => {
        const me = (res.data ?? []).find(u => u.id === userId)
        if (me) setUser(me)
      })
      .finally(() => setLoading(false))
    // Get ticket count
    api.get('/dashboard').then(res => {
      setTicketCount((res.data as any)?.total_tickets ?? 0)
    }).catch(() => {})
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
      <span className="material-symbols-outlined animate-spin text-primary text-4xl">refresh</span>
    </div>
  )

  if (!user) return (
    <div className="max-w-2xl mx-auto p-8 text-center">
      <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 block mb-3">person_off</span>
      <p className="text-on-surface-variant">Profile not found</p>
    </div>
  )

  const initial = user.full_name.charAt(0).toUpperCase()
  const memberSince = new Date(user.created_at)
  const daysSince = Math.floor((Date.now() - memberSince.getTime()) / (1000 * 60 * 60 * 24))

  return (
    <div className="max-w-5xl mx-auto p-6 md:p-8">
      {/* Profile Header Banner */}
      <div className="relative bg-gradient-to-r from-primary to-primary-container rounded-2xl overflow-hidden mb-8 shadow-lg">
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-[-20%] right-[-10%] w-[50%] h-[150%] bg-white/20 rounded-full blur-3xl" />
          <div className="absolute bottom-[-30%] left-[-5%] w-[40%] h-[120%] bg-white/10 rounded-full blur-3xl" />
        </div>
        <div className="relative px-8 py-10 flex flex-col md:flex-row items-center gap-6">
          {/* Avatar */}
          <div className="relative group">
            <div className="w-28 h-28 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center text-white text-5xl font-black border-4 border-white/30 shadow-xl">
              {initial}
            </div>
            <div className="absolute -bottom-1 -right-1 w-8 h-8 bg-emerald-500 rounded-full border-3 border-white flex items-center justify-center">
              <span className="material-symbols-outlined text-white text-[14px]">
                {user.is_active ? 'check' : 'close'}
              </span>
            </div>
          </div>

          {/* Name & Quick Info */}
          <div className="text-center md:text-left">
            <h1 className="text-3xl font-extrabold text-white tracking-tight">{user.full_name}</h1>
            <p className="text-white/70 text-sm mt-1 flex items-center justify-center md:justify-start gap-2">
              <span className="material-symbols-outlined text-[16px]">mail</span>
              {user.email}
            </p>
            <div className="flex flex-wrap items-center justify-center md:justify-start gap-2 mt-3">
              <span className={`text-xs px-3 py-1 rounded-full font-bold capitalize ${
                user.role === 'admin' ? 'bg-white/25 text-white' :
                user.role === 'approver' ? 'bg-amber-400/30 text-amber-100' :
                'bg-white/15 text-white/80'
              }`}>
                <span className="material-symbols-outlined text-[12px] mr-1">shield</span>
                {user.role}
              </span>
              {user.position && (
                <span className="text-xs px-3 py-1 rounded-full font-bold bg-white/15 text-white/90">
                  <span className="material-symbols-outlined text-[12px] mr-1">badge</span>
                  {POSITION_LABELS[user.position]}
                </span>
              )}
              <span className="text-xs px-3 py-1 rounded-full font-bold bg-white/15 text-white/80">
                <span className="material-symbols-outlined text-[12px] mr-1">calendar_today</span>
                Joined {memberSince.toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        {[
          { icon: 'confirmation_number', label: 'My Tickets', value: ticketCount, color: 'text-primary bg-primary-fixed' },
          { icon: 'schedule', label: 'Days Active', value: daysSince, color: 'text-tertiary bg-tertiary-fixed' },
          { icon: 'verified', label: 'Status', value: user.is_active ? 'Active' : 'Inactive', color: user.is_active ? 'text-emerald-700 bg-emerald-100' : 'text-red-700 bg-red-100' },
          { icon: 'shield', label: 'Access Level', value: user.role, color: 'text-amber-700 bg-amber-100' },
        ].map(({ icon, label, value, color }) => (
          <div key={label} className="bg-surface-container-lowest rounded-xl p-4 shadow-sm border border-outline-variant/10">
            <div className={`w-10 h-10 rounded-lg ${color} flex items-center justify-center mb-3`}>
              <span className="material-symbols-outlined text-[20px]">{icon}</span>
            </div>
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider">{label}</p>
            <p className="text-lg font-bold text-on-surface capitalize mt-0.5">{value}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Account Information */}
        <div className="lg:col-span-6">
          <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 overflow-hidden">
            <div className="px-6 py-4 border-b border-outline-variant/10 flex items-center gap-2">
              <span className="material-symbols-outlined text-primary text-[20px]">person</span>
              <h3 className="text-sm font-bold text-on-surface uppercase tracking-widest">Account Information</h3>
            </div>
            <div className="p-6 space-y-4">
              {[
                { icon: 'badge', label: 'Full Name', value: user.full_name },
                { icon: 'mail', label: 'Email Address', value: user.email },
                { icon: 'shield', label: 'Role', value: user.role },
                { icon: 'toggle_on', label: 'Account Status', value: user.is_active ? 'Active' : 'Inactive' },
                { icon: 'calendar_today', label: 'Member Since', value: memberSince.toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) },
                { icon: 'update', label: 'Last Updated', value: new Date(user.updated_at).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) },
              ].map(({ icon, label, value }) => (
                <div key={label} className="flex items-start gap-3">
                  <span className="material-symbols-outlined text-on-surface-variant/50 text-[18px] mt-0.5">{icon}</span>
                  <div>
                    <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider">{label}</p>
                    <p className="text-sm font-semibold text-on-surface capitalize">{value}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Organization */}
        <div className="lg:col-span-6">
          <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 overflow-hidden">
            <div className="px-6 py-4 border-b border-outline-variant/10 flex items-center gap-2">
              <span className="material-symbols-outlined text-primary text-[20px]">account_tree</span>
              <h3 className="text-sm font-bold text-on-surface uppercase tracking-widest">Organization</h3>
            </div>
            <div className="p-6 space-y-4">
              {[
                { icon: 'domain', label: 'Department', value: deptName || '—' },
                { icon: 'workspaces', label: 'Division', value: divName || '—' },
                { icon: 'groups', label: 'Team', value: teamName || '—' },
                { icon: 'badge', label: 'Position', value: user.position ? POSITION_LABELS[user.position] : '—' },
              ].map(({ icon, label, value }) => (
                <div key={label} className="flex items-start gap-3">
                  <span className="material-symbols-outlined text-on-surface-variant/50 text-[18px] mt-0.5">{icon}</span>
                  <div>
                    <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider">{label}</p>
                    <p className="text-sm font-semibold text-on-surface">{value}</p>
                  </div>
                </div>
              ))}
            </div>

            {/* Org Visual */}
            {(deptName || divName || teamName) && (
              <div className="px-6 pb-6">
                <div className="bg-surface-container-low rounded-xl p-4">
                  <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-3">Organization Path</p>
                  <div className="flex items-center gap-2 text-sm flex-wrap">
                    {deptName && (
                      <>
                        <span className="px-2.5 py-1 bg-primary-fixed text-primary rounded-lg font-semibold text-xs">{deptName}</span>
                      </>
                    )}
                    {divName && (
                      <>
                        <span className="material-symbols-outlined text-on-surface-variant/40 text-[16px]">chevron_right</span>
                        <span className="px-2.5 py-1 bg-tertiary-fixed text-on-tertiary-fixed rounded-lg font-semibold text-xs">{divName}</span>
                      </>
                    )}
                    {teamName && (
                      <>
                        <span className="material-symbols-outlined text-on-surface-variant/40 text-[16px]">chevron_right</span>
                        <span className="px-2.5 py-1 bg-secondary-container text-on-secondary-container rounded-lg font-semibold text-xs">{teamName}</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
