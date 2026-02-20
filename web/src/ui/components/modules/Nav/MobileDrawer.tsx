'use client'
import { useEffect, useState } from 'react'

import { createPortal } from 'react-dom'

interface Props {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
}

export default function MobileDrawer({ isOpen, onClose, children }: Props) {
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
    <div className={`fixed inset-0 z-100 lg:hidden transition-all duration-400 ${isOpen ? 'visible' : 'invisible'}`}>
      <div
        className={`absolute inset-0 bg-black/60 transition-opacity duration-400 ease-in-out ${
          isOpen ? 'opacity-100' : 'opacity-0'
        }`}
        onClick={onClose}
      />

      <div
        className={`absolute top-0 right-0 h-full w-[calc(100%-70px)] max-w-[440px] bg-neutral-900 shadow-2xl p-6
          transition-transform duration-400 ease-in-out ${isOpen ? 'translate-x-0' : 'translate-x-full'}`}
      >
        <p className="text-white">
          Lorem ipsum dolor sit amet consectetur, adipisicing elit. Tempora cum architecto fugiat eius aliquam rem in
          accusamus fugit numquam, quos similique sed corporis placeat quis corrupti veniam est eligendi dolor?
        </p>
        <div className="mt-10">{children}</div>
      </div>
    </div>,
    document.body,
  )
}
