'use client'

import { MouseEvent } from 'react'

import TrashIcon from '@/assets/icons/ic_trash.svg'
import { useToggleFollowMutation } from '@/lib/hooks/api/user/useToggleFollow'

type Props = {
  targetId: string
  currentUserId?: string
}

export default function UnfollowButton({ targetId, currentUserId }: Props) {
  const { mutate: toggle, isPending } = useToggleFollowMutation(currentUserId)

  const handleToggle = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()

    if (isPending) return

    toggle(targetId)
  }

  return (
    <button className="relative group/trash cursor-pointer" onClick={handleToggle}>
      <TrashIcon
        className="w-4.5 h-4.5 stroke-[1.5px] text-red-500/80 group-hover/trash:text-red-500/90 transition-all
          duration-300"
      />
      <div
        className="absolute w-8 h-8 left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2 group-hover/trash:bg-red-500/20
          transition-all duration-300 rounded-full"
      />
    </button>
  )
}
