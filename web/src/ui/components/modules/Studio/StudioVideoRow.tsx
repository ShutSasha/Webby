'use client'

import { useEffect, useRef, useState } from 'react'

import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { DEFAULT_VIDEO_THUMBNAIL } from '@/lib/constants/url.constasts'
import { useCancelUploadVideo } from '@/lib/hooks/api/video/useCancelUploadVideo'
import { useCheckVideoUploadStatus } from '@/lib/hooks/api/video/useCheckVideoUploadStatus'
import { useDeleteVideo } from '@/lib/hooks/api/video/useDeleteVideo'
import { formatDate, formatVideoTime } from '@/lib/utils/date.utils'
import { cn, serverLog } from '@/lib/utils/general.utils'
import { formatViews } from '@/lib/utils/video.utils'
import { useVideoDraftStore } from '@/stores/video-draft.store'
import { Video } from '@/types/video.types'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  video: Video
}

export function StudioVideoRow({ video }: Props) {
  const { mutate, isPending } = useDeleteVideo()
  const { mutate: cancelUploadVideo, isPending: cancelUploadPending } = useCancelUploadVideo()
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const router = useRouter()
  const clearDraft = useVideoDraftStore(state => state.clearDraft)

  const isReady = video.videoUploadStatus === 'Ready'
  const isUploading = video.videoUploadStatus === 'Uploading'
  const isFailed = video.videoUploadStatus === 'Failed'

  useCheckVideoUploadStatus(video.videoId, isUploading)

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
          clearDraft()
          setIsMenuOpen(false)
        },
      },
    )
  }

  const handleEditRedirect = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsMenuOpen(false)
    router.push(`/videos/edit/${video.videoId}`)
  }

  const handleCancelUploadVideo = async () => {
    try {
      cancelUploadVideo(video.videoId)
      setIsMenuOpen(false)
    } catch (error) {
      serverLog('FAILED_CANCEL_UPLOAD_VIDEO', error, true)
    }
  }

  return (
    <div
      className={cn(
        'flex items-center border-b border-border py-3 px-4 transition-colors group',
        isReady ? 'hover:bg-background/40' : 'bg-surface/50',
      )}
    >
      <div className="flex-1 flex gap-4 min-w-[300px]">
        <div className="relative w-32 aspect-video bg-background rounded-lg overflow-hidden shrink-0">
          <Image
            src={video.previewUrl || DEFAULT_VIDEO_THUMBNAIL}
            alt="Video thumbnail"
            fill
            className={cn('object-cover transition-all', !isReady && 'opacity-40 grayscale-50')}
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />

          {isReady && (
            <div
              className="absolute bottom-1 right-1 bg-black/80 backdrop-blur-sm px-1 py-0.5 rounded text-[10px]
                font-medium text-neutral-50"
            >
              {video.duration ? formatVideoTime(video.duration) : '0:00'}
            </div>
          )}

          {!isReady && (
            <div className="absolute inset-0 flex flex-col items-center justify-center bg-surface/40
              backdrop-blur-[2px]">
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
              className="text-sm font-semibold text-foreground-secondary line-clamp-2 wrap-break-word
                hover:text-emerald-400 transition-colors"
            >
              {video.name}
            </Link>
          ) : (
            <span className="text-sm font-semibold text-foreground-muted line-clamp-2 wrap-break-word">
              {video.name}
            </span>
          )}

          {isReady ? (
            <p className="text-xs text-foreground-faint mt-1 line-clamp-1 wrap-break-word">
              {video.description || 'No description'}
            </p>
          ) : (
            <div className="flex items-center gap-2 mt-1.5">
              <span
                className={cn(
                  'text-[10px] uppercase tracking-wider font-bold px-1.5 py-0.5 rounded-sm',
                  isUploading && 'bg-emerald-500/10 text-emerald-400',
                  isFailed && 'bg-red-500/10 text-red-400',
                  video.videoUploadStatus === 'Canceled' && 'bg-surface-faint/20 text-foreground-muted',
                )}
              >
                {video.videoUploadStatus}
              </span>
              {isUploading && <span className="text-xs text-foreground-muted animate-pulse">Processing...</span>}
            </div>
          )}
        </div>
      </div>

      <div className="w-28 flex justify-center shrink-0">
        <div
          className={cn(
            'flex items-center gap-1.5 text-xs font-medium px-2 py-1 rounded-md',
            video.isPublished ? 'text-foreground-subtle' : 'bg-surface-faint/10 text-foreground-muted',
          )}
        >
          {video.isPublished ? 'Published' : 'Draft'}
        </div>
      </div>

      <div className="w-28 flex justify-center shrink-0">
        <div className="flex items-center gap-1.5 text-xs text-foreground-subtle">
          {video.isPrivate ? <span>Private</span> : <span className="text-emerald-500">Public</span>}
        </div>
      </div>

      <div className="w-32 flex flex-col items-center justify-center shrink-0 text-center">
        <p className="text-xs text-foreground-tertiary">{formatDate(video.createdAt)}</p>
        <p className="text-[10px] text-foreground-faint mt-0.5">Uploaded</p>
      </div>

      <div className="w-24 text-center text-xs text-foreground-subtle shrink-0">
        {isReady ? formatViews(video.views) : '-'}
      </div>

      <div className="w-12 flex justify-end shrink-0 relative" ref={menuRef}>
        {video.videoUploadStatus !== 'Canceled' && (
          <button
            onClick={toggleMenu}
            className={cn(
              'p-1.5 -mr-1.5 -mt-1 rounded-full transition-all duration-300 cursor-pointer z-20',
              'hover:bg-surface-faint/20 active:bg-surface-faint/40',
              isMenuOpen
                ? 'bg-surface-faint/20 text-foreground-subtle'
                : 'text-foreground-muted opacity-0 group-hover:opacity-100 md:opacity-100',
            )}
          >
            <MoreVertical className="size-5 text-foreground-subtle" />
          </button>
        )}

        {isMenuOpen && (
          <div
            className="absolute right-0 top-full mt-2 w-48 bg-background border border-neutral-700/60 shadow-xl
              shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
            onClick={e => e.preventDefault()}
          >
            {isUploading && (
              <button
                className="w-full text-left px-4 py-2 text-sm text-foreground-tertiary hover:bg-surface-tertiary/50
                  transition-colors flex items-center gap-3 cursor-pointer"
                onClick={handleCancelUploadVideo}
              >
                <span>{cancelUploadPending ? 'Canceling upload...' : 'Cancel upload'}</span>
              </button>
            )}
            <button
              className="w-full text-left px-4 py-2 text-sm text-foreground-tertiary hover:bg-surface-tertiary/50
                transition-colors flex items-center gap-3 cursor-pointer"
              onClick={handleEditRedirect}
            >
              <span>Edit</span>
            </button>
            {!isUploading && (
              <button
                className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 transition-colors flex
                  items-center gap-3 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                onClick={handleDelete}
                disabled={isPending}
              >
                <span>{isPending ? 'Deleting...' : 'Delete video'}</span>
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
