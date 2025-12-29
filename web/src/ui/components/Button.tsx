'use client'
import React from 'react'

type ViewType = 'Confirm' | 'Cancel'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  viewType: ViewType
  paddingClasses?: string
}

export default function Button({ children, viewType, paddingClasses = 'px-4 py-2', className, ...props }: ButtonProps) {
  const btnClasess =
    viewType === 'Confirm'
      ? 'bg-emerald-500 hover:bg-emerald-400 text-neutral-900'
      : 'hover:border-border border border-gray-50/0 bg-neutral-900'

  return (
    <button
      {...props}
      className={`${btnClasess} ${paddingClasses} ${className} w-fit cursor-pointer rounded transition-all duration-300
        ease-in-out`}
    >
      {children}
    </button>
  )
}
