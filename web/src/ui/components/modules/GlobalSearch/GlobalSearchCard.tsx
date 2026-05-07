import { useEffect, useRef, useState } from 'react'

import { Route } from 'next'
import Link from 'next/link'
import { useSession } from 'next-auth/react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { clog } from '@/lib/utils/general.utils'
import { BLUR_DATA_URLS } from '@/ui/images'

import SafeImage, { FallbackType } from '../../shared/SafeImage'
import SaveToPlaylistModal from '../Playlists/SaveToPlaylistModal'

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
  const { data: session } = useSession()
  const currentUserId = session?.user?.id

  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isPlaylistModalOpen, setIsPlaylistModalOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const isUser = type === 'User'
  const isRoom = type === 'Room'
  const isStream = type === 'Stream'
  const isYoutubeVideo = type === 'YouTube'
  const isWebbyVideo = type === 'Video'
  const isPlaylist = type === 'Playlist'
  const isExternal = type === 'Stream'

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false)
      }
    }
    if (isMenuOpen) document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isMenuOpen])

  const toggleMenu = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(!isMenuOpen)
  }

  const handleAddToPlaylistClick = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(false)
    setIsPlaylistModalOpen(true)
  }

  const handleAddToRoomClick = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(false)
    clog(`[Room] TODO: Add ${type} with ID to Room:`, id)
  }

  const routes: Record<EntityType, Route> = {
    Video: `/videos/${id}` as Route,
    Room: `/rooms/${id}` as Route,
    Playlist: `/playlists/${id}` as Route,
    Stream: `/streams/${id}` as Route,
    User: `/profile/${id}` as Route,
    YouTube: `/videos/${id}` as Route,
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

      <div className="flex items-center gap-1 shrink-0">
        {showAddButton && !isUser && !isRoom && (
          <div className="relative shrink-0" ref={menuRef}>
            <button
              onClick={toggleMenu}
              className="p-1.5 rounded-full text-neutral-500 hover:text-emerald-500 hover:bg-emerald-500/10
                transition-colors opacity-0 group-hover:opacity-100 shrink-0"
            >
              <PlusIcon className="w-5 h-5 stroke-2" />
            </button>
            {isMenuOpen && (
              <div
                className="absolute right-0 top-full mt-2 w-48 bg-neutral-800 border border-neutral-700/60 shadow-xl
                  shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
                onClick={e => e.preventDefault()}
              >
                {(isYoutubeVideo || isWebbyVideo) && (
                  <button
                    onClick={handleAddToPlaylistClick}
                    className="w-full text-left px-4 py-2 text-sm text-neutral-200 hover:bg-neutral-700/50
                      transition-colors flex items-center gap-3 cursor-pointer"
                  >
                    <span>Add to playlist</span>
                  </button>
                )}

                {(isYoutubeVideo || isWebbyVideo || isPlaylist || isStream) && (
                  <button
                    onClick={handleAddToRoomClick}
                    className="w-full text-left px-4 py-2 text-sm text-neutral-200 hover:bg-neutral-700/50
                      transition-colors flex items-center gap-3 cursor-pointer"
                  >
                    <span>Add to room</span>
                  </button>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </>
  )

  const wrapperClasses = `group flex items-center justify-between gap-3 p-2 rounded-xl hover:bg-neutral-800/40 transition-colors border border-transparent hover:border-neutral-800/60 ${
    isExternal ? '' : 'cursor-pointer'
  }`

  return (
    <>
      {isExternal ? (
        <div className={wrapperClasses}>{cardContent}</div>
      ) : (
        <Link href={targetUrl} className={wrapperClasses}>
          {cardContent}
        </Link>
      )}

      <SaveToPlaylistModal
        isOpen={isPlaylistModalOpen}
        onClose={() => setIsPlaylistModalOpen(false)}
        videoId={id}
        userId={currentUserId}
      />
    </>
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
