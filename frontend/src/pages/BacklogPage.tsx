import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { jiraService } from '../services/jira.service'
import { projectService } from '../services/project.service'
import { BacklogView } from '../components/project/BacklogView'
import { SearchFilterBar } from '../components/project/SearchFilterBar'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { Sprint, Label, WorkflowStatus } from '../types/jira'

export function BacklogPage() {
  const { id: projectId } = useParams()
  const [sprints, setSprints] = useState<Sprint[]>([])
  const [labels, setLabels] = useState<Label[]>([])
  const [statuses, setStatuses] = useState<WorkflowStatus[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showCreateRecord, setShowCreateRecord] = useState(false)
  const [recordTitle, setRecordTitle] = useState('')
  const [recordDesc, setRecordDesc] = useState('')
  const [creatingRecord, setCreatingRecord] = useState(false)

  useEffect(() => {
    const fetchBacklogData = async () => {
      if (!projectId) return
      try {
        setLoading(true)
        setError(null)

        // Get sprints (404 = no sprints yet)
        try {
          const sprintsRes = await jiraService.listSprints(projectId)
          setSprints(sprintsRes.data || [])
        } catch (e: any) {
          if (e?.response?.status !== 404) throw e
        }

        // Get labels (404 = no labels yet)
        try {
          const labelsRes = await jiraService.listLabels(projectId)
          setLabels(labelsRes.data || [])
        } catch (e: any) {
          if (e?.response?.status !== 404) throw e
        }

        // Get workflow statuses (404 = no workflow yet)
        try {
          const workflowRes = await jiraService.getWorkflow(projectId)
          if (workflowRes.data) {
            const statusesRes = await jiraService.listWorkflowStatuses(projectId, workflowRes.data.id)
            setStatuses(statusesRes.data || [])
          }
        } catch (e: any) {
          if (e?.response?.status !== 404) throw e
        }
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Gagal memuat backlog'
        setError(errorMessage)
      } finally {
        setLoading(false)
      }
    }

    fetchBacklogData()
  }, [projectId])

  if (loading) {
    return (
      <div className="flex flex-col h-[calc(100vh-64px)]">
        <div className="px-8 pt-8 pb-4">
          <div className="h-7 w-48 bg-surface-container-high rounded-lg animate-pulse mb-2" />
          <div className="h-4 w-64 bg-surface-container-high rounded-lg animate-pulse" />
        </div>
        <div className="flex-1 px-8 pb-8">
          <div className="grid grid-cols-4 gap-4">
            <div className="lg:col-span-3 bg-surface-container-low rounded-xl p-4 min-h-96 animate-pulse" />
            <div className="lg:col-span-1 bg-surface-container-low rounded-xl p-4 min-h-96 animate-pulse" />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Backlog</h1>
            <button onClick={() => setShowCreateRecord(true)}
              className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors">
              <span className="material-symbols-outlined text-[18px]">add</span>
              Create Issue
            </button>
          </div>
          <div className="mt-3">
            <SearchFilterBar 
              projectId={projectId || ''} 
              onSearchResults={() => {}} 
              issueTypes={[]}
              statuses={statuses}
              labels={labels}
            />
          </div>
        </div>

        {error && (
          <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
            <span className="material-symbols-outlined text-lg">error</span>
            <span>{error}</span>
          </div>
        )}

        <div className="flex-1 overflow-auto px-8 pb-8">
          {projectId && (
            <BacklogView projectId={projectId} sprints={sprints} labels={labels} statuses={statuses} />
          )}
        </div>
      </div>

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
              <input value={recordTitle} onChange={e => setRecordTitle(e.target.value)} placeholder="Judul issue..."
                className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none" autoFocus />
              <textarea value={recordDesc} onChange={e => setRecordDesc(e.target.value)} placeholder="Deskripsi (opsional)..."
                className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none resize-none" rows={3} />
              <div className="flex justify-end gap-3">
                <button onClick={() => setShowCreateRecord(false)} className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-xl">Batal</button>
                <button onClick={async () => {
                  if (!projectId || !recordTitle.trim()) return
                  setCreatingRecord(true)
                  try {
                    await projectService.createRecord(projectId, { title: recordTitle, description: recordDesc })
                    setShowCreateRecord(false)
                    setRecordTitle('')
                    setRecordDesc('')
                    window.location.reload()
                  } catch { setError('Gagal membuat issue') }
                  finally { setCreatingRecord(false) }
                }} disabled={!recordTitle.trim() || creatingRecord}
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
