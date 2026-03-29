'use client'

import { useState } from 'react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { cn } from '@/lib/utils/general.utils'

import Modal from './Modal'

type MediaButtonProps = {
  actionLabel: string
  children?: React.ReactNode
  isOpen?: boolean
  setIsOpen?: (isOpen: boolean) => void
}

export default function MediaButton({
  actionLabel,
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
          'flex items-center gap-2 shrink-0 cursor-pointer uppercase bg-emerald-500',
          'transition-all duration-300 ease-out hover:bg-emerald-400 text-neutral-900',
          'text-[13px] tracking-wide font-bold px-4 py-2 md:px-5 rounded-xl order-3',
          'shadow-[0_0_15px_rgba(16,185,129,0.2)] hover:shadow-[0_0_20px_rgba(16,185,129,0.4)]',
        )}
        onClick={() => setIsOpen(true)}
      >
        <PlusIcon className="h-4 w-4 text-neutral-900" aria-hidden="true" />
        {actionLabel}
      </button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)}>
        {children}
      </Modal>
    </>
  )
}
