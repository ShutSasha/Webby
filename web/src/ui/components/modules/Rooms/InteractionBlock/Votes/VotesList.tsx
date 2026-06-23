'use client'

import { useCheckNextVideoVotingQuery } from '@/lib/hooks/api/vote/useCheckNextVideoVoting'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { cn } from '@/lib/utils/general.utils'

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
        {isNextVideoActive && (
          <div
            onClick={() => onSelectVote('next-video-active', 'NEXT_VIDEO')}
            className="p-4 rounded-xl border transition-all cursor-pointer group bg-purple-500/10 border-purple-500/30
              hover:border-purple-500 shadow-[0_0_15px_rgba(168,85,247,0.1)]"
          >
            <div className="flex justify-between items-start mb-2">
              <span
                className="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md bg-purple-500/20
                  text-purple-400 flex items-center gap-2"
              >
                <span className="relative flex size-1.5">
                  <span
                    className="animate-ping absolute inline-flex h-full w-full rounded-full bg-purple-400 opacity-75"
                  />
                  <span className="relative inline-flex rounded-full size-1.5 bg-purple-500" />
                </span>
                Quick Vote
              </span>
            </div>
            <h4 className="font-medium text-foreground-tertiary group-hover:text-purple-400 transition-colors">
              Choose the next video!
            </h4>
          </div>
        )}

        {isLoading ? (
          <div className="text-foreground0 text-center py-10 animate-pulse">Loading votes...</div>
        ) : votes.length === 0 && !isNextVideoActive ? (
          <div className="text-foreground0 text-center py-10">No polls have been created yet.</div>
        ) : (
          votes.map(vote => {
            const isResolved = !!vote.rightChoice
            const isClosed = !isResolved && vote.isLocked
            const isActive = !isResolved && !vote.isLocked

            return (
              <div
                key={vote.id}
                onClick={() => onSelectVote(vote.id, 'RIGHT_CHOICE')}
                className={cn('p-4 rounded-xl border transition-all cursor-pointer group', {
                  'bg-neutral-800/30 border-neutral-800': isResolved,
                  'bg-neutral-800/80 border-neutral-700 hover:bg-neutral-700': isClosed,
                  'bg-neutral-800 border-emerald-500/30 hover:border-emerald-500': isActive,
                })}
              >
                <div className="flex justify-between items-start mb-2">
                  <span
                    className={cn('text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md', {
                      'bg-neutral-800 text-foreground0': isResolved,
                      'bg-neutral-700 text-foreground-muted': isClosed,
                      'bg-emerald-500/20 text-emerald-500': isActive,
                    })}
                  >
                    {isResolved ? 'Resolved' : isClosed ? 'Closed' : 'Active'}
                  </span>
                </div>
                <h4 className="font-medium text-foreground-tertiary group-hover:text-emerald-400 transition-colors">
                  {vote.voteText}
                </h4>
                <div className="flex justify-between items-center mt-2">
                  <p className="text-xs text-foreground0">{vote.choices.length} options</p>
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
