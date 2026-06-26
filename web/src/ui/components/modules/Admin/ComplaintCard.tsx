'use client'

import { useState } from 'react'

import Link from 'next/link'

import { Complaint } from '@/lib/actions/admin.actions'
import { useAcceptComplaintMutation } from '@/lib/hooks/api/admin/useAcceptComplaint'
import { useDenyComplaintMutation } from '@/lib/hooks/api/admin/useDenyComplaint'
import { formatTimeAgo } from '@/lib/utils/date.utils'

type Props = {
  complaint: Complaint
}

export default function ComplaintCard({ complaint }: Props) {
  const [isDenying, setIsDenying] = useState(false)
  const [denyReason, setDenyReason] = useState('')

  const { mutate: acceptComplaint } = useAcceptComplaintMutation()
  const { mutate: denyComplaint } = useDenyComplaintMutation()

  const handleAccept = () => {
    acceptComplaint({ complaintId: complaint.id })
  }

  const handleDenySubmit = () => {
    if (!denyReason.trim()) return
    denyComplaint({ complaintId: complaint.id, reason: denyReason.trim() })
  }

  return (
    <div className="bg-surface border border-border rounded-2xl p-5 flex flex-col gap-4">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <span
              className="px-2 py-0.5 bg-background text-foreground-subtle rounded text-xs font-bold uppercase
                tracking-wider"
            >
              {complaint.targetType}
            </span>
            <span className="text-red-400 text-sm font-bold">{complaint.reasonType}</span>
            <span className="text-foreground-disabled text-xs ml-2">{formatTimeAgo(complaint.createdAt)}</span>
          </div>
          <Link
            href={
              complaint.targetType === 'User' ? `/profile/${complaint.targetId}` : `/videos/wb_${complaint.targetId}`
            }
            className="text-lg font-bold text-foreground-secondary mt-1 flex items-center gap-2"
          >
            Target name: <span className="text-sm font-mono text-foreground-muted">{complaint.targetId}</span>
          </Link>
        </div>
        <div className="text-sm text-foreground-faint flex flex-col items-end">
          <span>Reported by:</span>
          <Link
            href={`/profile/${complaint.authorId}`}
            className="text-emerald-500 font-medium cursor-pointer hover:underline font-mono text-xs mt-0.5"
          >
            {complaint?.username || 'смайли'}
          </Link>
        </div>
      </div>

      {/* Body */}
      {complaint.additionalInfo && (
        <div className="bg-surface/50 rounded-xl p-4 text-foreground-subtle text-sm border-l-4 border-red-500/50">
          &quot;{complaint.additionalInfo}&quot;
        </div>
      )}

      {/* Actions */}
      <div className="flex flex-col gap-3 mt-2 pt-4">
        {isDenying ? (
          <div className="flex items-center gap-2 w-full animate-in slide-in-from-top-2 duration-200">
            <input
              type="text"
              autoFocus
              placeholder="Reason for dismissal..."
              value={denyReason}
              onChange={e => setDenyReason(e.target.value)}
              className="flex-1 bg-background border border-border text-foreground-tertiary rounded-xl px-4 py-2 text-sm
                outline-none focus:border-neutral-600"
            />
            <button
              onClick={() => setIsDenying(false)}
              className="px-4 py-2 text-foreground-muted hover:text-foreground-secondary text-sm font-medium
                transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleDenySubmit}
              disabled={!denyReason.trim()}
              className="px-4 py-2 bg-emerald-500 text-foreground-inverse font-semibold rounded-xl text-sm
                transition-colors disabled:opacity-50"
            >
              Confirm
            </button>
          </div>
        ) : (
          <div className="flex items-center gap-3">
            <button
              onClick={handleAccept}
              className="ml-auto px-5 py-2 bg-red-500 hover:bg-red-600 text-foreground-strong font-semibold rounded-xl
                text-sm transition-colors"
            >
              {complaint.targetType.toLowerCase() === 'user' ? 'Ban User' : 'Delete Content'}
            </button>

            <button
              onClick={() => setIsDenying(true)}
              className="px-5 py-2 bg-background hover:bg-surface-tertiary text-foreground-subtle font-semibold
                rounded-xl text-sm transition-colors"
            >
              Dismiss Report
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
