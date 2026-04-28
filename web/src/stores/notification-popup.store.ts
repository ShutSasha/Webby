import { create } from 'zustand'

import { Notification } from '@/types/notification.types'

type Popup = Notification & { popupId: string }

interface NotificationPopupState {
  popups: Popup[]
  unreadCount: number
  addPopup: (notification: Notification) => void
  removePopup: (popupId: string) => void
  setUnreadCount: (count: number) => void
}

export const useNotificationPopupStore = create<NotificationPopupState>(set => ({
  popups: [],
  unreadCount: 0,
  addPopup: notification => {
    const popupId = Math.random().toString(36).substring(2, 9)
    set(state => ({ popups: [...state.popups, { ...notification, popupId }] }))
  },
  removePopup: popupId => set(state => ({ popups: state.popups.filter(p => p.popupId !== popupId) })),
  setUnreadCount: count => set({ unreadCount: count }),
}))
