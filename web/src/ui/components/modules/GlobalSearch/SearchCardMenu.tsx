import { useEffect, useRef, useState } from 'react'

import { useParams, usePathname } from 'next/navigation'

import PlusIcon from '@/assets/icons/ic_plus_create.svg'
import { useAddQueueItemMutation } from '@/lib/hooks/api/room/useAddQueueItem'
import { EntityType } from '@/lib/utils/global-search-modal.utils'

type Props = {
  type: EntityType
  id: string
  onOpenPlaylistModal: () => void
}

export const SearchCardMenu = ({ type, id, onOpenPlaylistModal }: Props) => {
  const [isOpen, setIsOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  const params = useParams()
  const pathname = usePathname()

  const isRoomPage = pathname?.startsWith('/rooms/')
  const roomId = isRoomPage ? (params?.id as string | undefined) : undefined

  const { mutate: addQueueItem, isPending } = useAddQueueItemMutation()

  const canAddToPlaylist = ['Video', 'YouTube', 'Stream'].includes(type)
  const canAddToRoom = !!roomId && ['Video', 'YouTube', 'Playlist', 'Stream'].includes(type)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    if (isOpen) document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isOpen])

  const toggleMenu = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsOpen(prev => !prev)
  }

  const handleAddToPlaylist = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsOpen(false)
    onOpenPlaylistModal()
  }

  const handleAddToRoom = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    if (!roomId || isPending) return

    addQueueItem(
      { roomId, videoId: id },
      {
        onSuccess: () => {
          setIsOpen(false)
        },
      },
    )
  }

  if (!canAddToPlaylist && !canAddToRoom) {
    return null
  }

  return (
    <div className="relative shrink-0" ref={menuRef}>
      <button
        onClick={toggleMenu}
        className="p-1.5 rounded-full text-foreground-faint hover:text-emerald-500 hover:bg-emerald-500/10
          transition-colors opacity-0 group-hover:opacity-100 shrink-0"
      >
        <PlusIcon className="w-5 h-5 stroke-2" />
      </button>

      {isOpen && (
        <div
          className="absolute right-0 top-full mt-2 w-48 bg-neutral-800 border border-neutral-700/60 shadow-xl
            shadow-black/50 z-50 py-1.5 rounded-xl animate-in fade-in zoom-in-95 duration-200"
          onClick={e => e.preventDefault()}
        >
          {canAddToPlaylist && (
            <button
              onClick={handleAddToPlaylist}
              className="w-full text-left px-4 py-2 text-sm text-foreground-tertiary hover:bg-neutral-700/50
                transition-colors flex items-center gap-3 cursor-pointer"
            >
              <span>Add to playlist</span>
            </button>
          )}

          {canAddToRoom && (
            <button
              onClick={handleAddToRoom}
              disabled={isPending}
              className={`w-full text-left px-4 py-2 text-sm text-foreground-tertiary transition-colors flex
              items-center gap-3
              ${isPending ? 'opacity-50 cursor-not-allowed' : 'hover:bg-neutral-700/50 cursor-pointer'}`}
            >
              <span>{isPending ? 'Adding...' : 'Add to room'}</span>
            </button>
          )}
        </div>
      )}
    </div>
  )
}
