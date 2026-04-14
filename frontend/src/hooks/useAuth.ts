import { useSelector, useDispatch } from 'react-redux'
import type { RootState, AppDispatch } from '../store'
import { setTokens, clearTokens } from '../store/authSlice'
import { authService } from '../services/auth.service'

export function useAuth() {
  const dispatch = useDispatch<AppDispatch>()
  const { isAuthenticated, accessToken } = useSelector((s: RootState) => s.auth)

  const login = async (email: string, password: string) => {
    const res = await authService.login(email, password)
    dispatch(setTokens({ accessToken: res.data.access_token, refreshToken: res.data.refresh_token }))
  }

  const logout = async () => {
    const refreshToken = localStorage.getItem('refresh_token') ?? ''
    try { await authService.logout(refreshToken) } catch {}
    dispatch(clearTokens())
  }

  return { isAuthenticated, accessToken, login, logout }
}
