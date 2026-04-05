import { ReactNode } from 'react'

type Props = {
  children: ReactNode
}

export default function MainContainer({ children }: Props) {
  return <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">{children}</div>
}
