import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import logoPcs from '../assets/logo-pcs.png'

interface ServiceOverview {
  name: string
  version: string
  status: string
  type: string
}

interface OverviewData {
  total_staff: number
  total_devices: number
  total_tickets: number
  total_vpn: number
  services: ServiceOverview[]
}

export function LandingPage() {
  const [data, setData] = useState<OverviewData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/public/overview')
      .then(r => r.json())
      .then(d => setData(d))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'operational': return 'bg-emerald-500'
      case 'maintenance': return 'bg-amber-500'
      case 'down': return 'bg-red-500'
      default: return 'bg-slate-400'
    }
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'platform': return 'confirmation_number'
      case 'pos': return 'point_of_sale'
      case 'vpn': return 'vpn_lock'
      default: return 'cloud'
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-primary/20">
      {/* Header */}
      <header className="px-6 py-4 flex items-center justify-between max-w-6xl mx-auto">
        <div className="flex items-center gap-3">
          <img src={logoPcs} alt="POSe" className="h-8 w-auto" />
          <span className="text-xs uppercase tracking-widest text-white/50 font-bold">Enterprise Management</span>
        </div>
        <Link to="/login" className="px-5 py-2 bg-white/10 hover:bg-white/20 text-white text-sm font-semibold rounded-lg backdrop-blur-sm border border-white/10 transition-all">
          Sign In
        </Link>
      </header>

      {/* Hero */}
      <main className="max-w-6xl mx-auto px-6 py-16">
        <div className="text-center mb-16">
          <h1 className="text-4xl md:text-5xl font-black text-white tracking-tight mb-4">
            Service Overview
          </h1>
          <p className="text-lg text-white/60 max-w-2xl mx-auto">
            Monitor all organizational services, infrastructure, and operational metrics at a glance.
          </p>
        </div>

        {/* Stats */}
        {loading ? (
          <div className="flex justify-center py-12">
            <div className="animate-spin w-8 h-8 border-4 border-white/30 border-t-white rounded-full" />
          </div>
        ) : data && (
          <>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
              <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-5 text-center">
                <p className="text-3xl font-black text-white">{data.total_staff}</p>
                <p className="text-xs text-white/50 uppercase tracking-wider mt-1 font-bold">Staff Active</p>
              </div>
              <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-5 text-center">
                <p className="text-3xl font-black text-white">{data.total_tickets}</p>
                <p className="text-xs text-white/50 uppercase tracking-wider mt-1 font-bold">Total Tickets</p>
              </div>
              <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-5 text-center">
                <p className="text-3xl font-black text-white">{data.total_vpn}</p>
                <p className="text-xs text-white/50 uppercase tracking-wider mt-1 font-bold">VPN Servers</p>
              </div>
              <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-5 text-center">
                <p className="text-3xl font-black text-emerald-400">{data.services.filter(s => s.status === 'operational').length}/{data.services.length}</p>
                <p className="text-xs text-white/50 uppercase tracking-wider mt-1 font-bold">Services Up</p>
              </div>
            </div>

            {/* Services */}
            <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl overflow-hidden">
              <div className="px-6 py-4 border-b border-white/10">
                <h2 className="text-sm font-bold text-white uppercase tracking-widest">Services</h2>
              </div>
              <div className="divide-y divide-white/5">
                {data.services.map((svc, i) => (
                  <div key={i} className="px-6 py-4 flex items-center justify-between hover:bg-white/5 transition-colors">
                    <div className="flex items-center gap-4">
                      <div className="w-10 h-10 rounded-lg bg-white/10 flex items-center justify-center">
                        <span className="material-symbols-outlined text-white/70 text-[20px]">{getTypeIcon(svc.type)}</span>
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-white">{svc.name}</p>
                        <p className="text-[10px] text-white/40">v{svc.version}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className={`w-2 h-2 rounded-full ${getStatusColor(svc.status)} ${svc.status === 'operational' ? 'animate-pulse' : ''}`} />
                      <span className="text-xs text-white/60 capitalize">{svc.status}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Footer */}
            <div className="text-center mt-12">
              <Link to="/login" className="inline-flex items-center gap-2 px-8 py-3 bg-primary text-white font-bold rounded-xl hover:opacity-90 transition-all shadow-lg shadow-primary/30">
                <span className="material-symbols-outlined text-[18px]">login</span>
                Go to Dashboard
              </Link>
              <p className="text-xs text-white/30 mt-4">© 2026 PCS Payments. All Rights Reserved.</p>
            </div>
          </>
        )}
      </main>
    </div>
  )
}
