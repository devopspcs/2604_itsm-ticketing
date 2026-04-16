import { useDraggable } from '@dnd-kit/core'
import type { ProjectRecord } from '../../types/project'

interface ProjectRecordCardProps {
  record: ProjectRecord
  onClick: () => void
}

export function ProjectRecordCard({ record, onClick }: ProjectRecordCardProps) {
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

  const handleClick = () => { if (!isDragging) onClick() }

  const isOverdue = record.due_date && !record.is_completed && new Date(record.due_date) < new Date()

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      onClick={handleClick}
      className={`bg-surface-container-lowest p-3 rounded-xl border border-outline-variant/5 cursor-pointer transition-all ${
        isDragging ? '' : 'hover:bg-surface-bright'
      } ${record.is_completed ? 'opacity-60' : ''}`}
    >
      <h3 className={`text-sm font-semibold text-on-surface leading-snug line-clamp-2 ${
        record.is_completed ? 'line-through' : ''
      }`}>
        {record.title}
      </h3>
      <div className="flex items-center justify-between mt-2">
        <div className="flex items-center gap-2">
          {record.due_date && (
            <span className={`text-[10px] font-medium flex items-center gap-1 ${
              isOverdue ? 'text-error' : 'text-on-surface-variant'
            }`}>
              <span className="material-symbols-outlined text-[12px]">schedule</span>
              {new Date(record.due_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}
            </span>
          )}
        </div>
        {record.assigned_to && !record.assignees?.length && (
          <div className="w-6 h-6 rounded-full bg-primary/20 flex items-center justify-center">
            <span className="material-symbols-outlined text-primary text-[14px]">person</span>
          </div>
        )}
        {record.assignees?.length > 0 && (
          <div className="flex -space-x-1.5">
            {record.assignees.slice(0, 3).map((_, i) => (
              <div key={i} className="w-6 h-6 rounded-full bg-primary/20 flex items-center justify-center ring-2 ring-surface-container-lowest">
                <span className="material-symbols-outlined text-primary text-[12px]">person</span>
              </div>
            ))}
            {record.assignees.length > 3 && (
              <div className="w-6 h-6 rounded-full bg-surface-container-high flex items-center justify-center ring-2 ring-surface-container-lowest">
                <span className="text-[9px] font-bold text-on-surface-variant">+{record.assignees.length - 3}</span>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
