'use client'

import { useState } from 'react'

import Image from 'next/image'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import { cn } from '@/lib/utils/utils'

interface Video {
  id: string
  title: string
  thumbnail: string
  isActive?: boolean
  isFolder?: boolean
  children?: Video[]
}

type Props = {
  video: Video
  isChild?: boolean
}

export default function PlaylistItem({ video, isChild = false }: Props) {
  const [isOpen, setOpen] = useState(false)

  return (
    <div className={cn('flex flex-col gap-1', isChild && 'ml-4 pl-2 border-l border-neutral-800')}>
      <div
        className={cn(
          `flex items-center justify-between p-2 rounded-xl bg-neutral-900/50 hover:bg-neutral-900 transition-all
          duration-300 border-b-2 border-transparent group cursor-pointer`,
          video.isActive && 'border-emerald-500 bg-neutral-900',
          video.isFolder && 'hover:bg-neutral-800/40',
        )}
      >
        <div className="flex items-center gap-3 overflow-hidden">
          <Image
            src={video.thumbnail}
            alt={video.title}
            width={96}
            height={96}
            className={cn(
              'object-cover size-10 rounded-lg',
              video.isActive ? 'opacity-100' : 'opacity-60 group-hover:opacity-100',
            )}
          />

          <span
            className={cn(
              'text-sm font-medium truncate pr-1 transition-colors',
              video.isActive ? 'text-neutral-300' : 'text-neutral-500 group-hover:text-neutral-300',
            )}
            title={video.title}
          >
            {video.title}
          </span>
        </div>

        <div className="flex items-center gap-1 shrink-0" onClick={e => e.stopPropagation()}>
          <button className="p-1.5 group/play hover:bg-neutral-300/10 rounded-full cursor-pointer">
            {video.isActive ? (
              <PauseIcon
                className={`size-4 ${video.isActive ? 'text-neutral-300' : 'text-neutral-500'}
                  group-hover/play:text-neutral-300`}
              />
            ) : (
              <PlayIcon
                className={`size-4 ${video.isActive ? 'text-neutral-300' : 'text-neutral-500'}
                  group-hover/play:text-neutral-300`}
              />
            )}
          </button>

          <button
            className="group/trash p-1.5 text-neutral-500 hover:text-red-500 transition-colors hover:bg-red-500/10
              cursor-pointer rounded-full"
          >
            <TrashIcon
              className={`size-4 transition-colors ${video.isActive ? 'text-neutral-300' : 'text-neutral-500'}
                group-hover/trash:text-red-500`}
            />
          </button>

          {video.isFolder && (
            <button
              className="p-1.5 text-neutral-500 hover:text-emerald-300 transition-colors hover:bg-neutral-500/10
                cursor-pointer rounded-full"
              onClick={() => video.isFolder && setOpen(prev => !prev)}
            >
              <svg
                className={cn('size-4 text-neutral-600 transition-transform duration-300', isOpen && 'rotate-180')}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
          )}
        </div>
      </div>

      {video.isFolder && isOpen && video.children && (
        <div className="flex flex-col gap-2 mt-1 animate-in fade-in slide-in-from-top-2 duration-300">
          {video.children.map(child => (
            <PlaylistItem key={child.id} video={child} isChild />
          ))}
        </div>
      )}
    </div>
  )
}
