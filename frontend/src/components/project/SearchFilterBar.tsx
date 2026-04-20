import { useState, useEffect } from 'react'
import { jiraService } from '../../services/jira.service'
import type { JiraProjectRecord, SearchFilters, SavedFilter, IssueType, Label, WorkflowStatus } from '../../types/jira'

interface SearchFilterBarProps {
  projectId: string
  onSearchResults: (records: JiraProjectRecord[]) => void
  issueTypes?: IssueType[]
  statuses?: WorkflowStatus[]
  labels?: Label[]
}

export function SearchFilterBar({
  projectId,
  onSearchResults,
  issueTypes = [],
  statuses = [],
  labels = [],
}: SearchFilterBarProps) {
  const [query, setQuery] = useState('')
  const [filters, setFilters] = useState<SearchFilters>({})
  const [savedFilters, setSavedFilters] = useState<SavedFilter[]>([])
  const [showFilters, setShowFilters] = useState(false)
  const [filterName, setFilterName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSearch = async () => {
    try {
      setLoading(true)
      setError(null)
      const res = await jiraService.searchRecords(projectId, query, filters)
      onSearchResults(res.data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Pencarian gagal'
      setError(errorMessage)
      setTimeout(() => setError(null), 4000)
    } finally {
      setLoading(false)
    }
  }

  const handleSaveFilter = async () => {
    if (!filterName.trim()) return
    try {
      setLoading(true)
      await jiraService.saveFilter(projectId, {
        name: filterName,
        filters,
      })
      setFilterName('')
      loadSavedFilters()
    } catch (err) {
      console.error('Failed to save filter:', err)
    } finally {
      setLoading(false)
    }
  }

  const loadSavedFilters = async () => {
    try {
      const res = await jiraService.listSavedFilters(projectId)
      setSavedFilters(res.data)
    } catch (err) {
      console.error('Failed to load saved filters:', err)
    }
  }

  useEffect(() => {
    if (showFilters) {
      loadSavedFilters()
    }
  }, [showFilters])

  const applySavedFilter = (filter: SavedFilter) => {
    setFilters(filter.filters)
    setQuery('')
  }

  const handleReset = () => {
    setFilters({})
    setQuery('')
    onSearchResults([])
  }

  return (
    <div className="space-y-3">
      {error && (
        <div className="text-error text-sm p-2 bg-error/10 rounded">{error}</div>
      )}

      {/* Search Bar */}
      <div className="flex gap-2">
        <input
          type="text"
          value={query}
          onChange={e => setQuery(e.target.value)}
          onKeyPress={e => e.key === 'Enter' && handleSearch()}
          placeholder="Cari records berdasarkan judul atau deskripsi..."
          className="flex-1 px-3 py-2 rounded border border-outline-variant text-sm"
        />
        <button
          onClick={handleSearch}
          disabled={loading}
          className="px-4 py-2 bg-primary text-on-primary rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-90"
        >
          {loading ? 'Mencari...' : 'Cari'}
        </button>
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="px-4 py-2 bg-surface-container-high text-on-surface rounded text-sm hover:bg-surface-container-highest"
        >
          Filter {Object.keys(filters).length > 0 && `(${Object.keys(filters).length})`}
        </button>
      </div>

      {/* Advanced Filters */}
      {showFilters && (
        <div className="bg-surface-container-low border border-outline-variant rounded-lg p-4 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            {/* Issue Type Filter */}
            {issueTypes.length > 0 && (
              <div>
                <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                  Issue Type
                </label>
                <select
                  value={filters.issue_type_id || ''}
                  onChange={e =>
                    setFilters({
                      ...filters,
                      issue_type_id: e.target.value || undefined,
                    })
                  }
                  className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
                >
                  <option value="">All Types</option>
                  {issueTypes.map(type => (
                    <option key={type.id} value={type.id}>
                      {type.name}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {/* Status Filter */}
            {statuses.length > 0 && (
              <div>
                <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                  Status
                </label>
                <select
                  value={filters.status || ''}
                  onChange={e =>
                    setFilters({
                      ...filters,
                      status: e.target.value || undefined,
                    })
                  }
                  className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
                >
                  <option value="">All Statuses</option>
                  {statuses.map(status => (
                    <option key={status.id} value={status.id}>
                      {status.status_name}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {/* Label Filter */}
            {labels.length > 0 && (
              <div>
                <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                  Label
                </label>
                <select
                  value={filters.label_id || ''}
                  onChange={e =>
                    setFilters({
                      ...filters,
                      label_id: e.target.value || undefined,
                    })
                  }
                  className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
                >
                  <option value="">All Labels</option>
                  {labels.map(label => (
                    <option key={label.id} value={label.id}>
                      {label.name}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {/* Assignee Filter */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                Assignee
              </label>
              <input
                type="text"
                placeholder="Nama assignee"
                value={filters.assignee_id || ''}
                onChange={e =>
                  setFilters({
                    ...filters,
                    assignee_id: e.target.value || undefined,
                  })
                }
                className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
              />
            </div>

            {/* Due Date From */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                Due Date From
              </label>
              <input
                type="date"
                value={filters.due_date_from || ''}
                onChange={e =>
                  setFilters({
                    ...filters,
                    due_date_from: e.target.value || undefined,
                  })
                }
                className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
              />
            </div>

            {/* Due Date To */}
            <div>
              <label className="text-[10px] font-medium text-on-surface-variant block mb-1">
                Due Date To
              </label>
              <input
                type="date"
                value={filters.due_date_to || ''}
                onChange={e =>
                  setFilters({
                    ...filters,
                    due_date_to: e.target.value || undefined,
                  })
                }
                className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
              />
            </div>
          </div>

          {/* Filter Actions */}
          <div className="flex gap-2 pt-3 border-t border-outline-variant">
            <button
              onClick={handleSearch}
              className="flex-1 px-3 py-2 bg-primary text-on-primary rounded text-sm hover:opacity-90"
            >
              Terapkan Filter
            </button>
            <button
              onClick={handleReset}
              className="flex-1 px-3 py-2 bg-surface-container-high text-on-surface rounded text-sm hover:bg-surface-container-highest"
            >
              Reset
            </button>
          </div>

          {/* Save Filter */}
          <div className="border-t border-outline-variant pt-3 space-y-2">
            <p className="text-[10px] font-medium text-on-surface-variant">Simpan Filter</p>
            <div className="flex gap-2">
              <input
                type="text"
                value={filterName}
                onChange={e => setFilterName(e.target.value)}
                placeholder="Nama filter"
                className="flex-1 px-3 py-1 rounded border border-outline-variant text-sm"
              />
              <button
                onClick={handleSaveFilter}
                disabled={!filterName.trim() || loading}
                className="px-3 py-1 bg-primary text-on-primary rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-90"
              >
                Simpan
              </button>
            </div>
          </div>

          {/* Saved Filters */}
          {savedFilters.length > 0 && (
            <div className="border-t border-outline-variant pt-3 space-y-2">
              <p className="text-[10px] font-medium text-on-surface-variant">Filter Tersimpan</p>
              <div className="space-y-1 max-h-32 overflow-y-auto">
                {savedFilters.map(filter => (
                  <button
                    key={filter.id}
                    onClick={() => applySavedFilter(filter)}
                    className="w-full text-left px-2 py-1 rounded hover:bg-surface-container-high text-sm text-on-surface transition-colors"
                  >
                    {filter.name}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
