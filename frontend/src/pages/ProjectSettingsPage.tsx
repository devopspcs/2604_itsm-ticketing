import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { jiraService } from '../services/jira.service'
import { projectService } from '../services/project.service'
import api from '../services/api'
import { ProjectBoardSidebar } from '../components/project/ProjectBoardSidebar'
import type { IssueType, CustomField, Workflow, Label } from '../types/jira'

type SettingsTab = 'members' | 'issue-types' | 'custom-fields' | 'workflows' | 'labels' | 'danger'

interface MemberInfo {
  project_id: string
  user_id: string
  role: string
  created_at: string
  name?: string
  email?: string
}

interface UserInfo {
  id: string
  name: string
  email: string
}

export function ProjectSettingsPage() {
  const { id: projectId } = useParams()
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<SettingsTab>('members')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Members
  const [members, setMembers] = useState<MemberInfo[]>([])
  const [allUsers, setAllUsers] = useState<UserInfo[]>([])
  const [showInvite, setShowInvite] = useState(false)
  const [inviteUserId, setInviteUserId] = useState('')
  const [inviting, setInviting] = useState(false)

  // Issue Types
  const [issueTypes, setIssueTypes] = useState<IssueType[]>([])

  // Custom Fields
  const [customFields, setCustomFields] = useState<CustomField[]>([])
  const [showAddField, setShowAddField] = useState(false)
  const [newFieldName, setNewFieldName] = useState('')
  const [newFieldType, setNewFieldType] = useState<'text' | 'textarea' | 'dropdown' | 'date' | 'number' | 'checkbox'>('text')
  const [newFieldRequired, setNewFieldRequired] = useState(false)

  // Workflows
  const [workflow, setWorkflow] = useState<Workflow | null>(null)

  // Labels
  const [labels, setLabels] = useState<Label[]>([])
  const [showAddLabel, setShowAddLabel] = useState(false)
  const [newLabelName, setNewLabelName] = useState('')
  const [newLabelColor, setNewLabelColor] = useState('#3B82F6')

  useEffect(() => {
    const fetchSettings = async () => {
      if (!projectId) return
      try {
        setLoading(true)
        setError(null)

        const [issueTypesRes, fieldsRes, workflowRes, labelsRes, membersRes, usersRes] = await Promise.all([
          jiraService.listIssueTypes(projectId),
          jiraService.listCustomFields(projectId),
          jiraService.getWorkflow(projectId).catch(() => ({ data: null })),
          jiraService.listLabels(projectId),
          projectService.listMembers(projectId).catch(() => ({ data: [] })),
          api.get<UserInfo[]>('/users/list').catch(() => ({ data: [] })),
        ])

        setIssueTypes(issueTypesRes.data || [])
        setCustomFields(fieldsRes.data || [])
        setWorkflow(workflowRes.data || null)
        setLabels(labelsRes.data || [])
        setMembers(membersRes.data || [])
        setAllUsers(usersRes.data || [])
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Gagal memuat pengaturan'
        setError(errorMessage)
      } finally {
        setLoading(false)
      }
    }

    fetchSettings()
  }, [projectId])

  const handleAddCustomField = async () => {
    if (!projectId || !newFieldName.trim()) return
    try {
      await jiraService.createCustomField(projectId, {
        name: newFieldName,
        field_type: newFieldType,
        is_required: newFieldRequired,
      })
      setNewFieldName('')
      setNewFieldType('text')
      setNewFieldRequired(false)
      setShowAddField(false)
      // Refresh fields
      const fieldsRes = await jiraService.listCustomFields(projectId)
      setCustomFields(fieldsRes.data || [])
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menambah field'
      setError(errorMessage)
    }
  }

  const handleDeleteCustomField = async (fieldId: string) => {
    if (!projectId) return
    if (!confirm('Apakah Anda yakin ingin menghapus field ini?')) return
    try {
      await jiraService.deleteCustomField(projectId, fieldId)
      setCustomFields(customFields.filter(f => f.id !== fieldId))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menghapus field'
      setError(errorMessage)
    }
  }

  const handleAddLabel = async () => {
    if (!projectId || !newLabelName.trim()) return
    try {
      await jiraService.createLabel(projectId, {
        name: newLabelName,
        color: newLabelColor,
      })
      setNewLabelName('')
      setNewLabelColor('#3B82F6')
      setShowAddLabel(false)
      // Refresh labels
      const labelsRes = await jiraService.listLabels(projectId)
      setLabels(labelsRes.data || [])
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menambah label'
      setError(errorMessage)
    }
  }

  const handleDeleteLabel = async (labelId: string) => {
    if (!confirm('Apakah Anda yakin ingin menghapus label ini?')) return
    try {
      await jiraService.deleteLabel(projectId!, labelId)
      setLabels(labels.filter(l => l.id !== labelId))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menghapus label'
      setError(errorMessage)
    }
  }

  const handleInviteMember = async () => {
    if (!projectId || !inviteUserId) return
    try {
      setInviting(true)
      await projectService.inviteMember(projectId, inviteUserId)
      const membersRes = await projectService.listMembers(projectId)
      setMembers(membersRes.data || [])
      setInviteUserId('')
      setShowInvite(false)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal mengundang member'
      setError(errorMessage)
    } finally {
      setInviting(false)
    }
  }

  const handleRemoveMember = async (userId: string) => {
    if (!projectId) return
    if (!confirm('Apakah Anda yakin ingin menghapus member ini?')) return
    try {
      await projectService.removeMember(projectId, userId)
      setMembers(members.filter(m => m.user_id !== userId))
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal menghapus member'
      setError(errorMessage)
    }
  }

  const getUserName = (userId: string) => {
    const user = allUsers.find(u => u.id === userId)
    return user?.name || user?.email || userId.slice(0, 8) + '...'
  }

  const getUserEmail = (userId: string) => {
    const user = allUsers.find(u => u.id === userId)
    return user?.email || ''
  }

  const nonMemberUsers = allUsers.filter(u => !members.find(m => m.user_id === u.id))

  if (loading) {
    return (
      <div className="flex flex-col h-[calc(100vh-64px)]">
        <div className="px-8 pt-8 pb-4">
          <div className="h-7 w-48 bg-surface-container-high rounded-lg animate-pulse" />
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-screen bg-surface">
      <ProjectBoardSidebar projectId={projectId || ''} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <div className="px-8 pt-8 pb-4">
          <h1 className="text-2xl font-extrabold text-on-surface tracking-tight font-headline">
            Project Settings
          </h1>
        </div>

      {error && (
        <div className="mx-8 mb-4 bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-3 text-sm">
          <span className="material-symbols-outlined text-lg">error</span>
          <span>{error}</span>
        </div>
      )}

      <div className="flex-1 overflow-auto px-8 pb-8">
        {/* Tabs */}
        <div className="flex gap-4 mb-6 border-b border-outline-variant/20">
          {(['members', 'issue-types', 'custom-fields', 'workflows', 'labels', 'danger'] as const).map(tab => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 font-medium text-sm transition-colors border-b-2 ${
                activeTab === tab
                  ? tab === 'danger' ? 'text-error border-error' : 'text-primary border-primary'
                  : 'text-on-surface-variant border-transparent hover:text-on-surface'
              }`}
            >
              {tab === 'members' && 'Members'}
              {tab === 'issue-types' && 'Issue Types'}
              {tab === 'custom-fields' && 'Custom Fields'}
              {tab === 'workflows' && 'Workflows'}
              {tab === 'labels' && 'Labels'}
              {tab === 'danger' && 'Danger Zone'}
            </button>
          ))}
        </div>

        {/* Members Tab */}
        {activeTab === 'members' && (
          <div className="bg-surface-container-low rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-on-surface">Project Members</h2>
                <p className="text-sm text-on-surface-variant mt-1">Hanya member yang bisa melihat dan mengakses project ini</p>
              </div>
              <button
                onClick={() => setShowInvite(!showInvite)}
                className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
              >
                <span className="material-symbols-outlined text-[18px]">person_add</span>
                Invite Member
              </button>
            </div>

            {showInvite && (
              <div className="bg-surface-container-highest p-4 rounded-lg mb-4 border border-outline-variant/20">
                <h3 className="text-sm font-semibold text-on-surface mb-3">Undang User ke Project</h3>
                <div className="space-y-3">
                  <select
                    value={inviteUserId}
                    onChange={e => setInviteUserId(e.target.value)}
                    className="w-full px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-on-surface focus:outline-none focus:border-primary"
                  >
                    <option value="">Pilih user...</option>
                    {nonMemberUsers.map(user => (
                      <option key={user.id} value={user.id}>
                        {user.name} ({user.email})
                      </option>
                    ))}
                  </select>
                  <div className="flex gap-2">
                    <button
                      onClick={handleInviteMember}
                      disabled={!inviteUserId || inviting}
                      className="px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50"
                    >
                      {inviting ? 'Mengundang...' : 'Undang'}
                    </button>
                    <button
                      onClick={() => setShowInvite(false)}
                      className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-lg transition-colors"
                    >
                      Batal
                    </button>
                  </div>
                </div>
              </div>
            )}

            <div className="space-y-2">
              {members.length === 0 ? (
                <p className="text-on-surface-variant text-center py-8">Belum ada member</p>
              ) : (
                members.map(member => (
                  <div key={member.user_id} className="bg-surface-container-highest p-4 rounded-lg border border-outline-variant/20 flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center">
                        <span className="material-symbols-outlined text-primary">person</span>
                      </div>
                      <div>
                        <h3 className="font-semibold text-on-surface">{getUserName(member.user_id)}</h3>
                        <p className="text-xs text-on-surface-variant">{getUserEmail(member.user_id)}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                        member.role === 'owner'
                          ? 'bg-primary/20 text-primary'
                          : 'bg-surface-container-high text-on-surface-variant'
                      }`}>
                        {member.role === 'owner' ? 'Owner' : 'Member'}
                      </span>
                      {member.role !== 'owner' && (
                        <button
                          onClick={() => handleRemoveMember(member.user_id)}
                          className="p-2 text-error hover:bg-error-container/20 rounded-lg transition-colors"
                          title="Hapus member"
                        >
                          <span className="material-symbols-outlined text-[18px]">close</span>
                        </button>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {/* Issue Types Tab */}
        {activeTab === 'issue-types' && (
          <div className="bg-surface-container-low rounded-lg p-6">
            <h2 className="text-lg font-semibold text-on-surface mb-4">Available Issue Types</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {issueTypes.map(type => (
                <div key={type.id} className="bg-surface-container-highest p-4 rounded-lg border border-outline-variant/20">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="material-symbols-outlined text-primary">{type.icon || 'task'}</span>
                    <h3 className="font-semibold text-on-surface">{type.name}</h3>
                  </div>
                  {type.description && (
                    <p className="text-sm text-on-surface-variant">{type.description}</p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Custom Fields Tab */}
        {activeTab === 'custom-fields' && (
          <div className="bg-surface-container-low rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-on-surface">Custom Fields</h2>
              <button
                onClick={() => setShowAddField(!showAddField)}
                className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
              >
                <span className="material-symbols-outlined text-[18px]">add</span>
                Add Field
              </button>
            </div>

            {showAddField && (
              <div className="bg-surface-container-highest p-4 rounded-lg mb-4 border border-outline-variant/20">
                <div className="space-y-3">
                  <input
                    type="text"
                    value={newFieldName}
                    onChange={e => setNewFieldName(e.target.value)}
                    placeholder="Field name"
                    className="w-full px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-on-surface placeholder-on-surface-variant focus:outline-none focus:border-primary"
                  />
                  <select
                    value={newFieldType}
                    onChange={e => setNewFieldType(e.target.value as any)}
                    className="w-full px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-on-surface focus:outline-none focus:border-primary"
                  >
                    <option value="text">Text</option>
                    <option value="textarea">Text Area</option>
                    <option value="dropdown">Dropdown</option>
                    <option value="date">Date</option>
                    <option value="number">Number</option>
                    <option value="checkbox">Checkbox</option>
                  </select>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={newFieldRequired}
                      onChange={e => setNewFieldRequired(e.target.checked)}
                      className="w-4 h-4"
                    />
                    <span className="text-sm text-on-surface">Required</span>
                  </label>
                  <div className="flex gap-2">
                    <button
                      onClick={handleAddCustomField}
                      className="px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
                    >
                      Add
                    </button>
                    <button
                      onClick={() => setShowAddField(false)}
                      className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-lg transition-colors"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            )}

            <div className="space-y-2">
              {customFields.length === 0 ? (
                <p className="text-on-surface-variant text-center py-8">No custom fields yet</p>
              ) : (
                customFields.map(field => (
                  <div key={field.id} className="bg-surface-container-highest p-4 rounded-lg border border-outline-variant/20 flex items-center justify-between">
                    <div>
                      <h3 className="font-semibold text-on-surface">{field.name}</h3>
                      <p className="text-sm text-on-surface-variant">
                        {field.field_type} {field.is_required && '• Required'}
                      </p>
                    </div>
                    <button
                      onClick={() => handleDeleteCustomField(field.id)}
                      className="p-2 text-error hover:bg-error-container/20 rounded-lg transition-colors"
                    >
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {/* Workflows Tab */}
        {activeTab === 'workflows' && (
          <div className="bg-surface-container-low rounded-lg p-6">
            <h2 className="text-lg font-semibold text-on-surface mb-4">Workflow Configuration</h2>
            {workflow ? (
              <div className="bg-surface-container-highest p-4 rounded-lg border border-outline-variant/20">
                <h3 className="font-semibold text-on-surface mb-2">{workflow.name}</h3>
                <p className="text-sm text-on-surface-variant">
                  Initial Status: <span className="font-medium">{workflow.initial_status}</span>
                </p>
                <p className="text-sm text-on-surface-variant mt-2">
                  Workflow configuration is managed through the API. Contact your administrator to modify workflows.
                </p>
              </div>
            ) : (
              <p className="text-on-surface-variant text-center py-8">No workflow configured</p>
            )}
          </div>
        )}

        {/* Labels Tab */}
        {activeTab === 'labels' && (
          <div className="bg-surface-container-low rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-on-surface">Labels</h2>
              <button
                onClick={() => setShowAddLabel(!showAddLabel)}
                className="flex items-center gap-2 px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
              >
                <span className="material-symbols-outlined text-[18px]">add</span>
                Add Label
              </button>
            </div>

            {showAddLabel && (
              <div className="bg-surface-container-highest p-4 rounded-lg mb-4 border border-outline-variant/20">
                <div className="space-y-3">
                  <input
                    type="text"
                    value={newLabelName}
                    onChange={e => setNewLabelName(e.target.value)}
                    placeholder="Label name"
                    className="w-full px-3 py-2 bg-surface-container-low border border-outline-variant/20 rounded-lg text-on-surface placeholder-on-surface-variant focus:outline-none focus:border-primary"
                  />
                  <div className="flex items-center gap-3">
                    <input
                      type="color"
                      value={newLabelColor}
                      onChange={e => setNewLabelColor(e.target.value)}
                      className="w-12 h-10 rounded-lg cursor-pointer"
                    />
                    <span className="text-sm text-on-surface-variant">{newLabelColor}</span>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={handleAddLabel}
                      className="px-4 py-2 text-sm font-semibold text-on-primary bg-primary rounded-lg hover:bg-primary/90 transition-colors"
                    >
                      Add
                    </button>
                    <button
                      onClick={() => setShowAddLabel(false)}
                      className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-lg transition-colors"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            )}

            <div className="space-y-2">
              {labels.length === 0 ? (
                <p className="text-on-surface-variant text-center py-8">No labels yet</p>
              ) : (
                labels.map(label => (
                  <div key={label.id} className="bg-surface-container-highest p-4 rounded-lg border border-outline-variant/20 flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div
                        className="w-4 h-4 rounded-full"
                        style={{ backgroundColor: label.color }}
                      />
                      <h3 className="font-semibold text-on-surface">{label.name}</h3>
                    </div>
                    <button
                      onClick={() => handleDeleteLabel(label.id)}
                      className="p-2 text-error hover:bg-error-container/20 rounded-lg transition-colors"
                    >
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {/* Danger Zone Tab */}
        {activeTab === 'danger' && (
          <div className="bg-surface-container-low rounded-lg p-6 border-2 border-error/20">
            <h2 className="text-lg font-semibold text-error mb-2">Danger Zone</h2>
            <p className="text-sm text-on-surface-variant mb-6">Tindakan di bawah ini tidak dapat dibatalkan. Harap berhati-hati.</p>

            <div className="bg-error/5 rounded-lg p-4 border border-error/20 flex items-center justify-between">
              <div>
                <h3 className="font-semibold text-on-surface">Hapus Project</h3>
                <p className="text-sm text-on-surface-variant mt-1">Menghapus project beserta semua data di dalamnya (records, sprints, labels, dll)</p>
              </div>
              <button
                onClick={async () => {
                  if (!projectId) return
                  const confirmText = prompt('Ketik "HAPUS" untuk mengkonfirmasi penghapusan project:')
                  if (confirmText !== 'HAPUS') return
                  try {
                    await projectService.delete(projectId)
                    navigate('/projects')
                  } catch (err) {
                    setError(err instanceof Error ? err.message : 'Gagal menghapus project')
                  }
                }}
                className="px-4 py-2 text-sm font-semibold text-white bg-error rounded-lg hover:bg-error/90 transition-colors flex-shrink-0"
              >
                Hapus Project
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
    </div>
  )
}
