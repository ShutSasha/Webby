'use client'

import SyncIcon from '@/assets/icons/Player/sync.svg'
import { cn } from '@/lib/utils/general.utils'

type Props = {
  className?: string
  // onClick?: () => void
}

export default function PlayerSyncButton({ className }: Props) {
  return (
    <button
      type="button"
      className={cn(
        'group/sync cursor-pointer p-1.5 -ml-1.5 rounded-full transition-colors hover:bg-white/10',
        className,
      )}
      title="Synchronize playback"
      // onClick={onClick}
    >
      <SyncIcon
        className="w-5 h-5 md:w-6 md:h-6 text-neutral-300 transition-colors duration-300
          group-hover/sync:text-emerald-500 stroke-[1.5px]"
      />
    </button>
  )
}
