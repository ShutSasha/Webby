'use client'

import { useEffect, useRef, useState } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { formatTimeAgo, formatVideoTime } from '@/lib/utils/date.utils'
import { cn } from '@/lib/utils/general.utils'
import { formatViews } from '@/lib/utils/video.utils'
import { BLUR_DATA_URLS } from '@/ui/images'

import ComplaintModal from '../../shared/ComplaintModal'

type Props = {
  videoId: string
  previewUrl: string
  duration: number
  title: string
  creator: string
  views: number
  createAt: string
  userAvatar: string
  currentUserId: string | undefined
}

export default function VideoCard({
  videoId,
  previewUrl,
  duration,
  title,
  creator,
  views,
  createAt,
  userAvatar,
  currentUserId,
}: Props) {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const [isComplaintOpen, setIsComplaintOpen] = useState(false)

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

  return (
    <>
      <Link href={`/videos/${videoId}`} className="group flex flex-col gap-3 cursor-pointer">
        <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-neutral-800">
          <Image
            src={previewUrl}
            alt="Video thumbnail"
            fill
            loading="lazy"
            className="object-cover z-0 transition-transform duration-500 ease-out group-hover:scale-105"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />

          <div
            className="absolute bottom-2 right-2 z-10 bg-black/80 px-1.5 py-0.5 rounded text-[11px] font-medium
              text-neutral-100 tracking-wide"
          >
            {formatVideoTime(duration)}
          </div>
        </div>

        <div className="flex gap-3 items-start px-1">
          <Image
            src={userAvatar}
            alt="Creator avatar"
            width={36}
            height={36}
            className="size-9 rounded-full object-cover shrink-0 mt-0.5"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />

          <div className="flex flex-1 justify-between items-start gap-2">
            <div className="flex flex-col overflow-hidden">
              <h3
                className="text-neutral-100 text-sm font-semibold leading-snug line-clamp-2 break-all transition-colors
                  duration-200 group-hover:text-emerald-400"
              >
                {title}
              </h3>

              <div className="flex flex-col mt-1">
                <p className="text-neutral-400 text-xs truncate hover:text-neutral-300 transition-colors">{creator}</p>
                <p className="text-neutral-500 text-[11px] truncate mt-0.5">
                  {formatViews(views)} • {formatTimeAgo(createAt)}
                </p>
              </div>
            </div>

            <div className="relative shrink-0" ref={menuRef}>
              <button
                onClick={toggleMenu}
                className={cn(
                  'p-1.5 -mr-1.5 -mt-1.5 rounded-full transition-all duration-300 cursor-pointer z-20',
                  'hover:bg-neutral-500/20 active:bg-neutral-500/40',
                  isMenuOpen
                    ? 'bg-neutral-500/20 text-neutral-300'
                    : 'text-neutral-400 opacity-0 group-hover:opacity-100 md:opacity-100',
                )}
              >
                <MoreVertical className="size-5 text-neutral-300" />
              </button>

              {isMenuOpen && (
                <div
                  className="absolute right-0 top-full mt-2 w-48 bg-neutral-800 border border-neutral-700/60 shadow-xl
                    shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
                  onClick={e => e.preventDefault()}
                >
                  <button
                    className="w-full text-left px-4 py-2 text-sm text-neutral-200 hover:bg-neutral-700/50
                      transition-colors flex items-center gap-3 cursor-pointer"
                    onClick={e => {
                      e.preventDefault()
                      e.stopPropagation()
                      setIsMenuOpen(false)
                    }}
                  >
                    <span>Add to playlist</span>
                  </button>
                  <button
                    className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 transition-colors
                      flex items-center gap-3 cursor-pointer"
                    onClick={e => {
                      e.preventDefault()
                      e.stopPropagation()
                      setIsMenuOpen(false)
                      setIsComplaintOpen(true)
                    }}
                  >
                    <span>Report</span>
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </Link>
      <ComplaintModal
        isOpen={isComplaintOpen}
        onClose={() => setIsComplaintOpen(false)}
        authorId={currentUserId}
        targetId={videoId}
        targetType="Video"
      />
    </>
  )
}

export function VideoCardSkeleton() {
  return (
    <div className="flex flex-col gap-3 w-full">
      {/* Thumbnail Skeleton */}
      <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-neutral-800/80 animate-pulse" />

      {/* Info Section Skeleton */}
      <div className="flex gap-3 items-start px-1">
        {/* Avatar Skeleton */}
        <div className="size-9 rounded-full bg-neutral-800/80 shrink-0 mt-0.5 animate-pulse" />

        {/* Text Content Skeleton */}
        <div className="flex flex-1 flex-col justify-start gap-3 mt-0.5">
          {/* Title Skeleton (2 lines to match line-clamp-2) */}
          <div className="flex flex-col gap-1.5">
            <div className="h-3.5 bg-neutral-800/80 rounded w-[90%] animate-pulse" />
            <div className="h-3.5 bg-neutral-800/80 rounded w-[70%] animate-pulse" />
          </div>

          {/* Metadata Skeleton (Creator & Views/Time) */}
          <div className="flex flex-col gap-1.5 mt-1">
            <div className="h-2.5 bg-neutral-800/60 rounded w-[40%] animate-pulse" />
            <div className="h-2.5 bg-neutral-800/60 rounded w-[50%] animate-pulse" />
          </div>
        </div>

        {/* Placeholder for the MoreVertical button so the width stays exact */}
        <div className="size-5 shrink-0" />
      </div>
    </div>
  )
}
