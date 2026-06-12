'use client'

import { useEffect, useRef } from 'react'

import { createPortal } from 'react-dom'

type Props = {
  isOpen: boolean
  x: number
  y: number
  isOwnMessage: boolean
  onClose: () => void
  onCopy: () => void
  onDelete: () => void
}

export default function MessageContextMenu({ isOpen, x, y, isOwnMessage, onClose, onCopy, onDelete }: Props) {
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!isOpen) return

    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        onClose()
      }
    }

    const handleScroll = () => onClose()

    document.addEventListener('mousedown', handleClickOutside)
    window.addEventListener('scroll', handleScroll, true)

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      window.removeEventListener('scroll', handleScroll, true)
    }
  }, [isOpen, onClose])

  if (!isOpen) return null

  const safeX = Math.min(x, typeof window !== 'undefined' ? window.innerWidth - 160 : x)
  const safeY = Math.min(y, typeof window !== 'undefined' ? window.innerHeight - 100 : y)

  return createPortal(
    <div
      ref={menuRef}
      style={{ top: safeY, left: safeX }}
      className="fixed z-100 w-40 bg-neutral-800 border border-neutral-700/60 shadow-xl shadow-black/50 py-1.5
        rounded-xl animate-in fade-in zoom-in-95 duration-150"
      onContextMenu={e => e.preventDefault()}
    >
      <button
        onClick={() => {
          onCopy()
          onClose()
        }}
        className="w-full text-left px-4 py-2 text-sm text-neutral-200 hover:bg-neutral-700/50 transition-colors"
      >
        Copy text
      </button>

      {isOwnMessage && (
        <button
          onClick={() => {
            onDelete()
            onClose()
          }}
          className="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-neutral-700/50 transition-colors"
        >
          Delete message
        </button>
      )}
    </div>,
    document.body,
  )
}
