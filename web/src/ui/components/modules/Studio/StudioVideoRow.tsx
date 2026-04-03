'use client'

import { useEffect, useRef, useState } from 'react'

import Image from 'next/image'
import Link from 'next/link'

import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { useDeleteVideo } from '@/lib/hooks/api/video/useDeleteVideo'
import { formatDate, formatVideoTime } from '@/lib/utils/date.utils'
import { cn } from '@/lib/utils/general.utils'
import { formatViews } from '@/lib/utils/video.utils'
import { Video } from '@/types/video.types'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  video: Video
}

export function StudioVideoRow({ video }: Props) {
  const { mutate, isPending } = useDeleteVideo()
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const isReady = video.videoUploadStatus === 'Ready'
  const isUploading = video.videoUploadStatus === 'Uploading'
  const isFailed = video.videoUploadStatus === 'Failed'

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

  const handleDelete = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    if (isPending) return

    mutate(
      { videoId: video.videoId },
      {
        onSuccess: () => {
          setIsMenuOpen(false)
        },
      },
    )
  }

  return (
    <div
      className={cn(
        'flex items-center border-b border-neutral-800 py-3 px-4 transition-colors group',
        isReady ? 'hover:bg-neutral-800/40' : 'bg-neutral-900/50',
      )}
    >
      <div className="flex-1 flex gap-4 min-w-[300px]">
        <div className="relative w-32 aspect-video bg-neutral-800 rounded-lg overflow-hidden shrink-0">
          <Image
            src={
              video.previewUrl ||
              'https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/videos/default_video_thumbnail.png'
            }
            alt="Video thumbnail"
            fill
            className={cn('object-cover transition-all', !isReady && 'opacity-40 grayscale-50')}
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />

          {isReady && (
            <div
              className="absolute bottom-1 right-1 bg-black/80 px-1 py-0.5 rounded text-[10px] font-medium text-white"
            >
              {video.duration ? formatVideoTime(video.duration) : '0:00'}
            </div>
          )}

          {!isReady && (
            <div
              className="absolute inset-0 flex flex-col items-center justify-center bg-neutral-950/40
                backdrop-blur-[2px]"
            >
              {isUploading && (
                <div className="size-6 border-2 border-emerald-500/30 border-t-emerald-500 rounded-full animate-spin" />
              )}
              {isFailed && <span className="text-red-500 font-bold text-lg">!</span>}
            </div>
          )}
        </div>

        <div className="flex flex-col justify-center overflow-hidden min-w-0">
          {isReady ? (
            <Link
              href={`/videos/${video.videoId}`}
              className="text-sm font-semibold text-neutral-100 line-clamp-2 wrap-break-word hover:text-emerald-400
                transition-colors"
            >
              {video.name}
            </Link>
          ) : (
            <span className="text-sm font-semibold text-neutral-400 line-clamp-2 wrap-break-word">{video.name}</span>
          )}

          {isReady ? (
            <p className="text-xs text-neutral-500 mt-1 line-clamp-1 wrap-break-word">
              {video.description || 'No description'}
            </p>
          ) : (
            <div className="flex items-center gap-2 mt-1.5">
              <span
                className={cn(
                  'text-[10px] uppercase tracking-wider font-bold px-1.5 py-0.5 rounded-sm',
                  isUploading && 'bg-emerald-500/10 text-emerald-400',
                  isFailed && 'bg-red-500/10 text-red-400',
                  video.videoUploadStatus === 'Canceled' && 'bg-neutral-500/20 text-neutral-400',
                )}
              >
                {video.videoUploadStatus}
              </span>
              {isUploading && <span className="text-xs text-neutral-400 animate-pulse">Processing...</span>}
            </div>
          )}
        </div>
      </div>

      <div className="w-28 flex justify-center shrink-0">
        <div className="flex items-center gap-1.5 text-xs text-neutral-300">
          {video.isPrivate ? <span>Private</span> : <span className="text-emerald-500">Public</span>}
        </div>
      </div>

      <div className="w-32 flex flex-col items-center justify-center shrink-0 text-center">
        <p className="text-xs text-neutral-200">{formatDate(video.createdAt)}</p>
        <p className="text-[10px] text-neutral-500 mt-0.5">Uploaded</p>
      </div>

      <div className="w-24 text-center text-xs text-neutral-300 shrink-0">
        {isReady ? formatViews(video.views) : '-'}
      </div>

      <div className="w-12 flex justify-end shrink-0 relative" ref={menuRef}>
        <button
          onClick={toggleMenu}
          className={cn(
            'p-1.5 -mr-1.5 -mt-1 rounded-full transition-all duration-300 cursor-pointer z-20',
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
              className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 transition-colors flex
                items-center gap-3 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
              onClick={handleDelete}
              disabled={isPending}
            >
              <span>{isPending ? 'Deleting...' : 'Delete video'}</span>
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
