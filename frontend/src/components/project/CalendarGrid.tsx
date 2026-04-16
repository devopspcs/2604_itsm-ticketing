import type { ProjectRecord } from '../../types/project'

interface CalendarGridProps {
  year: number
  month: number // 0-indexed
  records: ProjectRecord[]
  onClickRecord: (record: ProjectRecord) => void
}

const DAY_NAMES = ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min']

export function CalendarGrid({ year, month, records, onClickRecord }: CalendarGridProps) {
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDow = (firstDay.getDay() + 6) % 7 // Monday = 0
  const totalDays = lastDay.getDate()

  const today = new Date()
  const isToday = (day: number) =>
    today.getFullYear() === year && today.getMonth() === month && today.getDate() === day

  // Group records by due_date day
  const recordsByDay: Record<number, ProjectRecord[]> = {}
  records.forEach(r => {
    if (!r.due_date) return
    const d = new Date(r.due_date)
    if (d.getFullYear() === year && d.getMonth() === month) {
      const day = d.getDate()
      if (!recordsByDay[day]) recordsByDay[day] = []
      recordsByDay[day].push(r)
    }
  })

  const cells: (number | null)[] = []
  for (let i = 0; i < startDow; i++) cells.push(null)
  for (let d = 1; d <= totalDays; d++) cells.push(d)
  while (cells.length % 7 !== 0) cells.push(null)

  return (
    <div>
      {/* Day headers */}
      <div className="grid grid-cols-7 gap-px mb-1">
        {DAY_NAMES.map(d => (
          <div key={d} className="text-center text-[10px] font-bold text-on-surface-variant uppercase tracking-widest py-2">
            {d}
          </div>
        ))}
      </div>

      {/* Calendar cells */}
      <div className="grid grid-cols-7 gap-px bg-outline-variant/10 rounded-xl overflow-hidden">
        {cells.map((day, i) => (
          <div
            key={i}
            className={`bg-white min-h-[100px] p-2 ${
              day === null ? 'bg-surface-container-low/50' : ''
            }`}
          >
            {day !== null && (
              <>
                <span className={`text-xs font-medium inline-flex items-center justify-center w-6 h-6 rounded-full ${
                  isToday(day) ? 'bg-primary text-on-primary font-bold' : 'text-on-surface-variant'
                }`}>
                  {day}
                </span>
                <div className="flex flex-col gap-1 mt-1">
                  {(recordsByDay[day] ?? []).slice(0, 3).map(r => (
                    <button
                      key={r.id}
                      onClick={() => onClickRecord(r)}
                      className="text-[10px] font-medium text-white px-1.5 py-0.5 rounded truncate text-left hover:opacity-80 transition-opacity"
                      style={{ backgroundColor: '#3b82f6' }}
                      title={r.title}
                    >
                      {r.title}
                    </button>
                  ))}
                  {(recordsByDay[day]?.length ?? 0) > 3 && (
                    <span className="text-[10px] text-on-surface-variant">
                      +{(recordsByDay[day]?.length ?? 0) - 3} lagi
                    </span>
                  )}
                </div>
              </>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
