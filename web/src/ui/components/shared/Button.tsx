'use client'
import React from 'react'

import { cn } from '@/lib/utils/general.utils'

type ViewType = 'confirm' | 'cancel' | 'loading'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  viewType: ViewType
  paddingClasses?: string
}

export default function Button({
  children,
  viewType,
  paddingClasses = 'px-4 py-2.5',
  className,
  ...props
}: ButtonProps) {
  return (
    <button
      {...props}
      className={cn(
        'flex items-center justify-center rounded-full font-semibold transition-all duration-300 ease-out',
        paddingClasses,
        {
          'bg-emerald-500 hover:bg-emerald-400 text-neutral-950 cursor-pointer': viewType === 'confirm',
          [`bg-transparent border border-skeleton-pulse text-foreground hover:border-muted hover:text-foreground-strong
          hover:bg-surface-subtle cursor-pointer`]: viewType === 'cancel',
          'bg-card text-muted cursor-not-allowed': viewType === 'loading',
        },
        className,
      )}
    >
      {children}
    </button>
  )
}
