'use client'

import { useMemo } from 'react'

import { useRoomQueueQuery } from '@/lib/hooks/api/room/useRoomQueueQuery'
import { useRoomStore } from '@/stores/room.store'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'

type Props = {
  roomId: string
}

export default function RoomPlayerContainer({ roomId }: Props) {
  const { data, isLoading } = useRoomQueueQuery(roomId)

  const optimisticPendingId = useRoomStore(state => state.optimisticPendingId)

  const activeVideoUrl = useMemo(() => {
    const items = data?.pages.flatMap(page => page?.data?.items || []) || []

    if (optimisticPendingId) {
      const optimisticVideo = items.find(v => v.id === optimisticPendingId)
      if (optimisticVideo?.videoUrl) return optimisticVideo.videoUrl
    }

    const activeVideo = items.find(v => v.isActive)
    if (activeVideo?.videoUrl) return activeVideo.videoUrl

    return ''
  }, [data, optimisticPendingId])

  if (isLoading && !activeVideoUrl) {
    return (
      <div className="aspect-video bg-neutral-900 w-full rounded-2xl flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
      </div>
    )
  }

  if (!activeVideoUrl) {
    return (
      <div className="aspect-video bg-neutral-900 w-full rounded-2xl flex flex-col items-center justify-center gap-2">
        <p className="text-neutral-500 font-medium">No video is currently playing</p>
        <p className="text-neutral-600 text-sm">Select a video from the queue to start</p>
      </div>
    )
  }

  return <CustomPlayer videoUrl={activeVideoUrl} isRoom roomId={roomId} />
}
