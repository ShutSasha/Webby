'use client'

import { MouseEvent } from 'react'

import { motion } from 'framer-motion'
import Image from 'next/image'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import { useTogglePlaylistMediaMutation } from '@/lib/hooks/api/playlist/useTogglePlaylistMedia'
import { cn } from '@/lib/utils/general.utils'
import { usePlayerPlayStore } from '@/stores/player.store'
import { usePlaylistStore } from '@/stores/playlist.store'
import { useToastStore } from '@/stores/toast-store'
import { MediaType } from '@/types/general.types'

interface Props {
  id: string
  title: string
  thumbnail: string
  playlistId: string
  isOwner: boolean
  userId: string
  mediaType: MediaType
}

export default function VideoItem({ id, title, thumbnail, playlistId, isOwner, userId, mediaType }: Props) {
  const playing = usePlayerPlayStore(state => state.playing)
  const togglePlay = usePlayerPlayStore(state => state.togglePlay)
  const addToast = useToastStore(state => state.addToast)

  const activeResourceId = usePlaylistStore(state => state.activeResourceId)
  const setActiveResourceId = usePlaylistStore(state => state.setActiveResourceId)
  const setActiveMediaType = usePlaylistStore(state => state.setActiveMediaType)

  const isActive = activeResourceId === id

  const { mutate: toggleMedia, isPending } = useTogglePlaylistMediaMutation(userId)

  const handleToggle = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    e.stopPropagation()
    if (isPending) return

    toggleMedia(
      { playlistId, mediaId: id, mediaType },
      {
        onError: () => addToast(`Failed to delete media from playlist`, 'error'),
      },
    )
  }

  const handlePlayPause = (e: MouseEvent<HTMLDivElement | HTMLButtonElement>) => {
    e.preventDefault()

    if (!isActive) {
      setActiveResourceId(id)
      setActiveMediaType(mediaType)
      if (!playing) togglePlay()
      return
    }

    togglePlay()
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 15 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: 'easeOut' }}
      layout
    >
      <div
        onClick={handlePlayPause}
        className={cn(
          `flex items-center justify-between p-2 rounded-xl bg-background/60 hover:bg-neutral-200 dark:bg-neutral-800/70
          dark:hover:bg-neutral-300/10 transition-all duration-300 border-b-2 border-transparent group cursor-pointer`,
          isActive && 'border-emerald-500 bg-background/80 dark:bg-neutral-800/70',
          isPending && 'opacity-50 pointer-events-none',
          !isOwner && 'gap-2',
        )}
      >
        <div className="flex items-center gap-3 overflow-hidden">
          <Image
            src={thumbnail}
            alt={title}
            width={96}
            height={96}
            className={cn(
              'object-cover size-10 rounded-lg shadow-sm ring-1 ring-border/50 shrink-0 transition-opacity duration-300',
            )}
          />

          <span
            className={cn(
              'text-sm truncate pr-1 transition-colors duration-200',
              isActive
                ? 'text-foreground-strong font-semibold'
                : 'text-foreground-secondary font-medium group-hover:text-foreground-strong',
            )}
            title={title}
          >
            {title}
          </span>
        </div>

        <div className="flex items-center gap-1 shrink-0" onClick={e => e.stopPropagation()}>
          <button
            className={cn(
              'group/play p-1.5 transition-colors duration-200 cursor-pointer rounded-full shrink-0',
              isActive ? 'text-foreground-muted' : 'text-muted group-hover:text-foreground-muted',
              !isPending && 'hover:text-emerald-500 hover:bg-emerald-500/10',
            )}
            onClick={handlePlayPause}
          >
            {playing && isActive ? (
              <PauseIcon
                className={cn(
                  'size-4 stroke-[1.5px] transition-transform duration-200',
                  !isPending && 'group-hover/play:scale-110',
                )}
              />
            ) : (
              <PlayIcon
                className={cn(
                  'size-4 stroke-[1.5px] transition-transform duration-200',
                  !isPending && 'group-hover/play:scale-110',
                )}
              />
            )}
          </button>

          {isOwner && (
            <button
              className={cn(
                'group/trash p-1.5 transition-colors duration-200 cursor-pointer rounded-full shrink-0',
                isActive ? 'text-foreground-muted' : 'text-muted group-hover:text-foreground-muted',
                !isPending && 'hover:text-red-500 hover:bg-red-500/10',
              )}
              onClick={handleToggle}
              disabled={isPending}
            >
              <TrashIcon
                className={cn(
                  'size-4 stroke-[1.5px] transition-transform duration-200',
                  !isPending && 'group-hover/trash:scale-110',
                  isPending && 'animate-pulse text-red-500',
                )}
              />
            </button>
          )}
        </div>
      </div>
    </motion.div>
  )
}

export function VideoItemSkeleton() {
  return (
    <div className="flex items-center justify-between p-2 rounded-xl bg-background/30 border-b-2 border-transparent">
      <div className="flex items-center gap-3 overflow-hidden w-full">
        <div className="size-10 bg-skeleton/80 rounded-lg shrink-0 animate-pulse" />
        <div className="h-4 bg-skeleton/80 rounded-md w-3/4 animate-pulse" />
      </div>
    </div>
  )
}
