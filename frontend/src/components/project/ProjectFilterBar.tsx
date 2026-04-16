import { useEffect, useState } from 'react'
import api from '../../services/api'
import type { User } from '../../types'

interface FilterState {
  search: string
  assignee: string
  dueDateFrom: string
  dueDateTo: string
}

interface ProjectFilterBarProps {
  filters: FilterState
  onChange: (filters: FilterState) => void
}

export function ProjectFilterBar({ filters, onChange }: ProjectFilterBarProps) {
  const [users, setUsers] = useState<User[]>([])

  useEffect(() => {
    api.get<User[]>('/users/list').then(res => setUsers(res.data ?? [])).catch(() => {})
  }, [])

  const activeCount = [filters.search, filters.assignee, filters.dueDateFrom, filters.dueDateTo].filter(Boolean).length

  const clearAll = () => onChange({ search: '', assignee: '', dueDateFrom: '', dueDateTo: '' })

  return (
    <div className="flex items-center gap-3 flex-wrap">
      <div className="relative">
        <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-[16px]">search</span>
        <input
          value={filters.search}
          onChange={e => onChange({ ...filters, search: e.target.value })}
          placeholder="Cari record..."
          className="pl-9 pr-3 py-1.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none w-48"
        />
      </div>

      <select
        value={filters.assignee}
        onChange={e => onChange({ ...filters, assignee: e.target.value })}
        className="px-3 py-1.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none"
      >
        <option value="">Semua Assignee</option>
        {users.map(u => (
          <option key={u.id} value={u.id}>{u.full_name}</option>
        ))}
      </select>

      <input
        type="date"
        value={filters.dueDateFrom}
        onChange={e => onChange({ ...filters, dueDateFrom: e.target.value })}
        className="px-3 py-1.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none"
        placeholder="Dari"
      />
      <input
        type="date"
        value={filters.dueDateTo}
        onChange={e => onChange({ ...filters, dueDateTo: e.target.value })}
        className="px-3 py-1.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none"
        placeholder="Sampai"
      />

      {activeCount > 0 && (
        <button onClick={clearAll} className="flex items-center gap-1 px-3 py-1.5 text-xs font-bold text-error hover:bg-error-container/20 rounded-xl transition-colors">
          <span className="material-symbols-outlined text-[14px]">close</span>
          Hapus filter ({activeCount})
        </button>
      )}
    </div>
  )
}
