'use client'

import { useState } from 'react'

import { leaveComplaint } from '@/app/api/user'
import ComplaintIcon from '@/assets/icons/Profile/ic_complaint.svg'
import { cn, extractServerMessage } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'

import ActionButton from '../../shared/ActionButton'
import Modal from '../../shared/Modal'

type Step = 'REASON' | 'DETAILS' | 'SUCCESS'

type TargetType = 'Video' | 'User'

type Props = {
  authorId: string | undefined
  targetId: string
  targetType: TargetType
}

const COMPLAINT_REASONS: Record<TargetType, string[]> = {
  User: ['Spam', 'Harassment', 'Inappropriate content', 'Scam / Fraud', 'Other'],
  Video: ['Harassment or bullying', 'Violence', 'Explicit content', 'Scam or fraud', 'Misinformation', 'Other'],
}

export default function ComplaintButton({ authorId, targetId, targetType }: Props) {
  const addToast = useToastStore(state => state.addToast)
  const [step, setStep] = useState<Step>('REASON')
  const [selectedReason, setSelectedReason] = useState('')
  const [description, setDescription] = useState('')
  const [loading, setLoading] = useState(false)
  const [isOpen, setIsOpen] = useState(false)

  const currentReasons = COMPLAINT_REASONS[targetType] || []

  const handleClose = () => {
    setIsOpen(false)

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

  const handleSubmit = async () => {
    setLoading(true)
    try {
      if (!authorId) {
        addToast('The user is not authenticated', 'error')
        return
      }

      const response = await leaveComplaint(targetId, selectedReason, targetType, description)

      if (response?.success) {
        setStep('SUCCESS')
      } else {
        const msg = extractServerMessage(response.errors)
        addToast(msg ? msg : 'error while leaving complaint occured', 'error')
      }
    } catch (error) {
      console.error('Failed to send report', error)
      addToast('An unexpected error occurred', 'error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <ActionButton onClick={() => setIsOpen(true)} label="Leave complaint" btnClassName="self-end">
        <ComplaintIcon className="w-4 h-4" />
      </ActionButton>
      <Modal isOpen={isOpen} onClose={handleClose}>
        <div className="flex flex-col w-full gap-3">
          {step === 'REASON' && (
            <>
              <h2 className="text-xl font-semibold text-neutral-300 mb-2 text-center">Why are you reporting?</h2>
              <div className="flex flex-col gap-3">
                {currentReasons.map(r => (
                  <ReasonButton key={r} reason={r} onClick={() => handleSelectReason(r)} />
                ))}
              </div>
            </>
          )}

          {step === 'DETAILS' && (
            <>
              <h2 className="text-xl font-semibold text-neutral-300 text-center">Details</h2>
              <p className="text-sm text-neutral-400 text-center -mt-2">Reason: {selectedReason}</p>

              <textarea
                autoFocus
                className="w-full h-32 bg-neutral-900 border border-neutral-700 rounded-xl p-3 text-neutral-300
                  focus:outline-none focus:border-emerald-500 transition-colors resize-none"
                placeholder="Describe the issue..."
                value={description}
                onChange={e => setDescription(e.target.value)}
              />

              <div className="flex gap-2">
                <button
                  onClick={() => setStep('REASON')}
                  className="flex-1 py-2 text-neutral-400 hover:text-neutral-300 transition-colors border border-border
                    rounded-xl cursor-pointer"
                >
                  Back
                </button>
                <button
                  disabled={loading || !description.trim()}
                  onClick={handleSubmit}
                  className={cn(
                    `flex-2 py-2 bg-emerald-500 hover:bg-emerald-400 rounded-xl text-neutral-900 font-medium
                    transition-all `,
                    {
                      'cursor-pointer': loading === false,
                      'cursor-not-allowed disabled:bg-neutral-700': loading === true,
                    },
                  )}
                >
                  {loading ? 'Sending...' : 'Send Report'}
                </button>
              </div>
            </>
          )}

          {step === 'SUCCESS' && (
            <div className="text-center py-6">
              <div className="text-emerald-500 text-4xl mb-4">✓</div>
              <h2 className="text-xl font-semibold text-neutral-300">Thank you!</h2>
              <p className="text-neutral-400 mt-2">We will review your report shortly.</p>
              <button
                onClick={handleClose}
                className="mt-6 px-8 py-2 rounded-full text-neutral-300 cursor-pointer border border-border
                  hover:bg-neutral-900/70 transition-all duration-300"
              >
                Close
              </button>
            </div>
          )}
        </div>
      </Modal>
    </>
  )
}

function ReasonButton({ reason, onClick }: { reason: string; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className="w-full border border-neutral-700 rounded-lg py-3 px-4 text-left cursor-pointer bg-neutral-900/50
        hover:bg-neutral-800 hover:border-emerald-500 transition-all text-neutral-200"
    >
      {reason}
    </button>
  )
}
