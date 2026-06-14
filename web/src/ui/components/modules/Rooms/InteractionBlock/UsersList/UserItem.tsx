'use client'

import { useState, useRef, useEffect } from 'react'

import Link from 'next/link'

import DotsIcon from '@/assets/icons/shared/more-horizontal.svg'
import { DEFAULT_USER_THUMBNAIL } from '@/lib/constants/url.constamts'
import { useRemoveRoomMemberMutation } from '@/lib/hooks/api/room/useRemoveRoomMember'
import { extractServerMessage } from '@/lib/utils/general.utils'
import { useToastStore } from '@/stores/toast-store'
import { RoomMember } from '@/types/room.types'
import SafeImage from '@/ui/components/shared/SafeImage'

export default function UserItem({ user, roomId }: { user: RoomMember; roomId: string }) {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const addToast = useToastStore(state => state.addToast)

  const { mutate: removeMember, isPending } = useRemoveRoomMemberMutation()

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false)
      }
    }
    if (isMenuOpen) document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isMenuOpen])

  const handleKickOut = () => {
    if (isPending) return

    removeMember(
      { roomId, memberId: user.userId },
      {
        onSuccess: res => {
          if (res.success) {
            setIsMenuOpen(false)
            addToast('User removed successfully', 'success')
          } else {
            const errorMessage = extractServerMessage(res.errors)
            addToast(errorMessage || 'An unexpected error occurred', 'error')
          }
        },
        onError: error => {
          addToast(error.message || 'An unexpected error occurred', 'error')
        },
      },
    )
  }

  return (
    <div
      className={`flex items-center justify-between p-2 rounded-lg hover:bg-neutral-900 group transition-colors
        duration-300 relative ${isPending ? 'opacity-50 pointer-events-none' : ''}`}
    >
      <div className="flex items-center gap-3">
        <SafeImage
          src={user.avatarUrl || DEFAULT_USER_THUMBNAIL}
          alt={user.username}
          width={50}
          height={50}
          className="object-cover w-8 h-8 rounded-full"
          fallbackType="user"
        />

        <span
          className="text-sm font-medium text-neutral-300 group-hover:text-white transition-colors truncate
            max-w-[120px]"
        >
          {user.username}
        </span>
      </div>

      <div className="relative shrink-0" ref={menuRef}>
        <button
          onClick={() => setIsMenuOpen(!isMenuOpen)}
          className="p-1 hover:bg-neutral-800 rounded-md text-neutral-500 hover:text-white transition-all
            cursor-pointer"
        >
          <DotsIcon className="w-5 h-5" />
        </button>

        {isMenuOpen && (
          <div
            className="absolute right-0 top-full mt-3 w-40 bg-neutral-900 border border-neutral-800 rounded-lg shadow-xl
              z-50 py-1 animate-in fade-in zoom-in duration-150"
          >
            <Link
              href={`/profile/${user.userId}`}
              className="block w-full text-left px-3 py-2 text-xs hover:bg-neutral-800 transition-colors cursor-pointer"
            >
              View Profile
            </Link>
            <div className="h-px bg-neutral-800 my-1" />
            <button
              onClick={handleKickOut}
              disabled={isPending}
              className="w-full text-left px-3 py-2 text-xs hover:bg-red-500/20 text-red-500 transition-colors
                cursor-pointer disabled:opacity-50"
            >
              {isPending ? 'Kicking...' : 'Kick Out'}
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
