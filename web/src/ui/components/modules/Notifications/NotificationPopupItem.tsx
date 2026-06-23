'use client'
import { useEffect, useRef, useState } from 'react'

import { motion, useAnimation } from 'framer-motion'
import Link from 'next/link'

import PlayCircleIcon from '@/assets/icons/Notifications/circle-play.svg'
import InfoIcon from '@/assets/icons/Notifications/info.svg'
import ListIcon from '@/assets/icons/Notifications/list.svg'
import UsersIcon from '@/assets/icons/Notifications/users.svg'
import XIcon from '@/assets/icons/shared/x.svg'
import { useNotificationPopupStore } from '@/stores/notification-popup.store'
import { Notification } from '@/types/notification.types'

type Props = {
  popup: Notification & { popupId: string }
}

function TypeIcon({ type }: { type: Notification['targetType'] }) {
  switch (type) {
    case 'Video':
      return <PlayCircleIcon className="w-5 h-5 stroke-2" />
    case 'Playlist':
      return <ListIcon className="w-5 h-5 stroke-2" />
    case 'Room':
      return <UsersIcon className="w-5 h-5 stroke-2" />
    case 'System':
    default:
      return <InfoIcon className="w-5 h-5 stroke-2" />
  }
}

export default function NotificationPopupItem({ popup }: Props) {
  const removePopup = useNotificationPopupStore(state => state.removePopup)
  const [remaining, setRemaining] = useState(5000)
  const timerRef = useRef<NodeJS.Timeout | null>(null)
  const startTimeRef = useRef(Date.now())

  const controls = useAnimation()

  const startTimer = () => {
    startTimeRef.current = Date.now()
    timerRef.current = setTimeout(() => removePopup(popup.popupId), remaining)

    controls.start({
      height: '0%',
      transition: { duration: remaining / 1000, ease: 'linear' },
    })
  }

  const pauseTimer = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current)
      setRemaining(prev => prev - (Date.now() - startTimeRef.current))
    }

    controls.stop()
  }

  useEffect(() => {
    startTimer()
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [])

  const isRoomInvite = popup.targetType === 'Room' && popup.title === 'Room invitation'

  return (
    <motion.div
      layout
      initial={{ opacity: 0, x: 100, scale: 0.9 }}
      animate={{ opacity: 1, x: 0, scale: 1 }}
      exit={{ opacity: 0, scale: 0.9, transition: { duration: 0.2 } }}
      onMouseEnter={pauseTimer}
      onMouseLeave={startTimer}
      className="pointer-events-auto relative flex gap-3 p-4 rounded-xl bg-neutral-900 border border-neutral-800
        shadow-2xl shadow-black/50 overflow-hidden cursor-default"
    >
      <div className="flex items-center shrink-0">
        <div className="w-10 h-10 flex items-center justify-center rounded-full bg-emerald-500/10 text-emerald-500">
          <TypeIcon type={popup.targetType} />
        </div>
      </div>

      <div className="flex flex-col flex-1 gap-1 pr-6">
        <h4 className="text-sm font-semibold text-foreground-secondary">{popup.title}</h4>
        <p className="text-[13px] text-foreground-muted leading-snug line-clamp-2">{popup.message}</p>

        {isRoomInvite && popup.targetIdentifier && (
          <Link
            href={`/rooms/${popup.targetIdentifier}`}
            onClick={() => removePopup(popup.popupId)}
            className="mt-1.5 w-max px-4 py-1.5 bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500
              hover:text-foreground-strong text-xs font-semibold rounded-lg transition-colors cursor-pointer"
          >
            Join Room
          </Link>
        )}
      </div>

      <button
        onClick={() => removePopup(popup.popupId)}
        className="absolute top-3 right-3 text-foreground-faint hover:text-foreground-subtle transition-colors
          cursor-pointer"
      >
        <XIcon className="w-4 h-4" />
      </button>

      <motion.div
        initial={{ height: '100%' }}
        animate={controls}
        className="absolute left-0 bottom-0 top-0 w-1 bg-emerald-500 rounded-r-full"
      />
    </motion.div>
  )
}
