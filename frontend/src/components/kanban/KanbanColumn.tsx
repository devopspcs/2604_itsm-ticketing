import { useDroppable } from '@dnd-kit/core'
import type { Ticket, TicketStatus } from '../../types'
import { TicketCard } from './TicketCard'

interface KanbanColumnProps {
  status: TicketStatus
  title: string
  tickets: Ticket[]
  loading: boolean
  isOver?: boolean
}

export function KanbanColumn({ status, title, tickets, loading }: KanbanColumnProps) {
  const { setNodeRef, isOver } = useDroppable({ id: status })

  return (
    <div className="w-[320px] min-w-[300px] flex flex-col gap-4 shrink-0">
      {/* Column Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-1 h-6 bg-primary rounded-full" />
          <h2 className="text-lg font-bold tracking-tight text-on-surface font-headline">{title}</h2>
          <span className="text-xs font-medium text-on-surface-variant bg-surface-container-high px-2 py-0.5 rounded-full">
            {tickets.length}
          </span>
        </div>
      </div>

      {/* Droppable Area */}
      <div
        ref={setNodeRef}
        className={`flex-1 bg-surface-container-low rounded-xl p-4 flex flex-col gap-4 overflow-y-auto transition-all min-h-[200px] ${
          isOver ? 'ring-2 ring-primary/20 bg-surface-container-low/80' : ''
        }`}
      >
        {loading ? (
          <>
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-surface-container-lowest p-4 rounded-xl animate-pulse">
                <div className="h-3 bg-surface-container-high rounded w-16 mb-3" />
                <div className="h-4 bg-surface-container-high rounded w-full mb-2" />
                <div className="h-3 bg-surface-container-high rounded w-2/3" />
              </div>
            ))}
          </>
        ) : tickets.length === 0 ? (
          <div className="flex-1 flex items-center justify-center">
            <p className="text-xs text-on-surface-variant/50 font-medium">No tickets</p>
          </div>
        ) : (
          tickets.map((ticket) => <TicketCard key={ticket.id} ticket={ticket} />)
        )}
      </div>
    </div>
  )
}
