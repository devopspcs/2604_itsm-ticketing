import { useCallback, useEffect, useState } from 'react'
import { jiraService } from '../services/jira.service'
import type { Attachment } from '../types/jira'

interface AttachmentsState {
  attachments: Attachment[]
  loading: boolean
  error: string | null
  uploading: boolean
  uploadProgress: number
  deletingAttachmentId: string | null
}

export function useAttachments(projectId: string | undefined, recordId: string | undefined) {
  const [state, setState] = useState<AttachmentsState>({
    attachments: [],
    loading: true,
    error: null,
    uploading: false,
    uploadProgress: 0,
    deletingAttachmentId: null,
  })

  const fetchAttachments = useCallback(async () => {
    if (!recordId || !projectId) return
    try {
      setState(prev => ({ ...prev, loading: true, error: null }))
      const res = await jiraService.listAttachments(projectId, recordId)
      setState(prev => ({
        ...prev,
        attachments: res.data,
        loading: false,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal memuat lampiran'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        loading: false,
      }))
    }
  }, [projectId, recordId])

  useEffect(() => {
    fetchAttachments()
  }, [fetchAttachments])

  const uploadAttachment = async (file: File) => {
    if (!recordId || !projectId) return
    try {
      setState(prev => ({ ...prev, uploading: true, error: null, uploadProgress: 0 }))
      const res = await jiraService.uploadAttachment(projectId, recordId, file)
      setState(prev => ({
        ...prev,
        attachments: [...prev.attachments, res.data],
        uploading: false,
        uploadProgress: 100,
      }))
      setTimeout(() => setState(prev => ({ ...prev, uploadProgress: 0 })), 1000)
      return res.data
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal mengunggah lampiran'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        uploading: false,
        uploadProgress: 0,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const deleteAttachment = async (attachmentId: string) => {
    try {
      setState(prev => ({ ...prev, deletingAttachmentId: attachmentId, error: null }))
      await jiraService.deleteAttachment(projectId!, attachmentId)
      setState(prev => ({
        ...prev,
        attachments: prev.attachments.filter(a => a.id !== attachmentId),
        deletingAttachmentId: null,
      }))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menghapus lampiran'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        deletingAttachmentId: null,
      }))
      setTimeout(() => setState(prev => ({ ...prev, error: null })), 4000)
    }
  }

  const isImageFile = (fileType: string): boolean => {
    return fileType.startsWith('image/')
  }

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    ...state,
    uploadAttachment,
    deleteAttachment,
    isImageFile,
    refresh: fetchAttachments,
    clearError,
  }
}
