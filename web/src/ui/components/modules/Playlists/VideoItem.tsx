'use client'

import { MouseEvent } from 'react'

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
}

export default function VideoItem({ id, title, thumbnail, playlistId }: Props) {
  const playing = usePlayerPlayStore(state => state.playing)
  const togglePlay = usePlayerPlayStore(state => state.togglePlay)
  const router = useRouter()
  const searchParams = useSearchParams()
  const videoId = searchParams.get('v')
  const isActive = videoId === id

  const handlePlayPause = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()

    if (!isActive) {
      router.push(`/playlists/${playlistId}?v=${id}`)

      if (!playing) togglePlay()
      return
    }

    togglePlay()
  }

  return (
    <Link
      href={`/playlists/${playlistId}?v=${id}`}
      className={cn(
        `flex items-center justify-between p-2 rounded-xl bg-neutral-800/50 hover:bg-neutral-800 transition-all
        duration-300 border-b-2 border-transparent group cursor-pointer`,
        isActive && 'border-emerald-500 bg-neutral-800',
      )}
    >
      <div className="flex items-center gap-3 overflow-hidden">
        <Image
          src={thumbnail}
          alt={title}
          width={96}
          height={96}
          className={cn(
            'object-cover size-10 rounded-lg',
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
          className="p-1.5 group/play hover:bg-neutral-300/10 rounded-full cursor-pointer"
          onClick={handlePlayPause}
        >
          {playing && isActive ? (
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
  )
}
