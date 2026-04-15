import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(email, password)
      navigate('/dashboard')
    } catch {
      setError('Invalid email or password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <body className="bg-surface text-on-surface flex items-center justify-center min-h-screen relative overflow-hidden">
      <div className="absolute inset-0 z-0 bg-gradient-to-br from-primary-fixed/20 via-surface-bright to-surface-bright" />
      <div className="absolute top-[-10%] left-[-5%] w-[40%] h-[40%] bg-tertiary-fixed/10 blur-[120px] rounded-full z-0" />

      <main className="relative z-10 w-full max-w-[440px] px-6">
        {/* Logo */}
        <div className="flex flex-col items-center mb-10">
          <div className="w-16 h-16 bg-gradient-to-br from-primary to-primary-container rounded-xl flex items-center justify-center mb-6 shadow-xl shadow-primary/10">
            <span className="material-symbols-outlined text-white text-3xl">confirmation_number</span>
          </div>
          <h1 className="text-3xl font-extrabold tracking-tight text-primary text-center">PCS Ticketing System</h1>
          <p className="text-on-surface-variant font-medium mt-2">IT Service Management</p>
        </div>

        {/* Card */}
        <div className="bg-surface-container-lowest p-8 md:p-10 rounded-xl shadow-[0px_24px_48px_rgba(25,28,29,0.04)] border border-outline-variant/10">
          <div className="mb-8">
            <h2 className="text-xl font-bold text-on-surface tracking-tight">Welcome Back</h2>
            <p className="text-on-surface-variant text-sm mt-1">Please enter your credentials to access the console.</p>
          </div>

          <form className="space-y-6" onSubmit={handleSubmit}>
            <div>
              <label className="block text-xs font-semibold text-on-surface-variant uppercase tracking-wider mb-2" htmlFor="email">
                Email or Username
              </label>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-outline group-focus-within:text-primary transition-colors">
                  <span className="material-symbols-outlined text-[20px]">alternate_email</span>
                </div>
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  placeholder="name@company.com"
                  className="w-full pl-11 pr-4 py-3 bg-surface-container-highest border-none rounded-xl text-on-surface placeholder:text-outline focus:ring-2 focus:ring-primary/20 focus:bg-surface-container-lowest transition-all duration-200 text-sm outline-none"
                />
              </div>
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <label className="block text-xs font-semibold text-on-surface-variant uppercase tracking-wider" htmlFor="password">
                  Password
                </label>
              </div>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-outline group-focus-within:text-primary transition-colors">
                  <span className="material-symbols-outlined text-[20px]">lock</span>
                </div>
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  placeholder="••••••••"
                  className="w-full pl-11 pr-12 py-3 bg-surface-container-highest border-none rounded-xl text-on-surface placeholder:text-outline focus:ring-2 focus:ring-primary/20 focus:bg-surface-container-lowest transition-all duration-200 text-sm outline-none"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute inset-y-0 right-0 pr-4 flex items-center text-outline hover:text-on-surface transition-colors"
                >
                  <span className="material-symbols-outlined text-[20px]">{showPassword ? 'visibility_off' : 'visibility'}</span>
                </button>
              </div>
            </div>

            {error && (
              <p role="alert" className="text-sm text-error bg-error-container/30 px-4 py-2 rounded-xl">{error}</p>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 px-4 bg-gradient-to-r from-primary to-primary-container text-white font-bold rounded-xl shadow-lg shadow-primary/20 hover:shadow-primary/30 active:scale-[0.98] transition-all duration-150 flex items-center justify-center gap-2 disabled:opacity-70"
            >
              {loading ? 'Signing in...' : 'Sign In'}
              {!loading && <span className="material-symbols-outlined text-[20px]">arrow_forward</span>}
            </button>
          </form>

          <div className="mt-10 flex flex-col items-center gap-4">
            <div className="flex items-center w-full gap-4">
              <div className="h-[1px] bg-outline-variant/30 flex-grow" />
              <span className="text-[10px] font-bold uppercase tracking-[0.2em] text-outline">Single Sign-On</span>
              <div className="h-[1px] bg-outline-variant/30 flex-grow" />
            </div>
            <a
              href="/api/v1/auth/sso/redirect"
              className="w-full py-3 px-4 bg-surface-container-high rounded-xl text-on-secondary-container font-semibold text-sm hover:bg-surface-container-highest transition-colors flex items-center justify-center gap-2 border border-outline-variant/10 no-underline"
            >
              <span className="material-symbols-outlined text-[18px]">shield_person</span>
              Login with Keycloak SSO
            </a>
          </div>
        </div>

        <footer className="mt-8 text-center">
          <p className="text-xs text-on-surface-variant/60 font-medium">
            © 2024 PCS Payments. All Rights Reserved.<br />
            Managed by IT Department.
          </p>
        </footer>
      </main>

      <div className="fixed bottom-8 right-8 hidden lg:block opacity-40">
        <div className="text-[140px] font-black text-primary/5 select-none pointer-events-none leading-none">PCS</div>
      </div>
    </body>
  )
}
