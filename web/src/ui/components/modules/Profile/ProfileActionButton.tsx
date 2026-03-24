'use client'

import { ReactNode } from 'react'

interface ProfileActionButtonProps {
  children: ReactNode
  label: string
  onClick?: () => void

  btnClassName?: string
}

export default function ProfileActionButton({ children, label, onClick, btnClassName }: ProfileActionButtonProps) {
  return (
    <button
      onClick={onClick}
      className={`flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-full transition-all duration-300
        hover:bg-neutral-700/40 cursor-pointer active:scale-90 group border border-transparent ${btnClassName}`}
    >
      {children}

      <span className="text-white text-sm font-medium leading-none">{label}</span>
    </button>
  )
}
