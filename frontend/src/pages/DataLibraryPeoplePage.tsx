import { useSearchParams } from 'react-router-dom'

// Lazy tab content imports
import { OrgChartPage } from './OrgChartPage'
import { UserManagementPage } from './UserManagementPage'

const tabs = [
  { id: 'overview', label: 'Overview', icon: 'dashboard' },
  { id: 'staff', label: 'Staff', icon: 'badge' },
  { id: 'monitoring', label: 'Monitoring', icon: 'monitoring' },
  { id: 'roles', label: 'Roles & Responsibilities', icon: 'assignment_ind' },
  { id: 'org-chart', label: 'Org Chart', icon: 'account_tree' },
  { id: 'configuration', label: 'Configuration', icon: 'settings' },
]

function OverviewTab() {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-primary">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Total Staff</p>
          <p className="text-2xl font-black text-on-surface">—</p>
        </div>
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-emerald-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Active</p>
          <p className="text-2xl font-black text-on-surface">—</p>
        </div>
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-amber-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Departments</p>
          <p className="text-2xl font-black text-on-surface">—</p>
        </div>
        <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-purple-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Teams</p>
          <p className="text-2xl font-black text-on-surface">—</p>
        </div>
      </div>
      <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
        <h3 className="text-lg font-bold text-on-surface mb-2">People Overview</h3>
        <p className="text-sm text-on-surface-variant">
          Manage your organization's people, roles, and structure. Use the tabs above to navigate between different views.
        </p>
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
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">settings</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">Configuration</h3>
        <p className="text-sm text-on-surface-variant">Configure people management settings, sync options, and integrations.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">Coming soon</p>
      </div>
    </div>
  )
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
        {activeTab === 'staff' && <UserManagementPage />}
        {activeTab === 'monitoring' && <MonitoringTab />}
        {activeTab === 'roles' && <RolesTab />}
        {activeTab === 'org-chart' && <OrgChartPage />}
        {activeTab === 'configuration' && <ConfigurationTab />}
      </div>
    </div>
  )
}
