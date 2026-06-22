'use client'

import { useParams } from 'next/navigation'
import { useSession } from 'next-auth/react'

import { useIsClient } from '@/lib/hooks/useIsClient'
import { useRoomStore } from '@/stores/room.store'

import ActiveVoteOverlay from './ActiveVoteOverlay'
import RoomChatContainer from './Chat/RoomChatContainer'
import RoomQueue from './Queue/RoomQueue'
import RoomInteractionHeader from './RoomInteractionHeader'
import RoomInteractionSkeleton from './RoomInteractionSkeleton'
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
  const isHost = session?.user.id === hostId

  if (!isClient) {
    return (
      <div
        className="flex items-center justify-center flex-none bg-black rounded-2xl py-1 px-2 w-full h-[500px]
          xl:w-[300px] 2xl:w-[340px] xl:h-auto"
      >
        <RoomInteractionSkeleton />
      </div>
    )
  }

  return (
    <div
      className="flex flex-col flex-none bg-black rounded-2xl pt-1 pb-3 px-2 relative w-full h-[500px] xl:w-[300px]
        2xl:w-[340px] xl:h-0 xl:min-h-full"
    >
      <RoomInteractionHeader isHost={isHost} />

      <ActiveVoteOverlay roomId={roomId} />

      {tab === 'queue' && <RoomQueue />}
      {tab === 'chat' && <RoomChatContainer chatId={chatId} />}
      {tab === 'users' && <RoomUsers />}
      {tab === 'settings' && <RoomSettings />}

      <RoomVotesModal roomId={roomId} isHost={isHost} />
    </div>
  )
}
