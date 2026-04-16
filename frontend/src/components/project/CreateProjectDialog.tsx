import { useState } from 'react'
import { projectService } from '../../services/project.service'

const COLORS = ['#3b82f6', '#ef4444', '#22c55e', '#f59e0b', '#8b5cf6', '#ec4899', '#06b6d4', '#f97316']

interface CreateProjectDialogProps {
  onClose: () => void
  onCreated: () => void
}

export function CreateProjectDialog({ onClose, onCreated }: CreateProjectDialogProps) {
  const [name, setName] = useState('')
  const [color, setColor] = useState(COLORS[0])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) { setError('Nama project tidak boleh kosong'); return }
    setLoading(true)
    setError('')
    try {
      await projectService.create({ name: name.trim(), icon_color: color })
      onCreated()
    } catch {
      setError('Gagal membuat project')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-bold text-on-surface font-headline">Project Baru</h2>
          <button onClick={onClose} className="p-1 hover:bg-surface-container-high rounded-lg transition-colors">
            <span className="material-symbols-outlined text-on-surface-variant">close</span>
          </button>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div>
            <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">
              Nama Project
            </label>
            <input
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Masukkan nama project..."
              className="w-full px-4 py-2.5 bg-surface-container-lowest border border-outline-variant/20 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none"
              autoFocus
            />
          </div>

          <div>
            <label className="text-xs font-bold text-on-surface-variant uppercase tracking-widest mb-1.5 block">
              Warna Ikon
            </label>
            <div className="flex gap-2">
              {COLORS.map(c => (
                <button
                  key={c}
                  type="button"
                  onClick={() => setColor(c)}
                  className={`w-8 h-8 rounded-full transition-all hover:scale-110 ${
                    color === c ? 'ring-2 ring-offset-2 ring-on-surface scale-110' : ''
                  }`}
                  style={{ backgroundColor: c }}
                />
              ))}
            </div>
          </div>

          {error && (
            <p className="text-sm text-error font-medium">{error}</p>
          )}

          <div className="flex justify-end gap-3 mt-2">
            <button type="button" onClick={onClose} className="px-4 py-2 text-sm font-medium text-on-surface-variant hover:bg-surface-container-high rounded-xl transition-colors">
              Batal
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-5 py-2 text-sm font-bold text-on-primary bg-primary rounded-xl hover:opacity-90 transition-opacity disabled:opacity-50"
            >
              {loading ? 'Membuat...' : 'Buat Project'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
