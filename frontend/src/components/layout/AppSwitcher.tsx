import { useState, useRef, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'

const APPS = [
  { id: 'ticketing', label: 'Ticketing System', icon: 'confirmation_number', path: '/dashboard' },
  { id: 'projects', label: 'Project Board', icon: 'view_kanban', path: '/projects' },
]

export function AppSwitcher() {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()
  const location = useLocation()

  const activeApp = location.pathname.startsWith('/projects') ? 'projects' : 'ticketing'
  const current = APPS.find(a => a.id === activeApp)!

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

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
          {APPS.map(app => (
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
