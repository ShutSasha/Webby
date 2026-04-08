'use client'

import { useState } from 'react'

import { cn } from '@/lib/utils/general.utils'

type NotificationStatus = 'Read' | 'Unread'

type Notification = {
  notificationId: string
  title: string
  message: string
  targetType: 'System' | 'Video' | 'Playlist' | 'Room'
  targetIdentifier: string
  createdAt: string
  notificationStatus: NotificationStatus
}

const MOCK_NOTIFICATIONS: Notification[] = [
  {
    notificationId: '1',
    title: 'Your video was successfully processed',
    message: 'The video "Survive 30 Days Stranded" is now available in high quality.',
    targetType: 'Video',
    targetIdentifier: 'vid-123',
    createdAt: '2026-04-08T10:30:00Z',
    notificationStatus: 'Unread',
  },
  {
    notificationId: '2',
    title: 'New subscriber!',
    message: 'User @qwerty123123 has subscribed to your channel.',
    targetType: 'System',
    targetIdentifier: 'user-456',
    createdAt: '2026-04-08T09:15:00Z',
    notificationStatus: 'Unread',
  },
  {
    notificationId: '3',
    title: 'Room invite',
    message: 'You have been invited to join the private watch room "Anime Night".',
    targetType: 'Room',
    targetIdentifier: 'room-789',
    createdAt: '2026-04-07T20:00:00Z',
    notificationStatus: 'Read',
  },
  {
    notificationId: '4',
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

  const unreadCount = MOCK_NOTIFICATIONS.filter(n => n.notificationStatus === 'Unread').length

  return (
    <>
      <div
        className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-neutral-800 pb-6
          mb-2"
      >
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-neutral-100">Notifications</h1>
          {unreadCount > 0 && (
            <span className="bg-emerald-500/20 text-emerald-500 text-xs font-bold px-2.5 py-0.5 rounded-full">
              {unreadCount} new
            </span>
          )}
        </div>

        <div className="flex items-center gap-4">
          <div className="flex bg-neutral-800/50 p-1 rounded-xl">
            <button
              onClick={() => setActiveTab('All')}
              className={cn(
                'px-4 py-1.5 text-sm font-medium rounded-lg transition-all',
                activeTab === 'All'
                  ? 'bg-neutral-700 text-neutral-100 shadow-sm'
                  : 'text-neutral-400 hover:text-neutral-200',
              )}
            >
              All
            </button>
            <button
              onClick={() => setActiveTab('Unread')}
              className={cn(
                'px-4 py-1.5 text-sm font-medium rounded-lg transition-all',
                activeTab === 'Unread'
                  ? 'bg-neutral-700 text-neutral-100 shadow-sm'
                  : 'text-neutral-400 hover:text-neutral-200',
              )}
            >
              Unread
            </button>
          </div>

          <button
            className="text-sm font-medium text-neutral-400 hover:text-emerald-500 transition-colors cursor-pointer"
          >
            Mark all as read
          </button>
        </div>
      </div>

      <div className="flex flex-col gap-2">
        {filteredNotifications.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-neutral-500">
            <BellIcon className="w-12 h-12 mb-4 opacity-20" />
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
                    <CheckIcon className="w-5 h-5" />
                  </button>
                )}
                <button
                  className="p-2 bg-neutral-700/50 hover:bg-red-500/20 text-neutral-300 hover:text-red-400 rounded-xl
                    transition-colors cursor-pointer"
                  title="Delete notification"
                >
                  <TrashIcon className="w-5 h-5" />
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </>
  )
}

function formatRelativeTime(dateString: string) {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

function TypeIcon({ type }: { type: Notification['targetType'] }) {
  switch (type) {
    case 'Video':
      return <PlayCircleIcon className="w-6 h-6" />
    case 'Playlist':
      return <ListIcon className="w-6 h-6" />
    case 'Room':
      return <UsersIcon className="w-6 h-6" />
    case 'System':
    default:
      return <InfoIcon className="w-6 h-6" />
  }
}

function CheckIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <polyline points="20 6 9 17 4 12" />
    </svg>
  )
}

function TrashIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M3 6h18" />
      <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
      <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
    </svg>
  )
}

function BellIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
      <path d="M13.73 21a2 2 0 0 1-3.46 0" />
    </svg>
  )
}

function PlayCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <circle cx="12" cy="12" r="10" />
      <polygon points="10 8 16 12 10 16 10 8" />
    </svg>
  )
}

function ListIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <line x1="8" y1="6" x2="21" y2="6" />
      <line x1="8" y1="12" x2="21" y2="12" />
      <line x1="8" y1="18" x2="21" y2="18" />
      <line x1="3" y1="6" x2="3.01" y2="6" />
      <line x1="3" y1="12" x2="3.01" y2="12" />
      <line x1="3" y1="18" x2="3.01" y2="18" />
    </svg>
  )
}

function UsersIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  )
}

function InfoIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="16" x2="12" y2="12" />
      <line x1="12" y1="8" x2="12.01" y2="8" />
    </svg>
  )
}
