import { ReactNode } from 'react'

import Link from 'next/link'

type Props = {
  title: string
  description: string
  icon: ReactNode
}

export default function AuthPlaceholder({ title, description, icon }: Props) {
  return (
    <div className="flex flex-col flex-1 items-center justify-center py-32 px-4 animate-in fade-in duration-500">
      <div
        className="size-20 bg-neutral-900/80 rounded-full flex items-center justify-center mb-6 border
          border-neutral-800"
      >
        {icon}
      </div>

      <h2 className="text-2xl font-bold text-foreground-tertiary mb-3 text-center">{title}</h2>
      <p className="text-foreground-muted text-center max-w-md mb-8 leading-relaxed">{description}</p>

      <Link
        href="/login"
        className="px-8 py-3 bg-emerald-500 hover:bg-emerald-400 text-foreground-inverse-subtle font-bold rounded-xl
          transition-colors duration-300"
      >
        Sign In
      </Link>
    </div>
  )
}
