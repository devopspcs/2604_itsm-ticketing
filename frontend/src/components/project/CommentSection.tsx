import { useState, useRef, useEffect } from 'react'
import { useComments } from '../../hooks/useComments'
import { formatDateTime, formatRelativeTime, parseMentions } from '../../utils/jira.utils'

interface CommentSectionProps {
  recordId: string
  projectId: string
  currentUserId: string
  projectMembers?: Array<{ id: string; name: string }>
}

export function CommentSection({ recordId, projectId, currentUserId, projectMembers = [] }: CommentSectionProps) {
  const { comments, loading, error, addComment, updateComment, deleteComment } = useComments(projectId, recordId)
  const [newComment, setNewComment] = useState('')
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editText, setEditText] = useState('')
  const [showMentions, setShowMentions] = useState(false)
  const [mentionSuggestions, setMentionSuggestions] = useState<Array<{ id: string; name: string }>>([])
  const [mentionIndex, setMentionIndex] = useState(0)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  const handleAddComment = async () => {
    if (!newComment.trim()) return
    await addComment(newComment)
    setNewComment('')
    setShowMentions(false)
  }

  const handleUpdateComment = async (commentId: string) => {
    if (!editText.trim()) return
    await updateComment(commentId, editText)
    setEditingId(null)
    setEditText('')
  }

  const handleDeleteComment = async (commentId: string) => {
    if (confirm('Hapus komentar ini?')) {
      await deleteComment(commentId)
    }
  }

  const handleCommentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = e.target.value
    setNewComment(text)

    // Check for @ mentions
    const lastAtIndex = text.lastIndexOf('@')
    if (lastAtIndex !== -1) {
      const afterAt = text.substring(lastAtIndex + 1)
      if (afterAt && !afterAt.includes(' ')) {
        const filtered = projectMembers.filter(m =>
          m.name.toLowerCase().includes(afterAt.toLowerCase())
        )
        setMentionSuggestions(filtered)
        setShowMentions(filtered.length > 0)
        setMentionIndex(0)
      } else {
        setShowMentions(false)
      }
    } else {
      setShowMentions(false)
    }
  }

  const handleMentionSelect = (member: { id: string; name: string }) => {
    const lastAtIndex = newComment.lastIndexOf('@')
    const beforeAt = newComment.substring(0, lastAtIndex)
    const afterAt = newComment.substring(lastAtIndex + 1)
    const afterSpace = afterAt.indexOf(' ')
    const afterText = afterSpace !== -1 ? afterAt.substring(afterSpace) : ''

    const newText = `${beforeAt}@${member.name}${afterText}`
    setNewComment(newText)
    setShowMentions(false)

    // Focus back to textarea
    setTimeout(() => textareaRef.current?.focus(), 0)
  }

  const mentions = parseMentions(newComment)

  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-on-surface">Komentar ({comments.length})</h3>

      {error && <div className="text-error text-sm p-2 bg-error/10 rounded">{error}</div>}

      {/* Comments List */}
      <div className="space-y-3 max-h-96 overflow-y-auto">
        {loading ? (
          <div className="text-on-surface-variant text-sm">Memuat komentar...</div>
        ) : comments.length === 0 ? (
          <div className="text-on-surface-variant text-sm">Belum ada komentar</div>
        ) : (
          comments.map(comment => (
            <div key={comment.id} className="bg-surface-container-low p-3 rounded-lg">
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center gap-2">
                  <div className="w-6 h-6 rounded-full bg-primary/20 flex items-center justify-center flex-shrink-0">
                    <span className="text-[10px] font-bold text-primary">
                      {(comment.author_name || 'U').charAt(0).toUpperCase()}
                    </span>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-on-surface">{comment.author_name || 'Unknown'}</p>
                    <p className="text-[10px] text-on-surface-variant" title={formatDateTime(comment.created_at)}>
                      {formatRelativeTime(comment.created_at)}
                    </p>
                  </div>
                </div>
                {comment.author_id === currentUserId && (
                  <div className="flex gap-1">
                    <button
                      onClick={() => {
                        setEditingId(comment.id)
                        setEditText(comment.text)
                      }}
                      className="text-[10px] text-primary hover:underline"
                      aria-label="Edit comment"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteComment(comment.id)}
                      className="text-[10px] text-error hover:underline"
                      aria-label="Delete comment"
                    >
                      Hapus
                    </button>
                  </div>
                )}
              </div>

              {editingId === comment.id ? (
                <div className="space-y-2">
                  <textarea
                    value={editText}
                    onChange={e => setEditText(e.target.value)}
                    className="w-full p-2 rounded border border-outline-variant text-sm"
                    rows={3}
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleUpdateComment(comment.id)}
                      className="px-3 py-1 bg-primary text-on-primary rounded text-sm hover:opacity-90"
                    >
                      Simpan
                    </button>
                    <button
                      onClick={() => setEditingId(null)}
                      className="px-3 py-1 bg-surface-container-high text-on-surface rounded text-sm hover:opacity-90"
                    >
                      Batal
                    </button>
                  </div>
                </div>
              ) : (
                <p className="text-sm text-on-surface whitespace-pre-wrap break-words">
                  {comment.text.split(/(@\w+)/g).map((part, i) =>
                    part.startsWith('@') ? (
                      <span key={i} className="bg-primary/20 text-primary px-1 rounded">
                        {part}
                      </span>
                    ) : (
                      part
                    )
                  )}
                </p>
              )}
            </div>
          ))
        )}
      </div>

      {/* Add Comment Form */}
      <div className="space-y-2 border-t border-outline-variant pt-3">
        <div className="relative">
          <textarea
            ref={textareaRef}
            value={newComment}
            onChange={handleCommentChange}
            placeholder="Tambah komentar... (gunakan @ untuk mention)"
            className="w-full p-3 rounded border border-outline-variant text-sm resize-none"
            rows={3}
          />

          {/* Mention Suggestions Dropdown */}
          {showMentions && mentionSuggestions.length > 0 && (
            <div className="absolute bottom-full left-0 mb-1 bg-surface-container-low border border-outline-variant rounded-lg shadow-lg z-10 w-full max-h-40 overflow-y-auto">
              {mentionSuggestions.map((member, idx) => (
                <button
                  key={member.id}
                  onClick={() => handleMentionSelect(member)}
                  className={`w-full text-left px-3 py-2 text-sm ${
                    idx === mentionIndex
                      ? 'bg-primary/20 text-primary'
                      : 'hover:bg-surface-container-high text-on-surface'
                  }`}
                >
                  @{member.name}
                </button>
              ))}
            </div>
          )}
        </div>

        {mentions.length > 0 && (
          <div className="text-[10px] text-on-surface-variant bg-surface-container-high p-2 rounded">
            Mentions: {mentions.map(m => `@${m}`).join(', ')}
          </div>
        )}

        <button
          onClick={handleAddComment}
          disabled={!newComment.trim()}
          className="w-full px-4 py-2 bg-primary text-on-primary rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-90"
        >
          Kirim Komentar
        </button>
      </div>
    </div>
  )
}
