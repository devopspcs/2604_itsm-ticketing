import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { setTokens } from '../store/authSlice'
import type { AppDispatch } from '../store'

export function SSOCallbackPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const dispatch = useDispatch<AppDispatch>()

  useEffect(() => {
    const accessToken = searchParams.get('access_token')
    const refreshToken = searchParams.get('refresh_token')

    if (accessToken && refreshToken) {
      dispatch(setTokens({ accessToken, refreshToken }))
      navigate('/dashboard', { replace: true })
    } else {
      // No tokens — redirect back to login
      navigate('/login', { replace: true })
    }
  }, [searchParams, dispatch, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-surface">
      <div className="text-center">
        <span className="material-symbols-outlined text-4xl text-primary animate-spin block mb-4">refresh</span>
        <p className="text-on-surface-variant font-medium">Completing SSO login...</p>
      </div>
    </div>
  )
}
