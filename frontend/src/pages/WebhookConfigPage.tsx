import { useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import api from '../services/api'
import type { RootState } from '../store'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import { ErrorMessage } from '../components/common/ErrorMessage'

interface WebhookConfig {
  id: string
  url: string
  events: string[]
  is_active: boolean
}

const ALL_EVENTS = ['ticket.created', 'ticket.status_changed', 'ticket.assigned', 'approval.requested', 'approval.decided']

export function WebhookConfigPage() {
  const [configs, setConfigs] = useState<WebhookConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ url: '', events: [] as string[], secret: '' })
  const role = useSelector((s: RootState) => s.auth.role)

  if (role !== 'admin') {
    return (
      <div className="max-w-2xl mx-auto p-8 pt-16 text-center">
        <span className="material-symbols-outlined text-6xl text-on-surface-variant/30 block mb-4">lock</span>
        <h2 className="text-2xl font-bold text-on-surface mb-2">Access Restricted</h2>
        <p className="text-on-surface-variant">Webhook configuration is only available to administrators.</p>
      </div>
    )
  }

  const fetchConfigs = () => {
    api.get<WebhookConfig[]>('/webhooks').then((r) => setConfigs(r.data ?? [])).catch(() => setError('Failed to load webhooks')).finally(() => setLoading(false))
  }

  useEffect(() => { fetchConfigs() }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.post('/webhooks', form)
      setShowForm(false)
      setForm({ url: '', events: [], secret: '' })
      fetchConfigs()
    } catch { setError('Failed to create webhook') }
  }

  const handleDelete = async (id: string) => {
    try { await api.delete(`/webhooks/${id}`); fetchConfigs() }
    catch { setError('Failed to delete webhook') }
  }

  const toggleEvent = (ev: string) => {
    setForm((f) => ({ ...f, events: f.events.includes(ev) ? f.events.filter((e) => e !== ev) : [...f.events, ev] }))
  }

  const inputStyle = { width: '100%', padding: '8px 12px', border: '1px solid #cbd5e1', borderRadius: 4, boxSizing: 'border-box' as const }

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h2 style={{ fontSize: 22, fontWeight: 700 }}>Webhook Configurations</h2>
        <button onClick={() => setShowForm(!showForm)} style={{ background: '#1e40af', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: 4, cursor: 'pointer' }}>
          + New Webhook
        </button>
      </div>

      {error && <ErrorMessage message={error} />}

      {showForm && (
        <form onSubmit={handleCreate} style={{ background: '#fff', padding: '1.5rem', borderRadius: 8, marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
          <input placeholder="URL" type="url" value={form.url} onChange={(e) => setForm({ ...form, url: e.target.value })} required style={inputStyle} />
          <input placeholder="Secret Key" value={form.secret} onChange={(e) => setForm({ ...form, secret: e.target.value })} required style={inputStyle} />
          <div>
            <p style={{ marginBottom: 8, fontWeight: 500 }}>Events:</p>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
              {ALL_EVENTS.map((ev) => (
                <label key={ev} style={{ display: 'flex', alignItems: 'center', gap: 4, cursor: 'pointer' }}>
                  <input type="checkbox" checked={form.events.includes(ev)} onChange={() => toggleEvent(ev)} />
                  {ev}
                </label>
              ))}
            </div>
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            <button type="submit" style={{ background: '#1e40af', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: 4, cursor: 'pointer' }}>Create</button>
            <button type="button" onClick={() => setShowForm(false)} style={{ background: '#e2e8f0', border: 'none', padding: '8px 16px', borderRadius: 4, cursor: 'pointer' }}>Cancel</button>
          </div>
        </form>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        {configs.map((c) => (
          <div key={c.id} style={{ background: '#fff', padding: '1rem 1.5rem', borderRadius: 8, border: '1px solid #e2e8f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div>
              <div style={{ fontWeight: 600 }}>{c.url}</div>
              <div style={{ fontSize: 13, color: '#64748b', marginTop: 4 }}>{c.events.join(', ')}</div>
            </div>
            <button onClick={() => handleDelete(c.id)} style={{ background: '#ef4444', color: '#fff', border: 'none', padding: '6px 12px', borderRadius: 4, cursor: 'pointer' }}>Delete</button>
          </div>
        ))}
      </div>
    </div>
  )
}
