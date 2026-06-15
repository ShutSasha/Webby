'use client'

import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { cn } from '@/lib/utils/general.utils'

import { VoteViewState } from './RoomVotesModal'

type Props = {
  roomId: string
  isHost: boolean
  onViewChange: (view: VoteViewState) => void
  onSelectVote: (voteId: string) => void
}

export default function VotesList({ roomId, isHost, onViewChange, onSelectVote }: Props) {
  const { data: votes = [], isLoading } = useGetRoomVotesQuery(roomId)

  return (
    <div className="flex flex-col h-full max-h-[60vh]">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-bold text-neutral-100">Polls & Quizzes</h3>
        {isHost && (
          <button
            onClick={() => onViewChange('CREATE_RIGHT_CHOICE')}
            className="text-sm px-3 py-1.5 rounded-lg bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20
              font-medium transition-colors"
          >
            + Create New
          </button>
        )}
      </div>

      <div className="flex-1 overflow-y-auto custom-scrollbar flex flex-col gap-3 pr-2">
        {isLoading ? (
          <div className="text-neutral-500 text-center py-10 animate-pulse">Loading votes...</div>
        ) : votes.length === 0 ? (
          <div className="text-neutral-500 text-center py-10">No polls have been created yet.</div>
        ) : (
          votes.map(vote => {
            const isResolved = !!vote.rightChoice
            const isClosed = !isResolved && vote.isLocked
            const isActive = !isResolved && !vote.isLocked

            return (
              <div
                key={vote.id}
                onClick={() => onSelectVote(vote.id)}
                className={cn('p-4 rounded-xl border transition-all cursor-pointer group', {
                  'bg-neutral-800/30 border-neutral-800': isResolved,
                  'bg-neutral-800/80 border-neutral-700 hover:bg-neutral-700': isClosed,
                  'bg-neutral-800 border-emerald-500/30 hover:border-emerald-500': isActive,
                })}
              >
                <div className="flex justify-between items-start mb-2">
                  <span
                    className={cn('text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md', {
                      'bg-neutral-800 text-neutral-500': isResolved,
                      'bg-neutral-700 text-neutral-400': isClosed,
                      'bg-emerald-500/20 text-emerald-500': isActive,
                    })}
                  >
                    {isResolved ? 'Resolved' : isClosed ? 'Closed' : 'Active'}
                  </span>
                </div>
                <h4 className="font-medium text-neutral-200 group-hover:text-emerald-400 transition-colors">
                  {vote.voteText}
                </h4>
                <div className="flex justify-between items-center mt-2">
                  <p className="text-xs text-neutral-500">{vote.choices.length} options</p>
                  {vote.myVote && <p className="text-[10px] text-emerald-500/80 font-medium">Voted</p>}
                </div>
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}