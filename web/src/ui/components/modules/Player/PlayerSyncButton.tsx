'use client'

import SyncIcon from '@/assets/icons/Player/sync.svg'
import { initiateSyncAction } from '@/lib/actions/room.actions'
import { cn } from '@/lib/utils/general.utils'
import { useRoomStore } from '@/stores/room.store'

type Props = {
  roomId: string
  className?: string
}

export default function PlayerSyncButton({ roomId, className }: Props) {
  const isSyncCooldown = useRoomStore(state => state.isSyncCooldown)
  const setSyncCooldown = useRoomStore(state => state.setSyncCooldown)

  const handleSyncClick = async () => {
    if (!roomId || isSyncCooldown) return

    setSyncCooldown(true)
    initiateSyncAction(roomId)

    setTimeout(() => {
      setSyncCooldown(false)
    }, 1100)
  }

  return (
    <button
      type="button"
      onClick={handleSyncClick}
      disabled={isSyncCooldown || !roomId}
      className={cn(
        'group flex items-center gap-2 shrink-0 px-3 py-2 md:px-4 md:py-2 rounded-full',
        'transition-all duration-300 ease-out text-nowrap',
        'bg-transparent text-foreground-muted',
        !isSyncCooldown && 'hover:bg-neutral-800 hover:text-foreground-secondary cursor-pointer',
        isSyncCooldown && 'opacity-60 cursor-not-allowed',
        className,
      )}
      title={isSyncCooldown ? 'Sync is on cooldown' : 'Synchronize playback'}
    >
      <SyncIcon
        className={cn(
          'size-4 stroke-[1.5px] transition-all duration-300',
          isSyncCooldown && 'animate-spin text-emerald-500',
        )}
      />
      <span className="text-sm font-medium leading-none">{isSyncCooldown ? 'syncing...' : 'sync video'}</span>
    </button>
  )
}
