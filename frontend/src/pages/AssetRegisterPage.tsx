import { useEffect, useState } from 'react'
import api from '../services/api'

interface Device {
  id: number
  name: string
  os: string
  serial_number: string
  is_online: boolean
  created_at: string
  updated_at: string
}

interface DeviceResponse {
  data: {
    data: Device[]
    page: number
    limit: number
    total_items: number
    total_pages: number
    has_next: boolean
    has_previous: boolean
  }
  meta: { code: number; message: string }
  success: boolean
}

interface DeviceStats {
  total_devices: number
  online_count: number
  offline_count: number
  os_breakdown: Record<string, number>
}

export function AssetRegisterPage() {
  const [devices, setDevices] = useState<Device[]>([])
  const [stats, setStats] = useState<DeviceStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalItems, setTotalItems] = useState(0)
  const [deviceType, setDeviceType] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const limit = 20

  const fetchDevices = () => {
    setLoading(true)
    let url = `/assets/devices?page=${page}&limit=${limit}`
    if (deviceType) url += `&device_type=${deviceType}`

    api.get<DeviceResponse>(url)
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
    api.get<DeviceStats>('/assets/devices/stats')
      .then(res => setStats(res.data as any))
      .catch(() => {})
  }

  useEffect(() => { fetchDevices() }, [page, deviceType])
  useEffect(() => { fetchStats() }, [])

  const filteredDevices = devices.filter(d => {
    if (!searchQuery) return true
    const q = searchQuery.toLowerCase()
    return d.name.toLowerCase().includes(q) || d.serial_number.toLowerCase().includes(q)
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
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-black text-on-surface font-headline">Asset Register</h1>
        <p className="text-sm text-on-surface-variant mt-1">
          Manage and monitor all registered devices across the organization.
        </p>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-primary">
            <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Total Devices</p>
            <p className="text-2xl font-black text-on-surface">{stats.total_devices}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-emerald-500">
            <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Online</p>
            <p className="text-2xl font-black text-emerald-600">{stats.online_count}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-red-500">
            <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Offline</p>
            <p className="text-2xl font-black text-red-600">{stats.offline_count}</p>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-5 shadow-sm border-l-4 border-purple-500">
            <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">OS Types</p>
            <p className="text-2xl font-black text-on-surface">{Object.keys(stats.os_breakdown).length}</p>
            <div className="flex gap-2 mt-1">
              {Object.entries(stats.os_breakdown).map(([os, count]) => (
                <span key={os} className="text-[10px] text-on-surface-variant capitalize">{os}: {count}</span>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Toolbar */}
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div className="flex items-center gap-3">
          <select value={deviceType} onChange={e => { setDeviceType(e.target.value); setPage(1) }}
            className="text-sm font-medium border border-outline-variant/30 rounded-lg px-3 py-2 bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30">
            <option value="">All Devices</option>
            <option value="laptop">Laptop</option>
            <option value="desktop">Desktop</option>
            <option value="mobile">Mobile</option>
            <option value="tablet">Tablet</option>
          </select>
          <span className="text-sm text-on-surface-variant">{totalItems} devices</span>
        </div>
        <div className="relative">
          <span className="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-[16px] text-on-surface-variant">search</span>
          <input
            type="text"
            placeholder="Search by name or serial..."
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
            className="pl-8 pr-3 py-2 text-sm border border-outline-variant/30 rounded-lg bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30 w-56"
          />
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
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Status</th>
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Device Name</th>
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">OS</th>
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Serial Number</th>
                <th className="px-4 py-3 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Registered</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/10">
              {filteredDevices.map(d => (
                <tr key={d.id} className="hover:bg-surface-container-low/30 transition-colors">
                  <td className="px-4 py-3">
                    <span className={`inline-flex items-center gap-1.5 text-[10px] font-bold px-2.5 py-1 rounded-full ${
                      d.is_online ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
                    }`}>
                      <span className={`w-2 h-2 rounded-full ${d.is_online ? 'bg-emerald-500 animate-pulse' : 'bg-red-400'}`} />
                      {d.is_online ? 'Online' : 'Offline'}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">{getOSIcon(d.os)}</span>
                      <span className="text-sm font-semibold text-on-surface">{d.name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-xs text-on-surface-variant capitalize">{d.os}</span>
                  </td>
                  <td className="px-4 py-3">
                    <code className="text-xs text-on-surface-variant font-mono bg-surface-container-high px-2 py-0.5 rounded">{d.serial_number}</code>
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-xs text-on-surface-variant">
                      {new Date(d.created_at).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })}
                    </span>
                  </td>
                </tr>
              ))}
              {filteredDevices.length === 0 && (
                <tr>
                  <td colSpan={5} className="text-center py-12 text-on-surface-variant text-sm">No devices found</td>
                </tr>
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
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="px-3 py-1.5 text-xs font-bold border border-outline-variant/30 rounded-lg disabled:opacity-30 hover:bg-surface-container-low transition-colors"
            >
              Previous
            </button>
            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="px-3 py-1.5 text-xs font-bold border border-outline-variant/30 rounded-lg disabled:opacity-30 hover:bg-surface-container-low transition-colors"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
