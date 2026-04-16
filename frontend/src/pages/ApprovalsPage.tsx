import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { ticketService } from '../services/ticket.service'
import { approvalService } from '../services/approval.service'
import type { Ticket } from '../types'
import type { RootState } from '../store'

export function ApprovalsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [loading, setLoading] = useState(true)
  const [actionMsg, setActionMsg] = useState('')
  const role = useSelector((s: RootState) => s.auth.role) ?? 'user'

  const load = () => {
    setLoading(true)
    if (role === 'user') {
      // User only sees their own tickets that can be submitted
      Promise.all([
        ticketService.list({ status: 'open', page: 1, page_size: 50 }),
        ticketService.list({ status: 'in_progress', page: 1, page_size: 50 }),
        ticketService.list({ status: 'pending_approval', page: 1, page_size: 50 }),
      ]).then(([openRes, inProgressRes, pendingRes]) => {
        setTickets([
          ...(pendingRes.data.tickets ?? []),
          ...(openRes.data.tickets ?? []),
          ...(inProgressRes.data.tickets ?? []),
        ])
      }).finally(() => setLoading(false))
    } else {
      Promise.all([
        ticketService.list({ status: 'pending_approval', page: 1, page_size: 50 }),
        ticketService.list({ status: 'open', page: 1, page_size: 50 }),
        ticketService.list({ status: 'in_progress', page: 1, page_size: 50 }),
      ]).then(([pendingRes, openRes, inProgressRes]) => {
        setTickets([
          ...(pendingRes.data.tickets ?? []),
          ...(openRes.data.tickets ?? []),
          ...(inProgressRes.data.tickets ?? []),
        ])
      }).finally(() => setLoading(false))
    }
  }

  useEffect(() => { load() }, [])

  const showMsg = (msg: string) => {
    setActionMsg(msg)
    setTimeout(() => setActionMsg(''), 3000)
  }

  const handleSubmit = async (ticketId: string) => {
    try {
      await ticketService.submit(ticketId)
      load()
      showMsg('Submitted for approval')
    } catch { showMsg('Failed to submit') }
  }

  const handleDecide = async (ticketId: string, decision: 'approved' | 'rejected') => {
    try {
      await approvalService.decide(ticketId, decision)
      load()
      showMsg(`Ticket ${decision}`)
    } catch { showMsg('Failed to record decision') }
  }

  const pendingCount = tickets.filter(t => t.status === 'pending_approval').length

  return (
    <div className="max-w-6xl mx-auto p-8">
      {actionMsg && (
        <div className="fixed top-20 right-6 z-50 px-5 py-3 rounded-xl shadow-lg font-semibold text-sm bg-emerald-100 text-emerald-800 flex items-center gap-2">
          <span className="material-symbols-outlined text-sm">check_circle</span>
          {actionMsg}
        </div>
      )}

      <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-10">
        <div>
          <h1 className="text-4xl font-extrabold tracking-tight text-on-surface mb-2">Approvals Inbox</h1>
          <p className="text-on-surface-variant text-lg max-w-2xl">
            Manage and review pending infrastructure changes and access requests.
            {pendingCount > 0 && <> You have <span className="font-bold text-primary">{pendingCount} pending approvals</span>.</>}
          </p>
        </div>
        <div className="flex gap-3">
          <div className="bg-surface-container-low px-4 py-2 rounded-full border border-outline-variant/10 flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-amber-500" />
            <span className="text-xs font-bold text-on-surface-variant uppercase tracking-tighter">{pendingCount} Pending</span>
          </div>
          <div className="bg-surface-container-low px-4 py-2 rounded-full border border-outline-variant/10 flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-red-500" />
            <span className="text-xs font-bold text-on-surface-variant uppercase tracking-tighter">{tickets.length - pendingCount} Open/In Progress</span>
          </div>
        </div>
      </div>

      <div className="bg-surface-container-lowest rounded-xl overflow-hidden shadow-sm border border-outline-variant/10">
        <div className="grid grid-cols-12 gap-4 px-6 py-4 bg-surface-container-low border-b border-outline-variant/5">
          <div className="col-span-4 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Request &amp; Details</div>
          <div className="col-span-2 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Type</div>
          <div className="col-span-2 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Status</div>
          <div className="col-span-2 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Urgency</div>
          <div className="col-span-2 text-xs font-bold text-on-surface-variant uppercase tracking-wider text-right">Actions</div>
        </div>

        <div className="divide-y divide-surface-container">
          {loading ? (
            <div className="px-6 py-12 text-center text-on-surface-variant">
              <span className="material-symbols-outlined animate-spin block mb-2">refresh</span>
              Loading...
            </div>
          ) : tickets.length === 0 ? (
            <div className="px-6 py-16 text-center">
              <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 block mb-3">inbox</span>
              <p className="text-on-surface-variant font-medium">No tickets found</p>
              <p className="text-sm text-on-surface-variant/60 mt-1">Create a ticket first from the Tickets page</p>
            </div>
          ) : tickets.map((t) => (
            <div key={t.id} className="grid grid-cols-12 gap-4 px-6 py-5 items-center hover:bg-slate-50 transition-colors group">
              <div className="col-span-4 flex items-center gap-3">
                <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 ${t.status === 'pending_approval' ? 'bg-amber-100' : 'bg-primary-fixed'}`}>
                  <span className={`material-symbols-outlined ${t.status === 'pending_approval' ? 'text-amber-600' : 'text-primary'}`}>
                    {t.type === 'incident' ? 'report' : t.type === 'change_request' ? 'published_with_changes' : 'contact_support'}
                  </span>
                </div>
                <div className="min-w-0">
                  <h3 className="text-sm font-bold text-on-surface leading-tight truncate">{t.title}</h3>
                  <p className="text-xs text-on-surface-variant mt-0.5">
                    {t.type === 'incident' ? 'INC' : t.type === 'change_request' ? 'CHG' : 'REQ'}-{t.id.slice(0, 6).toUpperCase()}
                  </p>
                </div>
              </div>
              <div className="col-span-2">
                <span className="text-xs px-2 py-1 bg-surface-container-highest text-on-surface-variant rounded font-semibold capitalize">
                  {t.type.replace(/_/g, ' ')}
                </span>
              </div>
              <div className="col-span-2">
                <span className={`text-xs px-2 py-1 rounded font-bold capitalize ${
                  t.status === 'pending_approval' ? 'bg-amber-100 text-amber-700' :
                  t.status === 'in_progress' ? 'bg-red-100 text-red-700' :
                  'bg-surface-container-high text-on-surface-variant'
                }`}>{t.status.replace(/_/g, ' ')}</span>
              </div>
              <div className="col-span-2">
                <div className="flex items-center gap-1.5">
                  <span className={`w-1.5 h-6 rounded-full ${
                    t.priority === 'critical' || t.priority === 'high' ? 'bg-error' :
                    t.priority === 'medium' ? 'bg-amber-500' : 'bg-red-400'
                  }`} />
                  <span className={`text-xs font-bold uppercase ${
                    t.priority === 'critical' || t.priority === 'high' ? 'text-error' :
                    t.priority === 'medium' ? 'text-amber-600' : 'text-red-500'
                  }`}>{t.priority}</span>
                </div>
              </div>
              <div className="col-span-2 flex justify-end gap-2">
                <Link to={`/tickets/${t.id}`} className="px-3 py-1.5 text-xs font-bold text-on-surface-variant hover:bg-surface-container-high rounded-xl transition-all">
                  View
                </Link>
                {t.status === 'pending_approval' && (role === 'admin' || role === 'approver') ? (
                  <>
                    <button onClick={() => handleDecide(t.id, 'rejected')}
                      className="px-3 py-1.5 text-xs font-bold text-error hover:bg-error-container rounded-xl transition-all">
                      Reject
                    </button>
                    <button onClick={() => handleDecide(t.id, 'approved')}
                      className="px-3 py-1.5 text-xs font-bold bg-primary text-white rounded-xl shadow-sm hover:shadow-md transition-all">
                      Approve
                    </button>
                  </>
                ) : t.status !== 'pending_approval' && t.status !== 'approved' && t.status !== 'rejected' && t.status !== 'done' ? (
                  <button onClick={() => handleSubmit(t.id)}
                    className="px-3 py-1.5 text-xs font-bold bg-amber-500 text-white rounded-xl shadow-sm hover:shadow-md transition-all">
                    Submit
                  </button>
                ) : null}
              </div>
            </div>
          ))}
        </div>

        <div className="px-6 py-4 bg-surface-container-lowest border-t border-outline-variant/10">
          <span className="text-xs text-on-surface-variant">
            Showing <span className="font-bold">{tickets.length}</span> tickets
          </span>
        </div>
      </div>

      <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-gradient-to-br from-primary to-primary-container p-6 rounded-2xl text-white shadow-xl flex flex-col justify-between h-48">
          <div className="flex justify-between items-start">
            <span className="material-symbols-outlined text-3xl">avg_time</span>
            <span className="text-[10px] font-bold uppercase tracking-widest bg-white/20 px-2 py-1 rounded">Efficiency</span>
          </div>
          <div>
            <p className="text-4xl font-extrabold mb-1">1.2h</p>
            <p className="text-xs text-primary-fixed/80 font-medium">Average approval time this month</p>
          </div>
        </div>
        <div className="bg-white/80 p-6 rounded-2xl shadow-sm flex flex-col justify-between h-48 border border-outline-variant/10">
          <div className="flex justify-between items-start">
            <span className="material-symbols-outlined text-3xl text-red-600">trending_up</span>
            <span className="text-[10px] font-bold uppercase tracking-widest bg-red-100 text-red-700 px-2 py-1 rounded">Volume</span>
          </div>
          <div>
            <p className="text-4xl font-extrabold text-on-surface mb-1">{pendingCount}</p>
            <p className="text-xs text-on-surface-variant font-medium">Pending requests requiring action</p>
          </div>
        </div>
        <div className="bg-surface-container-lowest p-6 rounded-2xl shadow-sm flex flex-col justify-between h-48 border border-outline-variant/10">
          <div className="flex justify-between items-start">
            <span className="material-symbols-outlined text-3xl text-emerald-600">verified</span>
            <span className="text-[10px] font-bold uppercase tracking-widest bg-slate-100 text-slate-600 px-2 py-1 rounded">Compliance</span>
          </div>
          <div>
            <p className="text-4xl font-extrabold text-on-surface mb-1">98%</p>
            <p className="text-xs text-on-surface-variant font-medium">SLA compliance for critical tasks</p>
          </div>
        </div>
      </div>
    </div>
  )
}
