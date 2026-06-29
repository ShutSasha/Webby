'use client'

import Link from 'next/link'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { cn } from '@/lib/utils/general.utils'
import { useVideoDraftStore } from '@/stores/video-draft.store'

export default function CreateVideoButton() {
  const clearDraft = useVideoDraftStore(state => state.clearDraft)

  const handleClick = () => {
    clearDraft()
  }

  return (
    <Link
      href={'/videos/create'}
      onClick={handleClick}
      className={cn(
        'flex items-center gap-2 shrink-0 cursor-pointer order-3',
        'bg-background text-foreground-tertiary',
        'text-sm font-semibold px-5 py-2.5 rounded-full',
        'transition-all duration-300 ease-out',
        'hover:bg-surface-tertiary/40',
        'shadow-none',
      )}
    >
      <PlusIcon className="h-4 w-4 text-foreground-tertiary" aria-hidden="true" />
      Create
    </Link>
  )
}
