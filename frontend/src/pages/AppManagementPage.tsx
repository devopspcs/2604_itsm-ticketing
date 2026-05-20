import { useEffect, useState } from 'react'
import { appService } from '../services/app.service'
import type { Application, User } from '../types'
import api from '../services/api'

interface UserWithAppAccess {
  user: User
  role: string
}

export function AppManagementPage() {
  const [apps, setApps] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedApp, setSelectedApp] = useState<Application | null>(null)
  const [appUsers, setAppUsers] = useState<UserWithAppAccess[]>([])
  const [allUsers, setAllUsers] = useState<User[]>([])
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [showAddUserModal, setShowAddUserModal] = useState(false)
  const [formData, setFormData] = useState({ name: '', code: '', description: '', icon: 'apps', color: '#1976d2' })

  useEffect(() => {
    loadApps()
    loadAllUsers()
  }, [])

  const loadApps = async () => {
    try {
      const res = await appService.listApps()
      setApps(res.data || [])
    } catch {
      // ignore
    } finally {
      setLoading(false)
    }
  }

  const loadAllUsers = async () => {
    try {
      const res = await api.get<User[]>('/users/list')
      setAllUsers(res.data || [])
    } catch {
      // ignore
    }
  }

  const loadAppUsers = async (appId: string) => {
    try {
      const res = await appService.getAppUsers(appId)
      setAppUsers(res.data || [])
    } catch {
      setAppUsers([])
    }
  }

  const handleSelectApp = (app: Application) => {
    setSelectedApp(app)
    loadAppUsers(app.id)
  }

  const handleCreate = async () => {
    try {
      await appService.createApp(formData)
      setShowCreateModal(false)
      setFormData({ name: '', code: '', description: '', icon: 'apps', color: '#1976d2' })
      loadApps()
    } catch {
      // ignore
    }
  }

  const handleUpdate = async () => {
    if (!selectedApp) return
    try {
      await appService.updateApp(selectedApp.id, formData)
      setShowEditModal(false)
      loadApps()
      setSelectedApp(null)
    } catch {
      // ignore
    }
  }

  const handleDelete = async (app: Application) => {
    if (!confirm(`Delete application "${app.name}"? This will remove all user access.`)) return
    try {
      await appService.deleteApp(app.id)
      loadApps()
      if (selectedApp?.id === app.id) {
        setSelectedApp(null)
        setAppUsers([])
      }
    } catch {
      // ignore
    }
  }

  const handleToggleActive = async (app: Application) => {
    try {
      await appService.updateApp(app.id, { is_active: !app.is_active })
      loadApps()
    } catch {
      // ignore
    }
  }

  const handleRevokeAccess = async (userId: string) => {
    if (!selectedApp) return
    try {
      await appService.revokeAccess(selectedApp.id, userId)
      loadAppUsers(selectedApp.id)
    } catch {
      // ignore
    }
  }

  const handleGrantAccess = async (userId: string, role: string) => {
    if (!selectedApp) return
    try {
      await appService.grantAccess(selectedApp.id, { user_id: userId, role })
      loadAppUsers(selectedApp.id)
      setShowAddUserModal(false)
    } catch {
      // ignore
    }
  }

  const handleBulkGrant = async (userIds: string[], role: string) => {
    if (!selectedApp) return
    try {
      await appService.bulkGrantAccess(selectedApp.id, { user_ids: userIds, role })
      loadAppUsers(selectedApp.id)
      setShowAddUserModal(false)
    } catch {
      // ignore
    }
  }

  const openEditModal = (app: Application) => {
    setFormData({ name: app.name, code: app.code, description: app.description, icon: app.icon, color: app.color })
    setSelectedApp(app)
    setShowEditModal(true)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      </div>
    )
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">App Management</h1>
          <p className="text-sm text-slate-500 mt-1">Manage applications and user access</p>
        </div>
        <button
          onClick={() => {
            setFormData({ name: '', code: '', description: '', icon: 'apps', color: '#1976d2' })
            setShowCreateModal(true)
          }}
          className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors"
        >
          <span className="material-symbols-outlined text-sm">add</span>
          New Application
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* App List */}
        <div className="lg:col-span-1 space-y-3">
          <h2 className="text-sm font-semibold text-slate-700 uppercase tracking-wide">Applications</h2>
          {apps.map(app => (
            <div
              key={app.id}
              onClick={() => handleSelectApp(app)}
              className={`p-4 rounded-xl border cursor-pointer transition-all ${
                selectedApp?.id === app.id
                  ? 'border-primary bg-primary/5 shadow-sm'
                  : 'border-slate-200 hover:border-slate-300 bg-white'
              }`}
            >
              <div className="flex items-start gap-3">
                <div
                  className="w-10 h-10 rounded-lg flex items-center justify-center"
                  style={{ backgroundColor: app.color + '20' }}
                >
                  <span className="material-symbols-outlined" style={{ color: app.color }}>
                    {app.icon}
                  </span>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-slate-900 truncate">{app.name}</h3>
                    {!app.is_active && (
                      <span className="text-[10px] px-1.5 py-0.5 bg-red-100 text-red-700 rounded font-medium">Inactive</span>
                    )}
                  </div>
                  <p className="text-xs text-slate-500 mt-0.5">{app.code}</p>
                  {app.description && (
                    <p className="text-xs text-slate-400 mt-1 line-clamp-2">{app.description}</p>
                  )}
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={(e) => { e.stopPropagation(); openEditModal(app) }}
                    className="p-1 text-slate-400 hover:text-primary rounded"
                    title="Edit"
                  >
                    <span className="material-symbols-outlined text-base">edit</span>
                  </button>
                  <button
                    onClick={(e) => { e.stopPropagation(); handleToggleActive(app) }}
                    className={`p-1 rounded ${app.is_active ? 'text-slate-400 hover:text-amber-600' : 'text-amber-600 hover:text-emerald-600'}`}
                    title={app.is_active ? 'Deactivate' : 'Activate'}
                  >
                    <span className="material-symbols-outlined text-base">
                      {app.is_active ? 'toggle_on' : 'toggle_off'}
                    </span>
                  </button>
                  <button
                    onClick={(e) => { e.stopPropagation(); handleDelete(app) }}
                    className="p-1 text-slate-400 hover:text-red-600 rounded"
                    title="Delete"
                  >
                    <span className="material-symbols-outlined text-base">delete</span>
                  </button>
                </div>
              </div>
            </div>
          ))}
          {apps.length === 0 && (
            <p className="text-sm text-slate-400 text-center py-8">No applications yet</p>
          )}
        </div>

        {/* App Users Panel */}
        <div className="lg:col-span-2">
          {selectedApp ? (
            <div className="bg-white rounded-xl border border-slate-200 p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-lg flex items-center justify-center"
                    style={{ backgroundColor: selectedApp.color + '20' }}
                  >
                    <span className="material-symbols-outlined" style={{ color: selectedApp.color }}>
                      {selectedApp.icon}
                    </span>
                  </div>
                  <div>
                    <h2 className="font-bold text-slate-900">{selectedApp.name}</h2>
                    <p className="text-xs text-slate-500">{appUsers.length} user(s) with access</p>
                  </div>
                </div>
                <button
                  onClick={() => setShowAddUserModal(true)}
                  className="flex items-center gap-1 px-3 py-1.5 bg-primary text-white text-sm rounded-lg hover:bg-primary/90"
                >
                  <span className="material-symbols-outlined text-sm">person_add</span>
                  Add Users
                </button>
              </div>

              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-slate-100">
                      <th className="text-left py-2 px-3 text-slate-500 font-medium">User</th>
                      <th className="text-left py-2 px-3 text-slate-500 font-medium">Email</th>
                      <th className="text-left py-2 px-3 text-slate-500 font-medium">Role</th>
                      <th className="text-right py-2 px-3 text-slate-500 font-medium">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {appUsers.map(({ user, role }) => (
                      <tr key={user.id} className="border-b border-slate-50 hover:bg-slate-50">
                        <td className="py-2 px-3 font-medium text-slate-900">{user.full_name}</td>
                        <td className="py-2 px-3 text-slate-500">{user.email}</td>
                        <td className="py-2 px-3">
                          <span className="px-2 py-0.5 bg-slate-100 text-slate-700 rounded text-xs font-medium">
                            {role}
                          </span>
                        </td>
                        <td className="py-2 px-3 text-right">
                          <button
                            onClick={() => handleRevokeAccess(user.id)}
                            className="text-red-500 hover:text-red-700 text-xs font-medium"
                          >
                            Revoke
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {appUsers.length === 0 && (
                  <p className="text-center text-slate-400 py-8 text-sm">No users have access to this app</p>
                )}
              </div>
            </div>
          ) : (
            <div className="flex items-center justify-center h-64 bg-white rounded-xl border border-slate-200">
              <div className="text-center">
                <span className="material-symbols-outlined text-4xl text-slate-300">apps</span>
                <p className="text-sm text-slate-400 mt-2">Select an application to manage access</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Create Modal */}
      {showCreateModal && (
        <Modal title="Create Application" onClose={() => setShowCreateModal(false)}>
          <AppForm
            formData={formData}
            setFormData={setFormData}
            onSubmit={handleCreate}
            submitLabel="Create"
          />
        </Modal>
      )}

      {/* Edit Modal */}
      {showEditModal && (
        <Modal title="Edit Application" onClose={() => setShowEditModal(false)}>
          <AppForm
            formData={formData}
            setFormData={setFormData}
            onSubmit={handleUpdate}
            submitLabel="Save Changes"
            disableCode
          />
        </Modal>
      )}

      {/* Add User Modal */}
      {showAddUserModal && selectedApp && (
        <AddUserModal
          allUsers={allUsers}
          existingUserIds={appUsers.map(u => u.user.id)}
          onGrant={handleGrantAccess}
          onBulkGrant={handleBulkGrant}
          onClose={() => setShowAddUserModal(false)}
        />
      )}
    </div>
  )
}

// --- Sub-components ---

function Modal({ title, onClose, children }: { title: string; onClose: () => void; children: React.ReactNode }) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={onClose}>
      <div className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-bold text-slate-900">{title}</h3>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>
        {children}
      </div>
    </div>
  )
}

interface AppFormProps {
  formData: { name: string; code: string; description: string; icon: string; color: string }
  setFormData: (data: { name: string; code: string; description: string; icon: string; color: string }) => void
  onSubmit: () => void
  submitLabel: string
  disableCode?: boolean
}

function AppForm({ formData, setFormData, onSubmit, submitLabel, disableCode }: AppFormProps) {
  return (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-slate-700 mb-1">Name</label>
        <input
          type="text"
          value={formData.name}
          onChange={e => setFormData({ ...formData, name: e.target.value })}
          className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
          placeholder="Application name"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-slate-700 mb-1">Code</label>
        <input
          type="text"
          value={formData.code}
          onChange={e => setFormData({ ...formData, code: e.target.value })}
          disabled={disableCode}
          className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 disabled:bg-slate-100"
          placeholder="app-code"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-slate-700 mb-1">Description</label>
        <textarea
          value={formData.description}
          onChange={e => setFormData({ ...formData, description: e.target.value })}
          className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
          rows={2}
          placeholder="Brief description"
        />
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-1">Icon</label>
          <input
            type="text"
            value={formData.icon}
            onChange={e => setFormData({ ...formData, icon: e.target.value })}
            className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
            placeholder="Material icon name"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-1">Color</label>
          <div className="flex gap-2">
            <input
              type="color"
              value={formData.color}
              onChange={e => setFormData({ ...formData, color: e.target.value })}
              className="w-10 h-10 rounded border border-slate-300 cursor-pointer"
            />
            <input
              type="text"
              value={formData.color}
              onChange={e => setFormData({ ...formData, color: e.target.value })}
              className="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
            />
          </div>
        </div>
      </div>
      <button
        onClick={onSubmit}
        disabled={!formData.name || !formData.code}
        className="w-full py-2 bg-primary text-white rounded-lg font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {submitLabel}
      </button>
    </div>
  )
}

interface AddUserModalProps {
  allUsers: User[]
  existingUserIds: string[]
  onGrant: (userId: string, role: string) => void
  onBulkGrant: (userIds: string[], role: string) => void
  onClose: () => void
}

function AddUserModal({ allUsers, existingUserIds, onGrant, onBulkGrant, onClose }: AddUserModalProps) {
  const [search, setSearch] = useState('')
  const [selectedIds, setSelectedIds] = useState<string[]>([])
  const [role, setRole] = useState('user')

  const available = allUsers.filter(u =>
    !existingUserIds.includes(u.id) &&
    (u.full_name.toLowerCase().includes(search.toLowerCase()) || u.email.toLowerCase().includes(search.toLowerCase()))
  )

  const toggleUser = (id: string) => {
    setSelectedIds(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])
  }

  const handleSubmit = () => {
    if (selectedIds.length === 1) {
      onGrant(selectedIds[0], role)
    } else if (selectedIds.length > 1) {
      onBulkGrant(selectedIds, role)
    }
  }

  return (
    <Modal title="Add Users" onClose={onClose}>
      <div className="space-y-4">
        <input
          type="text"
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
          placeholder="Search users..."
        />
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-1">Role</label>
          <select
            value={role}
            onChange={e => setRole(e.target.value)}
            className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
          >
            <option value="user">User</option>
            <option value="agent">Agent</option>
            <option value="approver">Approver</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div className="max-h-48 overflow-y-auto border border-slate-200 rounded-lg">
          {available.map(user => (
            <label
              key={user.id}
              className="flex items-center gap-3 px-3 py-2 hover:bg-slate-50 cursor-pointer border-b border-slate-50 last:border-0"
            >
              <input
                type="checkbox"
                checked={selectedIds.includes(user.id)}
                onChange={() => toggleUser(user.id)}
                className="rounded border-slate-300"
              />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-slate-900 truncate">{user.full_name}</p>
                <p className="text-xs text-slate-500 truncate">{user.email}</p>
              </div>
              <span className="text-[10px] px-1.5 py-0.5 bg-slate-100 text-slate-600 rounded">{user.role}</span>
            </label>
          ))}
          {available.length === 0 && (
            <p className="text-center text-slate-400 py-4 text-sm">No users available</p>
          )}
        </div>
        <button
          onClick={handleSubmit}
          disabled={selectedIds.length === 0}
          className="w-full py-2 bg-primary text-white rounded-lg font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          Grant Access ({selectedIds.length} selected)
        </button>
      </div>
    </Modal>
  )
}
