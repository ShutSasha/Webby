type Props = {
  onSelectVote: (voteId: string, type: 'RIGHT_CHOICE' | 'NEXT_VIDEO') => void
}

export default function NextVideoVoteCard({ onSelectVote }: Props) {
  return (
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
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-purple-400 opacity-75" />
            <span className="relative inline-flex rounded-full size-1.5 bg-purple-500" />
          </span>
          Quick Vote
        </span>
      </div>
      <h4 className="font-medium text-foreground-tertiary group-hover:text-purple-400 transition-colors">
        Choose the next video!
      </h4>
    </div>
  )
}
