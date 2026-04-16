import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { dashboardService } from '../services/dashboard.service'
import type { DashboardStats } from '../types'
import type { RootState } from '../store'

function StatCard({ label, value, icon, borderColor, delta, deltaLabel, onClick }: {
  label: string; value: number; icon: string; borderColor: string; delta?: string; deltaLabel?: string; onClick?: () => void
}) {
  return (
    <div
      onClick={onClick}
      className={`bg-surface-container-lowest p-6 rounded-xl shadow-sm border-l-4 ${borderColor} hover:shadow-md transition-all cursor-pointer active:scale-[0.98]`}
    >
      <div className="flex justify-between items-start mb-4">
        <span className="text-sm font-bold text-on-surface-variant uppercase tracking-wider">{label}</span>
        <span className="material-symbols-outlined text-on-surface-variant/30">{icon}</span>
      </div>
      <div className="flex items-baseline gap-2">
        <span className="text-3xl font-black text-on-surface">{value.toLocaleString()}</span>
        {delta && <span className={`text-xs font-bold ${deltaLabel === 'High Risk' ? 'text-amber-600' : delta.startsWith('+') ? 'text-emerald-600' : 'text-red-500'}`}>{delta}</span>}
        {deltaLabel && <span className="text-xs font-bold text-emerald-600">{deltaLabel}</span>}
      </div>
    </div>
  )
}

export function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const role = useSelector((s: RootState) => s.auth.role) ?? 'user'
  const navigate = useNavigate()

  useEffect(() => {
    dashboardService.getStats()
      .then((res) => setStats(res.data))
      .catch(() => setError('Failed to load dashboard stats'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="max-w-7xl mx-auto p-6 md:p-10">
      {/* Welcome Header */}
      <div className="mb-10">
        <h1 className="text-3xl font-extrabold text-on-surface tracking-tight mb-2">
          {role === 'user' ? 'My Service Requests' : 'Operational Overview'}
        </h1>
        <p className="text-on-surface-variant font-medium">
          {role === 'user'
            ? 'Track and manage your IT service requests.'
            : 'Monitoring architectural health and service tickets.'}
        </p>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20">
          <div className="flex items-center gap-3 text-on-surface-variant">
            <span className="material-symbols-outlined animate-spin">refresh</span>
            <span>Loading dashboard...</span>
          </div>
        </div>
      ) : error ? (
        <div className="bg-error-container/30 text-error px-6 py-4 rounded-xl flex items-center gap-3">
          <span className="material-symbols-outlined">error</span>
          <span>{error}</span>
        </div>
      ) : stats ? (
        <>
          {/* Summary Cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
            <StatCard label="Total Tickets" value={stats.total_tickets ?? 0} icon="confirmation_number" borderColor="border-primary" delta="+12%" onClick={() => navigate('/tickets')} />
            <StatCard label="Pending Approvals" value={(stats.by_status as Record<string,number>)?.['pending_approval'] ?? 0} icon="fact_check" borderColor="border-amber-500" delta="-5" onClick={() => navigate('/tickets?status=pending_approval')} />
            <StatCard label="Active Incidents" value={(stats.by_type as Record<string,number>)?.['incident'] ?? 0} icon="report" borderColor="border-error" deltaLabel="High Risk" onClick={() => navigate('/tickets?type=incident')} />
            <StatCard label="Completed" value={(stats.by_status as Record<string,number>)?.['done'] ?? 0} icon="check_circle" borderColor="border-emerald-500" deltaLabel="Daily Goal Hit" onClick={() => navigate('/tickets?status=done')} />
          </div>

          {/* Bento Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 mb-10">
            {/* Chart */}
            <div className="lg:col-span-8 bg-surface-container-lowest rounded-xl p-6 shadow-sm">
              <div className="flex justify-between items-center mb-6">
                <div>
                  <h3 className="text-xl font-bold text-on-surface">Tickets per Category</h3>
                  <p className="text-sm text-on-surface-variant">Distribution across ticket types</p>
                </div>
                <div className="flex gap-2">
                  <button className="px-3 py-1 text-xs font-bold bg-primary text-white rounded-full">Weekly</button>
                </div>
              </div>
              <div className="relative flex items-end justify-around gap-8 px-8" style={{ height: '280px' }}>
                {[
                  { label: 'Helpdesk', value: (stats.by_type as Record<string,number>)?.['helpdesk_request'] ?? 0, color: '#00307d', light: '#dae2ff' },
                  { label: 'Incident', value: (stats.by_type as Record<string,number>)?.['incident'] ?? 0, color: '#ba1a1a', light: '#ffdad6' },
                  { label: 'Change', value: (stats.by_type as Record<string,number>)?.['change_request'] ?? 0, color: '#f59e0b', light: '#fef3c7' },
                ].map((bar) => {
                  const byType = (stats.by_type as Record<string,number>) ?? {}
                  const max = Math.max(byType['helpdesk_request'] ?? 0, byType['incident'] ?? 0, byType['change_request'] ?? 0, 1)
                  const barHeight = Math.max(Math.round((bar.value / max) * 220), bar.value > 0 ? 50 : 12)
                  return (
                    <div key={bar.label} className="flex flex-col items-center gap-3 flex-1 cursor-pointer group"
                      onClick={() => navigate(`/tickets?type=${bar.label === 'Helpdesk' ? 'helpdesk_request' : bar.label === 'Incident' ? 'incident' : 'change_request'}`)}>
                      <div className="w-full max-w-[120px] rounded-t-xl relative overflow-hidden transition-all group-hover:opacity-80"
                        style={{ height: `${barHeight}px`, backgroundColor: bar.light }}>
                        <div className="absolute bottom-0 w-full rounded-t-xl transition-all"
                          style={{ height: '60%', backgroundColor: bar.color }} />
                      </div>
                      <div className="text-center">
                        <div className="text-xs font-bold text-on-surface-variant">{bar.label}</div>
                        <div className="text-lg font-black text-on-surface">{bar.value}</div>
                      </div>
                    </div>
                  )
                })}
                {/* Y-axis grid lines */}
                <div className="absolute inset-0 flex flex-col justify-between pointer-events-none opacity-10" style={{ bottom: '50px' }}>
                  <div className="border-b border-outline w-full" />
                  <div className="border-b border-outline w-full" />
                  <div className="border-b border-outline w-full" />
                </div>
              </div>
            </div>

            {/* Recent Activity */}
            <div className="lg:col-span-4 bg-surface-container-lowest rounded-xl p-6 shadow-sm flex flex-col">
              <div className="flex justify-between items-center mb-6">
                <h3 className="text-lg font-bold text-on-surface">Recent Tickets</h3>
                <Link to="/tickets" className="text-xs font-bold text-primary hover:underline">View All</Link>
              </div>
              <div className="space-y-4 overflow-y-auto flex-1">
                {(stats.recent_tickets ?? []).length === 0 ? (
                  <p className="text-sm text-on-surface-variant text-center py-8">No tickets yet</p>
                ) : (
                  (stats.recent_tickets ?? []).filter(Boolean).map((t) => (
                    <Link to={`/tickets/${t.id}`} key={t.id} className="flex gap-3 group hover:bg-surface-container-low p-2 rounded-xl transition-colors">
                      <div className="flex-shrink-0 w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center text-blue-600 group-hover:bg-blue-600 group-hover:text-white transition-colors">
                        <span className="material-symbols-outlined text-lg">
                          {(t.type ?? '') === 'incident' ? 'report' : (t.type ?? '') === 'change_request' ? 'published_with_changes' : 'contact_support'}
                        </span>
                      </div>
                      <div className="flex-grow min-w-0">
                        <p className="text-sm font-bold text-on-surface truncate">{t.title}</p>
                        <div className="mt-1 flex gap-2">
                          <span className={`text-[9px] px-2 py-0.5 rounded-full font-bold uppercase ${
                            t.priority === 'critical' ? 'bg-red-100 text-red-700' :
                            t.priority === 'high' ? 'bg-amber-100 text-amber-700' :
                            'bg-surface-container-high text-on-surface-variant'
                          }`}>{t.priority}</span>
                          <span className="text-[9px] px-2 py-0.5 rounded-full bg-surface-container-high text-on-surface-variant font-bold uppercase">{(t.status ?? '').replace(/_/g, ' ')}</span>
                        </div>
                      </div>
                    </Link>
                  ))
                )}
              </div>
              <div className="mt-6 pt-4 border-t border-outline-variant/20">
                <Link
                  to="/tickets/new"
                  className="w-full py-3 bg-primary text-white rounded-xl font-bold text-sm hover:opacity-90 active:scale-[0.98] transition-all flex items-center justify-center gap-2 shadow-lg shadow-primary/20"
                >
                  <span className="material-symbols-outlined text-sm">add</span>
                  Create New Ticket
                </Link>
              </div>
            </div>
          </div>

          {/* SLA Footer */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center bg-primary-container/5 rounded-2xl p-8 border border-primary/5">
            <div className="flex items-center gap-6">
              <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                <span className="material-symbols-outlined text-3xl">query_stats</span>
              </div>
              <div>
                <h3 className="text-lg font-bold text-on-surface">SLA Compliance</h3>
                <p className="text-on-surface-variant text-sm">
                  Average resolution time is{' '}
                  <span className="font-bold text-primary">
                    {(stats.avg_resolution_hours ?? 0).toFixed(1)} hours
                  </span>
                </p>
              </div>
            </div>
            <div className="flex justify-end gap-10">
              <div className="text-center">
                <div className="text-2xl font-black text-on-surface">
                  {(stats.sla_compliance_rate ?? 0).toFixed(1)}%
                </div>
                <div className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest">Compliance</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-black text-emerald-600">{stats.on_time_count ?? 0}</div>
                <div className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest">On Time</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-black text-red-600">{stats.breached_count ?? 0}</div>
                <div className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest">Breached</div>
              </div>
            </div>
          </div>
        </>
      ) : null}

      {/* FAB Mobile */}
      <Link
        to="/tickets/new"
        className="fixed bottom-8 right-8 w-14 h-14 bg-primary text-white rounded-full shadow-2xl flex items-center justify-center hover:scale-110 active:scale-95 transition-transform md:hidden z-50"
      >
        <span className="material-symbols-outlined">add</span>
      </Link>
    </div>
  )
}
