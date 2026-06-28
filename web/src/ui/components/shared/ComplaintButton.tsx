'use client'

import { useState } from 'react'

import ComplaintIcon from '@/assets/icons/Profile/ic_complaint.svg'

import ActionButton from './ActionButton'
import ComplaintModal, { TargetType } from './ComplaintModal'

type Props = {
  targetId: string
  targetType: TargetType
}

export default function ComplaintButton({ targetId, targetType }: Props) {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      <ActionButton onClick={() => setIsOpen(true)} label="Leave complaint" btnClassName="self-end">
        <ComplaintIcon className="w-4 h-4" />
      </ActionButton>

      <ComplaintModal isOpen={isOpen} onClose={() => setIsOpen(false)} targetId={targetId} targetType={targetType} />
    </>
  )
}
