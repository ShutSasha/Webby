'use client'

import { ReactNode } from 'react'

import { cn } from '@/lib/utils/general.utils'

interface Props {
  children: ReactNode
  label: string
  onClick?: () => void
  btnClassName?: string
}

export default function ActionButton({ children, label, onClick, btnClassName }: Props) {
  return (
    <button
      onClick={onClick}
      className={cn(
        'group flex items-center gap-2 shrink-0 px-3 py-2 md:px-4 md:py-2 rounded-full',
        'transition-all duration-300 ease-out cursor-pointer text-nowrap',
        'bg-transparent text-neutral-400',
        'hover:bg-neutral-800 hover:text-neutral-100',
        btnClassName,
      )}
    >
      {children}
      <span className="text-sm font-medium leading-none">{label}</span>
    </button>
  )
}
