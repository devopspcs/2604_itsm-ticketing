import { useEffect, useState } from 'react'
import api from '../services/api'
import type { User } from '../types'

export function ACLDashboardPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filterStatus, setFilterStatus] = useState<'all' | 'active' | 'inactive'>('all')
  const [filterRole, setFilterRole] = useState<string>('all')
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set())
  const [actionMsg, setActionMsg] = useState('')
  const [confirmAction, setConfirmAction] = useState<{ type: 'deactivate' | 'activate'; ids: string[]; names: string[] } | null>(null)
  const [processing, setProcessing] = useState(false)

  const fetchUsers = () => {
    setLoading(true)
    api.get<User[]>('/users')
      .then((r) => setUsers(r.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [])

  const showMsg = (msg: string) => { setActionMsg(msg); setTimeout(() => setActionMsg(''), 4000) }

  // Stats
  const totalUsers = users.length
  const activeUsers = users.filter(u => u.is_active).length
  const inactiveUsers = users.filter(u => !u.is_active).length

  // Filtered users
  const filtered = users.filter(u => {
    if (filterStatus === 'active' && !u.is_active) return false
    if (filterStatus === 'inactive' && u.is_active) return false
    if (filterRole !== 'all' && u.role !== filterRole) return false
    if (search) {
      const q = search.toLowerCase()
      return u.full_name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q)
    }
    return true
  })

  // Selection
  const toggleSelect = (id: string) => {
    setSelectedIds(prev => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  const toggleSelectAll = () => {
    if (selectedIds.size === filtered.length) {
      setSelectedIds(new Set())
    } else {
      setSelectedIds(new Set(filtered.map(u => u.id)))
    }
  }

  // Actions
  const handleBulkDeactivate = () => {
    const ids = Array.from(selectedIds)
    const names = users.filter(u => ids.includes(u.id)).map(u => u.full_name)
    setConfirmAction({ type: 'deactivate', ids, names })
  }

  const handleBulkActivate = () => {
    const ids = Array.from(selectedIds)
    const names = users.filter(u => ids.includes(u.id)).map(u => u.full_name)
    setConfirmAction({ type: 'activate', ids, names })
  }

  const handleSingleAction = (user: User, type: 'deactivate' | 'activate') => {
    setConfirmAction({ type, ids: [user.id], names: [user.full_name] })
  }

  const executeAction = async () => {
    if (!confirmAction) return
    setProcessing(true)
    const endpoint = confirmAction.type === 'deactivate' ? 'deactivate' : 'activate'
    let success = 0
    for (const id of confirmAction.ids) {
      try {
        await api.patch(`/users/${id}/${endpoint}`)
        success++
      } catch { /* skip failed */ }
    }
    setProcessing(false)
    setConfirmAction(null)
    setSelectedIds(new Set())
    fetchUsers()
    showMsg(`${success} user(s) ${confirmAction.type === 'deactivate' ? 'deactivated' : 'activated'} successfully`)
  }

  return (
    <div className="max-w-7xl mx-auto p-6 md:p-10">
      {/* Toast */}
      {actionMsg && (
        <div className="fixed top-20 right-6 z-50 px-5 py-3 rounded-xl shadow-lg font-semibold text-sm bg-emerald-100 text-emerald-800 flex items-center gap-2 animate-in fade-in">
          <span className="material-symbols-outlined text-sm">check_circle</span>
          {actionMsg}
        </div>
      )}

      {/* Confirmation Modal */}
      {confirmAction && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 overflow-hidden">
            <div className="p-6">
              <div className={`w-12 h-12 rounded-full flex items-center justify-center mb-4 ${
                confirmAction.type === 'deactivate' ? 'bg-red-100' : 'bg-emerald-100'
              }`}>
                <span className={`material-symbols-outlined text-2xl ${
                  confirmAction.type === 'deactivate' ? 'text-red-600' : 'text-emerald-600'
                }`}>
                  {confirmAction.type === 'deactivate' ? 'person_off' : 'person_check'}
                </span>
              </div>
              <h3 className="text-lg font-bold text-on-surface mb-2">
                {confirmAction.type === 'deactivate' ? 'Revoke Access' : 'Grant Access'}
              </h3>
              <p className="text-sm text-on-surface-variant mb-4">
                {confirmAction.type === 'deactivate'
                  ? 'User berikut akan kehilangan akses ke sistem. Mereka tidak bisa login sampai diaktifkan kembali.'
                  : 'User berikut akan mendapatkan kembali akses ke sistem.'}
              </p>
              <div className="bg-surface-container-low rounded-xl p-3 max-h-32 overflow-y-auto mb-4">
                {confirmAction.names.map((name, i) => (
                  <div key={i} className="flex items-center gap-2 py-1">
                    <span className={`w-2 h-2 rounded-full ${confirmAction.type === 'deactivate' ? 'bg-red-500' : 'bg-emerald-500'}`} />
                    <span className="text-sm font-medium text-on-surface">{name}</span>
                  </div>
                ))}
              </div>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => setConfirmAction(null)}
                  disabled={processing}
                  className="px-5 py-2.5 rounded-xl font-semibold text-on-surface-variant hover:bg-surface-container-high transition-colors"
                >
                  Batal
                </button>
                <button
                  onClick={executeAction}
                  disabled={processing}
                  className={`px-6 py-2.5 rounded-xl font-bold text-white transition-all disabled:opacity-60 ${
                    confirmAction.type === 'deactivate'
                      ? 'bg-red-600 hover:bg-red-700'
                      : 'bg-emerald-600 hover:bg-emerald-700'
                  }`}
                >
                  {processing ? 'Processing...' : confirmAction.type === 'deactivate' ? 'Revoke Access' : 'Grant Access'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-extrabold tracking-tight text-on-surface mb-2">Access Control</h1>
        <p className="text-on-surface-variant">Kelola akses user ke sistem. Cabut akses user yang resign atau tidak aktif.</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-8">
        <div className="bg-surface-container-lowest p-6 rounded-xl shadow-sm border-l-4 border-primary">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Total Users</p>
          <p className="text-3xl font-black text-on-surface">{totalUsers}</p>
        </div>
        <div className="bg-surface-container-lowest p-6 rounded-xl shadow-sm border-l-4 border-emerald-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Active</p>
          <div className="flex items-baseline gap-2">
            <p className="text-3xl font-black text-emerald-600">{activeUsers}</p>
            <span className="text-xs text-on-surface-variant">({totalUsers > 0 ? Math.round((activeUsers / totalUsers) * 100) : 0}%)</span>
          </div>
        </div>
        <div className="bg-surface-container-lowest p-6 rounded-xl shadow-sm border-l-4 border-red-500">
          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Inactive / Revoked</p>
          <div className="flex items-baseline gap-2">
            <p className="text-3xl font-black text-red-600">{inactiveUsers}</p>
            <span className="text-xs text-on-surface-variant">({totalUsers > 0 ? Math.round((inactiveUsers / totalUsers) * 100) : 0}%)</span>
          </div>
        </div>
      </div>

      {/* Toolbar */}
      <div className="bg-surface-container-lowest rounded-xl shadow-sm p-4 mb-6">
        <div className="flex flex-wrap items-center gap-4">
          {/* Search */}
          <div className="flex-1 min-w-[200px] relative">
            <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg">search</span>
            <input
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Cari nama atau email..."
              className="w-full pl-10 pr-4 py-2.5 bg-surface-container-high border-none rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          {/* Filter Status */}
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value as 'all' | 'active' | 'inactive')}
            className="bg-surface-container-high border-none rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-primary appearance-none"
          >
            <option value="all">Semua Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>

          {/* Filter Role */}
          <select
            value={filterRole}
            onChange={(e) => setFilterRole(e.target.value)}
            className="bg-surface-container-high border-none rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-primary appearance-none"
          >
            <option value="all">Semua Role</option>
            <option value="user">User</option>
            <option value="agent">Agent</option>
            <option value="approver">Approver</option>
            <option value="admin">Admin</option>
          </select>

          {/* Bulk Actions */}
          {selectedIds.size > 0 && (
            <div className="flex items-center gap-2 ml-auto">
              <span className="text-xs font-bold text-on-surface-variant">{selectedIds.size} selected</span>
              <button
                onClick={handleBulkDeactivate}
                className="px-4 py-2 bg-red-600 text-white rounded-xl text-xs font-bold hover:bg-red-700 transition-colors flex items-center gap-1"
              >
                <span className="material-symbols-outlined text-sm">person_off</span>
                Revoke Access
              </button>
              <button
                onClick={handleBulkActivate}
                className="px-4 py-2 bg-emerald-600 text-white rounded-xl text-xs font-bold hover:bg-emerald-700 transition-colors flex items-center gap-1"
              >
                <span className="material-symbols-outlined text-sm">person_check</span>
                Grant Access
              </button>
            </div>
          )}
        </div>
      </div>

      {/* User Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-sm overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-20">
            <span className="material-symbols-outlined animate-spin text-on-surface-variant">refresh</span>
          </div>
        ) : (
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-surface-container-low/50">
                <th className="px-4 py-4 w-10">
                  <input
                    type="checkbox"
                    checked={selectedIds.size === filtered.length && filtered.length > 0}
                    onChange={toggleSelectAll}
                    className="w-4 h-4 rounded border-outline-variant text-primary focus:ring-primary"
                  />
                </th>
                <th className="px-4 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">User</th>
                <th className="px-4 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Role</th>
                <th className="px-4 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Status</th>
                <th className="px-4 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Registered</th>
                <th className="px-4 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-surface-container-low">
              {filtered.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center text-on-surface-variant">
                    <span className="material-symbols-outlined text-4xl block mb-2 opacity-30">group_off</span>
                    Tidak ada user ditemukan
                  </td>
                </tr>
              ) : (
                filtered.map((u) => (
                  <tr key={u.id} className={`hover:bg-surface-container-low/30 transition-colors ${!u.is_active ? 'opacity-60' : ''}`}>
                    <td className="px-4 py-4">
                      <input
                        type="checkbox"
                        checked={selectedIds.has(u.id)}
                        onChange={() => toggleSelect(u.id)}
                        className="w-4 h-4 rounded border-outline-variant text-primary focus:ring-primary"
                      />
                    </td>
                    <td className="px-4 py-4">
                      <div className="flex items-center gap-3">
                        <div className={`w-9 h-9 rounded-full flex items-center justify-center font-bold text-sm ${
                          u.is_active ? 'bg-primary-fixed text-primary' : 'bg-red-100 text-red-600'
                        }`}>
                          {u.full_name.charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <p className="text-sm font-semibold text-on-surface">{u.full_name}</p>
                          <p className="text-xs text-on-surface-variant">{u.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-4">
                      <span className={`text-xs px-2.5 py-1 rounded-full font-bold capitalize ${
                        u.role === 'admin' ? 'bg-primary-fixed text-primary' :
                        u.role === 'agent' ? 'bg-blue-100 text-blue-700' :
                        u.role === 'approver' ? 'bg-amber-100 text-amber-700' :
                        'bg-surface-container-high text-on-surface-variant'
                      }`}>{u.role}</span>
                    </td>
                    <td className="px-4 py-4">
                      <div className="flex items-center gap-2">
                        <span className={`w-2.5 h-2.5 rounded-full ${u.is_active ? 'bg-emerald-500' : 'bg-red-500'}`} />
                        <span className={`text-xs font-bold ${u.is_active ? 'text-emerald-700' : 'text-red-700'}`}>
                          {u.is_active ? 'Active' : 'Revoked'}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-4 text-xs text-on-surface-variant">
                      {new Date(u.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
                    </td>
                    <td className="px-4 py-4 text-right">
                      {u.is_active ? (
                        <button
                          onClick={() => handleSingleAction(u, 'deactivate')}
                          className="inline-flex items-center gap-1 px-3 py-1.5 text-xs font-bold text-red-700 bg-red-50 hover:bg-red-100 rounded-lg transition-colors"
                        >
                          <span className="material-symbols-outlined text-sm">block</span>
                          Revoke
                        </button>
                      ) : (
                        <button
                          onClick={() => handleSingleAction(u, 'activate')}
                          className="inline-flex items-center gap-1 px-3 py-1.5 text-xs font-bold text-emerald-700 bg-emerald-50 hover:bg-emerald-100 rounded-lg transition-colors"
                        >
                          <span className="material-symbols-outlined text-sm">check_circle</span>
                          Grant
                        </button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        )}
      </div>

      {/* Info Footer */}
      <div className="mt-6 p-4 bg-surface-container-low rounded-xl flex items-start gap-3">
        <span className="material-symbols-outlined text-primary text-lg mt-0.5">info</span>
        <div className="text-xs text-on-surface-variant leading-relaxed">
          <p className="font-bold text-on-surface mb-1">Bagaimana access control bekerja?</p>
          <p>User yang di-revoke tidak bisa login ke sistem. Token yang sudah ada akan expired secara otomatis dan tidak bisa di-refresh. Untuk mengembalikan akses, klik "Grant Access".</p>
        </div>
      </div>
    </div>
  )
}
