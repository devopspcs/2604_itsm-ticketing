import { useCallback, useEffect, useState } from 'react'
import type { DragEndEvent, DragStartEvent } from '@dnd-kit/core'
import { projectService } from '../services/project.service'
import type { ProjectDetail, ProjectRecord } from '../types/project'

interface FilterState {
  search: string
  assignee: string
  dueDateFrom: string
  dueDateTo: string
}

export function useProjectBoard(projectId: string | undefined) {
  const [project, setProject] = useState<ProjectDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [dragError, setDragError] = useState<string | null>(null)
  const [activeRecord, setActiveRecord] = useState<ProjectRecord | null>(null)
  const [filters, setFilters] = useState<FilterState>({ search: '', assignee: '', dueDateFrom: '', dueDateTo: '' })

  const fetchProject = useCallback(async () => {
    if (!projectId) return
    try {
      setLoading(true)
      const res = await projectService.get(projectId)
      setProject(res.data)
    } catch {
      setError('Gagal memuat project')
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => { fetchProject() }, [fetchProject])

  const handleDragStart = (event: DragStartEvent) => {
    const record = event.active.data.current?.record as ProjectRecord | undefined
    setActiveRecord(record ?? null)
    setDragError(null)
  }

  const handleDragEnd = async (event: DragEndEvent) => {
    setActiveRecord(null)
    const { active, over } = event
    if (!over || !project) return

    const record = active.data.current?.record as ProjectRecord | undefined
    if (!record) return

    const targetColumnId = over.id as string
    if (record.column_id === targetColumnId) return

    // Optimistic update
    const prevProject = { ...project, columns: project.columns.map(c => ({ ...c, records: [...c.records] })) }
    setProject(prev => {
      if (!prev) return prev
      const cols = prev.columns.map(c => {
        if (c.id === record.column_id) {
          return { ...c, records: c.records.filter(r => r.id !== record.id) }
        }
        if (c.id === targetColumnId) {
          return { ...c, records: [...c.records, { ...record, column_id: targetColumnId }] }
        }
        return c
      })
      return { ...prev, columns: cols }
    })

    try {
      const targetCol = project.columns.find(c => c.id === targetColumnId)
      const position = targetCol ? targetCol.records.length : 0
      await projectService.moveRecord(project.id, record.id, { target_column_id: targetColumnId, position })
    } catch {
      setProject(prevProject)
      setDragError(`Gagal memindahkan record "${record.title}"`)
      setTimeout(() => setDragError(null), 4000)
    }
  }

  const addRecord = async (columnId: string, title: string) => {
    if (!project) return
    try {
      const res = await projectService.createRecord(project.id, { column_id: columnId, title })
      setProject(prev => {
        if (!prev) return prev
        return {
          ...prev,
          columns: prev.columns.map(c =>
            c.id === columnId ? { ...c, records: [...c.records, res.data] } : c
          ),
        }
      })
    } catch {
      setError('Gagal menambah record')
    }
  }

  const addColumn = async (name: string) => {
    if (!project) return
    try {
      const res = await projectService.createColumn(project.id, { name })
      setProject(prev => prev ? { ...prev, columns: [...prev.columns, { ...res.data, records: [] }] } : prev)
    } catch {
      setError('Gagal menambah kolom')
    }
  }

  const editColumn = async (columnId: string, name: string) => {
    if (!project) return
    try {
      await projectService.updateColumn(project.id, columnId, { name })
      setProject(prev => {
        if (!prev) return prev
        return { ...prev, columns: prev.columns.map(c => c.id === columnId ? { ...c, name } : c) }
      })
    } catch {
      setError('Gagal mengubah kolom')
    }
  }

  const deleteColumn = async (columnId: string) => {
    if (!project) return
    try {
      await projectService.deleteColumn(project.id, columnId)
      setProject(prev => prev ? { ...prev, columns: prev.columns.filter(c => c.id !== columnId) } : prev)
    } catch {
      setError('Kolom harus dikosongkan terlebih dahulu')
      setTimeout(() => setError(null), 4000)
    }
  }

  // Apply client-side filters
  const getFilteredRecords = (records: ProjectRecord[]) => {
    return records.filter(r => {
      if (filters.search && !r.title.toLowerCase().includes(filters.search.toLowerCase())) return false
      if (filters.assignee && r.assigned_to !== filters.assignee) return false
      if (filters.dueDateFrom && (!r.due_date || r.due_date < filters.dueDateFrom)) return false
      if (filters.dueDateTo && (!r.due_date || r.due_date > filters.dueDateTo)) return false
      return true
    })
  }

  return {
    project, loading, error, dragError, activeRecord, filters,
    setFilters, handleDragStart, handleDragEnd,
    addRecord, addColumn, editColumn, deleteColumn,
    getFilteredRecords, refresh: fetchProject,
  }
}
