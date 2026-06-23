'use client'
import { useEffect, useState } from 'react'

import { createPortal } from 'react-dom'
import { RemoveScroll } from 'react-remove-scroll'

import { useIsClient } from '@/lib/hooks/useIsClient'

interface Props {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
  modalClasses?: string
}

export default function Modal({ isOpen, onClose, children, modalClasses = '' }: Props) {
  const mounted = useIsClient()

  const [shouldRender, setShouldRender] = useState(isOpen)

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose()
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown)
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen, onClose])

  useEffect(() => {
    if (isOpen) {
      setShouldRender(true)
    } else {
      const timer = setTimeout(() => setShouldRender(false), 400)
      return () => clearTimeout(timer)
    }
  }, [isOpen])

  if (!mounted || !shouldRender) return null

  return createPortal(
    <RemoveScroll enabled={shouldRender}>
      <div
        className={`fixed inset-0 z-100 flex items-center justify-center transition-all duration-400 ${
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
          <div className="text-foreground-subtle">{children}</div>
        </div>
      </div>
    </RemoveScroll>,
    document.body,
  )
}
