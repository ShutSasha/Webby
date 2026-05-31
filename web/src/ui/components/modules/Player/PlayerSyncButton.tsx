'use client'

import { useParams } from 'next/navigation'

import SyncIcon from '@/assets/icons/Player/sync.svg'
import { initiateSyncAction } from '@/lib/actions/room.actions'
import { cn } from '@/lib/utils/general.utils'
import { useRoomStore } from '@/stores/room.store'

type Props = {
  className?: string
}

export default function PlayerSyncButton({ className }: Props) {
  const isSyncCooldown = useRoomStore(state => state.isSyncCooldown)
  const setSyncCooldown = useRoomStore(state => state.setSyncCooldown)
  const params = useParams()
  const roomId = params?.id as string | undefined

  const handleSyncClick = async () => {
    if (!roomId) {
      console.error('roomId is undefined')
      return
    }

    if (isSyncCooldown) return

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
        'group/sync p-1.5 -ml-1.5 rounded-full transition-all duration-300',
        'active:scale-90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50',
        isSyncCooldown || !roomId ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer hover:bg-white/10',
        className,
      )}
      title={isSyncCooldown ? 'Sync is on cooldown' : 'Synchronize playback'}
    >
      <SyncIcon
        className={cn(
          'w-5 h-5 md:w-6 md:h-6 transition-colors duration-300 stroke-[1.5px]',
          isSyncCooldown ? 'text-amber-500' : 'text-neutral-300 group-hover/sync:text-emerald-500',
        )}
      />
    </button>
  )
}
