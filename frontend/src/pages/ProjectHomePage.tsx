import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useSelector } from 'react-redux'
import type { RootState } from '../store'
import { projectService } from '../services/project.service'
import type { ProjectHomeData } from '../types/project'
import type { User } from '../types'
import api from '../services/api'

export function ProjectHomePage() {
  const userId = useSelector((s: RootState) => s.auth.userId)
  const [data, setData] = useState<ProjectHomeData | null>(null)
  const [loading, setLoading] = useState(true)
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [users, setUsers] = useState<User[]>([])

  useEffect(() => {
    projectService.getHome()
      .then(res => setData(res.data))
      .catch(() => {})
      .finally(() => setLoading(false))

    api.get<User[]>('/users/list').then(res => {
      const list = res.data ?? []
      setUsers(list)
      const me = list.find(u => u.id === userId)
      if (me) setCurrentUser(me)
    }).catch(() => {})
  }, [userId])

  if (loading) {
    return (
      <div className="p-8">
        <div className="h-8 w-72 bg-surface-container-high rounded-lg animate-pulse mb-6" />
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 flex flex-col gap-4">
            <div className="h-16 bg-surface-container-high rounded-xl animate-pulse" />
            {[1, 2, 3].map(i => (
              <div key={i} className="h-14 bg-surface-container-high rounded-xl animate-pulse" />
            ))}
          </div>
          <div className="flex flex-col gap-3">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-12 bg-surface-container-high rounded-xl animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline mb-1">
        Selamat datang, {currentUser?.full_name ?? 'User'}
      </h1>
      <p className="text-sm text-on-surface-variant mb-6">Ringkasan aktivitas project board Anda</p>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 flex flex-col gap-4">
          {/* Overdue alert */}
          {data && data.overdue_count > 0 && (
            <div className="bg-error-container/20 border border-error/10 rounded-xl px-5 py-4 flex items-center gap-3">
              <span className="material-symbols-outlined text-error text-xl">warning</span>
              <div>
                <p className="text-sm font-bold text-error">
                  Anda memiliki {data.overdue_count} record yang overdue
                </p>
                <Link to="/projects/calendar" className="text-xs font-bold text-error underline">
                  Lihat di kalender
                </Link>
              </div>
            </div>
          )}

          {/* Recent activities */}
          <div className="bg-white rounded-xl border border-outline-variant/10 p-5">
            <h2 className="text-sm font-bold text-on-surface uppercase tracking-widest mb-4">Aktivitas Terbaru</h2>
            {!data?.recent_activities?.length ? (
              <p className="text-sm text-on-surface-variant/60">Belum ada aktivitas</p>
            ) : (
              <div className="flex flex-col gap-3">
                {data.recent_activities.map(a => (
                  <div key={a.id} className="flex items-start gap-3 py-2 border-b border-outline-variant/5 last:border-0">
                    <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                      <span className="material-symbols-outlined text-primary text-[16px]">person</span>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm text-on-surface">
                        <span className="font-semibold">{a.action}</span>
                        {a.detail && <span className="text-on-surface-variant ml-1">{a.detail}</span>}
                      </p>
                      <p className="text-[10px] text-on-surface-variant/60 mt-0.5">
                        {new Date(a.created_at).toLocaleString('id-ID')}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
            <Link to="/projects/calendar" className="text-xs font-bold text-primary mt-3 inline-block hover:underline">
              LIHAT SEMUA
            </Link>
          </div>
        </div>

        {/* Active people panel */}
        <div className="bg-white rounded-xl border border-outline-variant/10 p-5 h-fit">
          <h2 className="text-sm font-bold text-on-surface uppercase tracking-widest mb-4">Orang Aktif Terakhir</h2>
          {users.length === 0 ? (
            <p className="text-sm text-on-surface-variant/60">Tidak ada data</p>
          ) : (
            <div className="flex flex-col gap-3">
              {users.slice(0, 8).map(u => (
                <div key={u.id} className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-primary-container flex items-center justify-center text-white text-xs font-bold shrink-0">
                    {u.full_name.charAt(0).toUpperCase()}
                  </div>
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-on-surface truncate">{u.full_name}</p>
                    <p className="text-[10px] text-on-surface-variant/60">{u.email}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
