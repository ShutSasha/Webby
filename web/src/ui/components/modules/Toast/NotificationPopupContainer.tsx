'use client'

import { AnimatePresence } from 'framer-motion'

import { useNotificationSocket } from '@/lib/hooks/useNotificationSocket'
import { useNotificationPopupStore } from '@/stores/notification-popup.store'

import NotificationPopupItem from '../Notifications/NotificationPopupItem'

export default function NotificationPopupContainer() {
  useNotificationSocket()

  const popups = useNotificationPopupStore(state => state.popups)

  return (
    <div className="fixed bottom-6 right-6 z-9999 flex flex-col gap-3 w-full max-w-xs md:max-w-sm pointer-events-none">
      <AnimatePresence mode="popLayout">
        {popups.map(popup => (
          <NotificationPopupItem key={popup.popupId} popup={popup} />
        ))}
      </AnimatePresence>
    </div>
  )
}
