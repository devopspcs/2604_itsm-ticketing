import { useRef, useState } from 'react'
import { useAttachments } from '../../hooks/useAttachments'
import { formatFileSize, formatDateTime, formatRelativeTime } from '../../utils/jira.utils'

interface AttachmentSectionProps {
  recordId: string
  projectId: string
  currentUserId: string
}

export function AttachmentSection({ recordId, projectId, currentUserId }: AttachmentSectionProps) {
  const { attachments, loading, error, uploading, uploadAttachment, deleteAttachment, isImageFile } =
    useAttachments(projectId, recordId)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [previewId, setPreviewId] = useState<string | null>(null)
  const [uploadProgress, setUploadProgress] = useState(0)

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      // Validate file size (50MB max)
      if (file.size > 50 * 1024 * 1024) {
        alert('File terlalu besar. Maksimal 50MB.')
        return
      }

      // Validate file type
      const allowedTypes = [
        'image/jpeg', 'image/png', 'image/gif',
        'application/pdf',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        'application/vnd.ms-excel',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        'application/zip',
        'application/x-rar-compressed',
        'text/plain',
      ]

      if (!allowedTypes.includes(file.type)) {
        alert('Tipe file tidak didukung.')
        return
      }

      await uploadAttachment(file)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
      setUploadProgress(0)
    }
  }

  const handleDeleteAttachment = async (attachmentId: string) => {
    if (confirm('Hapus lampiran ini?')) {
      await deleteAttachment(attachmentId)
    }
  }

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
  }

  const handleDrop = async (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()

    const files = e.dataTransfer.files
    if (files.length > 0) {
      const file = files[0]
      if (file.size > 50 * 1024 * 1024) {
        alert('File terlalu besar. Maksimal 50MB.')
        return
      }
      await uploadAttachment(file)
    }
  }

  const getAttachmentUrl = (attachmentId: string) => {
    return `/api/v1/projects/${projectId}/attachments/${attachmentId}/download`
  }

  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-on-surface">Lampiran ({attachments.length})</h3>

      {error && <div className="text-error text-sm p-2 bg-error/10 rounded">{error}</div>}

      {/* Attachments List */}
      <div className="space-y-2 max-h-64 overflow-y-auto">
        {loading ? (
          <div className="text-on-surface-variant text-sm">Memuat lampiran...</div>
        ) : attachments.length === 0 ? (
          <div className="text-on-surface-variant text-sm">Belum ada lampiran</div>
        ) : (
          attachments.map(attachment => (
            <div key={attachment.id} className="bg-surface-container-low p-3 rounded-lg">
              <div className="flex items-start justify-between gap-2">
                <div className="flex items-start gap-2 flex-1 min-w-0">
                  <span className="material-symbols-outlined text-on-surface-variant flex-shrink-0 mt-0.5">
                    attach_file
                  </span>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-on-surface truncate">{attachment.file_name}</p>
                    <p className="text-[10px] text-on-surface-variant">
                      {formatFileSize(attachment.file_size)} • {attachment.uploader_name} •{' '}
                      <span title={formatDateTime(attachment.created_at)}>
                        {formatRelativeTime(attachment.created_at)}
                      </span>
                    </p>
                  </div>
                </div>
                <div className="flex gap-1 flex-shrink-0">
                  {isImageFile(attachment.file_type) && (
                    <button
                      onClick={() => setPreviewId(previewId === attachment.id ? null : attachment.id)}
                      className="text-[10px] text-primary hover:underline"
                      aria-label="Preview image"
                    >
                      👁️
                    </button>
                  )}
                  <a
                    href={getAttachmentUrl(attachment.id)}
                    download={attachment.file_name}
                    className="text-[10px] text-primary hover:underline"
                    aria-label="Download file"
                  >
                    ⬇️
                  </a>
                  {attachment.uploader_id === currentUserId && (
                    <button
                      onClick={() => handleDeleteAttachment(attachment.id)}
                      className="text-[10px] text-error hover:underline"
                      aria-label="Delete attachment"
                    >
                      🗑️
                    </button>
                  )}
                </div>
              </div>

              {/* Image Preview */}
              {previewId === attachment.id && isImageFile(attachment.file_type) && (
                <div className="mt-2 pt-2 border-t border-outline-variant">
                  <img
                    src={getAttachmentUrl(attachment.id)}
                    alt={attachment.file_name}
                    className="max-w-xs max-h-48 rounded"
                  />
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* Upload Area */}
      <div
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        className="border-2 border-dashed border-outline-variant rounded-lg p-4 text-center hover:border-primary transition-colors"
      >
        <input
          ref={fileInputRef}
          type="file"
          onChange={handleFileSelect}
          disabled={uploading}
          className="hidden"
          aria-label="Upload file"
        />

        {uploading ? (
          <div className="space-y-2">
            <p className="text-sm text-on-surface-variant">Mengunggah...</p>
            <div className="w-full bg-surface-container-high rounded-full h-2 overflow-hidden">
              <div
                className="bg-primary h-full transition-all"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
            <p className="text-[10px] text-on-surface-variant">{uploadProgress}%</p>
          </div>
        ) : (
          <>
            <button
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
              className="text-sm text-primary hover:underline disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Klik untuk mengunggah
            </button>
            <p className="text-[10px] text-on-surface-variant mt-1">atau drag & drop file di sini</p>
            <p className="text-[10px] text-on-surface-variant mt-2">
              Maksimal 50MB • Tipe: JPG, PNG, GIF, PDF, DOC, XLS, ZIP, TXT
            </p>
          </>
        )}
      </div>
    </div>
  )
}
