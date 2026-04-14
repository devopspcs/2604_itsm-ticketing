import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { Notification } from '../types'

interface NotificationState {
  items: Notification[]
  unreadCount: number
}

const initialState: NotificationState = {
  items: [],
  unreadCount: 0,
}

const notificationSlice = createSlice({
  name: 'notifications',
  initialState,
  reducers: {
    setNotifications(state, action: PayloadAction<Notification[]>) {
      const items = action.payload ?? []
      state.items = items
      state.unreadCount = items.filter((n) => !n.is_read).length
    },
    markRead(state, action: PayloadAction<string>) {
      const notif = state.items.find((n) => n.id === action.payload)
      if (notif && !notif.is_read) {
        notif.is_read = true
        state.unreadCount = Math.max(0, state.unreadCount - 1)
      }
    },
  },
})

export const { setNotifications, markRead } = notificationSlice.actions
export default notificationSlice.reducer
