'use client'

import { useEffect, useState } from 'react'

import Image from 'next/image'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import { activateQueueItemAction } from '@/lib/actions/room.actions'
import { DEFAULT_VIDEO_THUMBNAIL } from '@/lib/constants/url.constamts'
import { useRemoveQueueItemMutation } from '@/lib/hooks/api/room/useRemoveQueueItem'
import { cn } from '@/lib/utils/general.utils'
import { useRoomStore } from '@/stores/room.store'

interface Video {
  id: string
  title: string
  thumbnail: string
  isActive?: boolean
  isFolder?: boolean
  children?: Video[]
}

type Props = {
  roomId: string | undefined
  video: Video
  isChild?: boolean
}

export default function PlaylistItem({ roomId, video, isChild = false }: Props) {
  const [isOpen, setOpen] = useState(false)

  const optimisticPendingId = useRoomStore(state => state.optimisticPendingId)
  const setOptimisticPendingId = useRoomStore(state => state.setOptimisticPendingId)

  const { mutate: removeQueueItem, isPending: isRemoving } = useRemoveQueueItemMutation()

  useEffect(() => {
    if (optimisticPendingId === video.id && video.isActive) {
      setOptimisticPendingId(null)
    }
  }, [video.isActive, video.id, optimisticPendingId, setOptimisticPendingId])

  const title = video.title || 'Untitled Video'
  const thumbnail = video.thumbnail || DEFAULT_VIDEO_THUMBNAIL

  const handleActivate = async () => {
    if (!roomId || !video.id || video.isActive || video.isFolder || isRemoving) return
    setOptimisticPendingId(video.id)
    const res = await activateQueueItemAction(roomId, video.id)

    if (!res.success) {
      if (useRoomStore.getState().optimisticPendingId === video.id) {
        setOptimisticPendingId(null)
      }
    }
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation()
    if (!roomId || !video.id || isRemoving) return

    removeQueueItem({ roomId, itemId: video.id })
  }

  const isPending = optimisticPendingId === video.id
  const isVisuallyActive = isPending || (video.isActive && !optimisticPendingId)

  return (
    <div className={cn('flex flex-col gap-1', isChild && 'ml-4 pl-2 border-l border-neutral-800')}>
      <div
        onClick={handleActivate}
        className={cn(
          `flex items-center justify-between p-2 rounded-xl bg-neutral-900/50 hover:bg-neutral-900 transition-all
          duration-300 border-b-2 border-transparent group cursor-pointer`,
          isVisuallyActive && !isPending && 'border-emerald-500 bg-neutral-900',
          isPending && 'border-amber-500 bg-neutral-900',
          video.isFolder && 'hover:bg-neutral-800/40',
          isRemoving && 'opacity-50 pointer-events-none',
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
              isVisuallyActive ? 'opacity-100' : 'opacity-60 group-hover:opacity-100',
            )}
          />

          <span
            className={cn(
              'text-sm font-medium truncate pr-1 transition-colors',
              isVisuallyActive ? 'text-neutral-300' : 'text-neutral-500 group-hover:text-neutral-300',
            )}
            title={title}
          >
            {title}
          </span>
        </div>

        <div className="flex items-center gap-1 shrink-0" onClick={e => e.stopPropagation()}>
          <button
            onClick={handleDelete}
            disabled={isRemoving}
            className={cn(
              'group/trash p-1.5 text-neutral-500 transition-colors cursor-pointer rounded-full',
              !isRemoving && 'hover:text-red-500 hover:bg-red-500/10',
            )}
          >
            <TrashIcon
              className={cn(
                'size-4 stroke-[1.5px] transition-colors',
                isVisuallyActive ? 'text-neutral-300' : 'text-neutral-500',
                !isRemoving && 'group-hover/trash:text-red-500',
                isRemoving && 'animate-pulse',
              )}
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
            <PlaylistItem key={child.id} video={child} roomId={roomId} isChild />
          ))}
        </div>
      )}
    </div>
  )
}
