'use client'

import { useEffect, useState } from 'react'

import Image from 'next/image'

import { ActiveReaction, useRoomStore } from '@/stores/room.store'

function Bubble({ reaction, onRemove }: { reaction: ActiveReaction; onRemove: (uid: string) => void }) {
  const [style, setStyle] = useState<React.CSSProperties>({})

  useEffect(() => {
    const size = Math.floor(Math.random() * 24) + 36
    const leftPos = Math.floor(Math.random() * 80) + 10
    const wobbleDuration = (Math.random() * 2 + 2).toFixed(1)
    const floatDuration = 5

    setStyle({
      width: `${size}px`,
      height: `${size}px`,
      left: `${leftPos}%`,
      animation: `float-up ${floatDuration}s cubic-bezier(0.25, 1, 0.5, 1) forwards, wobble ${wobbleDuration}s ease-in-out infinite alternate`,
    })

    const timer = setTimeout(() => {
      onRemove(reaction.uid)
    }, floatDuration * 1000)

    return () => clearTimeout(timer)
  }, [reaction, onRemove])

  if (!style.width) return null

  return (
    <div className="absolute bottom-0 z-51 pointer-events-none" style={style}>
      <Image
        src={reaction.stickerUrl}
        alt={reaction.name}
        sizes="60px"
        fill
        className="object-contain drop-shadow-[0_0_15px_rgba(16,185,129,0.3)]"
      />
    </div>
  )
}

export default function ReactionBubbles() {
  const activeReactions = useRoomStore(state => state.activeReactions)
  const removeReaction = useRoomStore(state => state.removeReaction)

  if (activeReactions.length === 0) return null

  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none z-51 rounded-2xl">
      <style>{`
        @keyframes float-up {
          0% { transform: translateY(50px) scale(0.5); opacity: 0; }
          10% { transform: translateY(0px) scale(1.1); opacity: 1; }
          20% { transform: translateY(-50px) scale(1); opacity: 1; }
          80% { opacity: 0.9; }
          100% { transform: translateY(-500px) scale(1); opacity: 0; }
        }
        @keyframes wobble {
          0% { margin-left: -15px; }
          100% { margin-left: 15px; }
        }
      `}</style>

      {activeReactions.map(reaction => (
        <Bubble key={reaction.uid} reaction={reaction} onRemove={removeReaction} />
      ))}
    </div>
  )
}
