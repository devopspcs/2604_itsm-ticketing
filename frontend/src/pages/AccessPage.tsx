import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import api from '../services/api'
import type { User } from '../types'

interface SyncStatus {
  last_synced_at: string | null
  status: string
  total_records: number
}

interface ServiceAccessEntry {
  service: string
  service_type: string
  has_access: boolean
  roles: string[]
  account_name: string
  status: string
}

interface ExternalService {
  id: string
  name: string
  type: string
  is_active: boolean
}

const tabs = [
  { id: 'staff-access', label: 'Staff Access', icon: 'badge' },
  { id: 'services', label: 'Services', icon: 'dns' },
]

function StaffAccessTab() {
  const [users, setUsers] = useState<User[]>([])
  const [services, setServices] = useState<ExternalService[]>([])
  const [accessMap, setAccessMap] = useState<Record<string, ServiceAccessEntry[]>>({})
  const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null)
  const [syncing, setSyncing] = useState(false)
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')

  const fetchData = async () => {
    setLoading(true)
    try {
      const [usersRes, servicesRes, syncRes] = await Promise.all([
        api.get<User[]>('/users'),
        api.get('/access/services'),
        api.get('/access/sync-status'),
      ])
      setUsers((usersRes.data as any) ?? [])
      setServices((servicesRes.data as any) ?? [])
      setSyncStatus(syncRes.data as any)

      // Fetch access for all users
      const activeUsers = ((usersRes.data as any) ?? []).filter((u: User) => u.is_active)
      const map: Record<string, ServiceAccessEntry[]> = {}
      
      // Batch: fetch access per user (from cache, should be fast)
      await Promise.all(
        activeUsers.map(async (u: User) => {
          try {
            const res = await api.get(`/access/check?email=${encodeURIComponent(u.email)}`)
            map[u.email] = (res.data as any) ?? []
          } catch {
            map[u.email] = []
          }
        })
      )
      setAccessMap(map)
    } catch {}
    setLoading(false)
  }

  useEffect(() => { fetchData() }, [])

  const handleSync = async () => {
    setSyncing(true)
    try {
      await api.post('/access/sync')
      // Refresh data after sync
      await fetchData()
    } catch {}
    setSyncing(false)
  }

  const filteredUsers = users.filter(u => u.is_active).filter(u => {
    if (!searchQuery) return true
    const q = searchQuery.toLowerCase()
    return u.full_name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q)
  }).sort((a, b) => a.full_name.localeCompare(b.full_name))

  const activeServices = services.filter(s => s.is_active)

  const getAccessIcon = (email: string, serviceName: string) => {
    const entries = accessMap[email] || []
    const entry = entries.find(e => e.service === serviceName)
    if (!entry) return <span className="text-on-surface-variant/30">—</span>
    if (entry.has_access) {
      return (
        <span title={`${entry.roles.join(', ')}`} className="inline-flex items-center gap-0.5">
          <span className="w-5 h-5 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center">
            <span className="material-symbols-outlined text-[14px]">check</span>
          </span>
        </span>
      )
    }
    return <span className="text-on-surface-variant/30">—</span>
  }

  const getAccessCount = (email: string) => {
    const entries = accessMap[email] || []
    return entries.filter(e => e.has_access).length
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
          <div className="relative">
            <span className="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-[16px] text-on-surface-variant">search</span>
            <input type="text" placeholder="Search staff..." value={searchQuery} onChange={e => setSearchQuery(e.target.value)}
              className="pl-8 pr-3 py-2 text-sm border border-outline-variant/30 rounded-lg bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30 w-56" />
          </div>
          <span className="text-xs text-on-surface-variant">{filteredUsers.length} staff</span>
        </div>
        <div className="flex items-center gap-3">
          {syncStatus && (
            <span className="text-[10px] text-on-surface-variant">
              Last synced: {syncStatus.last_synced_at
                ? new Date(syncStatus.last_synced_at).toLocaleString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
                : 'Never'}
            </span>
          )}
          <button onClick={handleSync} disabled={syncing}
            className="flex items-center gap-1.5 px-4 py-2 text-sm font-bold text-white bg-primary rounded-lg hover:opacity-90 disabled:opacity-50 transition-all">
            <span className={`material-symbols-outlined text-[16px] ${syncing ? 'animate-spin' : ''}`}>
              {syncing ? 'refresh' : 'sync'}
            </span>
            {syncing ? 'Syncing...' : 'Sync'}
          </button>
        </div>
      </div>

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-surface-container-low/50">
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant w-8"></th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Staff ↕</th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Role ↕</th>
              <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Status</th>
              {activeServices.map(s => (
                <th key={s.id} className="px-3 py-3 text-[9px] font-black uppercase tracking-widest text-on-surface-variant text-center" title={s.name}>
                  {s.name.length > 12 ? s.name.slice(0, 12) + '…' : s.name}
                </th>
              ))}
              <th className="px-3 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant text-center">Total</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-outline-variant/10">
            {filteredUsers.map(u => {
              const initial = u.full_name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()
              return (
                <tr key={u.id} className="hover:bg-surface-container-low/30 transition-colors">
                  <td className="px-4 py-2.5">
                    <div className="w-7 h-7 rounded-full bg-surface-container-high flex items-center justify-center text-[9px] font-bold text-on-surface-variant overflow-hidden">
                      {u.avatar_url ? <img src={u.avatar_url} className="w-full h-full object-cover" referrerPolicy="no-referrer" /> : initial}
                    </div>
                  </td>
                  <td className="px-4 py-2.5">
                    <p className="text-sm font-semibold text-on-surface">{u.full_name}</p>
                    <p className="text-[10px] text-on-surface-variant">{u.email}</p>
                  </td>
                  <td className="px-4 py-2.5">
                    <span className="text-xs text-on-surface-variant">{u.position || '—'}</span>
                  </td>
                  <td className="px-4 py-2.5">
                    <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-700">Active</span>
                  </td>
                  {activeServices.map(s => (
                    <td key={s.id} className="px-3 py-2.5 text-center">
                      {getAccessIcon(u.email, s.name)}
                    </td>
                  ))}
                  <td className="px-3 py-2.5 text-center">
                    {getAccessCount(u.email) > 0 ? (
                      <span className="text-xs font-bold text-primary">+{getAccessCount(u.email)}</span>
                    ) : (
                      <span className="text-xs text-on-surface-variant/30">0</span>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function ServicesTab() {
  const [services, setServices] = useState<ExternalService[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/access/services').then(r => setServices((r.data as any) ?? [])).finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="flex justify-center py-12"><div className="animate-spin w-8 h-8 border-4 border-primary border-t-transparent rounded-full" /></div>

  return (
    <div className="bg-surface-container-lowest rounded-xl shadow-sm border border-outline-variant/10 p-6">
      <h3 className="text-sm font-bold text-on-surface mb-4 uppercase tracking-wider">Registered Services</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {services.map(s => (
          <div key={s.id} className="p-4 rounded-lg border border-outline-variant/20 flex items-center gap-3">
            <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${s.type === 'pritunl' ? 'bg-blue-100 text-blue-600' : 'bg-amber-100 text-amber-600'}`}>
              <span className="material-symbols-outlined text-[20px]">{s.type === 'pritunl' ? 'vpn_lock' : 'point_of_sale'}</span>
            </div>
            <div>
              <p className="text-sm font-semibold text-on-surface">{s.name}</p>
              <p className="text-[10px] text-on-surface-variant capitalize">{s.type}</p>
            </div>
            <span className={`ml-auto text-[10px] font-bold px-2 py-0.5 rounded-full ${s.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'}`}>
              {s.is_active ? 'Active' : 'Disabled'}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

export function AccessPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'staff-access'

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-black text-on-surface font-headline">Access</h1>
        <p className="text-sm text-on-surface-variant mt-1">Monitor staff access across all organizational services and applications.</p>
      </div>

      <div className="border-b border-outline-variant/20">
        <nav className="flex gap-1 overflow-x-auto">
          {tabs.map(tab => (
            <button key={tab.id} onClick={() => setSearchParams({ tab: tab.id })}
              className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium whitespace-nowrap border-b-2 transition-colors ${
                activeTab === tab.id ? 'border-primary text-primary' : 'border-transparent text-on-surface-variant hover:text-on-surface'
              }`}>
              <span className="material-symbols-outlined text-[18px]">{tab.icon}</span>
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      <div>
        {activeTab === 'staff-access' && <StaffAccessTab />}
        {activeTab === 'services' && <ServicesTab />}
      </div>
    </div>
  )
}
