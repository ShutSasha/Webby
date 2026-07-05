'use client'

import { useEffect, useRef } from 'react'

import Image from 'next/image'

import { useRoomStore } from '@/stores/room.store'

export type Reaction = {
  id: string
  name: string
  cost: number
  stickerUrl: string
}

const MOCK_REACTIONS: Reaction[] = [
  { id: '1', name: 'Like', cost: 10, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Like' },
  { id: '2', name: 'Heart', cost: 50, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Heart' },
  { id: '3', name: 'Laugh', cost: 100, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Laugh' },
  { id: '4', name: 'Wow', cost: 250, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Wow' },
  { id: '5', name: 'Like', cost: 10, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Like' },
  { id: '6', name: 'Heart', cost: 50, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Heart' },
  { id: '7', name: 'Laugh', cost: 100, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Laugh' },
  { id: '8', name: 'Wow', cost: 250, stickerUrl: 'https://api.dicebear.com/7.x/notionists/svg?seed=Wow' },
]

export default function ActiveReactionsOverlay() {
  const isOpen = useRoomStore(state => state.isReactionsModalOpen)
  const setIsOpen = useRoomStore(state => state.setReactionsModalOpen)
  const overlayRef = useRef<HTMLDivElement>(null)

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

  const handleReactionClick = (reaction: Reaction) => {
    setIsOpen(false)
  }

  return (
    <div
      ref={overlayRef}
      className="absolute bottom-[100px] left-2 right-2 z-50 bg-surface border border-border shadow-xl shadow-black/20
        rounded-xl p-3 animate-in slide-in-from-bottom-2 fade-in duration-200"
    >
      <p className="text-xs font-semibold text-foreground-subtle px-1 pb-1 border-b border-border/50">Send Reaction</p>

      <div className="grid grid-cols-4 gap-1.5 pt-2">
        {MOCK_REACTIONS.map(reaction => (
          <button
            key={reaction.id}
            onClick={() => handleReactionClick(reaction)}
            className="flex flex-col items-center justify-center gap-1.5 p-2 rounded-lg hover:bg-surface-tertiary
              transition-colors border border-transparent hover:border-border group/btn"
          >
            <div className="size-12 relative transition-transform duration-200 group-hover/btn:scale-110">
              <Image src={reaction.stickerUrl} alt={reaction.name} fill className="object-contain" />
            </div>
            <span className="text-[10px] font-bold text-emerald-500">{reaction.cost}</span>
          </button>
        ))}
      </div>
    </div>
  )
}
