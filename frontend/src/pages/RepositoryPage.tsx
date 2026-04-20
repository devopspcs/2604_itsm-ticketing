import { useState, useEffect, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { projectService } from '../services/project.service'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { ActivityLogEntry, ActivityLogFilterParams } from '../types/project'

export function RepositoryPage() {
  const { id: projectId } = useParams()
  const [logs, setLogs] = useState<ActivityLogEntry[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [actionFilter, setActionFilter] = useState('')
  const pageSize = 20

  const fetchLogs = useCallback(async () => {
    if (!projectId) return
    try {
      setLoading(true)
      const params: ActivityLogFilterParams = { page, page_size: pageSize }
      if (actionFilter) params.action_type = actionFilter
      const res = await projectService.listActivityLog(projectId, params)
      setLogs(res.data?.logs ?? [])
      setTotal(res.data?.total ?? 0)
    } catch (err) {
      console.error('Failed to fetch activity log:', err)
    } finally {
      setLoading(false)
    }
  }, [projectId, page, actionFilter])

  useEffect(() => { fetchLogs() }, [fetchLogs])
  useEffect(() => { setPage(1) }, [actionFilter])

  const totalPages = Math.max(1, Math.ceil(total / pageSize))

  const getActionColor = (action: string) => {
    if (action.includes('create') || action.includes('add')) return 'bg-green-500/20 text-green-700'
    if (action.includes('delete') || action.includes('remove')) return 'bg-red-500/20 text-red-700'
    if (action.includes('update') || action.includes('move') || action.includes('edit')) return 'bg-blue-500/20 text-blue-700'
    if (action.includes('complete')) return 'bg-purple-500/20 text-purple-700'
    return 'bg-gray-500/20 text-gray-700'
  }

  const getActionIcon = (action: string) => {
    if (action.includes('create') || action.includes('add')) return 'add_circle'
    if (action.includes('delete') || action.includes('remove')) return 'delete'
    if (action.includes('update') || action.includes('edit')) return 'edit'
    if (action.includes('move')) return 'swap_horiz'
    if (action.includes('complete')) return 'check_circle'
    return 'history'
  }

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Activity Log</h1>
          <p className="text-sm text-on-surface-variant mt-1">Project activity history</p>
        </div>

        <div className="px-8 py-4 border-b border-outline-variant/10">
          <div className="flex items-center gap-3">
            <span className="text-sm text-on-surface-variant">Filter:</span>
            <select value={actionFilter} onChange={e => setActionFilter(e.target.value)}
              className="px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-sm text-on-surface focus:outline-none focus:border-primary">
              <option value="">All Actions</option>
              <option value="create_record">Create Record</option>
              <option value="update_record">Update Record</option>
              <option value="delete_record">Delete Record</option>
              <option value="move_record">Move Record</option>
              <option value="complete_record">Complete Record</option>
              <option value="add_comment">Add Comment</option>
            </select>
            <span className="text-sm text-on-surface-variant ml-auto">{total} activit{total !== 1 ? 'ies' : 'y'}</span>
          </div>
        </div>

        <div className="flex-1 overflow-auto px-8 pb-8">
          <div className="mt-6 space-y-3">
            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3, 4].map(i => (
                  <div key={i} className="bg-surface-container-low rounded-lg p-4 border border-outline-variant/10 animate-pulse">
                    <div className="flex items-start gap-4">
                      <div className="w-10 h-10 bg-surface-container-high rounded-full" />
                      <div className="flex-1">
                        <div className="h-4 w-48 bg-surface-container-high rounded mb-2" />
                        <div className="h-3 w-32 bg-surface-container-high rounded" />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : logs.length === 0 ? (
              <div className="text-center py-12">
                <span className="material-symbols-outlined text-6xl text-on-surface-variant/30 block mb-4">history</span>
                <p className="text-on-surface-variant">No activity yet</p>
              </div>
            ) : (
              logs.map((log, index) => (
                <div key={log.id}
                  className="bg-surface-container-low rounded-lg p-4 border border-outline-variant/10 hover:border-outline-variant/20 transition-colors">
                  <div className="flex items-start gap-4">
                    <div className="flex flex-col items-center">
                      <div className={`p-2 rounded-full ${getActionColor(log.action)}`}>
                        <span className="material-symbols-outlined text-[18px]">{getActionIcon(log.action)}</span>
                      </div>
                      {index < logs.length - 1 && <div className="w-0.5 h-8 bg-outline-variant/20 my-1" />}
                    </div>
                    <div className="flex-1 pt-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className={`px-2 py-0.5 rounded text-xs font-medium ${getActionColor(log.action)}`}>{log.action}</span>
                      </div>
                      <p className="text-sm text-on-surface mb-1">{log.detail}</p>
                      <div className="flex items-center gap-3 text-xs text-on-surface-variant">
                        <span>Actor: {log.actor_id.slice(0, 8)}...</span>
                        <span>{new Date(log.created_at).toLocaleString()}</span>
                      </div>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>

          {/* Pagination */}
          {!loading && totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-6">
              <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1}
                className="px-3 py-2 text-sm rounded-lg bg-surface-container-low border border-outline-variant/20 text-on-surface disabled:opacity-40 hover:bg-surface-container-high transition-colors">
                Previous
              </button>
              <span className="text-sm text-on-surface-variant">Page {page} of {totalPages}</span>
              <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page >= totalPages}
                className="px-3 py-2 text-sm rounded-lg bg-surface-container-low border border-outline-variant/20 text-on-surface disabled:opacity-40 hover:bg-surface-container-high transition-colors">
                Next
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
