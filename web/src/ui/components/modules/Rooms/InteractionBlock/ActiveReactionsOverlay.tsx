'use client'

import { useEffect, useRef } from 'react'

import Image from 'next/image'

import { Reaction } from '@/lib/actions/reaction.actions'
import { useGetReactionsQuery } from '@/lib/hooks/api/reactions/useGetReactionsQuery'
import { useSendReactionMutation } from '@/lib/hooks/api/reactions/useSendReactionMutation'
import { useRoomStore } from '@/stores/room.store'

type Props = {
  roomId: string
}

export default function ActiveReactionsOverlay({ roomId }: Props) {
  const isOpen = useRoomStore(state => state.isReactionsModalOpen)
  const setIsOpen = useRoomStore(state => state.setReactionsModalOpen)
  const overlayRef = useRef<HTMLDivElement>(null)

  const { data: reactions = [], isLoading } = useGetReactionsQuery()
  const { mutate: sendReaction, isPending: isSending } = useSendReactionMutation()

  useEffect(() => {
    if (!isOpen) return

    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as Element

      if (target.closest('.reaction-trigger')) return

      if (overlayRef.current && !overlayRef.current.contains(target)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isOpen, setIsOpen])

  if (!isOpen) return null

  const handleReactionClick = (reaction: Reaction) => {
    if (isSending) return
    sendReaction({ roomId, reactionId: reaction.id })
  }

  return (
    <div
      ref={overlayRef}
      className="absolute bottom-[100px] left-2 right-2 z-50 bg-surface border border-border shadow-xl shadow-black/20
        rounded-xl p-3 animate-in slide-in-from-bottom-2 fade-in duration-200"
    >
      <p className="text-xs font-semibold text-foreground-subtle px-1 pb-1 border-b border-border/50">Send Reaction</p>

      {isLoading ? (
        <div className="flex justify-center items-center py-6">
          <div className="size-6 border-2 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      ) : (
        <div
          className={`grid grid-cols-4 gap-1.5 pt-2 transition-opacity duration-200
            ${isSending ? 'opacity-50 pointer-events-none' : ''}`}
        >
          {reactions.map(reaction => (
            <button
              key={reaction.id}
              onClick={() => handleReactionClick(reaction)}
              disabled={isSending}
              className="flex flex-col items-center justify-center gap-1.5 p-2 rounded-lg hover:bg-surface-tertiary
                transition-colors border border-transparent hover:border-border group/btn disabled:cursor-wait"
            >
              <div className="size-10 relative transition-transform duration-200 group-hover/btn:scale-110">
                <Image
                  src={reaction.stickerUrl}
                  alt={reaction.name}
                  width={80}
                  height={80}
                  className="object-contain"
                />
              </div>
              <span className="text-[10px] font-bold text-emerald-500">{reaction.cost}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
