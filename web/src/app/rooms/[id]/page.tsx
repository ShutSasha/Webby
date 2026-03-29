import { testVideoUrls } from '@/lib/placeholder-data/player'
import MainLayout from '@/ui/components/layouts/MainLayout'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import RoomInteractionContainer from '@/ui/components/modules/Rooms/InteractionBlock/RoomInteractionContainer'

export default function RoomPage() {
  const testVideoUrl =
    testVideoUrls.find(it => it.type === 'mp4')?.link ?? 'https://www.youtube.com/watch?v=Zmrj90wYt4c'

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <div className="flex gap-5">
          <CustomPlayer videoUrl={testVideoUrl} />
          <RoomInteractionContainer />
        </div>
      </div>
    </MainLayout>
  )
}
