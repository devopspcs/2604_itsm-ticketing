import { useDraggable } from '@dnd-kit/core'
import type { JiraProjectRecord, IssueType } from '../../types/jira'
import { getIssueTypeIcon, formatDate, isOverdue } from '../../utils/jira.utils'

interface RecordCardProps {
  record: JiraProjectRecord
  issueType?: IssueType
  onClick: () => void
}

export function RecordCard({ record, issueType, onClick }: RecordCardProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: record.id,
    data: { record },
  })

  const style = transform
    ? {
        transform: `translate(${transform.x}px, ${transform.y}px)`,
        boxShadow: isDragging ? '0 20px 40px rgba(25, 28, 30, 0.06)' : undefined,
        zIndex: isDragging ? 50 : undefined,
        opacity: isDragging ? 0.9 : undefined,
      }
    : undefined

  const overdue = record.due_date && !record.is_completed && isOverdue(record.due_date)

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`bg-surface-container-lowest p-3 rounded-xl border border-outline-variant/5 transition-all ${
        isDragging ? '' : 'hover:bg-surface-bright'
      } ${record.is_completed ? 'opacity-60' : ''} flex gap-2`}
    >
      {/* Drag Handle */}
      <div
        {...listeners}
        {...attributes}
        className="flex-shrink-0 cursor-grab active:cursor-grabbing pt-0.5 text-on-surface-variant/40 hover:text-on-surface-variant"
        title="Drag to move"
      >
        <span className="material-symbols-outlined text-[16px]">drag_indicator</span>
      </div>

      {/* Clickable Content */}
      <div className="flex-1 cursor-pointer min-w-0" onClick={onClick}>
        <div className="flex items-start gap-2 mb-2">
          <span className="text-lg flex-shrink-0">{getIssueTypeIcon(issueType?.name || 'Task')}</span>
          <h3 className={`text-sm font-semibold text-on-surface leading-snug line-clamp-2 flex-1 ${
            record.is_completed ? 'line-through' : ''
          }`}>
            {record.title}
          </h3>
        </div>

        {record.labels && record.labels.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-2">
            {record.labels.slice(0, 2).map(label => (
              <span key={label.id} className="text-[10px] px-2 py-0.5 rounded-full text-white"
                style={{ backgroundColor: label.color }}>{label.name}</span>
            ))}
            {record.labels.length > 2 && (
              <span className="text-[10px] px-2 py-0.5 rounded-full bg-surface-container-high text-on-surface-variant">
                +{record.labels.length - 2}
              </span>
            )}
          </div>
        )}

        <div className="flex items-center justify-between text-[10px]">
          <div className="flex items-center gap-2">
            {record.due_date && (
              <span className={`flex items-center gap-1 ${overdue ? 'text-error' : 'text-on-surface-variant'}`}>
                <span className="material-symbols-outlined text-[12px]">schedule</span>
                {formatDate(record.due_date)}
              </span>
            )}
          </div>
          <div className="flex items-center gap-1">
            {record.comments_count && record.comments_count > 0 && (
              <span className="flex items-center gap-0.5 text-on-surface-variant">
                <span className="material-symbols-outlined text-[12px]">chat</span>
                {record.comments_count}
              </span>
            )}
            {record.attachments_count && record.attachments_count > 0 && (
              <span className="flex items-center gap-0.5 text-on-surface-variant">
                <span className="material-symbols-outlined text-[12px]">attach_file</span>
                {record.attachments_count}
              </span>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
