'use client'

import { MouseEvent } from 'react'

import { motion } from 'framer-motion'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import { cn } from '@/lib/utils/utils'
import { usePlayerPlayStore } from '@/stores/player.store'

interface Props {
  id: string
  title: string
  thumbnail: string
  playlistId: string
  optimisticId: string | null
  onOptimisticClick: () => void
}

export default function VideoItem({ id, title, thumbnail, playlistId, optimisticId, onOptimisticClick }: Props) {
  const playing = usePlayerPlayStore(state => state.playing)
  const togglePlay = usePlayerPlayStore(state => state.togglePlay)
  const router = useRouter()
  const searchParams = useSearchParams()
  const actualVideoId = searchParams.get('v')

  const isActuallyActive = actualVideoId === id
  const isOptimisticallyActive = optimisticId === id

  const isActive = isActuallyActive || isOptimisticallyActive

  const isLoading = isOptimisticallyActive && !isActuallyActive

  const handlePlayPause = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()

    if (isLoading) return

    if (!isActive) {
      onOptimisticClick()

      router.push(`/playlists/${playlistId}?v=${id}`, { scroll: false })

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
      <Link
        href={`/playlists/${playlistId}?v=${id}`}
        onClick={() => {
          if (!isActive) onOptimisticClick()
        }}
        className={cn(
          `flex items-center justify-between p-2 rounded-xl bg-neutral-800/50 hover:bg-neutral-800 transition-all
          duration-300 border-b-2 border-transparent group cursor-pointer`,
          isActuallyActive && 'border-emerald-500 bg-neutral-800',
          isOptimisticallyActive && 'border-amber-500 bg-neutral-800',
        )}
        scroll={false}
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
            className={cn(
              'p-1.5 group/play rounded-full cursor-pointer flex items-center justify-center',
              !isLoading && 'hover:bg-neutral-300/10',
              isLoading && 'hover:bg-none',
            )}
            onClick={handlePlayPause}
          >
            {isLoading ? (
              <div className="size-4 border-2 border-amber-500/20 border-t-amber-500 rounded-full animate-spin" />
            ) : playing && isActuallyActive ? (
              <PauseIcon
                className={`size-4 ${isActive ? 'text-neutral-300' : 'text-neutral-500'}
                  group-hover/play:text-neutral-300`}
              />
            ) : (
              <PlayIcon
                className={`size-4 ${isActive ? 'text-neutral-300' : 'text-neutral-500'}
                  group-hover/play:text-neutral-300`}
              />
            )}
          </button>

          <button
            className="group/trash p-1.5 text-neutral-500 hover:text-red-500 transition-colors hover:bg-red-500/10
              cursor-pointer rounded-full"
          >
            <TrashIcon
              className={`size-4 transition-colors ${isActive ? 'text-neutral-300' : 'text-neutral-500'}
                group-hover/trash:text-red-500`}
            />
          </button>
        </div>
      </Link>
    </motion.div>
  )
}
