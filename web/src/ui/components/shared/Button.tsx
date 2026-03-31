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
          'bg-emerald-500 hover:bg-emerald-400 text-neutral-900 cursor-pointer': viewType === 'confirm',
          [`bg-transparent border border-neutral-700 text-neutral-300 hover:border-neutral-500 hover:text-neutral-100
          hover:bg-neutral-800/50 cursor-pointer`]: viewType === 'cancel',
          'bg-neutral-800 text-neutral-500 cursor-not-allowed': viewType === 'loading',
        },
        className,
      )}
    >
      {children}
    </button>
  )
}
