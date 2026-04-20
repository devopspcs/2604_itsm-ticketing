import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { projectService } from '../services/project.service'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { ReportsSummary, VelocityDataPoint, BurndownData } from '../types/project'

export function ReportsPage() {
  const { id: projectId } = useParams()
  const [loading, setLoading] = useState(true)
  const [summary, setSummary] = useState<ReportsSummary | null>(null)
  const [velocity, setVelocity] = useState<VelocityDataPoint[]>([])
  const [burndown, setBurndown] = useState<BurndownData | null>(null)

  useEffect(() => {
    if (!projectId) return
    const load = async () => {
      setLoading(true)
      try {
        const [sumRes, velRes, burnRes] = await Promise.allSettled([
          projectService.getReportsSummary(projectId),
          projectService.getReportsVelocity(projectId),
          projectService.getReportsBurndown(projectId),
        ])
        if (sumRes.status === 'fulfilled') setSummary(sumRes.value.data)
        if (velRes.status === 'fulfilled') setVelocity(velRes.value.data ?? [])
        if (burnRes.status === 'fulfilled') setBurndown(burnRes.value.data)
      } catch (err) {
        console.error('Failed to fetch reports:', err)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [projectId])

  const completionRate = summary && summary.total_records > 0
    ? Math.round((summary.completed_count / summary.total_records) * 100)
    : 0

  const maxVelocity = velocity.length > 0
    ? Math.max(...velocity.map(v => v.total_records), 1)
    : 1

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">
            Reports
          </h1>
          <p className="text-sm text-on-surface-variant mt-1">Project analytics and statistics</p>
        </div>

        <div className="flex-1 overflow-auto px-8 pb-8">
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 animate-pulse">
                  <div className="h-4 w-24 bg-surface-container-high rounded mb-3" />
                  <div className="h-8 w-16 bg-surface-container-high rounded" />
                </div>
              ))}
            </div>
          ) : (
            <>
              {/* Summary Cards */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
                <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-on-surface-variant mb-1">Total Issues</p>
                      <p className="text-3xl font-bold text-on-surface">{summary?.total_records ?? 0}</p>
                    </div>
                    <span className="material-symbols-outlined text-4xl text-primary/30">task</span>
                  </div>
                </div>
                <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-on-surface-variant mb-1">Completed</p>
                      <p className="text-3xl font-bold text-on-surface">{summary?.completed_count ?? 0}</p>
                    </div>
                    <span className="material-symbols-outlined text-4xl text-green-500/30">check_circle</span>
                  </div>
                </div>
                <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-on-surface-variant mb-1">Open</p>
                      <p className="text-3xl font-bold text-on-surface">{summary?.open_count ?? 0}</p>
                    </div>
                    <span className="material-symbols-outlined text-4xl text-orange-500/30">schedule</span>
                  </div>
                </div>
                <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-on-surface-variant mb-1">Completion Rate</p>
                      <p className="text-3xl font-bold text-on-surface">{completionRate}%</p>
                    </div>
                    <span className="material-symbols-outlined text-4xl text-purple-500/30">percent</span>
                  </div>
                </div>
              </div>

              {/* Status Breakdown */}
              {summary && Object.keys(summary.by_status).length > 0 && (
                <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 mb-8">
                  <h2 className="text-lg font-semibold text-on-surface mb-4">Status Breakdown</h2>
                  <div className="space-y-3">
                    {Object.entries(summary.by_status).map(([status, count]) => {
                      const pct = summary.total_records > 0 ? Math.round((count / summary.total_records) * 100) : 0
                      return (
                        <div key={status}>
                          <div className="flex items-center justify-between mb-1">
                            <span className="text-sm text-on-surface-variant">{status}</span>
                            <span className="text-sm font-semibold text-on-surface">{count} ({pct}%)</span>
                          </div>
                          <div className="w-full bg-surface-container-highest rounded-full h-2 overflow-hidden">
                            <div className="bg-primary h-full transition-all" style={{ width: `${pct}%` }} />
                          </div>
                        </div>
                      )
                    })}
                  </div>
                </div>
              )}

              {/* Velocity Chart */}
              <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 mb-8">
                <h2 className="text-lg font-semibold text-on-surface mb-4">Sprint Velocity</h2>
                {velocity.length === 0 ? (
                  <div className="text-center py-8">
                    <span className="material-symbols-outlined text-4xl text-on-surface-variant/30 block mb-2">bar_chart</span>
                    <p className="text-sm text-on-surface-variant">No completed sprints yet</p>
                  </div>
                ) : (
                  <div className="flex items-end gap-4 h-48">
                    {velocity.map(v => (
                      <div key={v.sprint_name} className="flex-1 flex flex-col items-center gap-1">
                        <div className="w-full flex gap-1 items-end" style={{ height: '160px' }}>
                          <div
                            className="flex-1 bg-primary/30 rounded-t"
                            style={{ height: `${(v.total_records / maxVelocity) * 100}%` }}
                            title={`Total: ${v.total_records}`}
                          />
                          <div
                            className="flex-1 bg-primary rounded-t"
                            style={{ height: `${(v.completed_count / maxVelocity) * 100}%` }}
                            title={`Completed: ${v.completed_count}`}
                          />
                        </div>
                        <span className="text-xs text-on-surface-variant truncate w-full text-center">{v.sprint_name}</span>
                      </div>
                    ))}
                  </div>
                )}
                {velocity.length > 0 && (
                  <div className="flex items-center gap-4 mt-4 text-xs text-on-surface-variant">
                    <span className="flex items-center gap-1"><span className="w-3 h-3 bg-primary/30 rounded inline-block" /> Total</span>
                    <span className="flex items-center gap-1"><span className="w-3 h-3 bg-primary rounded inline-block" /> Completed</span>
                  </div>
                )}
              </div>

              {/* Burndown */}
              <div className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10">
                <h2 className="text-lg font-semibold text-on-surface mb-4">Sprint Burndown</h2>
                {!burndown || !burndown.has_active ? (
                  <div className="text-center py-8">
                    <span className="material-symbols-outlined text-4xl text-on-surface-variant/30 block mb-2">trending_down</span>
                    <p className="text-sm text-on-surface-variant">No active sprint</p>
                  </div>
                ) : (
                  <div>
                    <p className="text-sm text-on-surface-variant mb-2">{burndown.sprint_name}</p>
                    <div className="flex items-center gap-6 mb-4 text-sm">
                      <span className="text-on-surface-variant">Total: <span className="font-semibold text-on-surface">{burndown.total_count}</span></span>
                      <span className="text-on-surface-variant">Done: <span className="font-semibold text-green-600">{burndown.done_count}</span></span>
                      <span className="text-on-surface-variant">Remaining: <span className="font-semibold text-orange-600">{burndown.total_count - burndown.done_count}</span></span>
                    </div>
                    <div className="w-full bg-surface-container-highest rounded-full h-4 overflow-hidden">
                      <div
                        className="bg-green-500 h-full transition-all"
                        style={{ width: `${burndown.total_count > 0 ? (burndown.done_count / burndown.total_count) * 100 : 0}%` }}
                      />
                    </div>
                    <div className="flex justify-between text-xs text-on-surface-variant mt-1">
                      <span>{burndown.start_date?.split('T')[0] ?? ''}</span>
                      <span>{burndown.end_date?.split('T')[0] ?? ''}</span>
                    </div>
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
