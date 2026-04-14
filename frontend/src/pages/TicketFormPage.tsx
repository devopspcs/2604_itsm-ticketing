import { useState, useRef, useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { ticketService } from '../services/ticket.service'
import api from '../services/api'
import type { TicketType, Priority, User } from '../types'

const steps = [
  { label: 'Request Type', icon: 'category' },
  { label: 'Basic Details', icon: 'edit_note' },
  { label: 'Attachments', icon: 'attach_file' },
  { label: 'Review', icon: 'fact_check' },
]

const typeOptions = [
  { value: 'incident' as TicketType, icon: 'bolt', label: 'Incident', desc: 'Something is broken or not working.' },
  { value: 'change_request' as TicketType, icon: 'published_with_changes', label: 'Change Request', desc: 'Modify existing systems or access.' },
  { value: 'helpdesk_request' as TicketType, icon: 'contact_support', label: 'Helpdesk', desc: 'General inquiries and guidance.' },
]

export function TicketFormPage() {
  const navigate = useNavigate()
  const [step, setStep] = useState(0)
  const [type, setType] = useState<TicketType>('helpdesk_request')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [category, setCategory] = useState('Network & Connectivity')
  const [priority, setPriority] = useState<Priority>('medium')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [files, setFiles] = useState<File[]>([])
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [assigneeId, setAssigneeId] = useState('')
  const [usersList, setUsersList] = useState<User[]>([])

  // Load users for assign dropdown
  useEffect(() => {
    api.get<User[]>('/users/list')
      .then(res => setUsersList(res.data ?? []))
      .catch(() => {})
  }, [])

  const canProceed = () => {
    if (step === 1) return title.trim().length > 0 && description.trim().length > 0
    return true
  }

  const goToStep = (i: number) => {
    if (i < step) { setError(''); setStep(i); return }
    if (i === step + 1 && canProceed()) { setError(''); setStep(i); return }
    if (i === step + 1 && !canProceed()) { setError('Please fill in all required fields') }
  }

  const handleSubmit = async () => {
    if (!title.trim() || !description.trim()) { setError('Title and description are required'); return }
    setLoading(true); setError('')
    try {
      const res = await ticketService.create({ title, description, type, category, priority })
      const ticketId = res.data.id

      // Upload attachments if any
      for (const file of files) {
        const formData = new FormData()
        formData.append('file', file)
        await api.post(`/tickets/${ticketId}/attachments`, formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
        }).catch(() => {}) // don't fail if attachment upload fails
      }

      // Assign ticket if assignee selected
      if (assigneeId) {
        await ticketService.assign(ticketId, assigneeId).catch(() => {})
      }

      navigate(`/tickets/${ticketId}`)
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message
      setError(msg ?? 'Failed to create ticket')
    } finally { setLoading(false) }
  }

  return (
    <div className="flex h-[calc(100vh-64px)] overflow-hidden">
      <aside className="bg-slate-100 w-64 flex flex-col p-4 space-y-2 h-full hidden md:flex">
        <div className="mb-8 px-2 py-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-primary flex items-center justify-center text-white">
              <span className="material-symbols-outlined">description</span>
            </div>
            <div>
              <p className="text-sm font-semibold text-primary">Service Console</p>
              <p className="text-[10px] uppercase tracking-widest text-slate-500">ITSM Workspace</p>
            </div>
          </div>
        </div>
        <nav className="space-y-1">
          {[
            { to: '/dashboard', icon: 'dashboard', label: 'Home' },
            { to: '/tickets', icon: 'description', label: 'My Requests', active: true },
            { to: '/approvals', icon: 'fact_check', label: 'Approvals' },
            { to: '/activity-logs', icon: 'history_edu', label: 'Activity Logs' },
          ].map(item => (
            <Link key={item.to} to={item.to}
              className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-semibold transition-all ${
                item.active ? 'bg-white text-primary shadow-sm' : 'text-slate-600 hover:bg-slate-200 hover:translate-x-1 duration-200'
              }`}>
              <span className="material-symbols-outlined">{item.icon}</span>
              {item.label}
            </Link>
          ))}
        </nav>
      </aside>

      <main className="flex-1 overflow-y-auto bg-surface p-8">
        <div className="max-w-5xl mx-auto">
          <div className="mb-10">
            <div className="flex items-center gap-2 text-on-surface-variant text-sm mb-4">
              <Link to="/tickets" className="hover:text-primary transition-colors">Requests</Link>
              <span className="material-symbols-outlined text-xs">chevron_right</span>
              <span className="text-primary font-medium">New Request</span>
            </div>
            <h1 className="text-4xl font-extrabold tracking-tight text-on-surface mb-2">Create New Request</h1>
            <p className="text-on-surface-variant max-w-2xl">Initialize a new ticket within the ITSM ecosystem.</p>
          </div>

          <div className="grid grid-cols-12 gap-8 items-start">
            {/* Stepper */}
            <div className="col-span-12 lg:col-span-4 space-y-6">
              <div className="bg-surface-container-low p-6 rounded-xl">
                <h3 className="text-sm font-bold uppercase tracking-widest text-on-surface-variant mb-6">Request Progress</h3>
                <div className="space-y-1 relative">
                  <div className="absolute left-[11px] top-3 bottom-3 w-0.5 bg-outline-variant/30" />
                  {steps.map((s, i) => (
                    <button key={s.label} onClick={() => goToStep(i)}
                      className={`relative flex items-center gap-4 w-full text-left px-2 py-2.5 rounded-xl transition-all ${
                        i === step ? 'bg-white shadow-sm' :
                        i < step ? 'hover:bg-white/60 cursor-pointer' :
                        'cursor-default opacity-50'
                      }`}>
                      <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs z-10 flex-shrink-0 transition-all ${
                        i < step ? 'bg-emerald-500 text-white' :
                        i === step ? 'bg-primary text-white' :
                        'bg-outline-variant text-on-surface-variant'
                      }`}>
                        {i < step
                          ? <span className="material-symbols-outlined text-sm" style={{ fontVariationSettings: "'FILL' 1" }}>check</span>
                          : i + 1}
                      </div>
                      <span className={`text-sm ${i === step ? 'font-bold text-primary' : i < step ? 'font-medium text-emerald-700' : 'text-on-surface-variant'}`}>
                        {s.label}
                      </span>
                    </button>
                  ))}
                </div>
              </div>

              <div className="bg-primary p-6 rounded-xl text-white overflow-hidden relative shadow-lg">
                <div className="relative z-10">
                  <h4 className="font-bold mb-2">Need immediate assistance?</h4>
                  <p className="text-sm opacity-80 mb-4">Critical system outages should be reported via the Emergency Hotline.</p>
                  <button className="bg-white/20 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-white/30 transition-colors">
                    Contact Support
                  </button>
                </div>
                <span className="material-symbols-outlined absolute -bottom-4 -right-4 text-9xl opacity-10"
                  style={{ fontVariationSettings: "'FILL' 1" }}>support_agent</span>
              </div>
            </div>

            {/* Form */}
            <div className="col-span-12 lg:col-span-8 space-y-6">

              {step === 0 && (
                <section className="bg-surface-container-lowest p-8 rounded-xl shadow-sm">
                  <h2 className="text-xl font-bold mb-1">Select Request Type</h2>
                  <p className="text-sm text-on-surface-variant mb-6">What kind of assistance do you need today?</p>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {typeOptions.map((opt) => (
                      <button key={opt.value} type="button" onClick={() => setType(opt.value)}
                        className={`flex flex-col items-start p-5 border-2 rounded-xl text-left transition-all ${
                          type === opt.value ? 'border-primary bg-primary-fixed/30 shadow-sm' : 'border-transparent bg-surface-container-low hover:bg-surface-container-high'
                        }`}>
                        <span className={`material-symbols-outlined mb-3 text-2xl ${type === opt.value ? 'text-primary' : 'text-secondary'}`}
                          style={type === opt.value ? { fontVariationSettings: "'FILL' 1" } : {}}>
                          {opt.icon}
                        </span>
                        <span className={`font-bold text-sm ${type === opt.value ? 'text-primary' : 'text-secondary'}`}>{opt.label}</span>
                        <span className="text-[11px] text-on-surface-variant leading-relaxed mt-1">{opt.desc}</span>
                      </button>
                    ))}
                  </div>
                </section>
              )}

              {step === 1 && (
                <section className="bg-surface-container-lowest p-8 rounded-xl shadow-sm">
                  <h2 className="text-xl font-bold mb-1">Basic Details</h2>
                  <p className="text-sm text-on-surface-variant mb-6">Provide information about your request.</p>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-5">
                    <div className="md:col-span-2">
                      <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Request Title *</label>
                      <input value={title} onChange={(e) => setTitle(e.target.value)} required
                        placeholder="e.g., Cannot access internal VPN gateway"
                        className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary outline-none text-sm" />
                    </div>
                    <div>
                      <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Category</label>
                      <select value={category} onChange={(e) => setCategory(e.target.value)}
                        className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary outline-none text-sm appearance-none">
                        <option>Network &amp; Connectivity</option>
                        <option>Software &amp; Applications</option>
                        <option>Hardware Replacement</option>
                        <option>Identity &amp; Access</option>
                        <option>Infrastructure</option>
                        <option>Security</option>
                        <option>Database</option>
                        <option>Other</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Priority</label>
                      <div className="flex items-center gap-2">
                        <select value={priority} onChange={(e) => setPriority(e.target.value as Priority)}
                          className="flex-1 bg-surface-container-highest border-none rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary outline-none text-sm appearance-none">
                          <option value="low">Low</option>
                          <option value="medium">Medium</option>
                          <option value="high">High</option>
                          <option value="critical">Critical</option>
                        </select>
                        <div className={`w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0 ${
                          priority === 'critical' || priority === 'high' ? 'bg-error-container text-on-error-container' : 'bg-surface-container-high text-on-surface-variant'
                        }`}>
                          <span className="material-symbols-outlined text-xl">priority_high</span>
                        </div>
                      </div>
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Description *</label>
                      <textarea value={description} onChange={(e) => setDescription(e.target.value)} required rows={5}
                        placeholder="Provide detailed information about the incident..."
                        className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary outline-none text-sm resize-none" />
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-2">Assign To (optional)</label>
                      <select
                        value={assigneeId}
                        onChange={(e) => setAssigneeId(e.target.value)}
                        className="w-full bg-surface-container-highest border-none rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary outline-none text-sm appearance-none"
                      >
                        <option value="">-- No assignment (assign later) --</option>
                        {usersList.filter(u => u.is_active).map(u => (
                          <option key={u.id} value={u.id}>{u.full_name} ({u.role})</option>
                        ))}
                      </select>
                      <p className="text-[10px] text-on-surface-variant mt-1">Select a user to immediately assign this ticket to.</p>
                    </div>
                  </div>
                </section>
              )}

              {step === 2 && (
                <section className="bg-surface-container-lowest p-8 rounded-xl shadow-sm">
                  <h2 className="text-xl font-bold mb-1">Attachments</h2>
                  <p className="text-sm text-on-surface-variant mb-6">Upload supporting files or screenshots (optional).</p>

                  {/* Hidden file input */}
                  <input
                    ref={fileInputRef}
                    type="file"
                    multiple
                    accept="image/*,.pdf,.doc,.docx,.txt,.csv,.xlsx"
                    className="hidden"
                    onChange={(e) => {
                      const newFiles = Array.from(e.target.files ?? [])
                      setFiles(prev => [...prev, ...newFiles])
                      if (fileInputRef.current) fileInputRef.current.value = ''
                    }}
                  />

                  {/* Drop zone */}
                  <div
                    className="border-2 border-dashed border-outline-variant/40 rounded-xl p-10 text-center hover:border-primary/40 transition-colors cursor-pointer"
                    onClick={() => fileInputRef.current?.click()}
                    onDragOver={(e) => { e.preventDefault(); e.currentTarget.classList.add('border-primary', 'bg-primary-fixed/10') }}
                    onDragLeave={(e) => { e.currentTarget.classList.remove('border-primary', 'bg-primary-fixed/10') }}
                    onDrop={(e) => {
                      e.preventDefault()
                      e.currentTarget.classList.remove('border-primary', 'bg-primary-fixed/10')
                      const droppedFiles = Array.from(e.dataTransfer.files)
                      setFiles(prev => [...prev, ...droppedFiles])
                    }}
                  >
                    <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 block mb-3">cloud_upload</span>
                    <p className="text-sm font-semibold text-on-surface-variant">Drag &amp; drop files here or click to browse</p>
                    <p className="text-xs text-on-surface-variant/60 mt-1">PNG, JPG, PDF, DOCX up to 10MB each</p>
                    <button
                      type="button"
                      onClick={(e) => { e.stopPropagation(); fileInputRef.current?.click() }}
                      className="mt-4 px-5 py-2 bg-primary text-white rounded-xl text-sm font-bold hover:opacity-90 transition-colors"
                    >
                      Browse Files
                    </button>
                  </div>

                  {/* File list */}
                  {files.length > 0 && (
                    <div className="mt-4 space-y-2">
                      <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider">{files.length} file(s) selected</p>
                      {files.map((f, i) => (
                        <div key={i} className="flex items-center justify-between p-3 bg-surface-container-low rounded-xl">
                          <div className="flex items-center gap-3 min-w-0">
                            <span className="material-symbols-outlined text-primary text-lg flex-shrink-0">
                              {f.type.startsWith('image/') ? 'image' : 'description'}
                            </span>
                            <div className="min-w-0">
                              <p className="text-sm font-medium text-on-surface truncate">{f.name}</p>
                              <p className="text-[10px] text-on-surface-variant">{(f.size / 1024).toFixed(1)} KB</p>
                            </div>
                          </div>
                          <button
                            type="button"
                            onClick={() => setFiles(prev => prev.filter((_, idx) => idx !== i))}
                            className="p-1 hover:bg-error-container rounded-lg transition-colors flex-shrink-0"
                          >
                            <span className="material-symbols-outlined text-error text-lg">close</span>
                          </button>
                        </div>
                      ))}
                    </div>
                  )}

                  <p className="text-xs text-on-surface-variant mt-4 text-center">
                    {files.length === 0 ? 'Attachments are optional. You can skip this step.' : 'Files will be uploaded when you submit the request.'}
                  </p>
                </section>
              )}

              {step === 3 && (
                <section className="bg-surface-container-lowest p-8 rounded-xl shadow-sm">
                  <h2 className="text-xl font-bold mb-1">Review &amp; Submit</h2>
                  <p className="text-sm text-on-surface-variant mb-6">Please review your request before submitting.</p>
                  <div className="space-y-4">
                    <div className="p-4 bg-surface-container-low rounded-xl flex justify-between items-center">
                      <div>
                        <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Request Type</p>
                        <p className="text-sm font-semibold text-on-surface capitalize">{type.replace(/_/g, ' ')}</p>
                      </div>
                      <button onClick={() => setStep(0)} className="text-xs text-primary font-bold hover:underline">Edit</button>
                    </div>
                    <div className="p-4 bg-surface-container-low rounded-xl flex justify-between items-center">
                      <div>
                        <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Title</p>
                        <p className="text-sm font-semibold text-on-surface">{title}</p>
                      </div>
                      <button onClick={() => setStep(1)} className="text-xs text-primary font-bold hover:underline">Edit</button>
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div className="p-4 bg-surface-container-low rounded-xl">
                        <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Category</p>
                        <p className="text-sm font-semibold text-on-surface">{category}</p>
                      </div>
                      <div className="p-4 bg-surface-container-low rounded-xl">
                        <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Priority</p>
                        <p className="text-sm font-semibold text-on-surface capitalize">{priority}</p>
                      </div>
                    </div>
                    <div className="p-4 bg-surface-container-low rounded-xl">
                      <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Description</p>
                      <p className="text-sm text-on-surface-variant leading-relaxed">{description}</p>
                    </div>
                    {assigneeId && (
                      <div className="p-4 bg-surface-container-low rounded-xl flex justify-between items-center">
                        <div>
                          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Assigned To</p>
                          <p className="text-sm font-semibold text-on-surface">
                            {usersList.find(u => u.id === assigneeId)?.full_name ?? assigneeId}
                          </p>
                        </div>
                        <button onClick={() => setStep(1)} className="text-xs text-primary font-bold hover:underline">Edit</button>
                      </div>
                    )}
                    {files.length > 0 && (
                      <div className="p-4 bg-surface-container-low rounded-xl flex justify-between items-center">
                        <div>
                          <p className="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-1">Attachments</p>
                          <p className="text-sm font-semibold text-on-surface">{files.length} file(s): {files.map(f => f.name).join(', ')}</p>
                        </div>
                        <button onClick={() => setStep(2)} className="text-xs text-primary font-bold hover:underline">Edit</button>
                      </div>
                    )}
                  </div>
                </section>
              )}

              {error && (
                <div className="bg-error-container/30 text-error px-4 py-3 rounded-xl flex items-center gap-2 text-sm">
                  <span className="material-symbols-outlined text-sm">error</span>
                  {error}
                </div>
              )}

              {/* Navigation buttons */}
              <div className="flex items-center justify-between pt-2 pb-12">
                <div>
                  {step > 0 && (
                    <button type="button" onClick={() => goToStep(step - 1)}
                      className="px-6 py-3 rounded-xl font-bold text-on-surface-variant hover:bg-surface-container-high transition-colors flex items-center gap-2">
                      <span className="material-symbols-outlined text-sm">arrow_back</span>
                      Back
                    </button>
                  )}
                </div>
                <div className="flex items-center gap-4">
                  <Link to="/tickets" className="px-6 py-3 rounded-xl font-bold text-on-surface-variant hover:bg-surface-container-high transition-colors">
                    Cancel
                  </Link>
                  {step < steps.length - 1 ? (
                    <button type="button" onClick={() => goToStep(step + 1)}
                      className="px-10 py-3 bg-gradient-to-r from-primary to-primary-container text-white font-bold rounded-xl shadow-lg shadow-primary/20 hover:opacity-90 active:scale-95 transition-all flex items-center gap-2">
                      Next Step
                      <span className="material-symbols-outlined text-sm">arrow_forward</span>
                    </button>
                  ) : (
                    <button type="button" onClick={handleSubmit} disabled={loading}
                      className="px-10 py-3 bg-gradient-to-r from-primary to-primary-container text-white font-bold rounded-xl shadow-lg shadow-primary/20 hover:opacity-90 active:scale-95 transition-all disabled:opacity-70 flex items-center gap-2">
                      {loading ? 'Submitting...' : 'Submit Request'}
                      {!loading && <span className="material-symbols-outlined text-sm">send</span>}
                    </button>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
