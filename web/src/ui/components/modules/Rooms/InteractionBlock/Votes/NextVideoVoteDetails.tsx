'use client'

import { useEffect, useState } from 'react'

import Image from 'next/image'

import { useRoomQueueQuery } from '@/lib/hooks/api/room/useRoomQueueQuery'
import { useCastNextVideoVoteMutation } from '@/lib/hooks/api/vote/useCastNextVideoVote'
import { useCheckNextVideoVotingQuery } from '@/lib/hooks/api/vote/useCheckNextVideoVoting'
import { useCountdown } from '@/lib/hooks/useCountdown'
import { cn } from '@/lib/utils/general.utils'

type Props = {
  roomId: string
  onBack: () => void
}

export default function NextVideoVoteDetails({ roomId, onBack }: Props) {
  const { data: hasNextVideoData } = useCheckNextVideoVotingQuery(roomId)
  const isNextVideoActive = hasNextVideoData?.exists

  const { data: queueData, isLoading } = useRoomQueueQuery(roomId)
  const { mutate: castVote, isPending } = useCastNextVideoVoteMutation()

  const [myVote, setMyVote] = useState<string | null>(null)

  const { progress, secondsLeft, isExpired } = useCountdown(
    isNextVideoActive ? hasNextVideoData?.expiresAt : undefined,
    isNextVideoActive ? hasNextVideoData?.duration : undefined,
  )

  useEffect(() => {
    if (hasNextVideoData && !isNextVideoActive) {
      onBack()
    }
  }, [isNextVideoActive, hasNextVideoData, onBack])

  const queueItems = queueData?.pages.flatMap(page => page?.data?.items || []) || []
  const availableItems = queueItems.filter(item => !item.isActive)

  return (
    <div className="flex flex-col max-h-[60vh]">
      <h3 className="text-xl font-bold text-neutral-100 mb-2">Vote for Next Video</h3>

      <div className="mb-6 flex flex-col gap-2">
        <div className="flex items-center gap-2">
          <span
            className="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md bg-purple-500/20
              text-purple-400"
          >
            Active • {secondsLeft}s
          </span>
        </div>
        {isNextVideoActive && !isExpired && (
          <div className="h-1 w-full bg-neutral-800 rounded-full overflow-hidden">
            <div className="h-full bg-purple-500 transition-none" style={{ width: `${progress}%` }} />
          </div>
        )}
      </div>

      <div className="flex-1 overflow-y-auto custom-scrollbar pr-2 flex flex-col gap-2">
        {isLoading ? (
          <div className="text-neutral-500 text-center py-6 animate-pulse">Loading queue...</div>
        ) : availableItems.length === 0 ? (
          <div className="text-neutral-500 text-center py-6">No more videos in the queue.</div>
        ) : (
          availableItems.map(item => {
            const isMyChoice = myVote === item.id
            const hasVoted = !!myVote

            return (
              <button
                key={item.id}
                onClick={() => {
                  if (!hasVoted) {
                    castVote({ roomId, queueItemId: item.id })
                    setMyVote(item.id)
                  }
                }}
                disabled={hasVoted || isPending}
                className={cn(
                  `w-full text-left p-3 rounded-xl border transition-all flex items-center gap-3
                    disabled:cursor-not-allowed`,
                  isMyChoice
                    ? 'bg-purple-500/10 border-purple-500/50'
                    : hasVoted
                      ? 'bg-neutral-800/50 border-neutral-800 opacity-60'
                      : 'bg-[#141414] border-neutral-800 hover:border-neutral-600 hover:bg-neutral-800',
                )}
              >
                <div className="w-16 h-9 bg-neutral-800 rounded-md shrink-0 overflow-hidden relative">
                  {item.thumbnail && (
                    <Image src={item.thumbnail} width={50} height={50} className="object-cover w-full h-full" alt="" />
                  )}
                </div>
                <div className="flex flex-col flex-1 overflow-hidden">
                  <span className="text-sm font-medium text-neutral-200 truncate">
                    {item.title || 'Video from queue'}
                  </span>
                  {isMyChoice && <span className="text-[10px] text-purple-400 font-semibold mt-0.5">Your vote</span>}
                </div>
              </button>
            )
          })
        )}
      </div>

      <div className="mt-6 pt-4 border-t border-neutral-800">
        <button onClick={onBack} className="text-neutral-400 hover:text-neutral-200 font-medium">
          Back to list
        </button>
      </div>
    </div>
  )
}
