import { RoomVote } from '@/lib/actions/vote.actions'
import { cn } from '@/lib/utils/general.utils'

type Props = {
  isActive: boolean
  isClosed: boolean
  isResolved: boolean
  vote: RoomVote
  onSelectVote: (voteId: string, type: 'RIGHT_CHOICE' | 'NEXT_VIDEO') => void
}

export default function RightChoiceVoteCard({ isActive, isClosed, isResolved, vote, onSelectVote }: Props) {
  return (
    <div
      key={vote.id}
      onClick={() => onSelectVote(vote.id, 'RIGHT_CHOICE')}
      className={cn('p-4 rounded-xl border transition-all cursor-pointer group', {
        'bg-background/30 border-border': isResolved,
        'bg-background/80 border-neutral-700 hover:bg-neutral-300/40 dark:hover:bg-neutral-700/30': isClosed,
        'bg-background border-emerald-500/30 hover:border-emerald-500': isActive,
      })}
    >
      <div className="flex justify-between items-start mb-2">
        <span
          className={cn('text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md', {
            'bg-background text-foreground-faint': isResolved,
            'bg-surface-tertiary text-foreground-muted': isClosed,
            'bg-emerald-500/20 text-emerald-500': isActive,
          })}
        >
          {isResolved ? 'Resolved' : isClosed ? 'Closed' : 'Active'}
        </span>
      </div>
      <h4 className="font-medium text-foreground-tertiary transition-colors">{vote.voteText}</h4>
      <div className="flex justify-between items-center mt-2">
        <p className="text-xs text-foreground-faint">{vote.choices.length} options</p>
        {vote.myVote && <p className="text-[10px] text-emerald-500/80 font-medium">Voted</p>}
      </div>
    </div>
  )
}
