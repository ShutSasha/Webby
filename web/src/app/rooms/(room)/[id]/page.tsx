import { getRoomById } from '@/lib/actions/room.actions'
import PlayerSyncButton from '@/ui/components/modules/Player/PlayerSyncButton'
import RoomInteractionContainer from '@/ui/components/modules/Rooms/InteractionBlock/RoomInteractionContainer'
import RoomInviteButton from '@/ui/components/modules/Rooms/RoomInviteButton'
import RoomPlayerContainer from '@/ui/components/modules/Rooms/RoomPlayerContainer'
import RoomWebSocketManager from '@/ui/components/modules/Rooms/RoomWebSocketManager'
import EmptyState from '@/ui/components/shared/EmptyState'

type Props = {
  params: Promise<{ id: string }>
}

export default async function RoomPage({ params }: Props) {
  const { id } = await params

  const roomResponse = await getRoomById(id)
  const room = roomResponse.data

  if (!room) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState title="Room not found" description="This room doesn't exist, is private, or has been deleted." />
      </div>
    )
  }

  return (
    <div className="flex flex-col">
      <RoomWebSocketManager roomId={id} chatId={room.chatId} />
      <div className="flex gap-5">
        <RoomPlayerContainer roomId={id} />

        <RoomInteractionContainer chatId={room.chatId} />
      </div>
      <div className="flex gap-5">
        <div className="w-full flex justify-between items-center mt-3 mb-2">
          <p className="text-neutral-300 text-[20px] font-bold">{room.name}</p>

          <div className="flex items-center gap-2">
            <PlayerSyncButton roomId={id} />

            <RoomInviteButton roomId={id} />
          </div>
        </div>
        <div className="hidden xl:block xl:w-[300px] 2xl:w-[340px] flex-none"></div>
      </div>
    </div>
  )
}
