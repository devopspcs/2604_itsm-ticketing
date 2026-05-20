import { NavLink, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import type { RootState } from '../../store'

interface NavItem {
  to: string
  icon: string
  label: string
  roles?: string[] // if undefined, visible to all
  children?: NavItem[]
}

const navItems: NavItem[] = [
  { to: '/dashboard', icon: 'dashboard', label: 'Dashboard' },
  { to: '/tickets', icon: 'confirmation_number', label: 'My Tickets', roles: ['user'] },
  { to: '/tickets', icon: 'confirmation_number', label: 'All Tickets', roles: ['agent', 'admin', 'approver'] },
  { to: '/kanban', icon: 'view_kanban', label: 'Kanban Board', roles: ['agent', 'admin', 'approver'] },
  { to: '/approvals', icon: 'fact_check', label: 'Approvals', roles: ['admin', 'approver'] },
  { to: '/activity-logs', icon: 'history_edu', label: 'Activity Logs', roles: ['agent', 'admin', 'approver'] },
  {
    to: '/users', icon: 'manage_accounts', label: 'User Management', roles: ['admin'],
    children: [
      { to: '/acl', icon: 'shield_person', label: 'Access Control' },
    ],
  },
  { to: '/webhooks', icon: 'webhook', label: 'Webhooks', roles: ['admin'] },
  { to: '/org-structure', icon: 'account_tree', label: 'Org Structure', roles: ['admin'] },
  { to: '/org-chart', icon: 'groups', label: 'Org Chart', roles: ['admin', 'approver'] },
  { to: '/settings', icon: 'settings', label: 'System Settings', roles: ['admin'] },
  { to: '/app-management', icon: 'apps', label: 'App Management', roles: ['admin'] },
]

export function Sidebar() {
  const role = useSelector((s: RootState) => s.auth.role) ?? 'user'
  const location = useLocation()
  const [apiStatus, setApiStatus] = useState<'ok' | 'error' | 'checking'>('checking')

  useEffect(() => {
    const check = () => {
      fetch('/health')
        .then(res => res.ok ? res.json() : Promise.reject())
        .then(data => setApiStatus(data?.status === 'ok' ? 'ok' : 'error'))
        .catch(() => setApiStatus('error'))
    }
    check()
    const interval = setInterval(check, 15000)
    return () => clearInterval(interval)
  }, [])

  const visible = navItems.filter(item =>
    !item.roles || item.roles.includes(role)
  )

  return (
    <aside className="h-screen w-64 fixed left-0 top-0 pt-16 bg-slate-50 flex flex-col gap-2 p-4 z-40 hidden md:flex">
      <div className="mb-6 px-2">
        <h2 className="text-lg font-black text-accent-900 font-headline">Service Console</h2>
        <p className="text-[10px] uppercase tracking-widest text-slate-500 font-bold">Enterprise Management</p>
      </div>

      <nav className="flex flex-col gap-1 flex-1">
        {visible.map((item) => {
          const isParentActive = location.pathname === item.to || item.children?.some(c => location.pathname === c.to)

          return (
            <div key={item.to + item.label}>
              <NavLink
                to={item.to}
                className={({ isActive }) =>
                  (isActive || isParentActive)
                    ? 'flex items-center gap-3 px-3 py-2.5 text-accent-900 font-bold bg-white rounded-lg text-sm shadow-sm'
                    : 'flex items-center gap-3 px-3 py-2.5 text-slate-600 hover:text-accent-700 hover:translate-x-1 transition-all text-sm font-medium'
                }
              >
                {({ isActive }) => (
                  <>
                    <span className="material-symbols-outlined"
                      style={(isActive || isParentActive) ? { fontVariationSettings: "'FILL' 1" } : {}}>
                      {item.icon}
                    </span>
                    <span>{item.label}</span>
                  </>
                )}
              </NavLink>

              {/* Sub-items */}
              {item.children && isParentActive && (
                <div className="ml-6 mt-1 flex flex-col gap-0.5 border-l-2 border-outline-variant/20 pl-3">
                  {item.children.map((child) => (
                    <NavLink
                      key={child.to}
                      to={child.to}
                      className={({ isActive }) =>
                        isActive
                          ? 'flex items-center gap-2 px-2 py-2 text-primary font-bold text-xs rounded-lg bg-primary-fixed/30'
                          : 'flex items-center gap-2 px-2 py-2 text-slate-500 hover:text-primary text-xs rounded-lg hover:bg-white/60 transition-all'
                      }
                    >
                      {({ isActive }) => (
                        <>
                          <span className="material-symbols-outlined text-base"
                            style={isActive ? { fontVariationSettings: "'FILL' 1" } : {}}>
                            {child.icon}
                          </span>
                          <span>{child.label}</span>
                        </>
                      )}
                    </NavLink>
                  ))}
                </div>
              )}
            </div>
          )
        })}
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
