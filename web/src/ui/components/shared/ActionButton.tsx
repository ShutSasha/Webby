'use client'

import { ReactNode } from 'react'

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
      className={`flex items-center gap-2 px-4 py-2 bg-black/40 rounded-full transition-all duration-500 ease-in-out
        cursor-pointer active:scale-90 group border border-transparent hover:border-neutral-700/70 text-nowrap
        ${btnClassName}`}
    >
      {children}

      <span className="text-neutral-200 text-sm leading-none">{label}</span>
    </button>
  )
}
