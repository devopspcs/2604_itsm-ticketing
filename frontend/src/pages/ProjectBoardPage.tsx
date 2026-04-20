import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { jiraService } from '../services/jira.service'
import { projectService } from '../services/project.service'
import { SprintBoard } from '../components/project/SprintBoard'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { WorkflowStatus, Sprint } from '../types/jira'

export function ProjectBoardPage() {
  const { id: projectId } = useParams()
  const navigate = useNavigate()
  const [sprint, setSprint] = useState<Sprint | null>(null)
  const [statuses, setStatuses] = useState<WorkflowStatus[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Create Sprint dialog
  const [showCreateSprint, setShowCreateSprint] = useState(false)
  const [sprintName, setSprintName] = useState('')
  const [sprintGoal, setSprintGoal] = useState('')
  const [sprintDuration, setSprintDuration] = useState(14)
  const [creatingSprint, setCreatingSprint] = useState(false)

  // Create Record dialog
  const [showCreateRecord, setShowCreateRecord] = useState(false)
  const [recordTitle, setRecordTitle] = useState('')
  const [recordDesc, setRecordDesc] = useState('')
  const [creatingRecord, setCreatingRecord] = useState(false)

  const fetchProjectData = async () => {
    if (!projectId) return
    try {
      setLoading(true)
      setError(null)
      try {
        const sprintRes = await jiraService.getActiveSprint(projectId)
        setSprint(sprintRes.data)
      } catch (e: any) {
        if (e?.response?.status === 404) setSprint(null)
        else throw e
      }
      try {
        const workflowRes = await jiraService.getWorkflow(projectId)
        if (workflowRes.data) {
          const statusesRes = await jiraService.listWorkflowStatuses(projectId, workflowRes.data.id)
          setStatuses(statusesRes.data || [])
        }
      } catch (e: any) {
        if (e?.response?.status === 404) setStatuses([])
        else throw e
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal memuat project board')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchProjectData() }, [projectId])

  const handleCreateSprint = async () => {
    if (!projectId || !sprintName.trim()) return
    try {
      setCreatingSprint(true)
      const startDate = new Date().toISOString()
      const endDate = new Date(Date.now() + sprintDuration * 86400000).toISOString()
      await jiraService.createSprint(projectId, {
        name: sprintName,
        goal: sprintGoal,
        start_date: startDate,
        end_date: endDate,
      })
      setShowCreateSprint(false)
      setSprintName('')
      setSprintGoal('')
      await fetchProjectData()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal membuat sprint')
    } finally {
      setCreatingSprint(false)
    }
  }

  const handleCreateRecord = async () => {
    if (!projectId || !recordTitle.trim()) return
    try {
      setCreatingRecord(true)
      const res = await projectService.createRecord(projectId, {
        title: recordTitle,
        description: recordDesc,
      })
      // Auto-assign to active sprint if exists
      if (sprint?.id && res.data?.id) {
        try {
          await jiraService.bulkAssignToSprint(projectId, {
            sprint_id: sprint.id,
            record_ids: [res.data.id],
          })
        } catch {
          // Non-critical: record created but not assigned to sprint
        }
      }
      setShowCreateRecord(false)
      setRecordTitle('')
      setRecordDesc('')
      window.location.reload()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal membuat record')
    } finally {
      setCreatingRecord(false)
    }
  }

  if (loading) {
    return (
      <div className="flex h-screen bg-surface">
        <ProjectBoardSidebar projectId={projectId || ''} />
        <div className="flex-1 flex flex-col">
          <div className="h-16 bg-surface-container-low border-b border-outline-variant/10" />
          <div className="flex-1 p-6">
            <div className="grid grid-cols-4 gap-4">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="bg-surface-container-low rounded-lg p-4 min-h-96 animate-pulse" />
              ))}
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <div className="bg-surface-container-low border-b border-outline-variant/10 px-8 py-4 flex items-center justify-between">
          <h1 className="text-2xl font-bold text-on-surface">Board</h1>
          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowCreateRecord(true)}
              className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
            >
              <span className="material-symbols-outlined text-[18px]">add</span>
              Create Issue
            </button>
          </div>
        </div>

        {/* Board Content */}
        <div className="flex-1 overflow-auto px-8 py-6">
          {error && (
            <div className="mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
              <span className="material-symbols-outlined text-lg">error</span>
              <span>{error}</span>
              <button onClick={() => setError(null)} className="ml-auto">
                <span className="material-symbols-outlined text-[16px]">close</span>
              </button>
            </div>
          )}

          {projectId && sprint?.id ? (
            <SprintBoard projectId={projectId} sprintId={sprint.id} statuses={statuses} />
          ) : (
            <div className="flex flex-col h-full items-center justify-center gap-4">
              <span className="material-symbols-outlined text-6xl text-on-surface-variant/30">sprint</span>
              <div className="text-center">
                <p className="text-lg font-semibold text-on-surface mb-2">No Active Sprint</p>
                <p className="text-sm text-on-surface-variant mb-6">Buat sprint untuk mulai mengelola board</p>
              </div>
              <div className="flex gap-3">
                <button
                  onClick={() => setShowCreateSprint(true)}
                  className="px-6 py-2 text-sm font-medium text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
                >
                  Create Sprint
                </button>
                <button
                  onClick={() => navigate(`/projects/${projectId}/backlog`)}
                  className="px-6 py-2 text-sm font-medium text-primary bg-primary/10 rounded-lg hover:bg-primary/20 transition-colors"
                >
                  View Backlog
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Create Sprint Dialog */}
      {showCreateSprint && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={() => setShowCreateSprint(false)}>
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6" onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-bold text-on-surface">Create Sprint</h2>
              <button onClick={() => setShowCreateSprint(false)} className="p-1 hover:bg-surface-container-high rounded-lg">
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
                <textarea value={sprintGoal} onChange={e => setSprintGoal(e.target.value)} placeholder="Tujuan sprint ini..."
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none resize-none" rows={2} />
              </div>
              <div>
                <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">Durasi (hari)</label>
                <select value={sprintDuration} onChange={e => setSprintDuration(Number(e.target.value))}
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none">
                  <option value={7}>1 Minggu</option>
                  <option value={14}>2 Minggu</option>
                  <option value={21}>3 Minggu</option>
                  <option value={28}>4 Minggu</option>
                </select>
              </div>
              <div className="flex justify-end gap-3 mt-2">
                <button onClick={() => setShowCreateSprint(false)} className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-xl">Batal</button>
                <button onClick={handleCreateSprint} disabled={!sprintName.trim() || creatingSprint}
                  className="px-5 py-2 text-sm font-bold text-on-primary bg-primary rounded-xl hover:opacity-90 disabled:opacity-50">
                  {creatingSprint ? 'Membuat...' : 'Buat Sprint'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Create Record Dialog */}
      {showCreateRecord && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={() => setShowCreateRecord(false)}>
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6" onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-bold text-on-surface">Create Issue</h2>
              <button onClick={() => setShowCreateRecord(false)} className="p-1 hover:bg-surface-container-high rounded-lg">
                <span className="material-symbols-outlined text-on-surface-variant">close</span>
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">Judul</label>
                <input value={recordTitle} onChange={e => setRecordTitle(e.target.value)} placeholder="Judul issue..."
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none" autoFocus />
              </div>
              <div>
                <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">Deskripsi (opsional)</label>
                <textarea value={recordDesc} onChange={e => setRecordDesc(e.target.value)} placeholder="Deskripsi issue..."
                  className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none resize-none" rows={3} />
              </div>
              <div className="flex justify-end gap-3 mt-2">
                <button onClick={() => setShowCreateRecord(false)} className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-xl">Batal</button>
                <button onClick={handleCreateRecord} disabled={!recordTitle.trim() || creatingRecord}
                  className="px-5 py-2 text-sm font-bold text-on-primary bg-primary rounded-xl hover:opacity-90 disabled:opacity-50">
                  {creatingRecord ? 'Membuat...' : 'Buat Issue'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
