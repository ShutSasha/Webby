import { useUnreadNotificationsCountQuery } from '@/lib/hooks/api/notifications/useUnreadNotificationsCount'
import { cn } from '@/lib/utils/general.utils'

type Props = {
  activeTab: 'All' | 'Unread'
  setActiveTab: (v: 'All' | 'Unread') => void
  onMarkAllAsRead: () => void
  hasUnreadToMark: boolean
}

export default function NotificationHeader({ activeTab, setActiveTab, onMarkAllAsRead, hasUnreadToMark }: Props) {
  const { data: unreadCount = 0 } = useUnreadNotificationsCountQuery()

  return (
    <div
      className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-neutral-800 pb-6 mb-2"
    >
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold text-foreground-secondary">Notifications</h1>
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
              'px-4 py-1.5 text-sm font-medium rounded-lg transition-all cursor-pointer',
              activeTab === 'All'
                ? 'bg-neutral-700 text-foreground-secondary shadow-sm'
                : 'text-foreground-muted hover:text-foreground-tertiary',
            )}
          >
            All
          </button>
          <button
            onClick={() => setActiveTab('Unread')}
            className={cn(
              'px-4 py-1.5 text-sm font-medium rounded-lg transition-all cursor-pointer',
              activeTab === 'Unread'
                ? 'bg-neutral-700 text-foreground-secondary shadow-sm'
                : 'text-foreground-muted hover:text-foreground-tertiary',
            )}
          >
            Unread
          </button>
        </div>

        <button
          onClick={onMarkAllAsRead}
          disabled={!hasUnreadToMark}
          className={cn(
            'text-sm font-medium transition-colors',
            hasUnreadToMark
              ? 'text-foreground-muted hover:text-emerald-500 cursor-pointer'
              : 'text-foreground-disabled cursor-not-allowed',
          )}
        >
          Mark all as read
        </button>
      </div>
    </div>
  )
}
