'use client'

import { useState } from 'react'

import UserPlusIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import UserMinusIcon from '@/assets/icons/shared/minus.svg'
import { useToggleFollowMutation } from '@/lib/hooks/api/user/useToggleFollow'
import { cn } from '@/lib/utils/general.utils'

interface Props {
  targetUserId: string
  initialIsFollowing: boolean
  currentUserId?: string
}

export default function FollowButton({ targetUserId, initialIsFollowing, currentUserId }: Props) {
  const [isFollowing, setIsFollowing] = useState(initialIsFollowing)

  const { mutate: toggleFollow, isPending } = useToggleFollowMutation(currentUserId)

  const handleFollow = () => {
    if (isPending) return

    const prevStatus = isFollowing
    setIsFollowing(!prevStatus)

    toggleFollow(targetUserId, {
      onError: _ => {
        setIsFollowing(prevStatus)
      },
    })
  }

  return (
    <button
      onClick={handleFollow}
      disabled={isPending}
      className={cn(
        'flex items-center gap-2 px-4 py-2 rounded-full font-medium transition-all duration-300 cursor-pointer',
        {
          'bg-neutral-800 text-foreground-muted hover:bg-red-500/10 hover:text-red-500': isFollowing && !isPending,
          'bg-emerald-500 text-foreground-inverse-subtle hover:bg-emerald-400': !isFollowing && !isPending,
          'opacity-50 grayscale pointer-events-none ': isPending,
        },
      )}
    >
      {isFollowing && !isPending && (
        <>
          <UserMinusIcon className="w-4 h-4" />
          <span className="text-sm leading-3.5">Unfollow</span>
        </>
      )}
      {!isFollowing && !isPending && (
        <>
          <UserPlusIcon className="w-4 h-4" />
          <span className="text-sm leading-3.5">Follow</span>
        </>
      )}
      {isPending && (
        <>
          <span className="text-sm leading-3.5">Thinking...</span>
        </>
      )}
    </button>
  )
}
