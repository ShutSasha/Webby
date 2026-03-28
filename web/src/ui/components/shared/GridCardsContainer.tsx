import { ReactNode } from 'react'

type Props = {
  children: ReactNode
}

export default function GridCardsContainer({ children }: Props) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3">
      {children}
    </div>
  )
}
