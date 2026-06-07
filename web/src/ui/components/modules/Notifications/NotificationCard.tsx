'use client'

import Link from 'next/link'

import CheckIcon from '@/assets/icons/Notifications/check.svg'
import PlayCircleIcon from '@/assets/icons/Notifications/circle-play.svg'
import InfoIcon from '@/assets/icons/Notifications/info.svg'
import ListIcon from '@/assets/icons/Notifications/list.svg'
import TrashIcon from '@/assets/icons/Notifications/trash.svg'
import UsersIcon from '@/assets/icons/Notifications/users.svg'
import { formatRelativeTime } from '@/lib/utils/date.utils'
import { cn } from '@/lib/utils/general.utils'
import { Notification } from '@/types/notification.types'

type NotificationCardProps = {
  notification: Notification
  onMarkAsRead?: (id: string) => void
  onDelete?: (id: string) => void
}

export default function NotificationCard({ notification, onMarkAsRead, onDelete }: NotificationCardProps) {
  const isUnread = notification.notificationStatus === 'Unread'

  const isRoomInvite = notification.targetType === 'Room' && notification.title === 'Room invitation'

  return (
    <div
      className={cn(
        'group relative flex gap-4 p-4 rounded-2xl transition-all duration-300 border border-transparent',
        isUnread ? 'bg-neutral-800/40 hover:bg-neutral-800 border-neutral-800/60' : 'hover:bg-neutral-800/40',
      )}
    >
      <div className="flex items-center shrink-0">
        <div
          className={cn(
            'w-12 h-12 flex items-center justify-center rounded-full',
            isUnread ? 'bg-emerald-500/10 text-emerald-500' : 'bg-neutral-800 text-neutral-400',
          )}
        >
          <TypeIcon type={notification.targetType} />
        </div>
      </div>

      <div className="flex gap-4 items-center pr-20">
        <div className="flex flex-col flex-1 gap-1">
          <div className="flex items-start justify-between">
            <h3 className={cn('text-sm font-semibold', isUnread ? 'text-neutral-100' : 'text-neutral-300')}>
              {notification.title}
            </h3>
          </div>
          <p className="text-sm text-neutral-400 leading-relaxed line-clamp-2">{notification.message}</p>
          <span className="text-[12px] font-medium text-neutral-500 mt-1">
            {formatRelativeTime(notification.createdAt)}
          </span>
        </div>
        {isRoomInvite && notification.targetIdentifier && (
          <div className="mt-2">
            <Link
              href={`/rooms/${notification.targetIdentifier}`}
              onClick={() => {
                if (isUnread && onMarkAsRead) {
                  onMarkAsRead(notification.notificationId)
                }
              }}
              className="inline-flex px-4 py-1.5 bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500
                hover:text-white text-xs font-semibold rounded-lg transition-colors cursor-pointer w-max"
            >
              Join Room
            </Link>
          </div>
        )}
      </div>

      {isUnread && <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-emerald-500 rounded-r-full" />}

      <div
        className="absolute right-4 top-1/2 -translate-y-1/2 flex items-center gap-2 opacity-0 group-hover:opacity-100
          transition-opacity"
      >
        {isUnread && (
          <button
            onClick={() => onMarkAsRead?.(notification.notificationId)}
            className="p-2 bg-neutral-700/50 hover:bg-emerald-500/20 text-neutral-300 hover:text-emerald-500 rounded-xl
              transition-colors cursor-pointer"
            title="Mark as read"
          >
            <CheckIcon className="w-5 h-5 stroke-2" />
          </button>
        )}
        <button
          onClick={() => onDelete?.(notification.notificationId)}
          className="p-2 bg-neutral-700/50 hover:bg-red-500/20 text-neutral-300 hover:text-red-400 rounded-xl
            transition-colors cursor-pointer"
          title="Delete notification"
        >
          <TrashIcon className="w-5 h-5 stroke-2" />
        </button>
      </div>
    </div>
  )
}

function TypeIcon({ type }: { type: Notification['targetType'] }) {
  switch (type) {
    case 'Video':
      return <PlayCircleIcon className="w-6 h-6 stroke-2" />
    case 'Playlist':
      return <ListIcon className="w-6 h-6 stroke-2" />
    case 'Room':
      return <UsersIcon className="w-6 h-6 stroke-2" />
    case 'System':
    default:
      return <InfoIcon className="w-6 h-6 stroke-2" />
  }
}
