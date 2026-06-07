import { useEffect, useRef, useState } from 'react'

import Image from 'next/image'

import { useTogglePlaylistMediaMutation } from '@/lib/hooks/api/playlist/useTogglePlaylistMedia'
import { useToastStore } from '@/stores/toast-store'
import { MediaType } from '@/types/general.types'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  playlistId: string
  videoId: string
  image: string
  name: string
  count: number
  isVideoAdded: boolean
  userId: string
  mediaType: MediaType | null
}

export default function PlaylistItem({
  videoId,
  playlistId,
  image,
  name,
  count,
  isVideoAdded,
  userId,
  mediaType,
}: Props) {
  const addToast = useToastStore(state => state.addToast)

  const [isAdded, setIsAdded] = useState(isVideoAdded)
  const [localCount, setLocalCount] = useState(count)

  useEffect(() => {
    setIsAdded(isVideoAdded)
    setLocalCount(count)
  }, [isVideoAdded, count])

  const { mutate: toggleMedia, isPending } = useTogglePlaylistMediaMutation(userId)

  const handleToggle = () => {
    if (isPending) return

    const newIsAdded = !isAdded
    setIsAdded(newIsAdded)
    setLocalCount(prev => (newIsAdded ? prev + 1 : prev - 1))

    toggleMedia(
      {
        playlistId,
        mediaId: videoId,
        mediaType: mediaType || 'Video',
      },
      {
        onError: () => {
          setIsAdded(!newIsAdded)
          setLocalCount(prev => (!newIsAdded ? prev + 1 : prev - 1))

          addToast(`Failed to update "${name}"`, 'error')
        },
      },
    )
  }

  return (
    <div
      className={`flex items-center justify-between py-2 px-2.5 mr-1 transition-all duration-300 ease-in-out rounded-xl
        group ${isPending ? 'cursor-default opacity-60' : 'cursor-pointer hover:bg-neutral-800/50'}`}
      onClick={handleToggle}
    >
      <div className="flex gap-4">
        <Image
          src={image}
          width={100}
          height={100}
          alt=""
          className="aspect-square size-10 object-cover rounded-lg"
          loading="lazy"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
        />
        <div className="flex flex-col">
          <p className="text-sm text-neutral-300 line-clamp-1" title={name}>
            {name}
          </p>
          <p className="text-sm text-neutral-500">{localCount === 1 ? '1 video' : `${localCount} videos`}</p>
        </div>
      </div>
      <div className="flex items-center justify-center size-6 shrink-0 ml-3">
        {isPending ? (
          <div className="size-4 border-2 border-emerald-500/30 border-t-emerald-500 rounded-full animate-spin" />
        ) : isAdded ? (
          <div
            className="size-5 bg-emerald-500 rounded-full flex items-center justify-center
              shadow-[0_0_8px_rgba(16,185,129,0.3)]"
          >
            <svg
              className="w-3.5 h-3.5 text-neutral-900"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={3}
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        ) : (
          <div
            className="size-5 border-2 border-neutral-600 rounded-full group-hover:border-neutral-400 transition-colors"
          />
        )}
      </div>
    </div>
  )
}
