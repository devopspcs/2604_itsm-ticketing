import { useState, useEffect, useCallback, useRef } from 'react'
import { useParams } from 'react-router-dom'
import { projectService } from '../services/project.service'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { ProjectRecord, IssuesFilterParams } from '../types/project'

export function IssuesPage() {
  const { id: projectId } = useParams()
  const [records, setRecords] = useState<ProjectRecord[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const pageSize = 20

  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [assigneeFilter, setAssigneeFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const [debouncedSearch, setDebouncedSearch] = useState('')

  // Debounce search input
  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      setDebouncedSearch(search)
      setPage(1)
    }, 300)
    return () => { if (debounceRef.current) clearTimeout(debounceRef.current) }
  }, [search])

  const fetchIssues = useCallback(async () => {
    if (!projectId) return
    try {
      setLoading(true)
      const params: IssuesFilterParams = { page, page_size: pageSize }
      if (debouncedSearch) params.search = debouncedSearch
      if (statusFilter) params.status_id = statusFilter
      if (assigneeFilter) params.assignee_id = assigneeFilter
      if (typeFilter) params.issue_type = typeFilter
      const res = await projectService.listIssues(projectId, params)
      setRecords(res.data?.records ?? [])
      setTotal(res.data?.total ?? 0)
    } catch (err) {
      console.error('Failed to fetch issues:', err)
    } finally {
      setLoading(false)
    }
  }, [projectId, page, debouncedSearch, statusFilter, assigneeFilter, typeFilter])

  useEffect(() => { fetchIssues() }, [fetchIssues])

  // Reset page when filters change
  useEffect(() => { setPage(1) }, [statusFilter, assigneeFilter, typeFilter])

  const totalPages = Math.max(1, Math.ceil(total / pageSize))

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Issues</h1>
          <p className="text-sm text-on-surface-variant mt-1">All project issues and tasks</p>
        </div>

        <div className="px-8 py-4 border-b border-outline-variant/10 space-y-4">
          {/* Search */}
          <div className="flex items-center gap-2 bg-surface-container-low px-4 py-2 rounded-lg border border-outline-variant/20">
            <span className="material-symbols-outlined text-on-surface-variant">search</span>
            <input type="text" value={search} onChange={e => setSearch(e.target.value)}
              placeholder="Search issues..."
              className="flex-1 bg-transparent outline-none text-on-surface placeholder-on-surface-variant" />
          </div>

          {/* Filters */}
          <div className="flex items-center gap-3 flex-wrap">
            <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
              className="px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-sm text-on-surface focus:outline-none focus:border-primary">
              <option value="">All Statuses</option>
            </select>
            <select value={assigneeFilter} onChange={e => setAssigneeFilter(e.target.value)}
              className="px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-sm text-on-surface focus:outline-none focus:border-primary">
              <option value="">All Assignees</option>
            </select>
            <select value={typeFilter} onChange={e => setTypeFilter(e.target.value)}
              className="px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-sm text-on-surface focus:outline-none focus:border-primary">
              <option value="">All Types</option>
            </select>
            <span className="text-sm text-on-surface-variant ml-auto">{total} issue{total !== 1 ? 's' : ''}</span>
          </div>
        </div>

        <div className="flex-1 overflow-auto px-8 pb-8">
          <div className="mt-6 space-y-2">
            {loading ? (
              <div className="space-y-2">
                {[1, 2, 3, 4, 5].map(i => (
                  <div key={i} className="bg-surface-container-low rounded-lg p-4 border border-outline-variant/10 animate-pulse">
                    <div className="h-4 w-64 bg-surface-container-high rounded mb-2" />
                    <div className="h-3 w-32 bg-surface-container-high rounded" />
                  </div>
                ))}
              </div>
            ) : records.length === 0 ? (
              <div className="text-center py-12">
                <span className="material-symbols-outlined text-6xl text-on-surface-variant/30 block mb-4">task</span>
                <p className="text-on-surface-variant">No issues found</p>
              </div>
            ) : (
              <div className="space-y-2">
                {records.map(record => (
                  <div key={record.id}
                    className="bg-surface-container-low rounded-lg p-4 border border-outline-variant/10 hover:border-outline-variant/20 hover:shadow-md transition-all cursor-pointer">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <h3 className="text-sm font-semibold text-on-surface">{record.title}</h3>
                          {record.is_completed && (
                            <span className="material-symbols-outlined text-[16px] text-green-500">check_circle</span>
                          )}
                        </div>
                        <div className="flex items-center gap-2 text-xs text-on-surface-variant">
                          {record.status && (
                            <span className="px-2 py-1 bg-surface-container-highest rounded">{record.status}</span>
                          )}
                          {record.assigned_to && (
                            <span className="px-2 py-1 bg-primary/10 text-primary rounded">Assigned</span>
                          )}
                          {record.due_date && (
                            <span>{record.due_date.split('T')[0]}</span>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
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
