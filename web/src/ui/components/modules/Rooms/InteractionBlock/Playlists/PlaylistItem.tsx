'use client'

import { useState } from 'react'

import Image from 'next/image'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import PauseIcon from '@/assets/icons/Player/pause.svg'
import PlayIcon from '@/assets/icons/Player/play.svg'
import { cn } from '@/lib/utils/utils'

type Props = {
  video: {
    id: string
    title: string
    thumbnail: string
    isActive?: boolean
    isFolder?: boolean
  }
}

export default function PlaylistItem({ video }: Props) {
  const [isOpen, setOpen] = useState(false)

  return (
    <>
      <div
        className={cn(
          `flex items-center justify-between p-2 rounded-xl bg-neutral-900/50 hover:bg-neutral-900 transition-all
          duration-300 border-b-2 border-transparent group cursor-pointer`,
          video.isActive && 'border-emerald-500 ',
        )}
      >
        <div className="flex items-center gap-3 overflow-hidden">
          <Image
            src={video.thumbnail}
            alt={video.title}
            width={96}
            height={96}
            className="object-cover size-12 rounded-xl"
          />

          <span
            className={cn(
              'text-sm font-medium truncate pr-1',
              video.isActive ? 'text-neutral-300' : 'text-neutral-500 group-hover:text-neutral-300',
            )}
            title={video.title}
          >
            {video.title}
          </span>
        </div>

        <div className="flex items-center gap-1 shrink-0">
          <button
            className="p-1.5 text-neutral-300 hover:text-neutral-300 transition-colors hover:bg-neutral-300/10
              cursor-pointer rounded-full"
          >
            {video.isActive ? <PauseIcon className="size-4" /> : <PlayIcon className="size-4" />}
          </button>

          <button
            className="p-1.5 text-neutral-300 hover:text-red-500 transition-colors hover:bg-red-500/10 cursor-pointer
              rounded-full"
          >
            <TrashIcon className="size-4" />
          </button>

          {video.isFolder && (
            <button
              className="p-1.5 hover:bg-neutral-300/10 cursor-pointer rounded-full"
              onClick={() => setOpen(prev => !prev)}
            >
              <svg
                className={cn('size-4 text-neutral-300 transition-transform', isOpen && 'rotate-180')}
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
      {isOpen && <p>text</p>}
    </>
  )
}
