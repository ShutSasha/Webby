'use client'

import { useState } from 'react'

import { useLeaveComplaintMutation } from '@/lib/hooks/api/user/useLeaveComplaintMutation'
import { cn } from '@/lib/utils/general.utils'

import Modal from './Modal'

export type Step = 'REASON' | 'DETAILS' | 'SUCCESS'
export type TargetType = 'Video' | 'User'

type ModalProps = {
  isOpen: boolean
  onClose: () => void
  targetId: string
  targetType: TargetType
}

const COMPLAINT_REASONS: Record<TargetType, string[]> = {
  User: ['Spam', 'Harassment', 'Inappropriate content', 'Scam / Fraud', 'Other'],
  Video: ['Harassment or bullying', 'Violence', 'Explicit content', 'Scam or fraud', 'Misinformation', 'Other'],
}

export default function ComplaintModal({ isOpen, onClose, targetId, targetType }: ModalProps) {
  const [step, setStep] = useState<Step>('REASON')
  const [selectedReason, setSelectedReason] = useState('')
  const [description, setDescription] = useState('')

  const { mutate: submitComplaint, isPending } = useLeaveComplaintMutation({
    onSuccess: () => setStep('SUCCESS'),
  })

  const currentReasons = COMPLAINT_REASONS[targetType] || []

  const handleClose = () => {
    onClose()

    setTimeout(() => {
      setStep('REASON')
      setSelectedReason('')
      setDescription('')
    }, 300)
  }

  const handleSelectReason = (reason: string) => {
    setSelectedReason(reason)
    setStep('DETAILS')
  }

  const handleSubmit = () => {
    submitComplaint({
      targetUserId: targetId,
      reasonType: selectedReason,
      targetType,
      additionalInfo: description.trim(),
    })
  }

  return (
    <Modal isOpen={isOpen} onClose={handleClose}>
      <div className="flex flex-col w-full gap-3">
        {step === 'REASON' && (
          <>
            <h2 className="text-xl font-semibold text-foreground-subtle mb-2 text-center">Why are you reporting?</h2>
            <div className="flex flex-col gap-3">
              {currentReasons.map(r => (
                <ReasonButton key={r} reason={r} onClick={() => handleSelectReason(r)} />
              ))}
            </div>
          </>
        )}

        {step === 'DETAILS' && (
          <>
            <h2 className="text-xl font-semibold text-foreground-subtle text-center">Details</h2>
            <p className="text-sm text-foreground0 text-center -mt-1 mb-2">
              You can skip this or add more info (optional)
            </p>

            <textarea
              autoFocus
              className="w-full h-32 bg-neutral-900 border border-neutral-700 rounded-xl p-3 text-foreground-subtle
                focus:outline-none focus:border-emerald-500 transition-colors resize-none"
              placeholder="Describe the issue..."
              value={description}
              onChange={e => setDescription(e.target.value)}
            />

            <div className="flex gap-2">
              <button
                onClick={() => setStep('REASON')}
                className="flex-1 py-2 text-foreground-muted hover:text-foreground-subtle transition-colors border
                  border-border rounded-xl cursor-pointer"
              >
                Back
              </button>
              <button
                disabled={isPending}
                onClick={handleSubmit}
                className={cn(
                  `flex-2 py-2 bg-emerald-500 hover:bg-emerald-400 rounded-xl text-foreground-inverse-subtle font-medium
                  transition-all `,
                  {
                    'cursor-pointer': !isPending,
                    'cursor-not-allowed disabled:bg-neutral-700': isPending,
                  },
                )}
              >
                {isPending ? 'Sending...' : 'Send Report'}
              </button>
            </div>
          </>
        )}

        {step === 'SUCCESS' && (
          <div className="text-center py-6">
            <div className="text-emerald-500 text-4xl mb-4">✓</div>
            <h2 className="text-xl font-semibold text-foreground-subtle">Thank you!</h2>
            <p className="text-foreground-muted mt-2">We will review your report shortly.</p>
            <button
              onClick={handleClose}
              className="mt-6 px-8 py-2 rounded-full text-foreground-subtle cursor-pointer border border-border
                hover:bg-neutral-900/70 transition-all duration-300"
            >
              Close
            </button>
          </div>
        )}
      </div>
    </Modal>
  )
}

function ReasonButton({ reason, onClick }: { reason: string; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className="w-full border border-neutral-700 rounded-2xl py-3 px-4 text-left cursor-pointer bg-neutral-900/50
        hover:bg-neutral-800 hover:border-emerald-500 transition-all text-foreground-tertiary"
    >
      {reason}
    </button>
  )
}
