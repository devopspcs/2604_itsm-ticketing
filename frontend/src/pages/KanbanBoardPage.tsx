import { useCallback, useEffect, useState } from 'react'
import { DndContext, DragEndEvent, DragOverlay, DragStartEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core'
import { ticketService } from '../services/ticket.service'
import { KanbanColumn } from '../components/kanban/KanbanColumn'
import { TicketCard } from '../components/kanban/TicketCard'
import { KANBAN_COLUMNS } from '../types/kanban'
import type { Ticket, TicketStatus } from '../types'

type ColumnsState = Record<TicketStatus, Ticket[]>
type LoadingState = Record<TicketStatus, boolean>

const emptyColumns = (): ColumnsState => ({
  open: [], in_progress: [], pending_approval: [], approved: [], rejected: [], done: [],
})
const allLoading = (): LoadingState => ({
  open: true, in_progress: true, pending_approval: true, approved: true, rejected: true, done: true,
})

export function KanbanBoardPage() {
  const [columns, setColumns] = useState<ColumnsState>(emptyColumns)
  const [loading, setLoading] = useState<LoadingState>(allLoading)
  const [error, setError] = useState<string | null>(null)
  const [dragError, setDragError] = useState<string | null>(null)
  const [activeTicket, setActiveTicket] = useState<Ticket | null>(null)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } })
  )

  const fetchColumn = useCallback(async (status: TicketStatus) => {
    try {
      setLoading((prev) => ({ ...prev, [status]: true }))
      const res = await ticketService.list({ status, page_size: 50 })
      setColumns((prev) => ({ ...prev, [status]: res.data.tickets ?? [] }))
    } catch {
      setError(`Failed to load ${status} tickets`)
    } finally {
      setLoading((prev) => ({ ...prev, [status]: false }))
    }
  }, [])

  useEffect(() => {
    KANBAN_COLUMNS.forEach((col) => fetchColumn(col.status))
  }, [fetchColumn])

  const handleDragStart = (event: DragStartEvent) => {
    const ticket = event.active.data.current?.ticket as Ticket | undefined
    setActiveTicket(ticket ?? null)
    setDragError(null)
  }

  const handleDragEnd = async (event: DragEndEvent) => {
    setActiveTicket(null)
    const { active, over } = event
    if (!over) return

    const ticket = active.data.current?.ticket as Ticket | undefined
    if (!ticket) return

    const sourceStatus = ticket.status
    const destStatus = over.id as TicketStatus
    if (sourceStatus === destStatus) return

    // Optimistic update
    const prevColumns = { ...columns }
    setColumns((prev) => {
      const source = prev[sourceStatus].filter((t) => t.id !== ticket.id)
      const dest = [...prev[destStatus], { ...ticket, status: destStatus }]
      return { ...prev, [sourceStatus]: source, [destStatus]: dest }
    })

    try {
      await ticketService.update(ticket.id, { status: destStatus })
    } catch {
      // Rollback
      setColumns(prevColumns)
      setDragError(`Failed to move ticket "${ticket.title}"`)
      setTimeout(() => setDragError(null), 4000)
    }
  }

  return (
    <div className="flex flex-col h-[calc(100vh-64px)]">
      {/* Header */}
      <div className="px-8 pt-8 pb-4">
        <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">Kanban Board</h1>
        <p className="text-sm text-on-surface-variant font-medium mt-1">Drag tickets between columns to update status</p>
      </div>

      {/* Error banners */}
      {error && (
        <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
          <span className="material-symbols-outlined text-lg">error</span>
          <span>{error}</span>
          <button onClick={() => { setError(null); KANBAN_COLUMNS.forEach((c) => fetchColumn(c.status)) }}
            className="ml-auto text-xs font-bold underline">Retry</button>
        </div>
      )}
      {dragError && (
        <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm animate-in fade-in">
          <span className="material-symbols-outlined text-lg">warning</span>
          <span>{dragError}</span>
        </div>
      )}

      {/* Kanban Board */}
      <DndContext sensors={sensors} onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
        <div className="flex-1 overflow-x-auto overflow-y-hidden px-8 pb-8">
          <div className="flex gap-8 h-full min-w-max">
            {KANBAN_COLUMNS.map((col) => (
              <KanbanColumn
                key={col.status}
                status={col.status}
                title={col.title}
                tickets={columns[col.status] ?? []}
                loading={loading[col.status] ?? false}
              />
            ))}
          </div>
        </div>
        <DragOverlay>
          {activeTicket ? <TicketCard ticket={activeTicket} /> : null}
        </DragOverlay>
      </DndContext>
    </div>
  )
}
