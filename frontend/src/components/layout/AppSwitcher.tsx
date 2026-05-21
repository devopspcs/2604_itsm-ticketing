import { useState, useRef, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { appService } from '../../services/app.service'
import type { RootState } from '../../store'

const ALL_APPS = [
  { id: 'ticketing', code: 'ticketing', label: 'Ticketing System', icon: 'confirmation_number', path: '/dashboard' },
  { id: 'projects', code: 'project-board', label: 'Project Board', icon: 'view_kanban', path: '/projects' },
]

export function AppSwitcher() {
  const [open, setOpen] = useState(false)
  const [accessibleCodes, setAccessibleCodes] = useState<string[]>([])
  const ref = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()
  const location = useLocation()
  const role = useSelector((s: RootState) => s.auth.role)

  const activeApp = location.pathname.startsWith('/projects') ? 'projects' : 'ticketing'

  // Fetch user's accessible apps
  useEffect(() => {
    if (role === 'admin') {
      // Admin sees all apps
      setAccessibleCodes(ALL_APPS.map(a => a.code))
      return
    }
    appService.getMyApps()
      .then(res => {
        const apps = res.data || []
        setAccessibleCodes(apps.map((a: { application: { code: string } }) => a.application.code))
      })
      .catch(() => {
        // On error, show all (fail-open for UX)
        setAccessibleCodes(ALL_APPS.map(a => a.code))
      })
  }, [role])

  // Filter apps by access
  const visibleApps = ALL_APPS.filter(app => accessibleCodes.includes(app.code))
  const current = visibleApps.find(a => a.id === activeApp) || visibleApps[0]

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  // Don't render switcher if user only has access to 1 or 0 apps
  if (!current || visibleApps.length <= 1) {
    if (current) {
      return (
        <div className="flex items-center gap-2 px-3 py-1.5 text-sm font-semibold text-on-surface">
          <span className="material-symbols-outlined text-[18px] text-primary">{current.icon}</span>
          <span className="hidden sm:inline">{current.label}</span>
        </div>
      )
    }
    return null
  }

  return (
    <div className="relative" ref={ref}>
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2 px-3 py-1.5 rounded-xl hover:bg-slate-100 transition-colors text-sm font-semibold text-on-surface"
      >
        <span className="material-symbols-outlined text-[18px] text-primary">{current.icon}</span>
        <span className="hidden sm:inline">{current.label}</span>
        <span className="material-symbols-outlined text-on-surface-variant text-[16px]">
          {open ? 'expand_less' : 'expand_more'}
        </span>
      </button>

      {open && (
        <div className="absolute left-0 top-full mt-1 w-56 bg-white rounded-xl shadow-xl border border-outline-variant/15 overflow-hidden z-50 animate-in fade-in slide-in-from-top-2">
          <p className="px-4 py-2 text-[10px] font-bold text-on-surface-variant uppercase tracking-widest">Aplikasi</p>
          {visibleApps.map(app => (
            <button
              key={app.id}
              onClick={() => { navigate(app.path); setOpen(false) }}
              className={`flex items-center gap-3 w-full px-4 py-2.5 text-sm transition-colors ${
                app.id === activeApp
                  ? 'bg-primary/5 text-primary font-bold'
                  : 'text-on-surface hover:bg-surface-container-low'
              }`}
            >
              <span className="material-symbols-outlined text-[20px]"
                style={app.id === activeApp ? { fontVariationSettings: "'FILL' 1" } : {}}>
                {app.icon}
              </span>
              <span>{app.label}</span>
              {app.id === activeApp && (
                <span className="material-symbols-outlined text-primary text-[16px] ml-auto">check</span>
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
