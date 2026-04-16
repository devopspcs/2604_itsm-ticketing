import { useEffect, useRef, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useAuth } from '../../hooks/useAuth'
import { useNotifications } from '../../hooks/useNotifications'
import { useTheme, THEMES } from '../../hooks/useTheme'
import { AppSwitcher } from './AppSwitcher'
import type { RootState } from '../../store'
import type { User } from '../../types'
import api from '../../services/api'

export function Header() {
  const { logout } = useAuth()
  const { unreadCount } = useNotifications()
  const { theme, setTheme } = useTheme()
  const navigate = useNavigate()
  const userId = useSelector((s: RootState) => s.auth.userId)
  const role = useSelector((s: RootState) => s.auth.role)

  const [showProfileMenu, setShowProfileMenu] = useState(false)
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!userId) return
    api.get<User[]>('/users/list')
      .then(res => {
        const me = (res.data ?? []).find(u => u.id === userId)
        if (me) setCurrentUser(me)
      })
      .catch(() => {})
  }, [userId])

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setShowProfileMenu(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleLogout = async () => {
    setShowProfileMenu(false)
    await logout()
    navigate('/login')
  }

  const displayName = currentUser?.full_name ?? 'User'
  const displayEmail = currentUser?.email ?? ''
  const initial = displayName.charAt(0).toUpperCase()

  return (
    <header className="fixed top-0 w-full z-50 bg-white/80 backdrop-blur-md shadow-sm flex justify-between items-center px-6 py-3">
      <div className="flex items-center gap-8">
        <Link to="/dashboard" className="text-xl font-bold text-accent-900 font-headline">
          PCS ITSM
        </Link>
        <AppSwitcher />
        <nav className="hidden md:flex items-center gap-4">
          <Link to="/dashboard" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Overview</Link>
          <Link to="/tickets" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Tickets</Link>
          <Link to="/approvals" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Approvals</Link>
          <Link to="/activity-logs" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Audit</Link>
        </nav>
      </div>

      <div className="flex items-center gap-3">
        <div className="relative hidden sm:block">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline text-lg">search</span>
          <input
            className="pl-10 pr-4 py-2 bg-surface-container-highest border-none rounded-xl text-sm focus:ring-2 focus:ring-primary/20 w-64 outline-none"
            placeholder="Global Search..."
            type="text"
          />
        </div>

        <Link to="/notifications" className="relative p-2 text-on-surface-variant hover:bg-slate-100 rounded-full active:scale-95 duration-150">
          <span className="material-symbols-outlined">notifications</span>
          {unreadCount > 0 && (
            <span className="absolute top-1 right-1 w-2 h-2 bg-error rounded-full" />
          )}
        </Link>

        {/* Profile Avatar + Dropdown */}
        <div className="relative" ref={menuRef}>
          <button
            onClick={() => setShowProfileMenu(!showProfileMenu)}
            className="flex items-center gap-2 p-1 pr-3 rounded-full hover:bg-slate-100 transition-colors active:scale-[0.97]"
          >
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-primary-container flex items-center justify-center text-white text-sm font-bold">
              {initial}
            </div>
            <span className="hidden sm:block text-sm font-semibold text-on-surface max-w-[120px] truncate">{displayName}</span>
            <span className="material-symbols-outlined text-on-surface-variant text-[18px]">
              {showProfileMenu ? 'expand_less' : 'expand_more'}
            </span>
          </button>

          {showProfileMenu && (
            <div className="absolute right-0 top-full mt-2 w-72 bg-white rounded-xl shadow-xl border border-outline-variant/15 overflow-hidden z-50 animate-in fade-in slide-in-from-top-2">
              {/* User Info Header */}
              <div className="px-5 py-4 bg-gradient-to-r from-primary/5 to-primary-container/5 border-b border-outline-variant/10">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary to-primary-container flex items-center justify-center text-white text-lg font-bold shrink-0">
                    {initial}
                  </div>
                  <div className="min-w-0">
                    <p className="text-sm font-bold text-on-surface truncate">{displayName}</p>
                    <p className="text-xs text-on-surface-variant truncate">{displayEmail}</p>
                    <div className="flex items-center gap-2 mt-1">
                      <span className={`text-[10px] px-2 py-0.5 rounded-full font-bold capitalize ${
                        role === 'admin' ? 'bg-primary-fixed text-primary' :
                        role === 'approver' ? 'bg-amber-100 text-amber-700' :
                        'bg-surface-container-high text-on-surface-variant'
                      }`}>{role}</span>
                      {currentUser?.position && (
                        <span className="text-[10px] px-2 py-0.5 rounded-full font-bold bg-tertiary-fixed text-on-tertiary-fixed capitalize">
                          {currentUser.position.replace('_', ' ')}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </div>

              {/* Menu Items */}
              <div className="py-1">
                <Link
                  to="/profile"
                  onClick={() => setShowProfileMenu(false)}
                  className="flex items-center gap-3 px-5 py-2.5 text-sm text-on-surface hover:bg-surface-container-low transition-colors"
                >
                  <span className="material-symbols-outlined text-[20px] text-on-surface-variant">person</span>
                  My Profile
                </Link>
                <Link
                  to="/notifications"
                  onClick={() => setShowProfileMenu(false)}
                  className="flex items-center gap-3 px-5 py-2.5 text-sm text-on-surface hover:bg-surface-container-low transition-colors"
                >
                  <span className="material-symbols-outlined text-[20px] text-on-surface-variant">notifications</span>
                  Notifications
                  {unreadCount > 0 && (
                    <span className="ml-auto text-[10px] font-bold bg-error text-white px-1.5 py-0.5 rounded-full">{unreadCount}</span>
                  )}
                </Link>
              </div>

              {/* Theme Picker */}
              <div className="border-t border-outline-variant/10 px-5 py-3">
                <p className="text-[10px] font-bold text-on-surface-variant uppercase tracking-widest mb-2">Theme</p>
                <div className="flex items-center gap-2">
                  {THEMES.map((t) => (
                    <button
                      key={t.id}
                      onClick={() => setTheme(t.id)}
                      title={t.label}
                      className={`w-7 h-7 rounded-full transition-all hover:scale-110 active:scale-95 ${
                        theme === t.id ? 'ring-2 ring-offset-2 ring-on-surface scale-110' : ''
                      }`}
                      style={{ backgroundColor: t.color }}
                    />
                  ))}
                </div>
              </div>

              <div className="border-t border-outline-variant/10">
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-3 px-5 py-2.5 text-sm text-error hover:bg-error-container/20 transition-colors w-full"
                >
                  <span className="material-symbols-outlined text-[20px]">logout</span>
                  Sign Out
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}
