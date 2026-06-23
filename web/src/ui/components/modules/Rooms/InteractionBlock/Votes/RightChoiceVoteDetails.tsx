'use client'

import { useState } from 'react'

import { useQueryClient } from '@tanstack/react-query'

import { useCastVoteMutation } from '@/lib/hooks/api/vote/useCastVote'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { useResolveVoteMutation } from '@/lib/hooks/api/vote/useResolveVote'
import { useCountdown } from '@/lib/hooks/useCountdown'
import { cn } from '@/lib/utils/general.utils'

type Props = {
  roomId: string
  voteId: string
  isHost: boolean
  onBack: () => void
}

export default function RightChoiceVoteDetails({ roomId, voteId, isHost, onBack }: Props) {
  const queryClient = useQueryClient()
  const [hostSelectedCorrect, setHostSelectedCorrect] = useState<string>('')

  const { data: votes = [] } = useGetRoomVotesQuery(roomId)
  const currentVote = votes.find(v => v.id === voteId)

  const { mutate: castVote, isPending: isCasting } = useCastVoteMutation()
  const { mutate: resolveVote, isPending: isResolving } = useResolveVoteMutation()

  const { progress, secondsLeft, isExpired } = useCountdown(
    currentVote && !currentVote.isLocked ? currentVote.expiresAt : undefined,
    currentVote && !currentVote.isLocked ? currentVote.duration : undefined,
  )

  if (!currentVote) return null

  const isResolvedDetails = !!currentVote.rightChoice
  const isClosedDetails = !isResolvedDetails && currentVote.isLocked
  const isActiveDetails = !isResolvedDetails && !currentVote.isLocked

  const handleCastVote = (choice: string) => {
    castVote({ roomId, voteId: currentVote.id, choice })

    queryClient.setQueryData(['room-votes', roomId], (oldData: any) => {
      if (!Array.isArray(oldData)) return oldData

      return oldData.map((vote: any) => {
        if (vote.id === currentVote.id) {
          return { ...vote, myVote: choice }
        }
        return vote
      })
    })
  }

  const handleResolve = () => {
    if (hostSelectedCorrect) {
      resolveVote({ roomId, voteId: currentVote.id, rightChoice: hostSelectedCorrect })
    }
  }

  return (
    <div className="flex flex-col max-h-[60vh]">
      <h3 className="text-xl font-bold text-foreground-secondary mb-2">{currentVote.voteText}</h3>

      <div className="mb-6 flex flex-col gap-2">
        <div className="flex items-center gap-2">
          <span
            className={cn('text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md', {
              'bg-neutral-800 text-foreground-faint': isResolvedDetails,
              'bg-neutral-700 text-foreground-muted': isClosedDetails,
              'bg-emerald-500/20 text-emerald-500': isActiveDetails,
            })}
          >
            {isResolvedDetails ? 'Resolved' : isClosedDetails ? 'Closed' : `Active • ${secondsLeft}s`}
          </span>
        </div>
        {isActiveDetails && !isExpired && (
          <div className="h-1 w-full bg-neutral-800 rounded-full overflow-hidden">
            <div className="h-full bg-emerald-500 transition-none" style={{ width: `${progress}%` }} />
          </div>
        )}
      </div>

      <div className="flex-1 overflow-y-auto custom-scrollbar pr-2 flex flex-col gap-2">
        {currentVote.choices.map((choice, index) => {
          const isWinner = currentVote.rightChoice === choice
          const isMyChoice = currentVote.myVote === choice
          const hasVoted = !!currentVote.myVote

          const getOptionClasses = () => {
            if (isWinner) return 'bg-emerald-500/10 border-emerald-500 text-emerald-400 font-semibold'
            if (isMyChoice && !isResolvedDetails) return 'bg-neutral-800 border-emerald-500/50 text-emerald-500'
            if (isResolvedDetails) return 'bg-neutral-900 border-neutral-800 text-foreground-disabled'
            if (currentVote.isLocked || hasVoted) return 'bg-neutral-800/50 border-neutral-800 text-foreground-faint'

            return 'bg-[#141414] border-neutral-800 text-foreground-tertiary hover:border-neutral-600 hover:bg-neutral-800'
          }

          return (
            <button
              key={index}
              onClick={() => {
                if (!currentVote.isLocked && !hasVoted) {
                  handleCastVote(choice)
                }
              }}
              disabled={currentVote.isLocked || isCasting || hasVoted}
              className={cn(
                'w-full text-left px-4 py-3 rounded-xl border transition-all cursor-pointer disabled:cursor-not-allowed',
                getOptionClasses(),
              )}
            >
              <div className="flex justify-between items-center">
                <span>{choice}</span>
                {isMyChoice && <span className="text-xs text-emerald-500 font-medium">Your vote</span>}
              </div>
            </button>
          )
        })}
      </div>

      {isHost && isClosedDetails && (
        <div className="mt-6 p-4 bg-emerald-500/10 border border-emerald-500/20 rounded-xl flex flex-col gap-3">
          <p className="text-sm text-emerald-400 font-medium">Select correct answer to publish results:</p>
          <select
            value={hostSelectedCorrect}
            onChange={e => setHostSelectedCorrect(e.target.value)}
            className="w-full bg-[#141414] border border-neutral-800 text-foreground-tertiary rounded-lg px-3 py-2
              outline-none"
          >
            <option value="" disabled hidden>
              Select option...
            </option>
            {currentVote.choices.map((choice, idx) => (
              <option key={idx} value={choice}>
                {choice}
              </option>
            ))}
          </select>
          <button
            onClick={handleResolve}
            disabled={!hostSelectedCorrect || isResolving}
            className="mt-2 w-full py-2 bg-emerald-500 text-foreground-inverse font-semibold rounded-lg
              hover:bg-emerald-400 disabled:opacity-50 transition-colors"
          >
            {isResolving ? 'Resolving...' : 'Publish Results'}
          </button>
        </div>
      )}

      {!isHost && isClosedDetails && (
        <div className="mt-4 p-3 bg-neutral-800/50 rounded-xl border border-neutral-800">
          <p className="text-sm text-foreground-muted text-center">
            Voting is closed. Waiting for host to publish results...
          </p>
        </div>
      )}

      <div className="mt-6 pt-4 border-t border-neutral-800">
        <button onClick={onBack} className="text-foreground-muted hover:text-foreground-tertiary font-medium">
          Back to list
        </button>
      </div>
    </div>
  )
}
