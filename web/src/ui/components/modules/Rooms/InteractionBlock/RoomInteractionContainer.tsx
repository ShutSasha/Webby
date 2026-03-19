'use client'

import { useIsClient } from '@/lib/hooks/useIsClient'
import { useRoomStore } from '@/stores/room.store'
import PageLoading from '@/ui/components/shared/PageLoading'

import RoomChatContainer from './Chat/RoomChatContainer'
import RoomInteractionHeader from './RoomInteractionHeader'
import RoomPlaylists from './RoomPlaylists'
import RoomSettings from './Settings/RoomSettings'
import RoomUsers from './RoomUsers'

export default function RoomInteractionContainer() {
  const tab = useRoomStore(state => state.tab)
  const isClient = useIsClient()

  if (!isClient) {
    return (
      <div
        className="hidden xl:block xl:w-[300px] 2xl:w-[340px] flex-none bg-black rounded-2xl py-1 px-2 items-center
          justify-center"
      >
        <PageLoading />
      </div>
    )
  }

  return (
    <div
      className="hidden xl:flex xl:flex-col xl:w-[300px] 2xl:w-[340px] flex-none bg-black rounded-2xl pt-1 pb-3 px-2 h-0
        min-h-full"
    >
      <RoomInteractionHeader />

      {tab === 'chat' && <RoomChatContainer />}
      {tab === 'playlist' && <RoomPlaylists />}
      {tab === 'users' && <RoomUsers />}
      {tab === 'settings' && <RoomSettings />}
    </div>
  )
}
