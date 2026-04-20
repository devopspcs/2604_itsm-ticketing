import { useCallback, useEffect, useState } from 'react'
import { jiraService } from '../services/jira.service'
import type { JiraProjectRecord, Sprint } from '../types/jira'

interface BacklogState {
  records: JiraProjectRecord[]
  sprints: Sprint[]
  loading: boolean
  error: string | null
  selectedRecords: Set<string>
  reordering: boolean
  assigning: boolean
}

export function useBacklog(projectId: string | undefined) {
  const [state, setState] = useState<BacklogState>({
    records: [],
    sprints: [],
    loading: true,
    error: null,
    selectedRecords: new Set(),
    reordering: false,
    assigning: false,
  })

  const fetchBacklogData = useCallback(async () => {
    if (!projectId) return
    try {
      setState(prev => ({ ...prev, loading: true, error: null }))
      const [backlogRes, sprintsRes] = await Promise.all([
        jiraService.getBacklog(projectId),
        jiraService.listSprints(projectId),
      ])
      setState(prev => ({
        ...prev,
        records: backlogRes.data,
        sprints: sprintsRes.data,
        loading: false,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal memuat backlog'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        loading: false,
      }))
    }
  }, [projectId])

  useEffect(() => {
    fetchBacklogData()
  }, [fetchBacklogData])

  const reorderRecords = async (recordIds: string[]) => {
    if (!projectId) return
    try {
      setState(prev => ({ ...prev, reordering: true, error: null }))
      await jiraService.reorderBacklog(projectId, { record_ids: recordIds })
      setState(prev => ({
        ...prev,
        records: recordIds
          .map(id => prev.records.find(r => r.id === id))
          .filter((r): r is JiraProjectRecord => r !== undefined),
        reordering: false,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal mengubah urutan'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        reordering: false,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const assignToSprint = async (sprintId: string, recordIds: string[]) => {
    if (!projectId) return
    try {
      setState(prev => ({ ...prev, assigning: true, error: null }))
      await jiraService.bulkAssignToSprint(projectId, { sprint_id: sprintId, record_ids: recordIds })
      setState(prev => ({
        ...prev,
        records: prev.records.filter(r => !recordIds.includes(r.id)),
        selectedRecords: new Set(),
        assigning: false,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menugaskan ke sprint'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        assigning: false,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const toggleRecordSelection = (recordId: string) => {
    setState(prev => {
      const newSelected = new Set(prev.selectedRecords)
      if (newSelected.has(recordId)) {
        newSelected.delete(recordId)
      } else {
        newSelected.add(recordId)
      }
      return { ...prev, selectedRecords: newSelected }
    })
  }

  const clearSelection = () => {
    setState(prev => ({ ...prev, selectedRecords: new Set() }))
  }

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    ...state,
    reorderRecords,
    assignToSprint,
    toggleRecordSelection,
    clearSelection,
    refresh: fetchBacklogData,
    clearError,
  }
}
