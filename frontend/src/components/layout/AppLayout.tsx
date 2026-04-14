import { Outlet, Navigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import type { RootState } from '../../store'
import { Header } from './Header'
import { Sidebar } from './Sidebar'

export function AppLayout() {
  const isAuthenticated = useSelector((s: RootState) => s.auth.isAuthenticated)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return (
    <div className="bg-surface min-h-screen">
      <Header />
      <Sidebar />
      <main className="md:ml-64 pt-16 min-h-screen bg-surface">
        <Outlet />
      </main>
    </div>
  )
}
