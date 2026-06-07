'use client'

import { MouseEvent } from 'react'

import { motion } from 'framer-motion'
import Image from 'next/image'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import { useTogglePlaylistVideoMutation } from '@/lib/hooks/api/playlist/useTogglePlaylistVideo'
import { cn } from '@/lib/utils/general.utils'
import { usePlayerPlayStore } from '@/stores/player.store'
import { usePlaylistStore } from '@/stores/playlist.store'
import { useToastStore } from '@/stores/toast-store'

interface Props {
  id: string
  title: string
  thumbnail: string
  playlistId: string
  isOwner: boolean
  userId: string
}

export default function VideoItem({ id, title, thumbnail, playlistId, isOwner, userId }: Props) {
  const playing = usePlayerPlayStore(state => state.playing)
  const togglePlay = usePlayerPlayStore(state => state.togglePlay)
  const addToast = useToastStore(state => state.addToast)

  const activeVideoId = usePlaylistStore(state => state.activeVideoId)
  const setActiveVideoId = usePlaylistStore(state => state.setActiveVideoId)

  const isActive = activeVideoId === id

  const { mutate: toggleVideo, isPending } = useTogglePlaylistVideoMutation(userId)

  const handleToggle = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    e.stopPropagation()
    if (isPending) return

    toggleVideo(
      { playlistId, videoId: id },
      {
        onError: () => addToast(`Failed to delete video in playlist`, 'error'),
      },
    )
  }

  const handlePlayPause = (e: MouseEvent<HTMLDivElement | HTMLButtonElement>) => {
    e.preventDefault()

    if (!isActive) {
      setActiveVideoId(id)
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
          `flex items-center justify-between p-2 rounded-xl bg-neutral-800/50 hover:bg-neutral-800 transition-all
          duration-300 border-b-2 border-transparent group cursor-pointer`,
          isActive && 'border-emerald-500 bg-neutral-800',
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
              'object-cover size-10 rounded-lg shrink-0',
              isActive ? 'opacity-100' : 'opacity-60 group-hover:opacity-100',
            )}
          />

          <span
            className={cn(
              'text-sm font-medium truncate pr-1 transition-colors',
              isActive ? 'text-neutral-300' : 'text-neutral-500 group-hover:text-neutral-300',
            )}
            title={title}
          >
            {title}
          </span>
        </div>

        <div className="flex items-center gap-1 shrink-0" onClick={e => e.stopPropagation()}>
          <button
            className="p-1.5 group/play hover:bg-neutral-300/10 rounded-full cursor-pointer flex items-center
              justify-center"
            onClick={handlePlayPause}
          >
            {playing && isActive ? (
              <PauseIcon
                className={`size-4 stroke-[1.5px] ${isActive ? 'text-neutral-300' : 'text-neutral-500'}
                  group-hover/play:text-neutral-300`}
              />
            ) : (
              <PlayIcon
                className={`size-4 stroke-[1.5px] ${isActive ? 'text-neutral-300' : 'text-neutral-500'}
                  group-hover/play:text-neutral-300`}
              />
            )}
          </button>

          {isOwner && (
            <button
              className="group/trash p-1.5 text-neutral-500 hover:text-red-500 transition-colors hover:bg-red-500/10
                cursor-pointer rounded-full"
              onClick={handleToggle}
            >
              {isPending ? (
                <div className="size-4 border-2 border-red-500/20 border-t-red-500 rounded-full animate-spin" />
              ) : (
                <TrashIcon
                  className={`size-4 stroke-[1.5px] transition-colors
                    ${isActive ? 'text-neutral-300' : 'text-neutral-500'} group-hover/trash:text-red-500`}
                />
              )}
            </button>
          )}
        </div>
      </div>
    </motion.div>
  )
}

export function VideoItemSkeleton() {
  return (
    <div className="flex items-center justify-between p-2 rounded-xl bg-neutral-800/30 border-b-2 border-transparent">
      <div className="flex items-center gap-3 overflow-hidden w-full">
        <div className="size-10 bg-neutral-700/50 rounded-lg shrink-0 animate-pulse" />

        <div className="h-4 bg-neutral-700/50 rounded-md w-3/4 animate-pulse" />
      </div>
    </div>
  )
}
