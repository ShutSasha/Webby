import { ReactNode } from 'react'

type Props = {
  children: ReactNode
  className?: string
}

export default function MainContainer({ children, className }: Props) {
  return (
    <div className={`w-full bg-neutral-900 rounded-[20px] p-5 box-border flex flex-col gap-4 ${className}`}>
      {children}
    </div>
  )
}
