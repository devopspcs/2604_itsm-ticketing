import { useEffect, useState } from 'react'
import { orgService } from '../services/org.service'
import type { OrgChartNode } from '../types'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import { ErrorMessage } from '../components/common/ErrorMessage'

function OrgChartNodeItem({ node, level = 0 }: { node: OrgChartNode; level?: number }) {
  const [expanded, setExpanded] = useState(level < 2)
  const hasChildren = node.children && node.children.length > 0

  const positionLabel = node.position
    ? node.position.charAt(0).toUpperCase() + node.position.slice(1)
    : node.role === 'admin' ? 'Admin' : ''

  const roleColor = {
    admin: 'text-purple-600',
    approver: 'text-blue-600',
    agent: 'text-green-600',
    user: 'text-on-surface-variant',
  }[node.role] || 'text-on-surface-variant'

  return (
    <div className={level > 0 ? 'ml-6 border-l border-outline-variant/40 pl-4' : ''}>
      <div className="flex items-center gap-2 py-1.5 group">
        {hasChildren ? (
          <button
            onClick={() => setExpanded(!expanded)}
            className="w-5 h-5 flex items-center justify-center rounded hover:bg-surface-container transition-colors"
          >
            <span className="material-symbols-outlined text-[16px] text-on-surface-variant">
              {expanded ? 'expand_more' : 'chevron_right'}
            </span>
          </button>
        ) : (
          <span className="w-5 h-5 flex items-center justify-center">
            <span className="material-symbols-outlined text-[14px] text-on-surface-variant/40">person</span>
          </span>
        )}

        <div className="flex items-center gap-2 flex-1 min-w-0">
          <div className="w-7 h-7 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
            <span className="text-xs font-bold text-primary">
              {node.full_name.split(' ').map(w => w[0]).slice(0, 2).join('').toUpperCase()}
            </span>
          </div>
          <span className="text-sm font-medium text-on-surface truncate">{node.full_name}</span>
          {positionLabel && (
            <span className={`text-xs font-medium ${roleColor}`}>{positionLabel}</span>
          )}
        </div>

        {hasChildren && (
          <span className="text-xs text-on-surface-variant/60 opacity-0 group-hover:opacity-100 transition-opacity">
            {node.children.length} direct report{node.children.length > 1 ? 's' : ''}
          </span>
        )}
      </div>

      {expanded && hasChildren && (
        <div>
          {node.children
            .sort((a, b) => a.full_name.localeCompare(b.full_name))
            .map(child => (
              <OrgChartNodeItem key={child.id} node={child} level={level + 1} />
            ))}
        </div>
      )}
    </div>
  )
}

export function OrgChartPage() {
  const [chart, setChart] = useState<OrgChartNode[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')

  useEffect(() => {
    orgService.getOrgChart()
      .then(r => setChart(r.data ?? []))
      .catch(() => setError('Failed to load org chart'))
      .finally(() => setLoading(false))
  }, [])

  // Flatten tree for search
  const flattenNodes = (nodes: OrgChartNode[]): OrgChartNode[] => {
    const result: OrgChartNode[] = []
    const traverse = (n: OrgChartNode) => {
      result.push(n)
      n.children?.forEach(traverse)
    }
    nodes.forEach(traverse)
    return result
  }

  // Filter tree based on search
  const filterTree = (nodes: OrgChartNode[], query: string): OrgChartNode[] => {
    if (!query) return nodes
    const lower = query.toLowerCase()
    const allNodes = flattenNodes(nodes)
    const matchingIds = new Set(
      allNodes
        .filter(n => n.full_name.toLowerCase().includes(lower) || n.email.toLowerCase().includes(lower))
        .map(n => n.id)
    )

    // Also include ancestors of matching nodes
    const includeAncestors = (nodes: OrgChartNode[]): OrgChartNode[] => {
      return nodes
        .map(node => {
          const filteredChildren = includeAncestors(node.children || [])
          if (matchingIds.has(node.id) || filteredChildren.length > 0) {
            return { ...node, children: filteredChildren }
          }
          return null
        })
        .filter(Boolean) as OrgChartNode[]
    }

    return includeAncestors(nodes)
  }

  const filteredChart = filterTree(chart, search)
  const totalPeople = flattenNodes(chart).length

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-xl font-bold text-on-surface">Org Chart</h2>
          <p className="text-sm text-on-surface-variant mt-0.5">{totalPeople} people in organization</p>
        </div>
      </div>

      {error && <ErrorMessage message={error} />}

      {/* Search */}
      <div className="mb-4">
        <div className="relative">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[18px] text-on-surface-variant">search</span>
          <input
            type="text"
            placeholder="Search people..."
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 border border-outline-variant rounded-xl text-sm bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30"
          />
        </div>
      </div>

      {/* Tree */}
      <div className="bg-surface-container-lowest rounded-xl p-4 shadow-sm">
        {filteredChart.length === 0 ? (
          <p className="text-sm text-on-surface-variant text-center py-8">
            {search ? 'No results found' : 'No org chart data. Assign "Reports To" to users to build the hierarchy.'}
          </p>
        ) : (
          filteredChart
            .sort((a, b) => a.full_name.localeCompare(b.full_name))
            .map(node => (
              <OrgChartNodeItem key={node.id} node={node} level={0} />
            ))
        )}
      </div>
    </div>
  )
}
