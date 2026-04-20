import { useState, useEffect, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { projectService } from '../services/project.service'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { ReleaseWithProgress } from '../types/project'

export function ReleasesPage() {
  const { id: projectId } = useParams()
  const [releases, setReleases] = useState<ReleaseWithProgress[]>([])
  const [loading, setLoading] = useState(true)
  const [showAddRelease, setShowAddRelease] = useState(false)
  const [newReleaseName, setNewReleaseName] = useState('')
  const [newReleaseVersion, setNewReleaseVersion] = useState('')

  const fetchReleases = useCallback(async () => {
    if (!projectId) return
    try {
      setLoading(true)
      const res = await projectService.listReleases(projectId)
      setReleases(res.data ?? [])
    } catch (err) {
      console.error('Failed to fetch releases:', err)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => { fetchReleases() }, [fetchReleases])

  const handleAddRelease = async () => {
    if (!projectId || !newReleaseName.trim() || !newReleaseVersion.trim()) return
    try {
      await projectService.createRelease(projectId, {
        name: newReleaseName,
        version: newReleaseVersion,
      })
      setNewReleaseName('')
      setNewReleaseVersion('')
      setShowAddRelease(false)
      fetchReleases()
    } catch (err) {
      console.error('Failed to create release:', err)
    }
  }

  const handleDeleteRelease = async (releaseId: string) => {
    if (!projectId || !confirm('Are you sure you want to delete this release?')) return
    try {
      await projectService.deleteRelease(projectId, releaseId)
      setReleases(prev => prev.filter(r => r.id !== releaseId))
    } catch (err) {
      console.error('Failed to delete release:', err)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'Planning': return 'bg-blue-500/20 text-blue-700'
      case 'In Progress': return 'bg-orange-500/20 text-orange-700'
      case 'Released': return 'bg-green-500/20 text-green-700'
      case 'Archived': return 'bg-gray-500/20 text-gray-700'
      default: return 'bg-gray-500/20 text-gray-700'
    }
  }

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Releases</h1>
              <p className="text-sm text-on-surface-variant mt-1">Manage project releases and versions</p>
            </div>
            <button
              onClick={() => setShowAddRelease(!showAddRelease)}
              className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
            >
              <span className="material-symbols-outlined text-[18px]">add</span>
              New Release
            </button>
          </div>
        </div>

        <div className="flex-1 overflow-auto px-8 pb-8">
          {showAddRelease && (
            <div className="bg-surface-container-low rounded-lg p-6 mb-6 border border-outline-variant/10">
              <h2 className="text-lg font-semibold text-on-surface mb-4">Create New Release</h2>
              <div className="space-y-4">
                <input type="text" value={newReleaseName} onChange={e => setNewReleaseName(e.target.value)}
                  placeholder="Release name (e.g., Version 1.0)"
                  className="w-full px-4 py-2 bg-surface-container-highest border border-outline-variant/20 rounded-lg text-on-surface placeholder-on-surface-variant focus:outline-none focus:border-primary" />
                <input type="text" value={newReleaseVersion} onChange={e => setNewReleaseVersion(e.target.value)}
                  placeholder="Version (e.g., 1.0.0)"
                  className="w-full px-4 py-2 bg-surface-container-highest border border-outline-variant/20 rounded-lg text-on-surface placeholder-on-surface-variant focus:outline-none focus:border-primary" />
                <div className="flex gap-2">
                  <button onClick={handleAddRelease}
                    className="px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors">Create</button>
                  <button onClick={() => setShowAddRelease(false)}
                    className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-lg transition-colors">Cancel</button>
                </div>
              </div>
            </div>
          )}

          {loading ? (
            <div className="space-y-4">
              {[1, 2].map(i => (
                <div key={i} className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 animate-pulse">
                  <div className="h-5 w-40 bg-surface-container-high rounded mb-3" />
                  <div className="h-4 w-24 bg-surface-container-high rounded mb-4" />
                  <div className="h-2 w-full bg-surface-container-high rounded" />
                </div>
              ))}
            </div>
          ) : releases.length === 0 ? (
            <div className="text-center py-12">
              <span className="material-symbols-outlined text-6xl text-on-surface-variant/30 block mb-4">tag</span>
              <p className="text-on-surface-variant">No releases yet</p>
            </div>
          ) : (
            <div className="space-y-4">
              {releases.map(release => (
                <div key={release.id} className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 hover:border-outline-variant/20 transition-colors">
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <div className="flex items-center gap-3 mb-2">
                        <h3 className="text-lg font-semibold text-on-surface">{release.name}</h3>
                        <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(release.status)}`}>{release.status}</span>
                      </div>
                      <p className="text-sm text-on-surface-variant">Version {release.version}</p>
                    </div>
                    <button onClick={() => handleDeleteRelease(release.id)}
                      className="p-2 text-error hover:bg-error-container/20 rounded-lg transition-colors">
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </div>
                  {release.description && <p className="text-sm text-on-surface-variant mb-4">{release.description}</p>}
                  <div className="grid grid-cols-3 gap-4 mb-4">
                    <div>
                      <p className="text-xs text-on-surface-variant mb-1">Release Date</p>
                      <p className="text-sm font-medium text-on-surface">{release.release_date?.split('T')[0] ?? '—'}</p>
                    </div>
                    <div>
                      <p className="text-xs text-on-surface-variant mb-1">Total Issues</p>
                      <p className="text-sm font-medium text-on-surface">{release.total_records}</p>
                    </div>
                    <div>
                      <p className="text-xs text-on-surface-variant mb-1">Completed</p>
                      <p className="text-sm font-medium text-on-surface">{release.completed_count}/{release.total_records}</p>
                    </div>
                  </div>
                  <div className="w-full bg-surface-container-highest rounded-full h-2 overflow-hidden">
                    <div className="bg-primary h-full transition-all" style={{ width: `${release.progress_percent}%` }} />
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
