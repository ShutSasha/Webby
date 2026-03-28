import { useState } from 'react'

import Image from 'next/image'

import { togglePlaylistVideo } from '@/app/api/playlists'
import { extractServerMessage, serverLog } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  playlistId: string
  videoId: string
  image: string
  name: string
  count: number
  isVideoAdded: boolean
}

export default function PlaylistItem({ videoId, playlistId, image, name, count, isVideoAdded }: Props) {
  const addToast = useToastStore(state => state.addToast)
  const [loading, setLoading] = useState(false)
  const [isAdded, setIsAdded] = useState(isVideoAdded)

  const [localCount, setLocalCount] = useState(count)

  const toggleVideoInPlaylsit = async () => {
    if (loading) return

    try {
      setLoading(true)
      const response = await togglePlaylistVideo(playlistId, videoId)

      if (response.success) {
        const newIsAdded = !isAdded
        setIsAdded(newIsAdded)

        setLocalCount(prev => (newIsAdded ? prev + 1 : prev - 1))
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg ? msg : `Failed to update "${name}"`, 'error')
      }
    } catch (error) {
      serverLog('toggle video in playlist error', error, true)
      addToast('Something went wrong while toggle video in playlist', 'error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className={`flex items-center justify-between py-2 px-2.5 mr-1 transition-all duration-300 ease-in-out rounded-xl
        group ${loading ? 'cursor-default opacity-60' : 'cursor-pointer hover:bg-neutral-800/50'}`}
      onClick={toggleVideoInPlaylsit}
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
        {loading ? (
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
