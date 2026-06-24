'use client'

import { useCheckNextVideoVotingQuery } from '@/lib/hooks/api/vote/useCheckNextVideoVoting'
import { useCreateNextVideoVoteMutation } from '@/lib/hooks/api/vote/useCreateNextVideoVote'

type Props = {
  roomId: string
  onBack: () => void
  onSuccess: () => void
}

export default function CreateNextVideoVote({ roomId, onBack, onSuccess }: Props) {
  const { data: hasNextVideoData } = useCheckNextVideoVotingQuery(roomId)
  const isNextVideoActive = hasNextVideoData?.exists

  const { mutate: createNextVideo, isPending } = useCreateNextVideoVoteMutation()

  const handleCreate = () => {
    createNextVideo(roomId, { onSuccess })
  }

  return (
    <div className="flex flex-col">
      <h3 className="text-xl font-bold text-foreground-secondary mb-4">Next Video Voting</h3>
      <p className="text-foreground-muted mb-6">
        This will immediately start a 20-second voting session. Users will be able to select any video currently in the
        queue. The video with the most votes will automatically play next.
      </p>

      <div className="flex items-center gap-3">
        <button
          onClick={onBack}
          className="px-5 py-2.5 text-foreground-muted hover:bg-background rounded-xl font-medium transition-colors"
        >
          Back
        </button>
        <button
          onClick={handleCreate}
          disabled={isPending || isNextVideoActive}
          className="flex-1 py-2.5 bg-purple-500 text-foreground-strong font-semibold rounded-xl hover:bg-purple-600
            disabled:opacity-50 transition-colors"
        >
          {isPending ? 'Starting...' : isNextVideoActive ? 'Already Active' : 'Start 20s Vote'}
        </button>
      </div>
    </div>
  )
}
