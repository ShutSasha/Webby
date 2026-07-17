import { ReactNode } from 'react'

import { cn } from '@/lib/utils/general.utils'

type Props = {
  children: ReactNode
  className?: string
}

export default function GridCardsContainer({ children, className }: Props) {
  return (
    <div
      className={cn('grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3', className)}
    >
      {children}
    </div>
  )
}
