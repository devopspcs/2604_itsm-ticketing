import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import { useNotifications } from '../../hooks/useNotifications'

export function Header() {
  const { logout } = useAuth()
  const { unreadCount } = useNotifications()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  return (
    <header className="fixed top-0 w-full z-50 bg-white/80 backdrop-blur-md shadow-sm flex justify-between items-center px-6 py-3">
      <div className="flex items-center gap-8">
        <Link to="/dashboard" className="text-xl font-bold text-blue-900 font-headline">
          PCS Ticketing System
        </Link>
        <nav className="hidden md:flex items-center gap-4">
          <Link to="/dashboard" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Overview</Link>
          <Link to="/tickets" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Tickets</Link>
          <Link to="/approvals" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Approvals</Link>
          <Link to="/activity-logs" className="text-slate-500 hover:bg-slate-100 transition-colors px-3 py-1 rounded-xl text-sm font-medium">Audit</Link>
        </nav>
      </div>

      <div className="flex items-center gap-4">
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

        <Link to="/profile" className="p-2 text-on-surface-variant hover:bg-slate-100 rounded-full active:scale-95 duration-150">
          <span className="material-symbols-outlined">account_circle</span>
        </Link>

        <button
          onClick={handleLogout}
          className="flex items-center gap-2 px-3 py-1.5 text-sm font-semibold text-on-surface-variant hover:bg-slate-100 rounded-xl transition-colors"
        >
          <span className="material-symbols-outlined text-[18px]">logout</span>
          <span className="hidden sm:inline">Logout</span>
        </button>
      </div>
    </header>
  )
}
