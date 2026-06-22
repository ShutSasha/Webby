'use client'

import { FC, SVGProps } from 'react'

import UsersIcon from '@/assets/icons/ic_users.svg'
import ChatsIcon from '@/assets/icons/Nav/chats.svg'
import PlaylistIcon from '@/assets/icons/Nav/list-video.svg'
import SettingsIcon from '@/assets/icons/shared/settings.svg'
import { cn } from '@/lib/utils/general.utils'
import { TabType, useRoomStore } from '@/stores/room.store'

type Props = {
  isHost: boolean
}

export default function RoomInteractionHeader({ isHost }: Props) {
  const tab = useRoomStore(state => state.tab)
  const setTab = useRoomStore(state => state.setTab)

  const toggleTab = (targetTab: TabType) => {
    setTab(targetTab)
  }

  const TabButton = ({ type, Icon }: { type: TabType; Icon: FC<SVGProps<SVGSVGElement>> }) => {
    const isActive = tab === type

    return (
      <Icon
        data-testid={`tab-${type}`}
        onClick={() => toggleTab(type)}
        className={cn(
          'size-6 cursor-pointer transition-all duration-300 ease-in-out select-none stroke-[1.5px]',
          isActive ? 'text-emerald-500' : 'text-neutral-300 hover:text-emerald-500',
        )}
      />
    )
  }

  return (
    <div className="border-b border-neutral-700/80 flex flex-row items-center justify-between py-2 px-3 mb-2">
      <TabButton type="queue" Icon={PlaylistIcon} />

      <div className="flex flex-row items-center gap-4">
        <TabButton type="chat" Icon={ChatsIcon} />
        <TabButton type="users" Icon={UsersIcon} />
        {isHost && <TabButton type="settings" Icon={SettingsIcon} />}
      </div>
    </div>
  )
}
