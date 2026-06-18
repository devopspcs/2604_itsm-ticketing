import { NavLink, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import type { RootState } from '../../store'
import logoPcs from '../../assets/logo-pcs.png'

interface NavItem {
  to: string
  icon: string
  label: string
  roles?: string[]
  children?: NavItem[]
}

interface NavSection {
  title?: string
  icon?: string
  roles?: string[]
  collapsible?: boolean
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    items: [
      { to: '/dashboard', icon: 'dashboard', label: 'Dashboard' },
      { to: '/tickets', icon: 'confirmation_number', label: 'My Tickets', roles: ['user'] },
      { to: '/tickets', icon: 'confirmation_number', label: 'All Tickets', roles: ['agent', 'admin', 'approver'] },
      { to: '/kanban', icon: 'view_kanban', label: 'Kanban Board', roles: ['agent', 'admin', 'approver'] },
      { to: '/approvals', icon: 'fact_check', label: 'Approvals', roles: ['admin', 'approver'] },
      { to: '/activity-logs', icon: 'history_edu', label: 'Activity Logs', roles: ['agent', 'admin', 'approver'] },
    ],
  },
  {
    title: 'Data Library',
    icon: 'database',
    roles: ['admin', 'approver', 'agent'],
    collapsible: true,
    items: [
      { to: '/data-library/people', icon: 'people', label: 'People', roles: ['admin', 'approver'] },
      { to: '/data-library/assets', icon: 'devices', label: 'Asset Register', roles: ['admin', 'approver', 'agent'] },
      { to: '/data-library/access', icon: 'shield_person', label: 'Access', roles: ['admin', 'approver'] },
      { to: '/data-library/vendors', icon: 'store', label: 'Vendors', roles: ['admin'] },
      { to: '/data-library/change-management', icon: 'change_circle', label: 'Change Management', roles: ['admin', 'approver'] },
      { to: '/data-library/incidents', icon: 'warning', label: 'Incidents', roles: ['admin', 'approver', 'agent'] },
      { to: '/data-library/knowledge-base', icon: 'menu_book', label: 'Knowledge Base', roles: ['admin', 'approver', 'agent'] },
      { to: '/data-library/sla-templates', icon: 'timer', label: 'SLA Templates', roles: ['admin'] },
      { to: '/data-library/categories', icon: 'category', label: 'Categories', roles: ['admin'] },
    ],
  },
  {
    title: 'Administration',
    icon: 'admin_panel_settings',
    roles: ['admin'],
    collapsible: true,
    items: [
      { to: '/app-management', icon: 'apps', label: 'App Management', roles: ['admin'] },
      { to: '/webhooks', icon: 'webhook', label: 'Webhooks', roles: ['admin'] },
    ],
  },
]

export function Sidebar() {
  const role = useSelector((s: RootState) => s.auth.role) ?? 'user'
  const location = useLocation()
  const [apiStatus, setApiStatus] = useState<'ok' | 'error' | 'checking'>('checking')
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({})

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

  // Auto-expand sections that contain the active route
  useEffect(() => {
    navSections.forEach(section => {
      if (section.title && section.collapsible) {
        const hasActive = section.items.some(item =>
          location.pathname === item.to || item.children?.some(c => location.pathname === c.to)
        )
        if (hasActive) {
          setExpandedSections(prev => ({ ...prev, [section.title!]: true }))
        }
      }
    })
  }, [location.pathname])

  const isItemVisible = (item: NavItem) => !item.roles || item.roles.includes(role)
  const isSectionVisible = (section: NavSection) => {
    if (section.roles && !section.roles.includes(role)) return false
    return section.items.some(isItemVisible)
  }

  const toggleSection = (title: string) => {
    setExpandedSections(prev => ({ ...prev, [title]: !prev[title] }))
  }

  return (
    <aside className="h-screen w-64 fixed left-0 top-0 pt-16 bg-slate-50 flex flex-col gap-2 p-4 z-40 hidden md:flex">
      <div className="mb-6 px-2">
      </div>

      <nav className="flex flex-col gap-1 flex-1 overflow-y-auto">
        {navSections.filter(isSectionVisible).map((section, sIdx) => {
          const isExpanded = section.title ? expandedSections[section.title] !== false : true

          return (
            <div key={sIdx}>
              {/* Section Header (collapsible) */}
              {section.title && (
                <button
                  onClick={() => toggleSection(section.title!)}
                  className="w-full flex items-center gap-2 px-3 pt-4 pb-2 group"
                >
                  <span className="material-symbols-outlined text-[16px] text-slate-400 group-hover:text-slate-600 transition-colors">
                    {isExpanded ? 'expand_more' : 'chevron_right'}
                  </span>
                  {section.icon && (
                    <span className="material-symbols-outlined text-[16px] text-slate-400">
                      {section.icon}
                    </span>
                  )}
                  <span className="text-[11px] uppercase tracking-widest text-slate-400 font-bold group-hover:text-slate-600 transition-colors">
                    {section.title}
                  </span>
                </button>
              )}

              {/* Section Items */}
              {isExpanded && (
                <div className={section.title ? 'ml-3 flex flex-col gap-0.5' : 'flex flex-col gap-0.5'}>
                  {section.items.filter(isItemVisible).map((item) => {
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
                              <span className="material-symbols-outlined text-[20px]"
                                style={(isActive || isParentActive) ? { fontVariationSettings: "'FILL' 1" } : {}}>
                                {item.icon}
                              </span>
                              <span>{item.label}</span>
                            </>
                          )}
                        </NavLink>

                        {/* Sub-items (nested submenu) */}
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
