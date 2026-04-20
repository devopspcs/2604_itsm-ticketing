import { useState, useEffect } from 'react'
import { jiraService } from '../../services/jira.service'
import { generateRandomColor } from '../../utils/jira.utils'
import type { Label } from '../../types/jira'

interface LabelManagerProps {
  recordId: string
  projectId: string
  selectedLabels: Label[]
  onLabelsChange: (labels: Label[]) => void
}

export function LabelManager({ recordId, projectId, selectedLabels, onLabelsChange }: LabelManagerProps) {
  const [availableLabels, setAvailableLabels] = useState<Label[]>([])
  const [loading, setLoading] = useState(false)
  const [showDropdown, setShowDropdown] = useState(false)
  const [newLabelName, setNewLabelName] = useState('')
  const [newLabelColor, setNewLabelColor] = useState(generateRandomColor())
  const [creatingLabel, setCreatingLabel] = useState(false)

  const loadLabels = async () => {
    try {
      setLoading(true)
      const res = await jiraService.listLabels(projectId)
      setAvailableLabels(res.data)
    } catch (error) {
      console.error('Failed to load labels:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (showDropdown) {
      loadLabels()
    }
  }, [showDropdown])

  const handleAddLabel = async (label: Label) => {
    if (!selectedLabels.find(l => l.id === label.id)) {
      try {
        await jiraService.addLabelToRecord(projectId, recordId, label.id)
        onLabelsChange([...selectedLabels, label])
      } catch (error) {
        console.error('Failed to add label:', error)
      }
    }
  }

  const handleRemoveLabel = async (labelId: string) => {
    try {
      await jiraService.removeLabelFromRecord(projectId, recordId, labelId)
      onLabelsChange(selectedLabels.filter(l => l.id !== labelId))
    } catch (error) {
      console.error('Failed to remove label:', error)
    }
  }

  const handleCreateLabel = async () => {
    if (!newLabelName.trim()) return
    try {
      setCreatingLabel(true)
      const res = await jiraService.createLabel(projectId, {
        name: newLabelName,
        color: newLabelColor,
      })
      const newLabel = res.data
      setAvailableLabels([...availableLabels, newLabel])
      await handleAddLabel(newLabel)
      setNewLabelName('')
      setNewLabelColor(generateRandomColor())
    } catch (error) {
      console.error('Failed to create label:', error)
    } finally {
      setCreatingLabel(false)
    }
  }

  const unselectedLabels = availableLabels.filter(l => !selectedLabels.find(sl => sl.id === l.id))

  return (
    <div className="space-y-2">
      <h4 className="text-sm font-medium text-on-surface">Labels</h4>

      {/* Selected Labels */}
      {selectedLabels.length > 0 ? (
        <div className="flex flex-wrap gap-2">
          {selectedLabels.map(label => (
            <div
              key={label.id}
              className="flex items-center gap-1 px-2 py-1 rounded-full text-white text-sm font-medium"
              style={{ backgroundColor: label.color }}
            >
              {label.name}
              <button
                onClick={() => handleRemoveLabel(label.id)}
                className="hover:opacity-80 ml-1"
                aria-label={`Remove ${label.name} label`}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-sm text-on-surface-variant">No labels selected</p>
      )}

      {/* Add Label Dropdown */}
      <div className="relative">
        <button
          onClick={() => setShowDropdown(!showDropdown)}
          className="text-sm text-primary hover:underline"
        >
          + Tambah Label
        </button>

        {showDropdown && (
          <div className="absolute top-full left-0 mt-2 bg-surface-container-low border border-outline-variant rounded-lg shadow-lg z-10 w-72 p-3">
            {loading ? (
              <div className="text-on-surface-variant text-sm">Memuat...</div>
            ) : (
              <>
                {/* Available Labels List */}
                {unselectedLabels.length > 0 && (
                  <div className="space-y-2 max-h-48 overflow-y-auto mb-3 pb-3 border-b border-outline-variant">
                    {unselectedLabels.map(label => (
                      <button
                        key={label.id}
                        onClick={() => handleAddLabel(label)}
                        className="w-full text-left px-2 py-1 rounded hover:bg-surface-container-high text-sm flex items-center gap-2 transition-colors"
                      >
                        <span
                          className="w-3 h-3 rounded-full flex-shrink-0"
                          style={{ backgroundColor: label.color }}
                        />
                        <span className="flex-1">{label.name}</span>
                      </button>
                    ))}
                  </div>
                )}

                {/* Create New Label */}
                <div className="space-y-2">
                  <p className="text-[10px] font-medium text-on-surface-variant">Buat Label Baru</p>
                  <input
                    type="text"
                    value={newLabelName}
                    onChange={e => setNewLabelName(e.target.value)}
                    placeholder="Nama label"
                    className="w-full px-2 py-1 rounded border border-outline-variant text-sm"
                  />
                  <div className="flex gap-2">
                    <input
                      type="color"
                      value={newLabelColor}
                      onChange={e => setNewLabelColor(e.target.value)}
                      className="w-10 h-8 rounded cursor-pointer"
                      aria-label="Label color"
                    />
                    <button
                      onClick={handleCreateLabel}
                      disabled={!newLabelName.trim() || creatingLabel}
                      className="flex-1 px-2 py-1 bg-primary text-on-primary rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-90"
                    >
                      {creatingLabel ? 'Membuat...' : 'Buat'}
                    </button>
                  </div>
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
