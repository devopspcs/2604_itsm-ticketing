import { useEffect, useState } from 'react'
import { NavLink, useParams } from 'react-router-dom'
import { projectService } from '../../services/project.service'
import type { Project } from '../../types/project'

interface ProjectSidebarProps {
  onCreateProject: () => void
}

export function ProjectSidebar({ onCreateProject }: ProjectSidebarProps) {
  const { id: activeId } = useParams()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    projectService.list()
      .then(res => setProjects(res.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  return (
    <aside className="w-64 h-[calc(100vh-64px)] fixed left-0 top-16 bg-slate-50 flex flex-col z-40 hidden md:flex">
      <div className="p-4 border-b border-outline-variant/10">
        <div className="flex items-center justify-between mb-1">
          <h2 className="text-sm font-bold text-on-surface uppercase tracking-widest">Projects</h2>
          <button
            onClick={onCreateProject}
            className="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-surface-container-high transition-colors"
          >
            <span className="material-symbols-outlined text-on-surface-variant text-[18px]">add</span>
          </button>
        </div>
      </div>

      <nav className="flex-1 overflow-y-auto p-2 flex flex-col gap-0.5">
        {/* Home & Calendar links */}
        <NavLink
          to="/projects"
          end
          className={({ isActive }) =>
            `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              isActive ? 'bg-white text-primary font-bold shadow-sm' : 'text-slate-600 hover:bg-white/60'
            }`
          }
        >
          <span className="material-symbols-outlined text-[18px]">home</span>
          <span>Home</span>
        </NavLink>
        <NavLink
          to="/projects/calendar"
          className={({ isActive }) =>
            `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              isActive ? 'bg-white text-primary font-bold shadow-sm' : 'text-slate-600 hover:bg-white/60'
            }`
          }
        >
          <span className="material-symbols-outlined text-[18px]">calendar_month</span>
          <span>Calendar</span>
        </NavLink>

        <div className="h-px bg-outline-variant/10 my-2" />

        {loading ? (
          <div className="flex flex-col gap-2 px-3">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-8 bg-surface-container-high rounded-lg animate-pulse" />
            ))}
          </div>
        ) : projects.length === 0 ? (
          <p className="px-3 py-4 text-xs text-on-surface-variant/60 text-center">Belum ada project</p>
        ) : (
          projects.map(project => (
            <NavLink
              key={project.id}
              to={`/projects/${project.id}`}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive ? 'bg-white text-on-surface font-bold shadow-sm' : 'text-slate-600 hover:bg-white/60'
                }`
              }
            >
              <span
                className="w-3 h-3 rounded-sm shrink-0"
                style={{ backgroundColor: project.icon_color }}
              />
              <span className="truncate">{project.name}</span>
            </NavLink>
          ))
        )}
      </nav>
    </aside>
  )
}
