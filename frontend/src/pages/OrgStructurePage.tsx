import { useEffect, useState } from 'react'
import { orgService } from '../services/org.service'
import type { Department, Division, Team } from '../types'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import { ErrorMessage } from '../components/common/ErrorMessage'

type Tab = 'divisions' | 'departments' | 'teams'

export function OrgStructurePage() {
  const [tab, setTab] = useState<Tab>('divisions')
  const [divisions, setDivisions] = useState<Division[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [teams, setTeams] = useState<Team[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Filters
  const [filterDivId, setFilterDivId] = useState('')
  const [filterDeptId, setFilterDeptId] = useState('')

  // Modal state
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null)

  // Form fields
  const [formName, setFormName] = useState('')
  const [formCode, setFormCode] = useState('')
  const [formDivId, setFormDivId] = useState('')
  const [formDeptId, setFormDeptId] = useState('')
  const [formEmail, setFormEmail] = useState('')

  const fetchDivisions = () => orgService.listDivisions().then(r => setDivisions(r.data ?? [])).catch(() => setError('Failed to load divisions'))
  const fetchDepartments = () => orgService.listDepartments(filterDivId || undefined).then(r => setDepartments(r.data ?? [])).catch(() => setError('Failed to load departments'))
  const fetchTeams = () => orgService.listTeams(filterDeptId || undefined).then(r => setTeams(r.data ?? [])).catch(() => setError('Failed to load teams'))

  useEffect(() => { fetchDivisions().finally(() => setLoading(false)) }, [])
  useEffect(() => { if (tab === 'departments') fetchDepartments() }, [tab, filterDivId])
  useEffect(() => { if (tab === 'teams') fetchTeams() }, [tab, filterDeptId])
  useEffect(() => { if (tab === 'departments' || tab === 'teams') fetchDivisions() }, [tab])

  // All departments for team filter dropdown (unfiltered)
  const [allDepartments, setAllDepartments] = useState<Department[]>([])
  useEffect(() => { if (tab === 'teams') orgService.listDepartments().then(r => setAllDepartments(r.data ?? [])) }, [tab])

  const resetForm = () => { setFormName(''); setFormCode(''); setFormDivId(''); setFormDeptId(''); setFormEmail(''); setEditingId(null) }

  const openAdd = () => { resetForm(); setShowModal(true) }
  const openEdit = (item: Division | Department | Team) => {
    setEditingId(item.id)
    setFormName(item.name)
    if ('code' in item) setFormCode(item.code)
    if ('division_id' in item) setFormDivId(item.division_id)
    if ('department_id' in item) setFormDeptId(item.department_id)
    if ('email' in item) setFormEmail(item.email ?? '')
    setShowModal(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    try {
      if (tab === 'divisions') {
        if (editingId) await orgService.updateDivision(editingId, { name: formName, code: formCode })
        else await orgService.createDivision({ name: formName, code: formCode })
        await fetchDivisions()
      } else if (tab === 'departments') {
        if (editingId) await orgService.updateDepartment(editingId, { division_id: formDivId, name: formName, code: formCode })
        else await orgService.createDepartment({ division_id: formDivId, name: formName, code: formCode })
        await fetchDepartments()
      } else {
        if (editingId) await orgService.updateTeam(editingId, { department_id: formDeptId, name: formName, email: formEmail || undefined })
        else await orgService.createTeam({ department_id: formDeptId, name: formName, email: formEmail || undefined })
        await fetchTeams()
      }
      setShowModal(false)
      resetForm()
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Operation failed'
      setError(msg)
    }
  }

  const handleDelete = async (id: string) => {
    setError('')
    try {
      if (tab === 'divisions') { await orgService.deleteDivision(id); await fetchDivisions() }
      else if (tab === 'departments') { await orgService.deleteDepartment(id); await fetchDepartments() }
      else { await orgService.deleteTeam(id); await fetchTeams() }
      setConfirmDeleteId(null)
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Delete failed'
      setError(msg)
      setConfirmDeleteId(null)
    }
  }

  const divName = (id: string) => divisions.find(d => d.id === id)?.name ?? id
  const deptName = (id: string) => allDepartments.find(d => d.id === id)?.name ?? departments.find(d => d.id === id)?.name ?? id

  const tabs: { key: Tab; label: string }[] = [
    { key: 'divisions', label: 'Divisions' },
    { key: 'departments', label: 'Departments' },
    { key: 'teams', label: 'Teams' },
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-on-surface">Organizational Structure</h2>
        <button onClick={openAdd} className="flex items-center gap-1.5 bg-primary text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors">
          <span className="material-symbols-outlined text-[18px]">add</span>
          Add {tab === 'divisions' ? 'Division' : tab === 'departments' ? 'Department' : 'Team'}
        </button>
      </div>

      {error && <ErrorMessage message={error} />}

      {/* Tabs */}
      <div className="flex gap-1 mb-4 bg-surface-container-lowest rounded-xl p-1">
        {tabs.map(t => (
          <button key={t.key} onClick={() => { setTab(t.key); setError('') }}
            className={`flex-1 py-2 px-4 rounded-lg text-sm font-semibold transition-colors ${tab === t.key ? 'bg-primary text-white shadow-sm' : 'text-on-surface-variant hover:bg-surface-container'}`}>
            {t.label}
          </button>
        ))}
      </div>

      {/* Filters */}
      {tab === 'departments' && (
        <div className="mb-4">
          <select value={filterDivId} onChange={e => setFilterDivId(e.target.value)}
            className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest">
            <option value="">All Divisions</option>
            {divisions.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
          </select>
        </div>
      )}
      {tab === 'teams' && (
        <div className="mb-4">
          <select value={filterDeptId} onChange={e => setFilterDeptId(e.target.value)}
            className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest">
            <option value="">All Departments</option>
            {allDepartments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
          </select>
        </div>
      )}

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-xl overflow-hidden shadow-sm">
        <table className="w-full">
          <thead>
            <tr className="bg-surface-container">
              <th className="text-left px-4 py-3 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Name</th>
              {(tab === 'divisions' || tab === 'departments') && <th className="text-left px-4 py-3 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Code</th>}
              {tab === 'departments' && <th className="text-left px-4 py-3 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Division</th>}
              {tab === 'teams' && <th className="text-left px-4 py-3 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Department</th>}
              {tab === 'teams' && <th className="text-left px-4 py-3 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Email</th>}
              <th className="text-right px-4 py-3 text-xs font-bold text-on-surface-variant uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody>
            {tab === 'divisions' && divisions.map(d => (
              <tr key={d.id} className="border-t border-outline-variant/30 hover:bg-surface-container/30">
                <td className="px-4 py-3 text-sm">{d.name}</td>
                <td className="px-4 py-3 text-sm text-on-surface-variant">{d.code}</td>
                <td className="px-4 py-3 text-right">
                  <button onClick={() => openEdit(d)} className="text-primary hover:underline text-sm mr-3">Edit</button>
                  <button onClick={() => setConfirmDeleteId(d.id)} className="text-error hover:underline text-sm">Delete</button>
                </td>
              </tr>
            ))}
            {tab === 'departments' && departments.map(d => (
              <tr key={d.id} className="border-t border-outline-variant/30 hover:bg-surface-container/30">
                <td className="px-4 py-3 text-sm">{d.name}</td>
                <td className="px-4 py-3 text-sm text-on-surface-variant">{d.code}</td>
                <td className="px-4 py-3 text-sm text-on-surface-variant">{divName(d.division_id)}</td>
                <td className="px-4 py-3 text-right">
                  <button onClick={() => openEdit(d)} className="text-primary hover:underline text-sm mr-3">Edit</button>
                  <button onClick={() => setConfirmDeleteId(d.id)} className="text-error hover:underline text-sm">Delete</button>
                </td>
              </tr>
            ))}
            {tab === 'teams' && teams.map(t => (
              <tr key={t.id} className="border-t border-outline-variant/30 hover:bg-surface-container/30">
                <td className="px-4 py-3 text-sm">{t.name}</td>
                <td className="px-4 py-3 text-sm text-on-surface-variant">{deptName(t.department_id)}</td>
                <td className="px-4 py-3 text-sm text-on-surface-variant">{t.email || <span className="italic text-on-surface-variant/50">No email</span>}</td>
                <td className="px-4 py-3 text-right">
                  <button onClick={() => openEdit(t)} className="text-primary hover:underline text-sm mr-3">Edit</button>
                  <button onClick={() => setConfirmDeleteId(t.id)} className="text-error hover:underline text-sm">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Delete Confirmation */}
      {confirmDeleteId && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-surface-container-lowest rounded-xl p-6 max-w-sm w-full mx-4 shadow-xl">
            <h3 className="text-lg font-bold text-on-surface mb-2">Confirm Delete</h3>
            <p className="text-sm text-on-surface-variant mb-6">Are you sure you want to delete this item? This action cannot be undone.</p>
            <div className="flex justify-end gap-2">
              <button onClick={() => setConfirmDeleteId(null)} className="px-4 py-2 rounded-lg text-sm font-semibold text-on-surface-variant hover:bg-surface-container transition-colors">Cancel</button>
              <button onClick={() => handleDelete(confirmDeleteId)} className="px-4 py-2 rounded-lg text-sm font-semibold bg-error text-white hover:bg-error/90 transition-colors">Delete</button>
            </div>
          </div>
        </div>
      )}

      {/* Add/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-surface-container-lowest rounded-xl p-6 max-w-md w-full mx-4 shadow-xl">
            <h3 className="text-lg font-bold text-on-surface mb-4">{editingId ? 'Edit' : 'Add'} {tab === 'divisions' ? 'Division' : tab === 'departments' ? 'Department' : 'Team'}</h3>
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              {tab === 'departments' && (
                <select value={formDivId} onChange={e => setFormDivId(e.target.value)} required
                  className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest">
                  <option value="">Select Division</option>
                  {divisions.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
                </select>
              )}
              {tab === 'teams' && (
                <select value={formDeptId} onChange={e => setFormDeptId(e.target.value)} required
                  className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest">
                  <option value="">Select Department</option>
                  {allDepartments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
                </select>
              )}
              <input placeholder="Name" value={formName} onChange={e => setFormName(e.target.value)} required
                className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest" />
              {(tab === 'divisions' || tab === 'departments') && (
                <input placeholder="Code" value={formCode} onChange={e => setFormCode(e.target.value)} required
                  className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest" />
              )}
              {tab === 'teams' && (
                <input placeholder="Team Email (for notifications)" value={formEmail} onChange={e => setFormEmail(e.target.value)} type="email"
                  className="border border-outline-variant rounded-lg px-3 py-2 text-sm bg-surface-container-lowest" />
              )}
              <div className="flex justify-end gap-2 mt-2">
                <button type="button" onClick={() => { setShowModal(false); resetForm() }}
                  className="px-4 py-2 rounded-lg text-sm font-semibold text-on-surface-variant hover:bg-surface-container transition-colors">Cancel</button>
                <button type="submit"
                  className="px-4 py-2 rounded-lg text-sm font-semibold bg-primary text-white hover:bg-primary/90 transition-colors">{editingId ? 'Update' : 'Create'}</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
