import { useState, useEffect, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { projectService } from '../services/project.service'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { ComponentWithCount } from '../types/project'

export function ComponentsPage() {
  const { id: projectId } = useParams()
  const [components, setComponents] = useState<ComponentWithCount[]>([])
  const [loading, setLoading] = useState(true)
  const [showAddComponent, setShowAddComponent] = useState(false)
  const [newComponentName, setNewComponentName] = useState('')
  const [newComponentDescription, setNewComponentDescription] = useState('')

  const fetchComponents = useCallback(async () => {
    if (!projectId) return
    try {
      setLoading(true)
      const res = await projectService.listComponents(projectId)
      setComponents(res.data ?? [])
    } catch (err) {
      console.error('Failed to fetch components:', err)
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => { fetchComponents() }, [fetchComponents])

  const handleAddComponent = async () => {
    if (!projectId || !newComponentName.trim()) return
    try {
      await projectService.createComponent(projectId, {
        name: newComponentName,
        description: newComponentDescription || undefined,
      })
      setNewComponentName('')
      setNewComponentDescription('')
      setShowAddComponent(false)
      fetchComponents()
    } catch (err) {
      console.error('Failed to create component:', err)
    }
  }

  const handleDeleteComponent = async (componentId: string) => {
    if (!projectId || !confirm('Are you sure you want to delete this component?')) return
    try {
      await projectService.deleteComponent(projectId, componentId)
      setComponents(prev => prev.filter(c => c.id !== componentId))
    } catch (err) {
      console.error('Failed to delete component:', err)
    }
  }

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Components</h1>
              <p className="text-sm text-on-surface-variant mt-1">Manage project components</p>
            </div>
            <button onClick={() => setShowAddComponent(!showAddComponent)}
              className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors">
              <span className="material-symbols-outlined text-[18px]">add</span>
              New Component
            </button>
          </div>
        </div>

        <div className="flex-1 overflow-auto px-8 pb-8">
          {showAddComponent && (
            <div className="bg-surface-container-low rounded-lg p-6 mb-6 border border-outline-variant/10">
              <h2 className="text-lg font-semibold text-on-surface mb-4">Create New Component</h2>
              <div className="space-y-4">
                <input type="text" value={newComponentName} onChange={e => setNewComponentName(e.target.value)}
                  placeholder="Component name"
                  className="w-full px-4 py-2 bg-surface-container-highest border border-outline-variant/20 rounded-lg text-on-surface placeholder-on-surface-variant focus:outline-none focus:border-primary" />
                <textarea value={newComponentDescription} onChange={e => setNewComponentDescription(e.target.value)}
                  placeholder="Component description"
                  className="w-full px-4 py-2 bg-surface-container-highest border border-outline-variant/20 rounded-lg text-on-surface placeholder-on-surface-variant focus:outline-none focus:border-primary resize-none"
                  rows={3} />
                <div className="flex gap-2">
                  <button onClick={handleAddComponent}
                    className="px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors">Create</button>
                  <button onClick={() => setShowAddComponent(false)}
                    className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-lg transition-colors">Cancel</button>
                </div>
              </div>
            </div>
          )}

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[1, 2, 3].map(i => (
                <div key={i} className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 animate-pulse">
                  <div className="h-5 w-32 bg-surface-container-high rounded mb-3" />
                  <div className="h-4 w-48 bg-surface-container-high rounded mb-4" />
                  <div className="h-4 w-20 bg-surface-container-high rounded" />
                </div>
              ))}
            </div>
          ) : components.length === 0 ? (
            <div className="text-center py-12">
              <span className="material-symbols-outlined text-6xl text-on-surface-variant/30 block mb-4">widgets</span>
              <p className="text-on-surface-variant">No components yet</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {components.map(component => (
                <div key={component.id} className="bg-surface-container-low rounded-lg p-6 border border-outline-variant/10 hover:border-outline-variant/20 transition-colors">
                  <div className="flex items-start justify-between mb-3">
                    <h3 className="text-lg font-semibold text-on-surface">{component.name}</h3>
                    <button onClick={() => handleDeleteComponent(component.id)}
                      className="p-2 text-error hover:bg-error-container/20 rounded-lg transition-colors">
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </div>
                  {component.description && <p className="text-sm text-on-surface-variant mb-4">{component.description}</p>}
                  <div className="space-y-2 text-sm">
                    <div className="flex items-center justify-between">
                      <span className="text-on-surface-variant">Lead:</span>
                      <span className="text-on-surface font-medium">{component.lead_user_id ?? 'Unassigned'}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-on-surface-variant">Records:</span>
                      <span className="text-on-surface font-medium">{component.record_count}</span>
                    </div>
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
