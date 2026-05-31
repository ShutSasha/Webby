import ShareIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import { getRoomById } from '@/lib/actions/room.actions'
import { testVideoUrls } from '@/lib/placeholder-data/player'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import RoomInteractionContainer from '@/ui/components/modules/Rooms/InteractionBlock/RoomInteractionContainer'
import RoomWebSocketManager from '@/ui/components/modules/Rooms/RoomWebSocketManager'
import ActionButton from '@/ui/components/shared/ActionButton'
import EmptyState from '@/ui/components/shared/EmptyState'

type Props = {
  params: Promise<{ id: string }>
}

export default async function RoomPage({ params }: Props) {
  const { id } = await params

  const roomResponse = await getRoomById(id)
  const room = roomResponse.data

  const testVideoUrl =
    testVideoUrls.find(it => it.type === 'mp4')?.link ?? 'https://www.youtube.com/watch?v=Zmrj90wYt4c'

  if (!room) {
    return (
      <div className="flex-1 flex items-center justify-center bg-neutral-900/20 rounded-[20px]">
        <EmptyState title="Room not found" description="This room doesn't exist, is private, or has been deleted." />
      </div>
    )
  }

  return (
    <div className="flex flex-col">
      <RoomWebSocketManager chatId={room.chatId} />
      <div className="flex gap-5">
        <CustomPlayer videoUrl={testVideoUrl} isRoom roomId={id} />

        <RoomInteractionContainer />
      </div>
      <div className="flex gap-5">
        <div className="w-full flex justify-between items-center mt-3 mb-2">
          <p className="text-neutral-300 text-[20px] font-bold">Room name</p>
          <ActionButton label="send an invite">
            <ShareIcon className="size-4 stroke-[1.5px]" />
          </ActionButton>
        </div>
        <div className="hidden xl:block xl:w-[300px] 2xl:w-[340px] flex-none"></div>
      </div>
    </div>
  )
}
