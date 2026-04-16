import { useEffect, useState } from 'react'
import { projectService } from '../services/project.service'
import { CalendarGrid } from '../components/project/CalendarGrid'
import { RecordDetailModal } from '../components/project/RecordDetailModal'
import type { ProjectRecord, ProjectDetail } from '../types/project'

const MONTH_NAMES = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

export function ProjectCalendarPage() {
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth())
  const [records, setRecords] = useState<ProjectRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedRecord, setSelectedRecord] = useState<ProjectRecord | null>(null)
  const [selectedProject, setSelectedProject] = useState<ProjectDetail | null>(null)

  useEffect(() => {
    setLoading(true)
    projectService.getCalendar(month + 1, year)
      .then(res => setRecords(res.data ?? []))
      .catch(() => setRecords([]))
      .finally(() => setLoading(false))
  }, [month, year])

  const prevMonth = () => {
    if (month === 0) { setMonth(11); setYear(y => y - 1) }
    else setMonth(m => m - 1)
  }

  const nextMonth = () => {
    if (month === 11) { setMonth(0); setYear(y => y + 1) }
    else setMonth(m => m + 1)
  }

  const goToday = () => {
    const t = new Date()
    setYear(t.getFullYear())
    setMonth(t.getMonth())
  }

  const handleClickRecord = async (record: ProjectRecord) => {
    setSelectedRecord(record)
    try {
      const res = await projectService.get(record.project_id)
      setSelectedProject(res.data)
    } catch {
      setSelectedProject(null)
    }
  }

  return (
    <div className="p-8">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">
            {MONTH_NAMES[month]} {year}
          </h1>
          <p className="text-sm text-on-surface-variant mt-1">Tampilan kalender record</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={goToday}
            className="px-3 py-1.5 text-xs font-bold text-primary border border-primary/20 rounded-xl hover:bg-primary/5 transition-colors"
          >
            HARI INI
          </button>
          <button onClick={prevMonth} className="p-1.5 hover:bg-surface-container-high rounded-lg transition-colors">
            <span className="material-symbols-outlined text-on-surface-variant">chevron_left</span>
          </button>
          <button onClick={nextMonth} className="p-1.5 hover:bg-surface-container-high rounded-lg transition-colors">
            <span className="material-symbols-outlined text-on-surface-variant">chevron_right</span>
          </button>
        </div>
      </div>

      {/* Calendar */}
      {loading ? (
        <div className="grid grid-cols-7 gap-px bg-outline-variant/10 rounded-xl overflow-hidden">
          {Array.from({ length: 35 }).map((_, i) => (
            <div key={i} className="bg-white min-h-[100px] p-2 animate-pulse">
              <div className="h-4 w-6 bg-surface-container-high rounded" />
            </div>
          ))}
        </div>
      ) : (
        <CalendarGrid
          year={year}
          month={month}
          records={records}
          onClickRecord={handleClickRecord}
        />
      )}

      {selectedRecord && selectedProject && (
        <RecordDetailModal
          record={selectedRecord}
          project={selectedProject}
          onClose={() => { setSelectedRecord(null); setSelectedProject(null) }}
          onUpdate={() => { setSelectedRecord(null); setSelectedProject(null) }}
        />
      )}
    </div>
  )
}
