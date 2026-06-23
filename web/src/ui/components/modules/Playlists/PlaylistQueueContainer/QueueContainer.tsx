'use client'

import { ReactNode } from 'react'

type Props = {
  children: ReactNode
}

export default function QueueContainer({ children }: Props) {
  return (
    <div
      className="flex-1 flex flex-col gap-2 overflow-y-auto pr-1 [&::-webkit-scrollbar]:w-1.5
        [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-background
        [&::-webkit-scrollbar-thumb]:border-0 [&::-webkit-scrollbar-thumb]:rounded-full
        hover:[&::-webkit-scrollbar-thumb]:bg-surface-tertiary/50"
    >
      {children}
    </div>
  )
}
