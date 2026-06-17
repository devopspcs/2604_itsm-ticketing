import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import api from '../services/api'

interface Device {
  id: number
  name: string
  os: string
  serial_number?: string
  is_online?: boolean
  brand?: string
  model?: string
  created_at: string
  updated_at: string
}

interface DeviceStats {
  total_devices: number
  laptop_count: number
  mobile_count: number
  online_count: number
  offline_count: number
  os_breakdown: Record<string, number>
}

const tabs = [
  { id: 'infrastructure', label: 'Infrastructure', icon: 'cloud' },
  { id: 'staff-devices', label: 'Staff Devices', icon: 'devices' },
  { id: 'code-repos', label: 'Code Repos', icon: 'code' },
  { id: 'people', label: 'People', icon: 'people' },
  { id: 'systems', label: 'Systems', icon: 'dns' },
]

// ============ Staff Devices Tab ============
function StaffDevicesTab() {
  const [devices, setDevices] = useState<Device[]>([])
  const [stats, setStats] = useState<DeviceStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalItems, setTotalItems] = useState(0)
  const [category, setCategory] = useState<'laptop' | 'mobile'>('laptop')
  const [deviceType, setDeviceType] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const limit = 20

  const fetchDevices = () => {
    setLoading(true)
    let url = `/assets/devices?page=${page}&limit=${limit}&category=${category}`
    if (deviceType) url += `&device_type=${deviceType}`

    api.get(url)
      .then(res => {
        const d = (res.data as any)?.data
        if (d) {
          setDevices(d.data || [])
          setTotalPages(d.total_pages || 1)
          setTotalItems(d.total_items || 0)
        }
      })
      .catch(() => setDevices([]))
      .finally(() => setLoading(false))
  }

  const fetchStats = () => {
    api.get('/assets/devices/stats')
      .then(res => setStats(res.data as any))
      .catch(() => {})
  }

  useEffect(() => { fetchDevices() }, [page, deviceType, category])
  useEffect(() => { fetchStats() }, [])

  const filteredDevices = devices.filter(d => {
    if (!searchQuery) return true
    const q = searchQuery.toLowerCase()
    return d.name.toLowerCase().includes(q) ||
      (d.serial_number || '').toLowerCase().includes(q) ||
      (d.brand || '').toLowerCase().includes(q) ||
      (d.model || '').toLowerCase().includes(q)
  })

  const getOSIcon = (os: string) => {
    switch (os.toLowerCase()) {
      case 'windows': return '🪟'
      case 'macos': case 'mac': return '🍎'
      case 'linux': return '🐧'
      case 'android': return '🤖'
      case 'ios': return '📱'
      default: return '💻'
    }
  }

  return (
    <div className="space-y-4">
      {/* Stats */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
          <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm border-l-4 border-primary">
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">Total</p>
            <p className="text-xl font-black text-on-surface">{stats.total_devices}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm border-l-4 border-blue-500">
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">Laptop/Desktop</p>
            <p className="text-xl font-black text-blue-600">{stats.laptop_count}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm border-l-4 border-purple-500">
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">Smartphone</p>
            <p className="text-xl font-black text-purple-600">{stats.mobile_count}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm border-l-4 border-emerald-500">
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">Online</p>
            <p className="text-xl font-black text-emerald-600">{stats.online_count}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm border-l-4 border-amber-500">
            <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1">OS Types</p>
            <p className="text-xl font-black text-on-surface">{Object.keys(stats.os_breakdown).length}</p>
            <div className="flex flex-wrap gap-1 mt-1">
              {Object.entries(stats.os_breakdown).map(([os, count]) => (
                <span key={os} className="text-[9px] text-on-surface-variant capitalize">{os}: {count}</span>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Category Tabs */}
      <div className="flex gap-1 border-b border-outline-variant/20">
        <button onClick={() => { setCategory('laptop'); setDeviceType(''); setPage(1) }}
          className={`px-4 py-2 text-xs font-medium border-b-2 transition-colors ${category === 'laptop' ? 'border-primary text-primary' : 'border-transparent text-on-surface-variant hover:text-on-surface'}`}>
          <span className="flex items-center gap-1.5"><span className="material-symbols-outlined text-[16px]">laptop</span>Laptop / Desktop</span>
        </button>
        <button onClick={() => { setCategory('mobile'); setDeviceType('smartphone'); setPage(1) }}
          className={`px-4 py-2 text-xs font-medium border-b-2 transition-colors ${category === 'mobile' ? 'border-primary text-primary' : 'border-transparent text-on-surface-variant hover:text-on-surface'}`}>
          <span className="flex items-center gap-1.5"><span className="material-symbols-outlined text-[16px]">smartphone</span>Smartphone</span>
        </button>
      </div>

      {/* Toolbar */}
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div className="flex items-center gap-3">
          {category === 'laptop' && (
            <select value={deviceType} onChange={e => { setDeviceType(e.target.value); setPage(1) }}
              className="text-xs font-medium border border-outline-variant/30 rounded-lg px-2.5 py-1.5 bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30">
              <option value="">All</option>
              <option value="laptop">Laptop</option>
              <option value="desktop">Desktop</option>
            </select>
          )}
          <span className="text-xs text-on-surface-variant">{totalItems} devices</span>
        </div>
        <div className="relative">
          <span className="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-[14px] text-on-surface-variant">search</span>
          <input type="text" placeholder="Search..." value={searchQuery} onChange={e => setSearchQuery(e.target.value)}
            className="pl-7 pr-3 py-1.5 text-xs border border-outline-variant/30 rounded-lg bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30 w-48" />
        </div>
      </div>

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 overflow-x-auto">
        {loading ? (
          <div className="flex justify-center py-12">
            <div className="animate-spin w-8 h-8 border-4 border-primary border-t-transparent rounded-full" />
          </div>
        ) : (
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-surface-container-low/50">
                {category === 'laptop' && <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Status</th>}
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Device Name</th>
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">OS</th>
                {category === 'laptop' && <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Serial Number</th>}
                {category === 'mobile' && <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Brand</th>}
                {category === 'mobile' && <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Model</th>}
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Registered</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/10">
              {filteredDevices.map(d => (
                <tr key={d.id} className="hover:bg-surface-container-low/30 transition-colors">
                  {category === 'laptop' && (
                    <td className="px-4 py-3">
                      <span className={`inline-flex items-center gap-1.5 text-[10px] font-bold px-2.5 py-1 rounded-full ${d.is_online ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'}`}>
                        <span className={`w-2 h-2 rounded-full ${d.is_online ? 'bg-emerald-500 animate-pulse' : 'bg-red-400'}`} />
                        {d.is_online ? 'Online' : 'Offline'}
                      </span>
                    </td>
                  )}
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">{getOSIcon(d.os)}</span>
                      <span className="text-sm font-semibold text-on-surface">{d.name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3"><span className="text-xs text-on-surface-variant capitalize">{d.os}</span></td>
                  {category === 'laptop' && (
                    <td className="px-4 py-3"><code className="text-xs text-on-surface-variant font-mono bg-surface-container-high px-2 py-0.5 rounded">{d.serial_number}</code></td>
                  )}
                  {category === 'mobile' && <td className="px-4 py-3"><span className="text-xs font-medium text-on-surface">{d.brand || '—'}</span></td>}
                  {category === 'mobile' && <td className="px-4 py-3"><span className="text-xs text-on-surface-variant">{d.model || '—'}</span></td>}
                  <td className="px-4 py-3">
                    <span className="text-xs text-on-surface-variant">{new Date(d.created_at).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })}</span>
                  </td>
                </tr>
              ))}
              {filteredDevices.length === 0 && (
                <tr><td colSpan={6} className="text-center py-12 text-on-surface-variant text-sm">No devices found</td></tr>
              )}
            </tbody>
          </table>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <span className="text-xs text-on-surface-variant">Page {page} of {totalPages}</span>
          <div className="flex gap-2">
            <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1}
              className="px-3 py-1.5 text-xs font-bold border border-outline-variant/30 rounded-lg disabled:opacity-30 hover:bg-surface-container-low transition-colors">Previous</button>
            <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page >= totalPages}
              className="px-3 py-1.5 text-xs font-bold border border-outline-variant/30 rounded-lg disabled:opacity-30 hover:bg-surface-container-low transition-colors">Next</button>
          </div>
        </div>
      )}
    </div>
  )
}

// ============ Placeholder Tabs ============
function InfrastructureTab() {
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">cloud</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">Infrastructure</h3>
        <p className="text-sm text-on-surface-variant">Cloud resources, servers, databases, and network assets.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">AWS, GCP, Cloudflare, and other integrations coming soon</p>
      </div>
    </div>
  )
}

function CodeReposTab() {
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">code</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">Code Repositories</h3>
        <p className="text-sm text-on-surface-variant">GitHub, GitLab, and Bitbucket repositories linked to your organization.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">Coming soon</p>
      </div>
    </div>
  )
}

function PeopleAssetsTab() {
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">people</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">People</h3>
        <p className="text-sm text-on-surface-variant">People-related assets: access accounts, security training, and compliance status.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">Coming soon</p>
      </div>
    </div>
  )
}

function SystemsTab() {
  return (
    <div className="bg-surface-container-lowest rounded-xl p-6 shadow-sm">
      <div className="text-center py-12">
        <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 mb-4">dns</span>
        <h3 className="text-lg font-bold text-on-surface mb-2">Systems</h3>
        <p className="text-sm text-on-surface-variant">Internal applications, SaaS tools, and system integrations.</p>
        <p className="text-xs text-on-surface-variant/60 mt-2">Coming soon</p>
      </div>
    </div>
  )
}

// ============ Main Page ============
export function AssetRegisterPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'staff-devices'

  const setTab = (tabId: string) => {
    setSearchParams({ tab: tabId })
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-black text-on-surface font-headline">Asset Register</h1>
        <p className="text-sm text-on-surface-variant mt-1">
          Manage and monitor all organizational assets — infrastructure, devices, code, and systems.
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
        {activeTab === 'infrastructure' && <InfrastructureTab />}
        {activeTab === 'staff-devices' && <StaffDevicesTab />}
        {activeTab === 'code-repos' && <CodeReposTab />}
        {activeTab === 'people' && <PeopleAssetsTab />}
        {activeTab === 'systems' && <SystemsTab />}
      </div>
    </div>
  )
}
