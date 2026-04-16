import { useEffect, useState } from 'react'
import { Link, useSearchParams, useNavigate } from 'react-router-dom'
import { ticketService } from '../services/ticket.service'
import type { Ticket, TicketStatus, TicketType } from '../types'
import { Pagination } from '../components/common/Pagination'

const priorityDot: Record<string, string> = {
  critical: 'bg-error', high: 'bg-error', medium: 'bg-secondary', low: 'bg-outline-variant',
}
const priorityText: Record<string, string> = {
  critical: 'text-error', high: 'text-error', medium: 'text-secondary', low: 'text-outline',
}
const statusStyle: Record<string, string> = {
  open: 'bg-error-container/20 text-error border-error',
  in_progress: 'bg-primary-container/10 text-primary border-primary',
  pending_approval: 'bg-on-tertiary-container/10 text-on-tertiary-fixed-variant border-on-tertiary-fixed-variant',
  approved: 'bg-green-100/50 text-green-700 border-green-600',
  rejected: 'bg-error-container/20 text-error border-error',
  done: 'bg-green-100/50 text-green-700 border-green-600',
}

export function TicketListPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState(searchParams.get('status') ?? '')
  const [typeFilter, setTypeFilter] = useState(searchParams.get('type') ?? '')
  const [priorityFilter, setPriorityFilter] = useState(searchParams.get('priority') ?? '')
  const [loading, setLoading] = useState(true)

  const fetchTickets = (p = 1) => {
    setLoading(true)
    const params: Record<string, string | number> = { page: p, page_size: 10 }
    if (search) params.search = search
    if (statusFilter) params.status = statusFilter
    if (typeFilter) params.type = typeFilter
    if (priorityFilter) params.priority = priorityFilter

    ticketService.list(params)
      .then((res) => { setTickets(res.data.tickets ?? []); setTotal(res.data.total ?? 0) })
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchTickets(page) }, [page, statusFilter, typeFilter, priorityFilter])

  const handleSearch = (e: React.FormEvent) => { e.preventDefault(); setPage(1); fetchTickets(1) }

  return (
    <div className="max-w-[1400px] mx-auto p-8 lg:p-12">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-end justify-between mb-10 gap-6">
        <div>
          <h1 className="text-4xl font-extrabold tracking-tight text-on-surface mb-2">Ticket Management</h1>
          <p className="text-on-surface-variant max-w-2xl leading-relaxed">
            Curate and monitor the infrastructure lifecycle across enterprise-grade incidents and requests.
          </p>
        </div>
        <Link
          to="/tickets/new"
          className="bg-gradient-to-r from-primary to-primary-container text-white px-6 py-3 rounded-xl font-bold flex items-center gap-2 shadow-lg shadow-primary/20 hover:opacity-90 active:scale-95 transition-all whitespace-nowrap"
        >
          <span className="material-symbols-outlined">add</span>
          Create New Ticket
        </Link>
      </div>

      {/* Filter Bar */}
      <div className="bg-surface-container-lowest rounded-xl p-2 mb-8 flex flex-wrap items-center gap-2 shadow-sm border border-outline-variant/10">
        <div className="flex items-center gap-2 px-4 border-r border-outline-variant/20 py-2">
          <span className="material-symbols-outlined text-outline text-lg">filter_list</span>
          <span className="text-sm font-bold text-on-surface">Filters</span>
        </div>
        <div className="flex flex-wrap gap-2 flex-grow px-2">
          <select
            value={typeFilter}
            onChange={(e) => { setTypeFilter(e.target.value); setPage(1) }}
            className="bg-surface-container-low border-none rounded-lg text-xs font-medium focus:ring-2 focus:ring-primary/20 px-3 py-2 cursor-pointer hover:bg-surface-container-high transition-colors outline-none"
          >
            <option value="">All Categories</option>
            {(['change_request', 'incident', 'helpdesk_request'] as TicketType[]).map(t => (
              <option key={t} value={t}>{t.replace(/_/g, ' ')}</option>
            ))}
          </select>
          <select
            value={priorityFilter}
            onChange={(e) => { setPriorityFilter(e.target.value); setPage(1) }}
            className="bg-surface-container-low border-none rounded-lg text-xs font-medium focus:ring-2 focus:ring-primary/20 px-3 py-2 cursor-pointer hover:bg-surface-container-high transition-colors outline-none"
          >
            <option value="">Priority: All</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
          <select
            value={statusFilter}
            onChange={(e) => { setStatusFilter(e.target.value); setPage(1) }}
            className="bg-surface-container-low border-none rounded-lg text-xs font-medium focus:ring-2 focus:ring-primary/20 px-3 py-2 cursor-pointer hover:bg-surface-container-high transition-colors outline-none"
          >
            <option value="">Status: All</option>
            {(['open', 'in_progress', 'pending_approval', 'approved', 'rejected', 'done'] as TicketStatus[]).map(s => (
              <option key={s} value={s}>{s.replace(/_/g, ' ')}</option>
            ))}
          </select>
        </div>
        <form onSubmit={handleSearch} className="relative flex items-center ml-auto px-4 py-2 bg-surface-container-low rounded-lg">
          <span className="material-symbols-outlined text-outline text-sm absolute left-3">search</span>
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="bg-transparent border-none focus:ring-0 text-xs w-48 pl-6 text-on-surface outline-none"
            placeholder="Quick search ID or user..."
          />
        </form>
      </div>

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-2xl shadow-sm overflow-hidden border border-outline-variant/10">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-surface-container-low/50">
                {['Ticket ID', 'Title', 'Category', 'Priority', 'Status', 'Assigned To', 'Created At', 'Actions'].map(h => (
                  <th key={h} className="px-6 py-4 text-[11px] font-black uppercase tracking-widest text-on-surface-variant">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-surface-container-low">
              {loading ? (
                <tr><td colSpan={8} className="px-6 py-12 text-center text-on-surface-variant">
                  <div className="flex items-center justify-center gap-2">
                    <span className="material-symbols-outlined animate-spin">refresh</span>
                    Loading tickets...
                  </div>
                </td></tr>
              ) : tickets.length === 0 ? (
                <tr><td colSpan={8} className="px-6 py-12 text-center text-on-surface-variant">
                  <span className="material-symbols-outlined text-4xl block mb-2 opacity-30">inbox</span>
                  No tickets found
                </td></tr>
              ) : tickets.map((t) => (
                <tr key={t.id} className="hover:bg-surface-container-low/30 transition-colors group">
                  <td className="px-6 py-5">
                    <span className="text-xs font-bold text-primary-fixed-dim bg-primary/5 px-2 py-1 rounded">
                      {t.type === 'incident' ? 'INC' : t.type === 'change_request' ? 'CHG' : 'REQ'}-{t.id.slice(0, 6).toUpperCase()}
                    </span>
                  </td>
                  <td className="px-6 py-5">
                    <p className="text-sm font-semibold text-on-surface truncate max-w-xs">{t.title}</p>
                    <p className="text-[10px] text-outline mt-0.5">{t.type.replace(/_/g, ' ')}</p>
                  </td>
                  <td className="px-6 py-5">
                    <span className="text-xs text-on-surface-variant capitalize">{t.category || t.type.replace(/_/g, ' ')}</span>
                  </td>
                  <td className="px-6 py-5">
                    <div className="flex items-center gap-2">
                      <div className={`w-2 h-2 rounded-full ${priorityDot[t.priority] ?? 'bg-outline'}`} />
                      <span className={`text-xs font-bold capitalize ${priorityText[t.priority] ?? 'text-outline'}`}>{t.priority}</span>
                    </div>
                  </td>
                  <td className="px-6 py-5">
                    <div className={`inline-flex items-center gap-2 px-3 py-1 rounded-full border-l-4 ${statusStyle[t.status] ?? 'bg-surface-container text-on-surface border-outline'}`}>
                      <span className="text-[10px] font-bold capitalize">{t.status.replace(/_/g, ' ')}</span>
                    </div>
                  </td>
                  <td className="px-6 py-5">
                    <span className="text-xs font-medium text-on-surface-variant">{t.assigned_to ? t.assigned_to.slice(0, 8) + '...' : <span className="italic text-outline">Unassigned</span>}</span>
                  </td>
                  <td className="px-6 py-5 text-xs text-outline">
                    {new Date(t.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                  </td>
                  <td className="px-6 py-5 text-center">
                    <Link to={`/tickets/${t.id}`} className="p-2 hover:bg-surface-container-high rounded-lg transition-colors inline-block">
                      <span className="material-symbols-outlined text-lg text-outline">open_in_new</span>
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="bg-surface-container-low px-6 py-4 flex items-center justify-between">
          <p className="text-xs text-on-surface-variant font-medium">
            Showing <span className="font-bold">{tickets.length}</span> of <span className="font-bold">{total}</span> records
          </p>
          <Pagination page={page} total={total} pageSize={10} onPageChange={setPage} />
        </div>
      </div>

      {/* Bottom Stats */}
      <div className="mt-12 grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="bg-gradient-to-br from-red-900 to-primary-container p-6 rounded-2xl text-white shadow-xl flex flex-col justify-between">
          <div>
            <h3 className="font-headline text-lg font-bold mb-1">Queue Health</h3>
            <p className="text-xs text-primary-fixed opacity-80 mb-6">Current throughput vs. weekly average</p>
            <div className="flex items-end gap-3 h-24">
              {[40, 65, 50, 85, 70].map((h, i) => (
                <div key={i} className="bg-white/20 w-full rounded-t" style={{ height: `${h}%` }} />
              ))}
            </div>
          </div>
          <div className="flex justify-between items-center mt-6">
            <span className="text-2xl font-black">+14%</span>
            <span className="text-[10px] uppercase font-bold tracking-widest bg-white/10 px-2 py-1 rounded">Efficiency</span>
          </div>
        </div>

        <div className="lg:col-span-2 bg-surface-container-lowest p-6 rounded-2xl border border-outline-variant/15 flex flex-col justify-between">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h3 className="font-headline text-lg font-bold text-on-surface">Urgent Focus</h3>
              <p className="text-xs text-on-surface-variant">Recommended priorities for your current session</p>
            </div>
            <span className="material-symbols-outlined text-primary">auto_awesome</span>
          </div>
          <div className="flex flex-col gap-4">
            <div className="flex items-center justify-between p-3 bg-surface-container-low rounded-xl cursor-pointer hover:bg-surface-container-high transition-colors"
              onClick={() => navigate('/tickets?status=open&priority=critical')}>
              <div className="flex items-center gap-3">
                <span className="material-symbols-outlined text-error">warning</span>
                <div>
                  <p className="text-sm font-bold">Unassigned Critical Tickets</p>
                  <p className="text-[10px] text-outline">Incidents needing immediate triage</p>
                </div>
              </div>
              <span className="text-xs font-bold text-primary">View All →</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-surface-container-low rounded-xl cursor-pointer hover:bg-surface-container-high transition-colors"
              onClick={() => navigate('/approvals')}>
              <div className="flex items-center gap-3">
                <span className="material-symbols-outlined text-secondary">schedule</span>
                <div>
                  <p className="text-sm font-bold">Pending Approvals</p>
                  <p className="text-[10px] text-outline">Tickets awaiting approval decision</p>
                </div>
              </div>
              <span className="text-xs font-bold text-primary">Review →</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
