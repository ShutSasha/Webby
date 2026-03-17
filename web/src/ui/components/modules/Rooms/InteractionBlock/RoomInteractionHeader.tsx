'use client'

import { FC, SVGProps } from 'react'

import PlaylistIcon from '@/assets/icons/ic_playlist.svg'
import UsersIcon from '@/assets/icons/ic_users.svg'
import SettingsIcon from '@/assets/icons/shared/settings.svg'
import XIcon from '@/assets/icons/shared/x.svg'
import { cn } from '@/lib/utils/utils'
import { TabType, useRoomStore } from '@/stores/room.store'

export default function RoomInteractionHeader() {
  const tab = useRoomStore(state => state.tab)
  const setTab = useRoomStore(state => state.setTab)

  const toggleTab = (targetTab: TabType) => {
    if (tab === targetTab) {
      setTab('chat')
    } else {
      setTab(targetTab)
    }
  }

  const TabButton = ({ type, Icon }: { type: TabType; Icon: FC<SVGProps<SVGSVGElement>> }) => {
    const isActive = tab === type
    const CurrentIcon = isActive ? XIcon : Icon

    return (
      <CurrentIcon
        onClick={() => toggleTab(type)}
        className={cn(
          'size-6 cursor-pointer transition-all duration-300 ease-in-out select-none',
          isActive ? 'text-emerald-500' : 'text-neutral-300 hover:text-emerald-500',
        )}
      />
    )
  }

  return (
    <div className="border-b border-emerald-500 rounded-lg flex flex-row items-center justify-between py-2 px-3 mb-2">
      <TabButton type="playlist" Icon={PlaylistIcon} />

      <div className="flex flex-row items-center gap-4">
        <TabButton type="users" Icon={UsersIcon} />
        <TabButton type="settings" Icon={SettingsIcon} />
      </div>
    </div>
  )
}
