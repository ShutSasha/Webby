'use client'

import { useParams } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { useIsClient } from '@/lib/hooks/useIsClient'
import { useRoomStore } from '@/stores/room.store'
import PageLoading from '@/ui/components/shared/PageLoading'

import ActiveVoteOverlay from './ActiveVoteOverlay'
import RoomChatContainer from './Chat/RoomChatContainer'
import RoomPlaylists from './Playlists/RoomPlaylists'
import RoomInteractionHeader from './RoomInteractionHeader'
import RoomSettings from './Settings/RoomSettings'
import RoomUsers from './UsersList/RoomUsers'
import RoomVotesModal from './Votes/RoomVotesModal'

type Props = {
  chatId: string
  hostId: string
}

export default function RoomInteractionContainer({ chatId, hostId }: Props) {
  const params = useParams()
  const roomId = params?.id as string
  const tab = useRoomStore(state => state.tab)
  const isClient = useIsClient()
  const { data: session } = useSession()

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
        min-h-full relative"
    >
      <RoomInteractionHeader />

      <ActiveVoteOverlay roomId={roomId} />

      {tab === 'chat' && <RoomChatContainer chatId={chatId} />}
      {tab === 'playlist' && <RoomPlaylists />}
      {tab === 'users' && <RoomUsers />}
      {tab === 'settings' && <RoomSettings />}

      <RoomVotesModal roomId={roomId} isHost={session?.user.id === hostId} />
    </div>
  )
}
