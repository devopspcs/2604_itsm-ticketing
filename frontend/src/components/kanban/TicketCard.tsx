import { useDraggable } from '@dnd-kit/core'
import { useNavigate } from 'react-router-dom'
import type { Ticket } from '../../types'

const priorityStyles: Record<string, string> = {
  critical: 'bg-error-container/30 text-on-error-container',
  high: 'bg-tertiary-fixed/30 text-on-tertiary-fixed-variant',
  medium: 'bg-secondary-fixed/30 text-on-secondary-fixed-variant',
  low: 'bg-surface-container-high text-on-surface-variant',
}

const typeIcons: Record<string, string> = {
  incident: 'report',
  change_request: 'published_with_changes',
  helpdesk_request: 'contact_support',
}

export function TicketCard({ ticket }: { ticket: Ticket }) {
  const navigate = useNavigate()
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: ticket.id,
    data: { ticket },
  })

  const style = transform
    ? {
        transform: `translate(${transform.x}px, ${transform.y}px)`,
        boxShadow: isDragging ? '0 20px 40px rgba(25, 28, 30, 0.06)' : undefined,
        zIndex: isDragging ? 50 : undefined,
        opacity: isDragging ? 0.9 : undefined,
      }
    : undefined

  const handleClick = () => {
    if (!isDragging) navigate(`/tickets/${ticket.id}`)
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      onClick={handleClick}
      className={`bg-surface-container-lowest p-4 rounded-xl border border-outline-variant/5 group cursor-pointer transition-all ${
        isDragging ? '' : 'hover:bg-surface-bright'
      }`}
    >
      <div className="flex justify-between items-start mb-3">
        <span className={`text-[10px] font-bold uppercase tracking-widest px-2 py-1 rounded-md ${priorityStyles[ticket.priority] ?? priorityStyles.low}`}>
          {ticket.priority}
        </span>
        <span className="material-symbols-outlined text-outline-variant/40 text-sm">drag_indicator</span>
      </div>
      <h3 className="text-sm font-semibold text-on-surface mb-2 leading-snug line-clamp-2">{ticket.title}</h3>
      <div className="flex items-center justify-between mt-3">
        <div className="flex items-center gap-2">
          <span className="material-symbols-outlined text-on-surface-variant text-[16px]">
            {typeIcons[ticket.type] ?? 'confirmation_number'}
          </span>
          <span className="text-[10px] font-medium text-on-surface-variant capitalize">{ticket.type.replace(/_/g, ' ')}</span>
        </div>
        {ticket.assigned_to && (
          <div className="w-6 h-6 rounded-full bg-primary/20 flex items-center justify-center">
            <span className="material-symbols-outlined text-primary text-[14px]">person</span>
          </div>
        )}
      </div>
    </div>
  )
}
