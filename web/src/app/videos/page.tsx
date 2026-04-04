import MainLayout from '@/ui/components/layouts/MainLayout'
import VideoCard from '@/ui/components/modules/Videos/VideoCard'
import VideoPageHeader from '@/ui/components/modules/Videos/VideoPageHeader'
import GridCardsContainer from '@/ui/components/shared/GridCardsContainer'

export default function VideosPage() {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <VideoPageHeader />

        {/* Room list goes here */}
        <GridCardsContainer>
          {Array.from({ length: 10 }).map((_, index) => (
            <VideoCard key={index} />
          ))}
        </GridCardsContainer>
      </div>
    </MainLayout>
  )
}
