import { useCallback, useEffect, useState } from 'react'
import { jiraService } from '../services/jira.service'
import type { Comment } from '../types/jira'

interface CommentsState {
  comments: Comment[]
  loading: boolean
  error: string | null
  addingComment: boolean
  updatingCommentId: string | null
  deletingCommentId: string | null
}

export function useComments(projectId: string | undefined, recordId: string | undefined) {
  const [state, setState] = useState<CommentsState>({
    comments: [],
    loading: true,
    error: null,
    addingComment: false,
    updatingCommentId: null,
    deletingCommentId: null,
  })

  const fetchComments = useCallback(async () => {
    if (!recordId || !projectId) return
    try {
      setState(prev => ({ ...prev, loading: true, error: null }))
      const res = await jiraService.listComments(projectId, recordId)
      setState(prev => ({
        ...prev,
        comments: res.data,
        loading: false,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal memuat komentar'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        loading: false,
      }))
    }
  }, [projectId, recordId])

  useEffect(() => {
    fetchComments()
  }, [fetchComments])

  const addComment = async (text: string) => {
    if (!recordId || !projectId) return
    try {
      setState(prev => ({ ...prev, addingComment: true, error: null }))
      const res = await jiraService.addComment(projectId, recordId, { text })
      setState(prev => ({
        ...prev,
        comments: [...prev.comments, res.data],
        addingComment: false,
      }))
      return res.data
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menambah komentar'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        addingComment: false,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const updateComment = async (commentId: string, text: string) => {
    try {
      setState(prev => ({ ...prev, updatingCommentId: commentId, error: null }))
      const res = await jiraService.updateComment(projectId!, commentId, { text })
      setState(prev => ({
        ...prev,
        comments: prev.comments.map(c => (c.id === commentId ? res.data : c)),
        updatingCommentId: null,
      }))
      return res.data
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal mengubah komentar'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        updatingCommentId: null,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const deleteComment = async (commentId: string) => {
    try {
      setState(prev => ({ ...prev, deletingCommentId: commentId, error: null }))
      await jiraService.deleteComment(projectId!, commentId)
      setState(prev => ({
        ...prev,
        comments: prev.comments.filter(c => c.id !== commentId),
        deletingCommentId: null,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menghapus komentar'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        deletingCommentId: null,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const parseMentions = (text: string): string[] => {
    const mentionRegex = /@(\w+)/g
    const matches = text.match(mentionRegex) || []
    return matches.map(m => m.substring(1))
  }

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    ...state,
    addComment,
    updateComment,
    deleteComment,
    parseMentions,
    refresh: fetchComments,
    clearError,
  }
}
