'use client'

import { useState } from 'react'

import ShareIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import ActionButton from '@/ui/components/shared/ActionButton'

import RoomInviteModal from './RoomInviteModal'

export default function RoomInviteButton({ roomId }: { roomId: string }) {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      <ActionButton label="send an invite" onClick={() => setIsOpen(true)}>
        <ShareIcon className="size-4 stroke-[1.5px]" />
      </ActionButton>

      {isOpen && <RoomInviteModal roomId={roomId} isOpen={isOpen} onClose={() => setIsOpen(false)} />}
    </>
  )
}
