import { Route } from 'next'
import Link from 'next/link'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { clog } from '@/lib/utils/general.utils'
import { BLUR_DATA_URLS } from '@/ui/images'

import SafeImage, { FallbackType } from '../../shared/SafeImage'

type EntityType = 'Video' | 'Room' | 'Playlist' | 'Stream' | 'User' | 'YouTube'

type Props = {
  id: string
  title: string
  subtitle: string
  thumbnail: string
  type: EntityType

  showAddButton?: boolean
}

const FALLBACK_MAP: Record<EntityType, FallbackType> = {
  User: 'user',
  Video: 'video',
  Room: 'room',
  Playlist: 'playlist',
  Stream: 'video',
  YouTube: 'video',
}

export default function GlobalSearchCard({ id, title, subtitle, thumbnail, type, showAddButton = true }: Props) {
  const isUser = type === 'User'
  const isRoom = type === 'Room'

  const isExternal = type === 'Stream'

  const handleAddToQueue = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    clog(`[Queue] Added ${type} with ID:`, id)
  }

  const routes: Record<EntityType, Route> = {
    Video: `/videos/${id}?source=Webby` as Route,
    Room: `/rooms/${id}` as Route,
    Playlist: `/playlists/${id}` as Route,
    Stream: `/streams/${id}` as Route,
    User: `/profile/${id}` as Route,
    YouTube: `/videos/${id}?source=YouTube` as Route,
  }

  const targetUrl = routes[type]

  const cardContent = (
    <>
      <div className="flex items-center gap-3 overflow-hidden">
        <div
          className={`relative shrink-0 overflow-hidden bg-neutral-800 flex items-center justify-center
            ${isUser ? 'w-12 h-12 rounded-full' : 'w-24 h-14 rounded-lg'}`}
        >
          <SafeImage
            src={thumbnail}
            alt={title}
            fallbackType={FALLBACK_MAP[type]}
            width={160}
            height={90}
            className="object-cover aspect-video h-full w-full"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral800']}
          />
        </div>
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
    </>
  )

  const wrapperClasses = `group flex items-center justify-between gap-3 p-2 rounded-xl hover:bg-neutral-800/40 transition-colors border border-transparent hover:border-neutral-800/60 ${
    isExternal ? '' : 'cursor-pointer'
  }`

  if (isExternal) {
    return <div className={wrapperClasses}>{cardContent}</div>
  }

  return (
    <Link href={targetUrl} className={wrapperClasses}>
      {cardContent}
    </Link>
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
