'use client'

import { useCheckNextVideoVotingQuery } from '@/lib/hooks/api/vote/useCheckNextVideoVoting'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'

import NextVideoVoteCard from './NextVideoVoteCard'
import RightChoiceVoteCard from './RightChoiceVoteCard'
import { VoteViewState } from './RoomVotesModal'

type Props = {
  roomId: string
  isHost: boolean
  onViewChange: (view: VoteViewState) => void
  onSelectVote: (voteId: string, type: 'RIGHT_CHOICE' | 'NEXT_VIDEO') => void
}

export default function VotesList({ roomId, isHost, onViewChange, onSelectVote }: Props) {
  const { data: votes = [], isLoading } = useGetRoomVotesQuery(roomId)

  const { data: hasNextVideoData } = useCheckNextVideoVotingQuery(roomId)

  const isNextVideoActive = hasNextVideoData?.exists

  return (
    <div className="flex flex-col h-full max-h-[60vh]">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-bold text-foreground-secondary">Polls & Quizzes</h3>
        {isHost && (
          <button
            onClick={() => onViewChange('SELECT_TYPE')}
            className="text-sm px-3 py-1.5 rounded-lg bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20
              font-medium transition-colors"
          >
            + Create New
          </button>
        )}
      </div>

      <div className="flex-1 overflow-y-auto custom-scrollbar flex flex-col gap-3 pr-2">
        {isNextVideoActive && <NextVideoVoteCard onSelectVote={onSelectVote} />}

        {isLoading ? (
          <div className="text-foreground-faint text-center py-10 animate-pulse">Loading votes...</div>
        ) : votes.length === 0 && !isNextVideoActive ? (
          <div className="text-foreground-faint text-center py-10">No polls have been created yet.</div>
        ) : (
          votes.map(vote => {
            const isResolved = !!vote.rightChoice
            const isClosed = !isResolved && vote.isLocked
            const isActive = !isResolved && !vote.isLocked

            return (
              <RightChoiceVoteCard
                key={vote.id}
                vote={vote}
                onSelectVote={onSelectVote}
                isActive={isActive}
                isClosed={isClosed}
                isResolved={isResolved}
              />
            )
          })
        )}
      </div>
    </div>
  )
}
