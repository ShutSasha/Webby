'use client'

import { useState } from 'react'

import { useCastVoteMutation } from '@/lib/hooks/api/vote/useCastVote'
import { useCreateRightChoiceVoteMutation } from '@/lib/hooks/api/vote/useCreateRightChoiceVote'
import { useGetRoomVotesQuery } from '@/lib/hooks/api/vote/useGetRoomVotes'
import { useResolveVoteMutation } from '@/lib/hooks/api/vote/useResolveVote'
import { useRoomStore } from '@/stores/room.store'
import Modal from '@/ui/components/shared/Modal'

type Props = {
  roomId: string
  isHost: boolean
}

type ViewState = 'LIST' | 'CREATE' | 'DETAILS'

export default function RoomVotesModal({ roomId, isHost }: Props) {
  const isVotesModalOpen = useRoomStore(state => state.isVotesModalOpen)
  const setVotesModalOpen = useRoomStore(state => state.setVotesModalOpen)

  const [view, setView] = useState<ViewState>('LIST')
  const [selectedVoteId, setSelectedVoteId] = useState<string | null>(null)

  const [voteText, setVoteText] = useState('')
  const [duration, setDuration] = useState(120)
  const [choices, setChoices] = useState(['', ''])

  const [hostSelectedCorrect, setHostSelectedCorrect] = useState<string | null>(null)

  const { data, isLoading } = useGetRoomVotesQuery(roomId)
  const votes = data || []

  const { mutate: createVote, isPending: isCreating } = useCreateRightChoiceVoteMutation()
  const { mutate: castVote, isPending: isCasting } = useCastVoteMutation()
  const { mutate: resolveVote, isPending: isResolving } = useResolveVoteMutation()

  const handleClose = () => {
    setVotesModalOpen(false)
    setTimeout(() => {
      setView('LIST')
      setSelectedVoteId(null)

      setVoteText('')
      setDuration(120)
      setChoices(['', ''])
    }, 300)
  }

  const handleCreateSubmit = () => {
    const validChoices = choices.map(c => c.trim()).filter(Boolean)
    if (!voteText.trim() || validChoices.length < 2) return

    createVote(
      { roomId, voteText: voteText.trim(), duration, choices: validChoices },
      {
        onSuccess: () => {
          setView('LIST')
          setVoteText('')
          setChoices(['', ''])
        },
      },
    )
  }

  const currentVote = votes.find(v => v.id === selectedVoteId)

  return (
    <Modal isOpen={isVotesModalOpen} onClose={handleClose} modalClasses="max-w-[480px]">
      {view === 'LIST' && (
        <div className="flex flex-col h-full max-h-[60vh]">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-bold text-neutral-100">Polls & Quizzes</h3>
            {isHost && (
              <button
                onClick={() => setView('CREATE')}
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
              votes.map(vote => (
                <div
                  key={vote.id}
                  onClick={() => {
                    setSelectedVoteId(vote.id)
                    setView('DETAILS')
                  }}
                  className={`p-4 rounded-xl border transition-all cursor-pointer group ${
                    vote.rightChoice
                      ? 'bg-neutral-800/30 border-neutral-800'
                      : vote.isLocked
                        ? 'bg-neutral-800/80 border-neutral-700 hover:bg-neutral-700'
                        : 'bg-neutral-800 border-emerald-500/30 hover:border-emerald-500'
                    }`}
                >
                  <div className="flex justify-between items-start mb-2">
                    <span
                      className={`text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md ${
                        vote.rightChoice
                          ? 'bg-neutral-800 text-neutral-500'
                          : vote.isLocked
                            ? 'bg-neutral-700 text-neutral-400'
                            : 'bg-emerald-500/20 text-emerald-500'
                        }`}
                    >
                      {vote.rightChoice ? 'Resolved' : vote.isLocked ? 'Closed' : 'Active'}
                    </span>
                  </div>
                  <h4 className="font-medium text-neutral-200 group-hover:text-emerald-400 transition-colors">
                    {vote.voteText}
                  </h4>
                  <p className="text-xs text-neutral-500 mt-2">{vote.choices.length} options</p>
                </div>
              ))
            )}
          </div>
        </div>
      )}

      {view === 'CREATE' && (
        <div className="flex flex-col max-h-[60vh]">
          <h3 className="text-xl font-bold text-neutral-100 mb-6">Create a Poll</h3>

          <div className="flex-1 overflow-y-auto custom-scrollbar pr-2 flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-sm font-medium text-neutral-400">Question</label>
              <input
                value={voteText}
                onChange={e => setVoteText(e.target.value)}
                placeholder="What is your favorite..."
                className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 rounded-xl px-4 py-3
                  outline-none focus:border-neutral-600"
              />
            </div>

            <div className="flex flex-col gap-2">
              <label className="text-sm font-medium text-neutral-400">Duration</label>
              <select
                value={duration}
                onChange={e => setDuration(Number(e.target.value))}
                className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 rounded-xl px-4 py-3
                  outline-none focus:border-neutral-600 appearance-none"
              >
                <option value={30}>30 seconds</option>
                <option value={60}>1 minute</option>
                <option value={120}>2 minutes</option>
                <option value={300}>5 minutes</option>
              </select>
            </div>

            <div className="flex flex-col gap-2 mt-2">
              <label className="text-sm font-medium text-neutral-400">Options</label>
              {choices.map((choice, index) => (
                <div key={index} className="flex gap-2">
                  <input
                    value={choice}
                    onChange={e => {
                      const newChoices = [...choices]
                      newChoices[index] = e.target.value
                      setChoices(newChoices)
                    }}
                    placeholder={`Option ${index + 1}`}
                    className="flex-1 bg-[#141414] border border-neutral-800 text-neutral-200 rounded-xl px-4 py-2.5
                      outline-none focus:border-neutral-600"
                  />
                  {choices.length > 2 && (
                    <button
                      onClick={() => setChoices(choices.filter((_, i) => i !== index))}
                      className="px-3 text-red-400 hover:bg-neutral-800 rounded-xl transition-colors"
                    >
                      ✕
                    </button>
                  )}
                </div>
              ))}

              {choices.length < 10 && (
                <button
                  onClick={() => setChoices([...choices, ''])}
                  className="mt-2 text-sm text-emerald-500 hover:text-emerald-400 font-medium self-start"
                >
                  + Add option
                </button>
              )}
            </div>
          </div>

          <div className="flex justify-between items-center mt-6 pt-4 border-t border-neutral-800">
            <button onClick={() => setView('LIST')} className="text-neutral-400 hover:text-neutral-200 font-medium">
              Cancel
            </button>
            <button
              onClick={handleCreateSubmit}
              disabled={isCreating || !voteText.trim() || choices.filter(c => c.trim()).length < 2}
              className="px-5 py-2.5 bg-emerald-500 text-neutral-950 font-semibold rounded-xl hover:bg-emerald-400
                disabled:opacity-50 disabled:hover:bg-emerald-500 transition-colors"
            >
              {isCreating ? 'Creating...' : 'Create Poll'}
            </button>
          </div>
        </div>
      )}

      {view === 'DETAILS' && currentVote && (
        <div className="flex flex-col max-h-[60vh]">
          <h3 className="text-xl font-bold text-neutral-100 mb-2">{currentVote.voteText}</h3>

          <div className="mb-6 flex items-center gap-2">
            <span
              className={`text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md ${
                currentVote.rightChoice
                  ? 'bg-neutral-800 text-neutral-500'
                  : currentVote.isLocked
                    ? 'bg-neutral-700 text-neutral-400'
                    : 'bg-emerald-500/20 text-emerald-500'
              }`}
            >
              {currentVote.rightChoice ? 'Resolved' : currentVote.isLocked ? 'Closed' : 'Active'}
            </span>
          </div>

          <div className="flex-1 overflow-y-auto custom-scrollbar pr-2 flex flex-col gap-2">
            {currentVote.choices.map((choice, index) => {
              const isResolved = !!currentVote.rightChoice
              const isWinner = currentVote.rightChoice === choice

              return (
                <button
                  key={index}
                  onClick={() => {
                    if (!currentVote.isLocked) {
                      castVote({ roomId, voteId: currentVote.id, choice })
                    }
                  }}
                  disabled={currentVote.isLocked || isCasting}
                  className={`w-full text-left px-4 py-3 rounded-xl border transition-all ${
                    isWinner
                      ? 'bg-emerald-500/10 border-emerald-500 text-emerald-400 font-semibold'
                      : isResolved
                        ? 'bg-neutral-900 border-neutral-800 text-neutral-600'
                        : currentVote.isLocked
                          ? 'bg-neutral-800 border-neutral-700 text-neutral-400 cursor-not-allowed'
                          : `bg-[#141414] border-neutral-800 text-neutral-200 hover:border-neutral-600
                            hover:bg-neutral-800`
                  }`}
                >
                  {choice}
                </button>
              )
            })}
          </div>

      
          {isHost && (
            <div className="mt-6 p-4 bg-emerald-500/10 border border-emerald-500/20 rounded-xl flex flex-col gap-3">
              <p className="text-sm text-emerald-400 font-medium">Select correct answer to publish results:</p>
              <select
                value={hostSelectedCorrect || ''}
                onChange={e => setHostSelectedCorrect(e.target.value)}
                className="w-full bg-[#141414] border border-neutral-800 text-neutral-200 rounded-lg px-3 py-2
                  outline-none"
              >
                <option value="" disabled>
                  Select option...
                </option>
                {currentVote.choices.map((choice, idx) => (
                  <option key={idx} value={choice}>
                    {choice}
                  </option>
                ))}
              </select>
              <button
                onClick={() => {
                  if (hostSelectedCorrect) {
                    resolveVote({ roomId, voteId: currentVote.id, rightChoice: hostSelectedCorrect })
                  }
                }}
                disabled={!hostSelectedCorrect || isResolving}
                className="mt-2 w-full py-2 bg-emerald-500 text-neutral-950 font-semibold rounded-lg
                  hover:bg-emerald-400 disabled:opacity-50 transition-colors"
              >
                {isResolving ? 'Resolving...' : 'Publish Results'}
              </button>
            </div>
          )}

          <div className="mt-6 pt-4 border-t border-neutral-800">
            <button onClick={() => setView('LIST')} className="text-neutral-400 hover:text-neutral-200 font-medium">
              Back to list
            </button>
          </div>
        </div>
      )}
    </Modal>
  )
}
