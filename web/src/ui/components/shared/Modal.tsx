'use client'
import { useEffect, useState } from 'react'

import { createPortal } from 'react-dom'

interface Props {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
  modalClasses?: string
}

export default function Modal({ isOpen, onClose, children, modalClasses }: Props) {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'unset'
    }
  }, [isOpen])

  if (!mounted) return null

  return createPortal(
    <div
      className={`fixed inset-0 z-100 flex items-center justify-center p-4 transition-all duration-400 ${
        isOpen ? 'visible' : 'invisible'
      }`}
    >
      <div
        className={`absolute inset-0 bg-black/70 transition-opacity duration-400 ease-in-out ${
          isOpen ? 'opacity-100' : 'opacity-0'
        }`}
        onClick={onClose}
      />

      <div
        className={`relative w-full max-w-[440px] border border-neutral-800 bg-neutral-900 shadow-2xl p-6 rounded-2xl
          transition-all duration-400 ease-out ${modalClasses}
          ${isOpen ? 'opacity-100 scale-100 translate-y-0' : 'opacity-0 scale-95 translate-y-8'}`}
      >
        <div className="text-neutral-300">{children}</div>
      </div>
    </div>,
    document.body,
  )
}
