import { useState } from 'react'
import { Outlet, Navigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import type { RootState } from '../../store'
import { Header } from './Header'
import { ProjectSidebar } from './ProjectSidebar'
import { ProjectNavigation } from './ProjectNavigation'
import { CreateProjectDialog } from '../project/CreateProjectDialog'

export function ProjectBoardLayout() {
  const isAuthenticated = useSelector((s: RootState) => s.auth.isAuthenticated)
  const [showCreate, setShowCreate] = useState(false)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return (
    <div className="bg-surface min-h-screen">
      <Header />
      <ProjectSidebar onCreateProject={() => setShowCreate(true)} />
      <ProjectNavigation />
      <main className="md:ml-64 pt-16 min-h-screen bg-surface">
        <Outlet />
      </main>
      {showCreate && (
        <CreateProjectDialog
          onClose={() => setShowCreate(false)}
          onCreated={() => { setShowCreate(false); window.location.reload() }}
        />
      )}
    </div>
  )
}
