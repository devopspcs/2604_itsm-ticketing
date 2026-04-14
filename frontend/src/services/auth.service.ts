import api from './api'
import type { TokenPair } from '../types'

export const authService = {
  login: (email: string, password: string) =>
    api.post<TokenPair>('/auth/login', { email, password }),

  refresh: (refreshToken: string) =>
    api.post<TokenPair>('/auth/refresh', { refresh_token: refreshToken }),

  logout: (refreshToken: string) =>
    api.post('/auth/logout', { refresh_token: refreshToken }),
}
