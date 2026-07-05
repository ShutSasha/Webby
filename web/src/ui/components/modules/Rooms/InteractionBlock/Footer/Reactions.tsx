'use client'

import ReactionSmileIcon from '@/assets/icons/Room/reaction-smile.svg'
import { useRoomStore } from '@/stores/room.store'

import InteractionButton from './InteractionButton'

export default function Reactions() {
  const isOpen = useRoomStore(state => state.isReactionsModalOpen)
  const setIsOpen = useRoomStore(state => state.setReactionsModalOpen)

  return (
    <InteractionButton
      onClick={() => setIsOpen(!isOpen)}
      icon={
        <ReactionSmileIcon
          className="size-5 text-foreground-subtle dark:group-hover:text-foreground-inverse-subtle transition-colors
            duration-300 ease-in-out"
        />
      }
    />
  )
}
