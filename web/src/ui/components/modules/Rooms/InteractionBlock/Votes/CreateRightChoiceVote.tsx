'use client'

import { useState } from 'react'

import { useCreateRightChoiceVoteMutation } from '@/lib/hooks/api/vote/useCreateRightChoiceVote'

type Props = {
  roomId: string
  onBack: () => void
}

export default function CreateRightChoiceVote({ roomId, onBack }: Props) {
  const [voteText, setVoteText] = useState('')
  const [duration, setDuration] = useState(120)
  const [choices, setChoices] = useState(['', ''])

  const { mutate: createVote, isPending } = useCreateRightChoiceVoteMutation()

  const handleCreateSubmit = () => {
    const validChoices = choices.map(c => c.trim()).filter(Boolean)
    if (!voteText.trim() || validChoices.length < 2) return

    createVote(
      { roomId, voteText: voteText.trim(), duration, choices: validChoices },
      { onSuccess: onBack }
    )
  }

  const handleChoiceChange = (index: number, value: string) => {
    const newChoices = [...choices]
    newChoices[index] = value
    setChoices(newChoices)
  }

  const removeChoice = (index: number) => {
    setChoices(choices.filter((_, i) => i !== index))
  }

  const isSubmitDisabled = isPending || !voteText.trim() || choices.filter(c => c.trim()).length < 2

  return (
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
                onChange={e => handleChoiceChange(index, e.target.value)}
                placeholder={`Option ${index + 1}`}
                className="flex-1 bg-[#141414] border border-neutral-800 text-neutral-200 rounded-xl px-4 py-2.5
                  outline-none focus:border-neutral-600"
              />
              {choices.length > 2 && (
                <button
                  onClick={() => removeChoice(index)}
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
        <button onClick={onBack} className="text-neutral-400 hover:text-neutral-200 font-medium">
          Cancel
        </button>
        <button
          onClick={handleCreateSubmit}
          disabled={isSubmitDisabled}
          className="px-5 py-2.5 bg-emerald-500 text-neutral-950 font-semibold rounded-xl hover:bg-emerald-400
            disabled:opacity-50 disabled:hover:bg-emerald-500 transition-colors"
        >
          {isPending ? 'Creating...' : 'Create Poll'}
        </button>
      </div>
    </div>
  )
}