'use client'

import { ReactNode } from 'react'

import { cn } from '@/lib/utils/general.utils'

interface Props {
  children: ReactNode
  label: string
  onClick?: () => void
  btnClassName?: string
  disabled?: boolean
}

export default function ActionButton({ children, label, onClick, btnClassName, disabled }: Props) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={cn(
        'group flex items-center gap-2 shrink-0 px-3 py-2 md:px-4 md:py-2 rounded-full',
        'transition-all duration-300 ease-out text-nowrap',
        'bg-transparent text-foreground-muted',
        'not-disabled:cursor-pointer hover:bg-neutral-800 hover:text-foreground-secondary',
        `disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent
        disabled:hover:text-foreground-muted`,
        btnClassName,
      )}
    >
      {children}
      <span className="text-sm font-medium leading-none">{label}</span>
    </button>
  )
}
