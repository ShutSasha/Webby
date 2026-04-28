import { motion } from 'framer-motion'

import { cn } from '@/lib/utils/general.utils'

export const SEARCH_TABS = ['Videos', 'Rooms', 'Playlists', 'Streams', 'Users'] as const
export type SearchTab = (typeof SEARCH_TABS)[number]

type Props = {
  activeTab: SearchTab
  onChange: (tab: SearchTab) => void
}

export default function GlobalSearchTabs({ activeTab, onChange }: Props) {
  return (
    <div className="flex items-center gap-1 border-b border-neutral-800 mb-4 overflow-x-auto no-scrollbar px-2">
      {SEARCH_TABS.map(tab => (
        <button
          key={tab}
          onClick={() => onChange(tab)}
          className={cn(
            'relative px-4 py-3 text-sm font-medium transition-colors shrink-0 outline-none',
            activeTab === tab ? 'text-emerald-500' : 'text-neutral-400 hover:text-neutral-200',
          )}
        >
          {tab}
          {activeTab === tab && (
            <motion.div
              layoutId="activeSearchTabIndicator"
              className="absolute bottom-0 left-0 right-0 h-0.5 bg-emerald-500 rounded-t-full"
              transition={{ type: 'spring', stiffness: 400, damping: 30 }}
            />
          )}
        </button>
      ))}
    </div>
  )
}
