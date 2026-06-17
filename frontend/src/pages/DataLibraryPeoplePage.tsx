import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import api from '../services/api'
import { orgService } from '../services/org.service'
import type { User, Department, Division, Team } from '../types'

// Tab content imports
import { OrgChartPage } from './OrgChartPage'
import { OrgStructurePage } from './OrgStructurePage'
import { StaffTable } from '../components/people/StaffTable'

const tabs = [
  { id: 'overview', label: 'Overview', icon: 'dashboard' },
  { id: 'staff', label: 'Staff', icon: 'badge' },
  { id: 'monitoring', label: 'Monitoring', icon: 'monitoring' },
  { id: 'roles', label: 'Roles & Responsibilities', icon: 'assignment_ind' },
  { id: 'org-chart', label: 'Org Chart', icon: 'account_tree' },
  { id: 'configuration', label: 'Configuration', icon: 'settings' },
]

function OverviewTab() {
  const [users, setUsers] = useState<User[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [divisions, setDivisions] = useState<Division[]>([])
  const [teams, setTeams] = useState<Team[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get<User[]>('/users'),
      orgService.listDepartments(),
      orgService.listDivisions(),
      orgService.listTeams(),
    ]).then(([usersRes, deptRes, divRes, teamsRes]) => {
      setUsers(usersRes.data ?? [])
      setDepartments(deptRes.data ?? [])
      setDivisions(divRes.data ?? [])
      setTeams(teamsRes.data ?? [])
    }).catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const activeUsers = users.filter(u => u.is_active)
  const inactiveUsers = users.filter(u => !u.is_active)

  const roleBreakdown = users.reduce((acc, u) => {
    acc[u.role] = (acc[u.role] || 0) + 1
    return acc
  }, {} as Record<string, number>)

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin w-8 h-8 border-4 border-primary border-t-transparent rounded-full" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Stat Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-primary">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Total Staff</p>
          <p className="text-2xl font-black text-on-surface">{users.length}</p>
          <p className="text-[10px] text-on-surface-variant mt-1">
            <span className="text-emerald-600 font-bold">{activeUsers.length} active</span>
            {inactiveUsers.length > 0 && <span className="text-red-500 ml-2">{inactiveUsers.length} inactive</span>}
          </p>
        </div>
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-emerald-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Departments</p>
          <p className="text-2xl font-black text-on-surface">{departments.length}</p>
          <p className="text-[10px] text-on-surface-variant mt-1">Across {divisions.length} divisions</p>
        </div>
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-amber-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Teams</p>
          <p className="text-2xl font-black text-on-surface">{teams.length}</p>
          <p className="text-[10px] text-on-surface-variant mt-1">Operational units</p>
        </div>
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-purple-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Divisions</p>
          <p className="text-2xl font-black text-on-surface">{divisions.length}</p>
          <p className="text-[10px] text-on-surface-variant mt-1">Business units</p>
        </div>
      </div>

      {/* Role Breakdown & Recent Staff */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Role Breakdown */}
        <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
          <h3 className="text-sm font-bold text-on-surface mb-4 uppercase tracking-wider">Role Distribution</h3>
          <div className="space-y-3">
            {Object.entries(roleBreakdown).map(([role, count]) => {
              const percentage = Math.round((count / users.length) * 100)
              const colorMap: Record<string, string> = {
                admin: 'bg-purple-500',
                approver: 'bg-blue-500',
                agent: 'bg-emerald-500',
                user: 'bg-slate-400',
              }
              return (
                <div key={role} className="flex items-center gap-3">
                  <span className="text-xs font-bold text-on-surface-variant capitalize w-20">{role}</span>
                  <div className="flex-1 h-2 bg-slate-200 rounded-full overflow-hidden">
                    <div className={`h-full rounded-full ${colorMap[role] || 'bg-slate-400'}`} style={{ width: `${percentage}%` }} />
                  </div>
                  <span className="text-xs font-bold text-on-surface w-8 text-right">{count}</span>
                </div>
              )
            })}
          </div>
        </div>

        {/* Departments List */}
        <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
          <h3 className="text-sm font-bold text-on-surface mb-4 uppercase tracking-wider">Departments</h3>
          <div className="space-y-2 max-h-48 overflow-y-auto">
            {departments.length === 0 ? (
              <p className="text-sm text-on-surface-variant">No departments configured</p>
            ) : (
              departments.map(dept => {
                const deptUsers = users.filter(u => u.department_id === dept.id)
                return (
                  <div key={dept.id} className="flex items-center justify-between py-2 border-b border-outline-variant/10 last:border-0">
                    <div>
                      <p className="text-sm font-medium text-on-surface">{dept.name}</p>
                      <p className="text-[10px] text-on-surface-variant">{dept.code}</p>
                    </div>
                    <span className="text-xs font-bold text-on-surface-variant bg-surface-container-high px-2 py-0.5 rounded-full">
                      {deptUsers.length} staff
                    </span>
                  </div>
                )
              })
            )}
          </div>
        </div>
      </div>

      {/* Divisions overview */}
      <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
        <h3 className="text-sm font-bold text-on-surface mb-4 uppercase tracking-wider">Organization Structure</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {divisions.map(div => {
            const divDepts = departments.filter(d => d.division_id === div.id)
            const divTeams = teams.filter(t => divDepts.some(d => d.id === t.department_id))
            return (
              <div key={div.id} className="p-4 rounded-lg border border-outline-variant/20 hover:shadow-sm transition-all">
                <p className="text-sm font-bold text-on-surface">{div.name}</p>
                <p className="text-[10px] text-on-surface-variant mb-2">{div.code}</p>
                <div className="flex gap-3 text-[10px]">
                  <span className="text-on-surface-variant"><strong className="text-on-surface">{divDepts.length}</strong> depts</span>
                  <span className="text-on-surface-variant"><strong className="text-on-surface">{divTeams.length}</strong> teams</span>
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}

function MonitoringTab() {
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">monitoring</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">Staff Monitoring</h3>
        <p className="text-sm text-on-surface-variant">Monitor staff activity, login sessions, and compliance status.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">Coming soon</p>
      </div>
    </div>
  )
}

function RolesTab() {
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">assignment_ind</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">Roles & Responsibilities</h3>
        <p className="text-sm text-on-surface-variant">Define and manage roles, permissions, and responsibilities across your organization.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">Coming soon</p>
      </div>
    </div>
  )
}

function ConfigurationTab() {
  return <OrgStructurePage />
}

export function DataLibraryPeoplePage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'overview'

  const setTab = (tabId: string) => {
    setSearchParams({ tab: tabId })
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-black text-on-surface font-headline">People</h1>
        <p className="text-sm text-on-surface-variant mt-1">
          Manage staff, organizational structure, roles, and responsibilities.
        </p>
      </div>

      {/* Tabs */}
      <div className="border-b border-outline-variant/20">
        <nav className="flex gap-1 overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setTab(tab.id)}
              className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium whitespace-nowrap border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-primary text-primary'
                  : 'border-transparent text-on-surface-variant hover:text-on-surface hover:border-outline-variant/30'
              }`}
            >
              <span className="material-symbols-outlined text-[18px]">{tab.icon}</span>
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div>
        {activeTab === 'overview' && <OverviewTab />}
        {activeTab === 'staff' && <StaffTable />}
        {activeTab === 'monitoring' && <MonitoringTab />}
        {activeTab === 'roles' && <RolesTab />}
        {activeTab === 'org-chart' && <OrgChartPage />}
        {activeTab === 'configuration' && <ConfigurationTab />}
      </div>
    </div>
  )
}
