'use client'

import { useState } from 'react'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'

import Modal from '../shared/Modal'

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
        className="flex items-center gap-2 shrink-0 order-3 lg:order-3 cursor-pointer uppercase bg-emerald-500
          transition-colors duration-300 ease-out hover:bg-emerald-400 text-neutral-900 text-[14px] font-semibold
          md:font-bold leading-5.5 px-4 py-1.5 md:py-2 md:px-5 rounded-xl"
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
