import { useParams } from 'next/navigation'

import PollIcon from '@/assets/icons/Room/poll-icon.svg'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { useRoomStore } from '@/stores/room.store'

import InteractionButton from './InteractionButton'

export default function Poll() {
  const params = useParams()
  const roomId = params?.id as string | undefined
  const setVotesModalOpen = useRoomStore(state => state.setVotesModalOpen)

  const { data: votesResponse } = useGetRoomVotesQuery(roomId || '')
  const votes = votesResponse || []
  const hasActiveVote = votes.some(v => !v.isLocked)

  return (
    <InteractionButton
      onClick={() => setVotesModalOpen(true)}
      icon={
        <PollIcon
          className="size-5 text-foreground-subtle dark:group-hover:text-foreground-inverse-subtle transition-colors
            duration-300 ease-in-out"
        />
      }
      text="Votes"
    >
      {hasActiveVote && (
        <span className="absolute -top-1 -right-1 flex size-3">
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span className="relative inline-flex rounded-full size-3 bg-emerald-500 border-2 border-surface"></span>
        </span>
      )}
    </InteractionButton>
  )
}
