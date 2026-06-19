'use client'

import { useMemo, useState } from 'react'

import { useSession } from 'next-auth/react'

import BellIcon from '@/assets/icons/Notifications/bell.svg'
import LockIcon from '@/assets/icons/shared/lock.svg'
import { useDeleteNotification } from '@/lib/hooks/api/notifications/useDeleteNotification'
import { useGetUnreadNotificationsQuery } from '@/lib/hooks/api/notifications/useGetUnreadNotificationsQuery'
import { useGetUserNotificationsQuery } from '@/lib/hooks/api/notifications/useGetUserNotificationsQuery'
import { useSetNotificationsReadStatus } from '@/lib/hooks/api/notifications/useSetNotificationsReadStatus'
import { useInfiniteScroll } from '@/lib/hooks/useInfiniteScroll'
import NotificationCard from '@/ui/components/modules/Notifications/NotificationCard'
import NotificationHeader from '@/ui/components/modules/Notifications/NotificationHeader'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'

export default function NotificationsPage() {
  const [activeTab, setActiveTab] = useState<'All' | 'Unread'>('All')

  const allQuery = useGetUserNotificationsQuery()
  const unreadQuery = useGetUnreadNotificationsQuery()

  const { mutate: deleteNotification } = useDeleteNotification()
  const { mutate: markAsRead } = useSetNotificationsReadStatus()

  const currentQuery = activeTab === 'All' ? allQuery : unreadQuery

  const notifications = useMemo(() => {
    return currentQuery.data?.pages.flatMap(page => page?.data?.items || []) || []
  }, [currentQuery.data])

  const lastElementRef = useInfiniteScroll({
    isLoading: currentQuery.isLoading,
    isFetchingNextPage: currentQuery.isFetchingNextPage,
    hasNextPage: currentQuery.hasNextPage,
    fetchNextPage: currentQuery.fetchNextPage,
  })

  const handleMarkAllAsRead = () => {
    const unreadIds = notifications.filter(n => n.notificationStatus === 'Unread').map(n => n.notificationId)

    if (unreadIds.length > 0) {
      markAsRead(unreadIds)
    }
  }

  const hasUnreadToMark = notifications.some(n => n.notificationStatus === 'Unread')

  const { data: session } = useSession()

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view your notifications"
        description="Please log in to check your recent activity, view updates from your channels, and manage your notification preferences all in one place."
        icon={<LockIcon className="size-10 text-neutral-500 stroke-1" />}
      />
    )
  }

  return (
    <>
      <NotificationHeader
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        onMarkAllAsRead={handleMarkAllAsRead}
        hasUnreadToMark={hasUnreadToMark}
      />

      <div className="flex flex-col gap-2">
        {currentQuery.isLoading && notifications.length === 0 ? (
          <div className="flex justify-center py-20">
            <div className="size-10 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
          </div>
        ) : notifications.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-neutral-500">
            <BellIcon className="w-12 h-12 mb-4 opacity-20 stroke-2" />
            <p>You have no {activeTab === 'Unread' ? 'unread ' : ''}notifications right now.</p>
          </div>
        ) : (
          <>
            {notifications.map((notification, index) => {
              const isLast = notifications.length === index + 1

              const card = (
                <NotificationCard
                  key={notification.notificationId}
                  notification={notification}
                  onMarkAsRead={id => markAsRead([id])}
                  onDelete={id => deleteNotification(id)}
                />
              )

              if (isLast) {
                return (
                  <div ref={lastElementRef} key={`last-${notification.notificationId}`}>
                    {card}
                  </div>
                )
              }

              return card
            })}
          </>
        )}

        {currentQuery.isFetchingNextPage && (
          <div className="flex justify-center py-6">
            <div className="size-6 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
          </div>
        )}
      </div>
    </>
  )
}
