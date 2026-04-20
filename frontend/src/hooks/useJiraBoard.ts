import { useCallback, useEffect, useState } from 'react'
import type { DragEndEvent } from '@dnd-kit/core'
import { jiraService } from '../services/jira.service'
import type { Sprint, JiraProjectRecord, WorkflowStatus } from '../types/jira'

interface SprintBoardState {
  sprint: Sprint | null
  records: JiraProjectRecord[]
  statuses: WorkflowStatus[]
  loading: boolean
  error: string | null
  transitioningRecordId: string | null
}

export function useJiraBoard(projectId: string | undefined, sprintId: string | undefined) {
  const [state, setState] = useState<SprintBoardState>({
    sprint: null,
    records: [],
    statuses: [],
    loading: true,
    error: null,
    transitioningRecordId: null,
  })

  const fetchSprintData = useCallback(async () => {
    if (!projectId || !sprintId) return
    try {
      setState(prev => ({ ...prev, loading: true, error: null }))
      const [sprintRes, recordsRes] = await Promise.all([
        jiraService.getActiveSprint(projectId),
        jiraService.getSprintRecords(projectId, sprintId),
      ])
      setState(prev => ({
        ...prev,
        sprint: sprintRes.data,
        records: recordsRes.data,
        loading: false,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal memuat sprint board'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        loading: false,
      }))
    }
  }, [projectId, sprintId])

  useEffect(() => {
    fetchSprintData()
  }, [fetchSprintData])

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event
    if (!over) return

    const record = state.records.find(r => r.id === active.id)
    if (!record) return

    const newStatus = over.id as string
    if (record.status === newStatus) return

    // Optimistic update
    const prevRecords = state.records
    setState(prev => ({
      ...prev,
      records: prev.records.map(r =>
        r.id === record.id ? { ...r, status: newStatus } : r
      ),
      transitioningRecordId: record.id,
    }))

    try {
      await jiraService.transitionRecord(projectId!, record.id, { to_status_id: newStatus })
      setState(prev => ({ ...prev, transitioningRecordId: null }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal mengubah status'
      setState(prev => ({
        ...prev,
        records: prevRecords,
        error: errorMessage,
        transitioningRecordId: null,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const getRecordsByStatus = (status: string) => {
    return state.records.filter(r => r.status === status)
  }

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    ...state,
    handleDragEnd,
    getRecordsByStatus,
    refresh: fetchSprintData,
    clearError,
  }
}
