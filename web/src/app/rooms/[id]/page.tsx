import { testVideoUrls } from '@/lib/placeholder-data/player'
import MainLayout from '@/ui/components/MainLayout'
import CustomPlayer from '@/ui/components/modules/Player/CustomPlayer'

export default function RoomPage() {
  const testVideoUrl =
    testVideoUrls.find(it => it.type === 'mp4')?.link ?? 'https://www.youtube.com/watch?v=Zmrj90wYt4c'

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <div className="flex gap-5">
          <CustomPlayer videoUrl={testVideoUrl} />

          <div className="hidden xl:block xl:w-[300px] 2xl:w-[340px] flex-none bg-amber-700 rounded-2xl" />
        </div>
        <p>some footer content</p>
        <p>some footer content</p>
        <p>some footer content</p>
        <p>some footer content</p>
      </div>
    </MainLayout>
  )
}
