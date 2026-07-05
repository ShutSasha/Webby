'use client'

import { useEffect, useRef } from 'react'

import Image from 'next/image'

import { Reaction } from '@/lib/actions/reaction.actions'
import { useGetReactionsQuery } from '@/lib/hooks/api/reactions/useGetReactionsQuery'
import { useRoomStore } from '@/stores/room.store'

export default function ActiveReactionsOverlay() {
  const isOpen = useRoomStore(state => state.isReactionsModalOpen)
  const setIsOpen = useRoomStore(state => state.setReactionsModalOpen)
  const overlayRef = useRef<HTMLDivElement>(null)

  const { data: reactions = [], isLoading } = useGetReactionsQuery()

  useEffect(() => {
    if (!isOpen) return

    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as Element

      if (overlayRef.current && !overlayRef.current.contains(target)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isOpen, setIsOpen])

  if (!isOpen) return null

  const handleReactionClick = (reaction: Reaction) => {}

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
        <div className="grid grid-cols-4 gap-1.5 pt-2">
          {reactions.map(reaction => (
            <button
              key={reaction.id}
              onClick={() => handleReactionClick(reaction)}
              className="flex flex-col items-center justify-center gap-1.5 p-2 rounded-lg hover:bg-surface-tertiary
                transition-colors border border-transparent hover:border-border group/btn"
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
