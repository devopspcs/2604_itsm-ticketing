import { NavLink, useParams } from 'react-router-dom'

export function ProjectNavigation() {
  const { id } = useParams()

  if (!id) return null

  return (
    <nav className="bg-surface-container-low border-b border-outline-variant/10 px-8 flex items-center gap-1">
      <NavLink
        to={`/projects/${id}`}
        end
        className={({ isActive }) =>
          `flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors border-b-2 ${
            isActive
              ? 'text-primary border-primary'
              : 'text-on-surface-variant border-transparent hover:text-on-surface'
          }`
        }
      >
        <span className="material-symbols-outlined text-[18px]">dashboard</span>
        Board
      </NavLink>

      <NavLink
        to={`/projects/${id}/sprint`}
        className={({ isActive }) =>
          `flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors border-b-2 ${
            isActive
              ? 'text-primary border-primary'
              : 'text-on-surface-variant border-transparent hover:text-on-surface'
          }`
        }
      >
        <span className="material-symbols-outlined text-[18px]">sprint</span>
        Sprint
      </NavLink>

      <NavLink
        to={`/projects/${id}/backlog`}
        className={({ isActive }) =>
          `flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors border-b-2 ${
            isActive
              ? 'text-primary border-primary'
              : 'text-on-surface-variant border-transparent hover:text-on-surface'
          }`
        }
      >
        <span className="material-symbols-outlined text-[18px]">list</span>
        Backlog
      </NavLink>

      <NavLink
        to={`/projects/${id}/settings`}
        className={({ isActive }) =>
          `flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors border-b-2 ${
            isActive
              ? 'text-primary border-primary'
              : 'text-on-surface-variant border-transparent hover:text-on-surface'
          }`
        }
      >
        <span className="material-symbols-outlined text-[18px]">settings</span>
        Settings
      </NavLink>
    </nav>
  )
}
