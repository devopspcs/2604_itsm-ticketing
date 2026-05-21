import { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'
import { appService } from '../../services/app.service'
import { useSelector } from 'react-redux'
import type { RootState } from '../../store'

interface Props {
  appCode: string
  children: React.ReactNode
}

export function AppAccessGuard({ appCode, children }: Props) {
  const [hasAccess, setHasAccess] = useState<boolean | null>(null)
  const role = useSelector((s: RootState) => s.auth.role)

  useEffect(() => {
    // Admin always has access
    if (role === 'admin') {
      setHasAccess(true)
      return
    }

    appService.getMyApps()
      .then(res => {
        const apps = res.data || []
        const found = apps.some((a: { application: { code: string } }) => a.application.code === appCode)
        setHasAccess(found)
      })
      .catch(() => setHasAccess(true)) // On error, allow access (fail-open for UX)
  }, [appCode, role])

  if (hasAccess === null) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      </div>
    )
  }

  if (!hasAccess) {
    return <Navigate to="/dashboard" replace />
  }

  return <>{children}</>
}
