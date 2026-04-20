import { NavLink } from 'react-router-dom'

interface ProjectBoardSidebarProps {
  projectId: string
}

const navItems = [
  { path: 'backlog', icon: 'list', label: 'Backlog' },
  { path: '', icon: 'dashboard', label: 'Board', end: true },
  { path: 'sprint', icon: 'sprint', label: 'Sprint' },
]

const projectItems = [
  { path: 'reports', icon: 'description', label: 'Reports' },
  { path: 'releases', icon: 'tag', label: 'Releases' },
  { path: 'components', icon: 'widgets', label: 'Components' },
  { path: 'issues', icon: 'bug_report', label: 'Issues' },
  { path: 'repository', icon: 'history', label: 'Activity' },
  { path: 'settings', icon: 'settings', label: 'Settings' },
]

const linkClass = ({ isActive }: { isActive: boolean }) =>
  `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
    isActive
      ? 'bg-primary/20 text-primary'
      : 'text-on-surface-variant hover:bg-surface-container-low'
  }`

export function ProjectBoardSidebar({ projectId }: ProjectBoardSidebarProps) {
  const base = `/projects/${projectId}`

  return (
    <div className="w-64 bg-primary/5 border-r border-outline-variant/10 flex flex-col flex-shrink-0">
      <div className="p-4 border-b border-outline-variant/10">
        <h2 className="text-sm font-bold text-on-surface mb-1">Project Board</h2>
        <p className="text-xs text-on-surface-variant">Jira-like board</p>
      </div>

      <nav className="flex-1 overflow-y-auto p-2 space-y-1">
        {navItems.map(item => (
          <NavLink
            key={item.path}
            to={item.path ? `${base}/${item.path}` : base}
            end={item.end}
            className={linkClass}
          >
            <span className="material-symbols-outlined text-[18px]">{item.icon}</span>
            {item.label}
          </NavLink>
        ))}

        <div className="px-3 py-2 text-xs font-semibold text-on-surface-variant uppercase tracking-wider mt-4 mb-2">
          Project
        </div>

        {projectItems.map(item => (
          <NavLink
            key={item.path}
            to={`${base}/${item.path}`}
            className={linkClass}
          >
            <span className="material-symbols-outlined text-[18px]">{item.icon}</span>
            {item.label}
          </NavLink>
        ))}
      </nav>
    </div>
  )
}
