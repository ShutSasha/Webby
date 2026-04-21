import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { clog } from '@/lib/utils/general.utils'

import SearchCardImage from './SearchCardImage'

type Props = {
  id: string
  title: string
  subtitle: string
  thumbnail: string
  type: 'Video' | 'Room' | 'Playlist' | 'Stream' | 'User'
  showAddButton?: boolean
}

export default function GlobalSearchCard({ id, title, subtitle, thumbnail, type, showAddButton = true }: Props) {
  const isUser = type === 'User'
  const isRoom = type === 'Room'

  const handleAddToQueue = (e: React.MouseEvent) => {
    e.stopPropagation()

    clog(`[Queue] Added ${type} with ID:`, id)
  }

  return (
    <div
      className="group flex items-center justify-between gap-3 p-2 rounded-xl hover:bg-neutral-800/40 transition-colors
        cursor-pointer border border-transparent hover:border-neutral-800/60"
    >
      <div className="flex items-center gap-3 overflow-hidden">
        <SearchCardImage src={thumbnail} alt={title} isUser={isUser} />
        <div className="flex flex-col overflow-hidden">
          <h4 className="text-sm font-semibold text-neutral-200 truncate">{title}</h4>
          <p className="text-xs text-neutral-500 truncate mt-0.5">{subtitle}</p>
        </div>
      </div>

      {showAddButton && !isUser && !isRoom && (
        <button
          onClick={handleAddToQueue}
          className="p-2 mr-1 rounded-full text-neutral-500 hover:text-emerald-500 hover:bg-emerald-500/10
            transition-colors opacity-0 group-hover:opacity-100 shrink-0"
          title="Add to Queue"
        >
          <PlusIcon className="w-5 h-5 stroke-2" />
        </button>
      )}
    </div>
  )
}

export function GlobalSearchCardSkeleton({ isUser = false }: { isUser?: boolean }) {
  return (
    <div className="flex items-center gap-3 p-2 w-full animate-pulse">
      <div className={`shrink-0 bg-neutral-800/80 ${isUser ? 'w-12 h-12 rounded-full' : 'w-24 h-14 rounded-lg'}`} />
      <div className="flex flex-col gap-2 w-full">
        <div className="h-3.5 bg-neutral-800/80 rounded w-[60%]" />
        <div className="h-3 bg-neutral-800/60 rounded w-[40%]" />
      </div>
    </div>
  )
}
