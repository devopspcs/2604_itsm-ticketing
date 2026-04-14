import { createSlice, type PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  role: string | null
  userId: string | null
}

function parseJwtPayload(token: string): { role?: string; user_id?: string } {
  try {
    const base64 = token.split('.')[1]
    const decoded = atob(base64.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(decoded)
  } catch {
    return {}
  }
}

function getStoredRole(): string | null {
  const token = localStorage.getItem('access_token')
  if (!token) return null
  return parseJwtPayload(token).role ?? null
}

function getStoredUserId(): string | null {
  const token = localStorage.getItem('access_token')
  if (!token) return null
  return parseJwtPayload(token).user_id ?? null
}

const initialState: AuthState = {
  accessToken: localStorage.getItem('access_token'),
  refreshToken: localStorage.getItem('refresh_token'),
  isAuthenticated: !!localStorage.getItem('access_token'),
  role: getStoredRole(),
  userId: getStoredUserId(),
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setTokens(state, action: PayloadAction<{ accessToken: string; refreshToken: string }>) {
      state.accessToken = action.payload.accessToken
      state.refreshToken = action.payload.refreshToken
      state.isAuthenticated = true
      const payload = parseJwtPayload(action.payload.accessToken)
      state.role = payload.role ?? null
      state.userId = payload.user_id ?? null
      localStorage.setItem('access_token', action.payload.accessToken)
      localStorage.setItem('refresh_token', action.payload.refreshToken)
    },
    clearTokens(state) {
      state.accessToken = null
      state.refreshToken = null
      state.isAuthenticated = false
      state.role = null
      state.userId = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    },
  },
})

export const { setTokens, clearTokens } = authSlice.actions
export default authSlice.reducer
