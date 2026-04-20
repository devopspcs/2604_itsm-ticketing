import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { jiraService } from '../services/jira.service'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { Sprint } from '../types/jira'

export function SprintBoardPage() {
  const { id: projectId } = useParams()
  const [sprints, setSprints] = useState<Sprint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [showCreate, setShowCreate] = useState(false)
  const [sprintName, setSprintName] = useState('')
  const [sprintGoal, setSprintGoal] = useState('')
  const [sprintDuration, setSprintDuration] = useState(14)
  const [creating, setCreating] = useState(false)
  const [actionId, setActionId] = useState<string | null>(null)

  const fetchSprints = async () => {
    if (!projectId) return
    try {
      setLoading(true)
      setError(null)
      const res = await jiraService.listSprints(projectId)
      setSprints(res.data || [])
    } catch (e: any) {
      if (e?.response?.status !== 404) setError('Gagal memuat sprints')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchSprints() }, [projectId])

  const handleCreate = async () => {
    if (!projectId || !sprintName.trim()) return
    try {
      setCreating(true)
      const start = new Date().toISOString()
      const end = new Date(Date.now() + sprintDuration * 86400000).toISOString()
      await jiraService.createSprint(projectId, { name: sprintName, goal: sprintGoal, start_date: start, end_date: end })
      setShowCreate(false)
      setSprintName('')
      setSprintGoal('')
      await fetchSprints()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal membuat sprint')
    } finally {
      setCreating(false)
    }
  }

  const handleStart = async (id: string) => {
    if (!projectId) return
    try {
      setActionId(id)
      await jiraService.startSprint(projectId, id)
      await fetchSprints()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal memulai sprint')
    } finally {
      setActionId(null)
    }
  }

  const handleComplete = async (id: string) => {
    if (!projectId || !confirm('Selesaikan sprint ini?')) return
    try {
      setActionId(id)
      await jiraService.completeSprint(projectId, id)
      await fetchSprints()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal menyelesaikan sprint')
    } finally {
      setActionId(null)
    }
  }

  const hasActive = sprints.some(s => s.status === 'Active')
  const active = sprints.filter(s => s.status === 'Active')
  const planned = sprints.filter(s => s.status === 'Planned')
  const completed = sprints.filter(s => s.status === 'Completed')

  const statusColor = (s: string) => {
    if (s === 'Active') return 'bg-green-500/20 text-green-700'
    if (s === 'Planned') return 'bg-blue-500/20 text-blue-700'
    return 'bg-gray-500/20 text-gray-700'
  }

  const formatDate = (d?: string) => d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-'

  if (loading) {
    return (
      <div className="flex h-screen bg-surface">
        <ProjectBoardSidebar projectId={projectId || ''} />
        <div className="flex-1 p-8 space-y-3">
          {[1, 2, 3].map(i => <div key={i} className="h-28 bg-surface-container-low rounded-xl animate-pulse" />)}
        </div>
      </div>
    )
  }

  const SprintCard = ({ sprint, actions }: { sprint: Sprint; actions?: React.ReactNode }) => (
    <div className="bg-surface-container-low rounded-xl p-5 border border-outline-variant/10 hover:border-outline-variant/20 transition-colors">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-3 mb-1">
            <h3 className="text-base font-semibold text-on-surface">{sprint.name}</h3>
            <span className={`px-2.5 py-0.5 rounded-full text-xs font-semibold ${statusColor(sprint.status)}`}>{sprint.status}</span>
          </div>
          {sprint.goal && <p className="text-sm text-on-surface-variant">{sprint.goal}</p>}
        </div>
        {actions}
      </div>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
        <div>
          <p className="text-xs text-on-surface-variant">Mulai</p>
          <p className="font-medium text-on-surface">{formatDate(sprint.start_date)}</p>
        </div>
        <div>
          <p className="text-xs text-on-surface-variant">Selesai</p>
          <p className="font-medium text-on-surface">{formatDate(sprint.end_date)}</p>
        </div>
        <div>
          <p className="text-xs text-on-surface-variant">Status</p>
          <p className="font-medium text-on-surface">{sprint.status}</p>
        </div>
        <div>
          <p className="text-xs text-on-surface-variant">Dibuat</p>
          <p className="font-medium text-on-surface">{formatDate(sprint.created_at)}</p>
        </div>
      </div>
    </div>
  )

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Sprint Management</h1>
            <p className="text-sm text-on-surface-variant mt-1">Kelola sprint project — buat, mulai, dan selesaikan sprint</p>
          </div>
          <button onClick={() => setShowCreate(true)}
            className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors">
            <span className="material-symbols-outlined text-[18px]">add</span>
            New Sprint
          </button>
        </div>

        {error && (
          <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
            <span className="material-symbols-outlined text-lg">error</span>
            <span>{error}</span>
            <button onClick={() => setError(null)} className="ml-auto"><span className="material-symbols-outlined text-[16px]">close</span></button>
          </div>
        )}

        <div className="flex-1 overflow-auto px-8 pb-8 space-y-8">
          {/* Active Sprint */}
          {active.length > 0 && (
            <div>
              <h2 className="text-xs font-bold text-green-600 uppercase tracking-widest mb-3 flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-green-500" />
                Active Sprint
              </h2>
              {active.map(s => (
                <SprintCard key={s.id} sprint={s} actions={
                  <button onClick={() => handleComplete(s.id)} disabled={actionId === s.id}
                    className="px-4 py-2 text-sm font-semibold text-on-surface-variant bg-surface-container-high rounded-lg hover:bg-surface-container-highest transition-colors disabled:opacity-50">
                    {actionId === s.id ? 'Menyelesaikan...' : 'Complete Sprint'}
                  </button>
                } />
              ))}
            </div>
          )}

          {/* Planned Sprints */}
          {planned.length > 0 && (
            <div>
              <h2 className="text-xs font-bold text-blue-600 uppercase tracking-widest mb-3 flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-blue-500" />
                Planned ({planned.length})
              </h2>
              <div className="space-y-3">
                {planned.map(s => (
                  <SprintCard key={s.id} sprint={s} actions={
                    <button onClick={() => handleStart(s.id)} disabled={hasActive || actionId === s.id}
                      className="px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50"
                      title={hasActive ? 'Selesaikan sprint aktif dulu' : ''}>
                      {actionId === s.id ? 'Memulai...' : 'Start Sprint'}
                    </button>
                  } />
                ))}
              </div>
            </div>
          )}

          {/* Completed Sprints */}
          {completed.length > 0 && (
            <div>
              <h2 className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-3 flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-gray-400" />
                Completed ({completed.length})
              </h2>
              <div className="space-y-3 opacity-70">
                {completed.map(s => <SprintCard key={s.id} sprint={s} />)}
              </div>
            </div>
          )}

          {/* Empty */}
          {sprints.length === 0 && (
            <div className="flex flex-col items-center justify-center py-20">
              <span className="material-symbols-outlined text-6xl text-on-surface-variant/30 mb-4">sprint</span>
              <p className="text-on-surface-variant text-lg mb-2">Belum ada sprint</p>
              <p className="text-on-surface-variant text-sm">Klik "New Sprint" untuk membuat sprint pertama</p>
            </div>
          )}
        </div>
      </div>

      {/* Create Sprint Dialog */}
      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={() => setShowCreate(false)}>
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6" onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-bold text-on-surface">Create Sprint</h2>
              <button onClick={() => setShowCreate(false)} className="p-1 hover:bg-surface-container-high rounded-lg">
                <span className="material-symbols-outlined text-on-surface-variant">close</span>
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">Nama Sprint</label>
                <input value={sprintName} onChange={e => setSprintName(e.target.value)} placeholder="Sprint 1"
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none" autoFocus />
              </div>
              <div>
                <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">Goal (opsional)</label>
                <textarea value={sprintGoal} onChange={e => setSprintGoal(e.target.value)} placeholder="Tujuan sprint..."
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none resize-none" rows={2} />
              </div>
              <div>
                <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">Durasi</label>
                <select value={sprintDuration} onChange={e => setSprintDuration(Number(e.target.value))}
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm outline-none">
                  <option value={7}>1 Minggu</option>
                  <option value={14}>2 Minggu</option>
                  <option value={21}>3 Minggu</option>
                  <option value={28}>4 Minggu</option>
                </select>
              </div>
              <div className="flex justify-end gap-3 mt-2">
                <button onClick={() => setShowCreate(false)} className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-xl">Batal</button>
                <button onClick={handleCreate} disabled={!sprintName.trim() || creating}
                  className="px-5 py-2 text-sm font-bold text-on-primary bg-primary rounded-xl hover:opacity-90 disabled:opacity-50">
                  {creating ? 'Membuat...' : 'Buat Sprint'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
