import LockIcon from '@/assets/icons/shared/lock.svg'
import { getRoomById } from '@/lib/actions/room.actions'
import PlayerSyncButton from '@/ui/components/modules/Player/PlayerSyncButton'
import RoomInteractionContainer from '@/ui/components/modules/Rooms/InteractionBlock/RoomInteractionContainer'
import RoomInviteButton from '@/ui/components/modules/Rooms/RoomInviteButton'
import RoomPlayerContainer from '@/ui/components/modules/Rooms/RoomPlayerContainer'
import RoomWebSocketManager from '@/ui/components/modules/Rooms/RoomWebSocketManager'
import AuthPlaceholder from '@/ui/components/shared/AuthPlaceholder'
import EmptyState from '@/ui/components/shared/EmptyState'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function RoomPage({ params }: Props) {
  const { id } = await params
  const session = await auth()
  const roomResponse = await getRoomById(id)
  const room = roomResponse.data

  if (!session) {
    return (
      <AuthPlaceholder
        title="Sign in to view rooms"
        description="Manage your custom rooms, adjust privacy settings, and host synchronized viewing sessions by logging into your account."
        icon={<LockIcon className="size-10 text-neutral-500 stroke-1" />}
      />
    )
  }

  if (!room) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState title="Room not found" description="This room doesn't exist, is private, or has been deleted." />
      </div>
    )
  }

  const hostId = room.hostId

  return (
    <div className="flex flex-col">
      <RoomWebSocketManager roomId={id} chatId={room.chatId} />

      <div className="flex flex-col xl:flex-row gap-5 w-full">
        <div className="flex flex-col flex-1 min-w-0">
          <RoomPlayerContainer roomId={id} />

          <div className="w-full flex flex-col sm:flex-row sm:justify-between sm:items-center mt-3 mb-2 gap-4">
            <p className="text-neutral-300 text-[20px] font-bold truncate">{room.name}</p>

            <div className="flex items-center gap-2 shrink-0">
              <PlayerSyncButton roomId={id} />
              <RoomInviteButton roomId={id} />
            </div>
          </div>
        </div>

        <RoomInteractionContainer chatId={room.chatId} hostId={hostId} />
      </div>
    </div>
  )
}
