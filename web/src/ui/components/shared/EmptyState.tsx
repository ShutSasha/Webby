import React from 'react'

import SadSmileIcon from '@/assets/icons/Errors/sad-smile.svg'
import { cn } from '@/lib/utils/general.utils'

type EmptyStateProps = {
  title: string
  description?: string
  className?: string
  icon?: React.ReactNode
  disableFlex?: boolean
}

export default function EmptyState({ title, description, className, icon, disableFlex }: EmptyStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center w-full py-20 px-4 animate-in fade-in zoom-in-95 duration-500',
        !disableFlex && 'flex-1',
        className,
      )}
    >
      <div
        className="size-20 bg-surface/50 rounded-full flex items-center justify-center mb-6 border border-border
          shadow-inner"
      >
        {icon ? icon : <SadSmileIcon className="size-9 text-foreground-disabled stroke-[1.5px]" />}
      </div>
      <h3 className="text-xl font-bold text-foreground-subtle text-center mb-2">{title}</h3>

      {description && (
        <p className="text-sm text-foreground-faint text-center max-w-sm leading-relaxed">{description}</p>
      )}
    </div>
  )
}
