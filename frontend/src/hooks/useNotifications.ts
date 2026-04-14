import { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import type { AppDispatch, RootState } from '../store'
import { setNotifications, markRead } from '../store/notificationSlice'
import api from '../services/api'
import type { Notification } from '../types'

export function useNotifications() {
  const dispatch = useDispatch<AppDispatch>()
  const { items, unreadCount } = useSelector((s: RootState) => s.notifications)
  const isAuthenticated = useSelector((s: RootState) => s.auth.isAuthenticated)

  useEffect(() => {
    if (!isAuthenticated) return
    api.get<Notification[]>('/notifications').then((res) => {
      dispatch(setNotifications(res.data))
    })
  }, [isAuthenticated, dispatch])

  const markAsRead = async (id: string) => {
    await api.patch(`/notifications/${id}/read`)
    dispatch(markRead(id))
  }

  return { notifications: items, unreadCount, markAsRead }
}
