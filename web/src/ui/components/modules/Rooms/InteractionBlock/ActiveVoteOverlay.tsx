'use client'

import { useEffect, useState } from 'react'

import CloseIcon from '@/assets/icons/shared/circle-xmark.svg'
import { useCheckNextVideoVotingQuery } from '@/lib/hooks/api/vote/useCheckNextVideoVoting'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { useRoomStore } from '@/stores/room.store'

type Props = {
  roomId: string
}

export default function ActiveVoteOverlay({ roomId }: Props) {
  const setVotesModalOpen = useRoomStore(state => state.setVotesModalOpen)

  const { data: votesData } = useGetRoomVotesQuery(roomId)
  const votes = votesData || []
  const activeRightChoiceVote = votes.find(v => !v.isLocked)

  const { data: nextVideoData } = useCheckNextVideoVotingQuery(roomId)

  const isNextVideoActive = nextVideoData?.exists

  const [dismissedRightChoiceId, setDismissedRightChoiceId] = useState<string | null>(null)
  const [isNextVideoDismissed, setIsNextVideoDismissed] = useState(false)


  useEffect(() => {
    if (!isNextVideoActive) {
      setIsNextVideoDismissed(false)
    }
  }, [isNextVideoActive])

  if (isNextVideoActive && !isNextVideoDismissed) {
    return (
      <div
        onClick={() => setVotesModalOpen(true)}
        className="absolute top-2 left-2 right-2 z-20 bg-neutral-800/95 backdrop-blur-sm border border-purple-500/30
          rounded-xl p-3 shadow-[0_0_20px_rgba(168,85,247,0.15)] cursor-pointer group hover:bg-neutral-800
          transition-colors animate-in slide-in-from-top-4 fade-in duration-300"
      >
        <button
          onClick={e => {
            e.stopPropagation()
            setIsNextVideoDismissed(true)
          }}
          className="absolute top-1/2 -translate-y-1/2 right-1.5 p-1 text-neutral-500 hover:text-neutral-300
            transition-colors rounded-full hover:bg-neutral-700/50"
        >
          <CloseIcon className="size-4.5" />
        </button>

        <div className="flex items-center justify-between gap-3 pr-6">
          <div className="flex flex-col overflow-hidden">
            <div className="flex items-center gap-2 mb-1">
              <span className="relative flex size-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-purple-400 opacity-75" />
                <span className="relative inline-flex rounded-full size-2 bg-purple-500" />
              </span>
              <span className="text-[10px] font-bold uppercase tracking-wider text-purple-400">15s Quick Vote</span>
            </div>
            <p className="text-sm font-medium text-neutral-200 truncate">Choose the next video!</p>
          </div>
          <button
            className="shrink-0 bg-purple-500 text-white text-xs font-semibold px-3 py-1.5 rounded-lg
              group-hover:bg-purple-600 transition-colors"
          >
            Vote
          </button>
        </div>
      </div>
    )
  }

  if (activeRightChoiceVote && activeRightChoiceVote.id !== dismissedRightChoiceId) {
    return (
      <div
        onClick={() => setVotesModalOpen(true)}
        className="absolute top-2 left-2 right-2 z-20 bg-neutral-800/95 backdrop-blur-sm border border-emerald-500/30
          rounded-xl p-3 shadow-lg shadow-black/50 cursor-pointer group hover:bg-neutral-800 transition-colors
          animate-in slide-in-from-top-4 fade-in duration-300"
      >
        <button
          onClick={e => {
            e.stopPropagation()
            setDismissedRightChoiceId(activeRightChoiceVote.id)
          }}
          className="absolute top-1/2 -translate-y-1/2 right-1.5 p-1 text-neutral-500 hover:text-neutral-300
            transition-colors rounded-full hover:bg-neutral-700/50"
        >
          <CloseIcon className="size-4.5" />
        </button>

        <div className="flex items-center justify-between gap-3 pr-6">
          <div className="flex flex-col overflow-hidden">
            <div className="flex items-center gap-2 mb-1">
              <span className="relative flex size-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
                <span className="relative inline-flex rounded-full size-2 bg-emerald-500" />
              </span>
              <span className="text-[10px] font-bold uppercase tracking-wider text-emerald-500">Active Poll</span>
            </div>
            <p className="text-sm font-medium text-neutral-200 truncate">{activeRightChoiceVote.voteText}</p>
          </div>
          <button
            className="shrink-0 bg-emerald-500 text-neutral-950 text-xs font-semibold px-3 py-1.5 rounded-lg
              group-hover:bg-emerald-400 transition-colors"
          >
            Vote
          </button>
        </div>
      </div>
    )
  }

  return null
}
