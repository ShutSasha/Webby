'use client'

import Link from 'next/link'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { cn } from '@/lib/utils/general.utils'

export default function CreateVideoButton() {
  return (
    <Link
      href={'/videos/create'}
      className={cn(
        'flex items-center gap-2 shrink-0 cursor-pointer order-3',
        'bg-neutral-800 text-neutral-200',
        'text-sm font-semibold px-5 py-2.5 rounded-full',
        'transition-all duration-300 ease-out',
        'hover:bg-neutral-700/40',
        'shadow-none',
      )}
    >
      <PlusIcon className="h-4 w-4 text-neutral-200" aria-hidden="true" />
      Create
    </Link>
  )
}
