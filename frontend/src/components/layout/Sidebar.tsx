import { NavLink } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import type { RootState } from '../../store'

interface NavItem {
  to: string
  icon: string
  label: string
  roles?: string[] // if undefined, visible to all
}

const navItems: NavItem[] = [
  { to: '/dashboard', icon: 'dashboard', label: 'Dashboard' },
  { to: '/tickets', icon: 'confirmation_number', label: 'All Tickets' },
  { to: '/approvals', icon: 'fact_check', label: 'Approvals', roles: ['admin', 'approver'] },
  { to: '/activity-logs', icon: 'history_edu', label: 'Activity Logs' },
  { to: '/users', icon: 'manage_accounts', label: 'User Management', roles: ['admin'] },
  { to: '/webhooks', icon: 'webhook', label: 'Webhooks', roles: ['admin'] },
  { to: '/org-structure', icon: 'account_tree', label: 'Org Structure', roles: ['admin'] },
  { to: '/settings', icon: 'settings', label: 'System Settings', roles: ['admin'] },
]

export function Sidebar() {
  const role = useSelector((s: RootState) => s.auth.role) ?? 'user'
  const [apiStatus, setApiStatus] = useState<'ok' | 'error' | 'checking'>('checking')

  useEffect(() => {
    const checkHealth = () => {
      fetch('/api/v1/../health', { method: 'GET' })
        .then(res => {
          if (res.ok) return res.json()
          throw new Error('not ok')
        })
        .then(data => setApiStatus(data.status === 'ok' ? 'ok' : 'error'))
        .catch(() => setApiStatus('error'))
    }
    // Use direct health endpoint
    const check = () => {
      fetch('/health')
        .then(res => res.ok ? res.json() : Promise.reject())
        .then(data => setApiStatus(data?.status === 'ok' ? 'ok' : 'error'))
        .catch(() => setApiStatus('error'))
    }
    check()
    const interval = setInterval(check, 15000) // check every 15s
    return () => clearInterval(interval)
  }, [])

  const visible = navItems.filter(item =>
    !item.roles || item.roles.includes(role)
  )

  return (
    <aside className="h-screen w-64 fixed left-0 top-0 pt-16 bg-slate-50 flex flex-col gap-2 p-4 z-40 hidden md:flex">
      <div className="mb-6 px-2">
        <h2 className="text-lg font-black text-red-900 font-headline">Service Console</h2>
        <p className="text-[10px] uppercase tracking-widest text-slate-500 font-bold">Enterprise Management</p>
      </div>

      <nav className="flex flex-col gap-1 flex-1">
        {visible.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              isActive
                ? 'flex items-center gap-3 px-3 py-2.5 text-red-900 font-bold bg-white rounded-lg text-sm shadow-sm'
                : 'flex items-center gap-3 px-3 py-2.5 text-slate-600 hover:text-red-700 hover:translate-x-1 transition-all text-sm font-medium'
            }
          >
            {({ isActive }) => (
              <>
                <span className="material-symbols-outlined"
                  style={isActive ? { fontVariationSettings: "'FILL' 1" } : {}}>
                  {item.icon}
                </span>
                <span>{item.label}</span>
              </>
            )}
          </NavLink>
        ))}
      </nav>

      <div className="mt-auto p-4 bg-primary-container/10 rounded-xl">
        <p className="text-xs font-semibold text-primary mb-1">System Status</p>
        <div className="flex items-center gap-2">
          <span className={`h-2 w-2 rounded-full ${
            apiStatus === 'ok' ? 'bg-emerald-500 animate-pulse' :
            apiStatus === 'error' ? 'bg-red-500' :
            'bg-amber-500 animate-pulse'
          }`} />
          <span className={`text-xs ${apiStatus === 'error' ? 'text-red-600 font-semibold' : 'text-on-surface-variant'}`}>
            {apiStatus === 'ok' ? 'All Systems Operational' :
             apiStatus === 'error' ? 'API Unreachable' :
             'Checking...'}
          </span>
        </div>
        <p className="text-[10px] text-on-surface-variant/60 mt-1 capitalize">Role: {role}</p>
      </div>
    </aside>
  )
}
