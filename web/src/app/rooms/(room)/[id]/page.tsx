import ShareIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import { testVideoUrls } from '@/lib/placeholder-data/player'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import RoomInteractionContainer from '@/ui/components/modules/Rooms/InteractionBlock/RoomInteractionContainer'
import ActionButton from '@/ui/components/shared/ActionButton'

export default function RoomPage() {
  const testVideoUrl =
    testVideoUrls.find(it => it.type === 'mp4')?.link ?? 'https://www.youtube.com/watch?v=Zmrj90wYt4c'

  return (
    <div className="flex flex-col">
      <div className="flex gap-5">
        <CustomPlayer videoUrl={testVideoUrl} />
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
