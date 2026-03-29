'use client'
import React from 'react'

import { cn } from '@/lib/utils/general.utils'

type ViewType = 'confirm' | 'cancel' | 'loading'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  viewType: ViewType
  paddingClasses?: string
}

export default function Button({ children, viewType, paddingClasses = 'px-4 py-2', className, ...props }: ButtonProps) {
  return (
    <button
      {...props}
      className={cn(
        paddingClasses,
        'w-fit rounded transition-all duration-300 ease-in-out',
        {
          'bg-emerald-500 hover:bg-emerald-400 text-neutral-900 cursor-pointer': viewType === 'confirm',
          'hover:border-border border border-gray-50/0 bg-neutral-900 cursor-pointer': viewType === 'cancel',
          'hover:border-border border border-gray-50/0 bg-neutral-900 cursor-not-allowed': viewType === 'loading',
        },
        className,
      )}
    >
      {children}
    </button>
  )
}
