'use client'

import { useState } from 'react'

import BellIcon from '@/assets/icons/Notifications/bell.svg'
import CheckIcon from '@/assets/icons/Notifications/check.svg'
import PlayCircleIcon from '@/assets/icons/Notifications/circle-play.svg'
import InfoIcon from '@/assets/icons/Notifications/info.svg'
import ListIcon from '@/assets/icons/Notifications/list.svg'
import TrashIcon from '@/assets/icons/Notifications/trash.svg'
import UsersIcon from '@/assets/icons/Notifications/users.svg'
import { formatRelativeTime } from '@/lib/utils/date.utils'
import { cn } from '@/lib/utils/general.utils'
import { Notification } from '@/types/notification.types'
import NotificationHeader from '@/ui/components/modules/Notifications/NotificationHeader'

const MOCK_NOTIFICATIONS: Notification[] = [
  {
    notificationId: '1',
    userId: '1',
    title: 'Your video was successfully processed',
    message: 'The video "Survive 30 Days Stranded" is now available in high quality.',
    targetType: 'Video',
    targetIdentifier: 'vid-123',
    createdAt: '2026-04-08T10:30:00Z',
    notificationStatus: 'Unread',
  },
  {
    notificationId: '2',
    userId: '1',
    title: 'New subscriber!',
    message: 'User @qwerty123123 has subscribed to your channel.',
    targetType: 'System',
    targetIdentifier: 'user-456',
    createdAt: '2026-04-08T09:15:00Z',
    notificationStatus: 'Unread',
  },
  {
    notificationId: '3',
    userId: '1',
    title: 'Room invite',
    message: 'You have been invited to join the private watch room "Anime Night".',
    targetType: 'Room',
    targetIdentifier: 'room-789',
    createdAt: '2026-04-07T20:00:00Z',
    notificationStatus: 'Read',
  },
  {
    notificationId: '4',
    userId: '1',
    title: 'Playlist updated',
    message: '3 new videos were added to the "My favorites" playlist.',
    targetType: 'Playlist',
    targetIdentifier: 'pl-101',
    createdAt: '2026-04-05T14:20:00Z',
    notificationStatus: 'Read',
  },
]

export default function NotificationsPage() {
  const [activeTab, setActiveTab] = useState<'All' | 'Unread'>('All')

  const filteredNotifications = MOCK_NOTIFICATIONS.filter(
    notif => activeTab === 'All' || notif.notificationStatus === 'Unread',
  )

  return (
    <>
      <NotificationHeader activeTab={activeTab} setActiveTab={setActiveTab} />

      <div className="flex flex-col gap-2">
        {filteredNotifications.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-neutral-500">
            <BellIcon className="w-12 h-12 mb-4 opacity-20 stroke-2" />
            <p>You have no {activeTab === 'Unread' ? 'unread ' : ''}notifications right now.</p>
          </div>
        ) : (
          filteredNotifications.map(notification => (
            <div
              key={notification.notificationId}
              className={cn(
                'group relative flex gap-4 p-4 rounded-2xl transition-all duration-300 border border-transparent',
                notification.notificationStatus === 'Unread'
                  ? 'bg-neutral-800/40 hover:bg-neutral-800 border-neutral-800/60'
                  : 'hover:bg-neutral-800/40',
              )}
            >
              <div className="flex items-center shrink-0">
                <div
                  className={cn(
                    'w-12 h-12 flex items-center justify-center rounded-full',
                    notification.notificationStatus === 'Unread'
                      ? 'bg-emerald-500/10 text-emerald-500'
                      : 'bg-neutral-800 text-neutral-400',
                  )}
                >
                  <TypeIcon type={notification.targetType} />
                </div>
              </div>

              <div className="flex flex-col flex-1 gap-1 pr-20">
                <div className="flex items-start justify-between">
                  <h3
                    className={cn(
                      'text-sm font-semibold',
                      notification.notificationStatus === 'Unread' ? 'text-neutral-100' : 'text-neutral-300',
                    )}
                  >
                    {notification.title}
                  </h3>
                </div>
                <p className="text-sm text-neutral-400 leading-relaxed line-clamp-2">{notification.message}</p>
                <span className="text-[12px] font-medium text-neutral-500 mt-1">
                  {formatRelativeTime(notification.createdAt)}
                </span>
              </div>

              {notification.notificationStatus === 'Unread' && (
                <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-emerald-500 rounded-r-full" />
              )}

              <div
                className="absolute right-4 top-1/2 -translate-y-1/2 flex items-center gap-2 opacity-0
                  group-hover:opacity-100 transition-opacity"
              >
                {notification.notificationStatus === 'Unread' && (
                  <button
                    className="p-2 bg-neutral-700/50 hover:bg-emerald-500/20 text-neutral-300 hover:text-emerald-500
                      rounded-xl transition-colors cursor-pointer"
                    title="Mark as read"
                  >
                    <CheckIcon className="w-5 h-5 stroke-2" />
                  </button>
                )}
                <button
                  className="p-2 bg-neutral-700/50 hover:bg-red-500/20 text-neutral-300 hover:text-red-400 rounded-xl
                    transition-colors cursor-pointer"
                  title="Delete notification"
                >
                  <TrashIcon className="w-5 h-5 stroke-2" />
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </>
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
