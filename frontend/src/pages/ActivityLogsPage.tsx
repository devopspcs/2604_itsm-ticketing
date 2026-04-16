import { useEffect, useState } from 'react'
import api from '../services/api'
import type { ActivityLog, User } from '../types'

const actionIcon: Record<string, string> = {
  ticket_created: 'post_add', status_changed: 'edit_note', assigned: 'person_add',
  reassigned: 'person_search', approval_requested: 'pending', approval_decided: 'verified',
  field_updated: 'edit',
}
const statusChip: Record<string, string> = {
  ticket_created: 'border-slate-500 bg-slate-500/10 text-slate-700',
  status_changed: 'border-accent-600 bg-accent-600/10 text-accent-700',
  assigned: 'border-green-600 bg-green-600/10 text-green-700',
  reassigned: 'border-green-600 bg-green-600/10 text-green-700',
  approval_requested: 'border-amber-500 bg-amber-500/10 text-amber-700',
  approval_decided: 'border-tertiary bg-tertiary/10 text-tertiary',
  field_updated: 'border-accent-600 bg-accent-600/10 text-accent-700',
}

export function ActivityLogsPage() {
  const [logs, setLogs] = useState<ActivityLog[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState('all')

  const userName = (id?: string) => {
    if (!id) return 'System'
    const u = users.find(u => u.id === id)
    return u ? u.full_name : id.slice(0, 8) + '...'
  }

  // Resolve UUID values to user names in activity details
  const resolveValue = (val?: string | null) => {
    if (!val) return ''
    // Check if value looks like a UUID
    if (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(val)) {
      return userName(val)
    }
    return val
  }

  useEffect(() => {
    // Load users for name resolution
    api.get<User[]>('/users/list').then(r => setUsers(r.data ?? [])).catch(() => {})
  }, [])

  useEffect(() => {
    // Fetch recent activity logs — we'll get from a few recent tickets
    api.get('/tickets', { params: { page: 1, page_size: 5 } })
      .then(async (res) => {
        const tickets = res.data.tickets ?? []
        const allLogs: ActivityLog[] = []
        for (const t of tickets) {
          try {
            const r = await api.get(`/tickets/${t.id}/activities`)
            allLogs.push(...(r.data ?? []))
          } catch {}
        }
        allLogs.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
        setLogs(allLogs)
      })
      .finally(() => setLoading(false))
  }, [])

  const filtered = filter === 'all' ? logs : logs.filter(l => l.action.includes(filter))

  return (
    <div className="max-w-[1600px] mx-auto p-8">
      {/* Header */}
      <div className="grid grid-cols-12 gap-6 mb-8">
        <div className="col-span-12 lg:col-span-8">
          <h2 className="text-3xl font-extrabold text-on-surface mb-2">Audit Trail &amp; Activity Logs</h2>
          <p className="text-on-surface-variant max-w-2xl">
            A high-fidelity record of system operations, user interactions, and status transitions across the ITSM ecosystem.
          </p>
        </div>
        <div className="col-span-12 lg:col-span-4 flex items-center justify-end gap-3">
          <button
            onClick={() => {
              // Generate CSV from current filtered logs
              const headers = ['Timestamp', 'Action', 'Ticket ID', 'Actor', 'Old Value', 'New Value']
              const rows = filtered.map(log => [
                new Date(log.created_at).toISOString(),
                log.action.replace(/_/g, ' '),
                log.ticket_id,
                userName(log.actor_id),
                resolveValue(log.old_value) || '',
                resolveValue(log.new_value) || '',
              ])
              const csv = [headers.join(','), ...rows.map(r => r.map(c => `"${(c ?? '').replace(/"/g, '""')}"`).join(','))].join('\n')
              const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
              const url = URL.createObjectURL(blob)
              const link = document.createElement('a')
              link.href = url
              link.download = `activity-logs-${new Date().toISOString().slice(0, 10)}.csv`
              link.click()
              URL.revokeObjectURL(url)
            }}
            className="bg-surface-container-high text-on-secondary-container px-4 py-2 rounded-xl font-semibold flex items-center gap-2 hover:bg-surface-container-highest transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-sm">download</span>
            Export Report
          </button>
          <button className="bg-gradient-to-r from-primary to-primary-container text-white px-5 py-2 rounded-xl font-bold flex items-center gap-2 shadow-sm">
            <span className="material-symbols-outlined text-sm">filter_list</span>
            Filters
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        {[
          { label: 'Total Actions (24h)', value: logs.length, delta: '+12%', border: 'border-primary' },
          { label: 'Critical Alerts', value: 0, delta: 'Priority Red', border: 'border-error', deltaClass: 'text-error' },
          { label: 'System Changes', value: logs.filter(l => l.action === 'status_changed').length, delta: 'Config updates', border: 'border-accent-600' },
          { label: 'Active Operators', value: 1, delta: 'Online now', border: 'border-tertiary' },
        ].map((s) => (
          <div key={s.label} className={`bg-surface-container-low p-6 rounded-xl border-l-4 ${s.border}`}>
            <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">{s.label}</p>
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-extrabold">{s.value}</span>
              <span className={`text-xs font-bold ${s.deltaClass ?? 'text-emerald-600'}`}>{s.delta}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Table */}
      <section className="bg-surface-container-lowest rounded-2xl shadow-sm overflow-hidden">
        <div className="px-8 py-6 flex flex-wrap items-center justify-between gap-4 border-b border-surface-container-low">
          <div className="flex items-center gap-4">
            {['all', 'ticket', 'approval', 'status'].map((f) => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`text-sm font-medium px-3 py-1 rounded-full transition-colors ${
                  filter === f ? 'font-bold text-primary bg-primary-fixed' : 'text-on-surface-variant hover:bg-surface-container-low'
                }`}
              >
                {f === 'all' ? 'All Activity' : f === 'ticket' ? 'Ticket Changes' : f === 'approval' ? 'Approvals' : 'Status Changes'}
              </button>
            ))}
          </div>
          <span className="text-xs text-on-surface-variant font-medium">Viewing {filtered.length} entries</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-surface-container-low/50">
                {['Timestamp', 'Action', 'Entity / Target', 'Details', 'Status'].map(h => (
                  <th key={h} className="px-8 py-4 text-[10px] font-bold text-on-surface-variant uppercase tracking-widest">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-surface-container-low">
              {loading ? (
                <tr><td colSpan={5} className="px-8 py-12 text-center text-on-surface-variant">
                  <span className="material-symbols-outlined animate-spin block mb-2">refresh</span>
                  Loading activity logs...
                </td></tr>
              ) : filtered.length === 0 ? (
                <tr><td colSpan={5} className="px-8 py-12 text-center text-on-surface-variant">
                  <span className="material-symbols-outlined text-4xl block mb-2 opacity-30">history</span>
                  No activity logs found
                </td></tr>
              ) : filtered.map((log) => (
                <tr key={log.id} className="hover:bg-surface-container-low/30 transition-colors">
                  <td className="px-8 py-5">
                    <div className="text-sm font-semibold text-on-surface">{new Date(log.created_at).toLocaleDateString()}</div>
                    <div className="text-[10px] text-on-surface-variant">{new Date(log.created_at).toLocaleTimeString()} UTC</div>
                  </td>
                  <td className="px-6 py-5">
                    <div className="flex items-center gap-2">
                      <span className="material-symbols-outlined text-sm text-primary">{actionIcon[log.action] ?? 'info'}</span>
                      <span className="text-sm font-medium capitalize">{log.action.replace(/_/g, ' ')}</span>
                    </div>
                  </td>
                  <td className="px-6 py-5">
                    <div className="text-sm font-mono text-on-primary-fixed-variant">{log.ticket_id.slice(0, 8).toUpperCase()}</div>
                    <div className="text-[10px] text-on-surface-variant">by {userName(log.actor_id)}</div>
                  </td>
                  <td className="px-6 py-5">
                    {log.old_value && <div className="text-[10px] text-on-surface-variant">{resolveValue(log.old_value)} → {resolveValue(log.new_value)}</div>}
                    {!log.old_value && log.new_value && <div className="text-xs text-on-surface-variant">{resolveValue(log.new_value)}</div>}
                  </td>
                  <td className="px-6 py-5 text-right">
                    <span className={`inline-flex items-center px-2 py-1 rounded border-l-2 text-[10px] font-bold tracking-tight ${statusChip[log.action] ?? 'border-slate-500 bg-slate-500/10 text-slate-700'}`}>
                      {log.action.replace(/_/g, ' ').toUpperCase()}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="px-8 py-6 border-t border-surface-container-low flex items-center justify-between">
          <p className="text-xs text-on-surface-variant">Showing {filtered.length} activities</p>
        </div>
      </section>

      {/* Footer */}
      <footer className="px-8 py-4 bg-surface-container-low flex items-center justify-between mt-6 rounded-xl">
        <div className="flex items-center gap-6">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
            <span className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider">Audit Service: Operational</span>
          </div>
        </div>
        <span className="text-[10px] text-on-surface-variant">© 2024 PCS Payments. All Rights Reserved.</span>
      </footer>
    </div>
  )
}
