import { testVideoUrls } from '@/lib/placeholder-data/player'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'
import RoomInteractionContainer from '@/ui/components/modules/Rooms/InteractionBlock/RoomInteractionContainer'

export default function RoomPage() {
  const testVideoUrl =
    testVideoUrls.find(it => it.type === 'mp4')?.link ?? 'https://www.youtube.com/watch?v=Zmrj90wYt4c'

  return (
    <div className="flex gap-5">
      <CustomPlayer videoUrl={testVideoUrl} />
      <RoomInteractionContainer />
    </div>
  )
}
