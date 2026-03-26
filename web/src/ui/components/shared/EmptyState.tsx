import React from 'react'

import SadSmileIcon from '@/assets/icons/Errors/sad-smile.svg'
import { cn } from '@/lib/utils/utils'

type EmptyStateProps = {
  title: string
  description?: string
  className?: string
  icon?: React.ReactNode
}

export default function EmptyState({ title, description, className, icon }: EmptyStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col flex-1 items-center justify-center w-full py-20 px-4 animate-in fade-in zoom-in-95 duration-500',
        className,
      )}
    >
      <div
        className="size-20 bg-neutral-900/50 rounded-full flex items-center justify-center mb-6 border
          border-neutral-800 shadow-inner"
      >
        {icon ? icon : <SadSmileIcon className="size-9 text-neutral-600 stroke-[1.5px]" />}
      </div>
      <h3 className="text-xl font-bold text-neutral-300 text-center mb-2">{title}</h3>

      {description && <p className="text-sm text-neutral-500 text-center max-w-sm leading-relaxed">{description}</p>}
    </div>
  )
}
