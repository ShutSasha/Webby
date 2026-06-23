'use client'

import { VoteViewState } from './RoomVotesModal'

type Props = {
  onSelect: (view: VoteViewState) => void
  onBack: () => void
}

export default function SelectVoteType({ onSelect, onBack }: Props) {
  return (
    <div className="flex flex-col">
      <h3 className="text-xl font-bold text-foreground-secondary mb-6">Select Voting Type</h3>

      <div className="flex flex-col gap-3">
        <button
          onClick={() => onSelect('CREATE_NEXT_VIDEO')}
          className="flex flex-col items-start p-4 bg-[#141414] border border-neutral-800 hover:border-purple-500
            hover:bg-neutral-800 rounded-xl transition-all group text-left"
        >
          <span className="font-semibold text-foreground-tertiary group-hover:text-purple-400 transition-colors">
            Next Video Vote
          </span>
          <span className="text-sm text-foreground-faint mt-1">
            Start a quick 20-second vote to let the room choose the next video from the queue.
          </span>
        </button>

        <button
          onClick={() => onSelect('CREATE_RIGHT_CHOICE')}
          className="flex flex-col items-start p-4 bg-[#141414] border border-neutral-800 hover:border-emerald-500
            hover:bg-neutral-800 rounded-xl transition-all group text-left"
        >
          <span className="font-semibold text-foreground-tertiary group-hover:text-emerald-400 transition-colors">
            Quiz / Custom Poll
          </span>
          <span className="text-sm text-foreground-faint mt-1">
            Create a custom question with multiple choices and determine the correct answer later.
          </span>
        </button>
      </div>

      <button
        onClick={onBack}
        className="mt-6 text-foreground-muted hover:text-foreground-tertiary font-medium self-start"
      >
        Cancel
      </button>
    </div>
  )
}
