import { cn } from '@/lib/utils/general.utils'

const unreadCount = 2

type Props = {
  activeTab: 'All' | 'Unread'
  setActiveTab: (v: 'All' | 'Unread') => void
}

export default function NotificationHeader({ activeTab, setActiveTab }: Props) {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-neutral-800 pb-6
      mb-2">
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

        <button className="text-sm font-medium text-neutral-400 hover:text-emerald-500 transition-colors cursor-pointer">
          Mark all as read
        </button>
      </div>
    </div>
  )
}
