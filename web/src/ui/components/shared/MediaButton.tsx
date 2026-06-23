'use client'

import { useState } from 'react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { cn } from '@/lib/utils/general.utils'

import Modal from './Modal'

type MediaButtonProps = {
  actionLabel?: string
  children?: React.ReactNode
  isOpen?: boolean
  setIsOpen?: (isOpen: boolean) => void
}

export default function MediaButton({
  actionLabel = 'Create',
  children,
  isOpen: externalIsOpen,
  setIsOpen: externalSetIsOpen,
}: MediaButtonProps) {
  const [internalIsOpen, setInternalIsOpen] = useState(false)
  const isControlled = externalIsOpen !== undefined && externalSetIsOpen !== undefined

  const isOpen = isControlled ? externalIsOpen : internalIsOpen
  const setIsOpen = isControlled ? externalSetIsOpen : setInternalIsOpen

  return (
    <>
      <button
        className={cn(
          'flex items-center gap-2 shrink-0 cursor-pointer order-3',
          'bg-background text-foreground-tertiary',
          'text-sm font-semibold px-5 py-2.5 rounded-full',
          'transition-all duration-300 ease-out',
          'hover:bg-surface-tertiary/40',
          'shadow-none',
        )}
        onClick={() => setIsOpen(true)}
      >
        <PlusIcon className="h-4 w-4 text-foreground-tertiary" aria-hidden="true" />
        <p>{actionLabel}</p>
      </button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)}>
        {children}
      </Modal>
    </>
  )
}
