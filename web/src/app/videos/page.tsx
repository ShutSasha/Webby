import MediaHeader from '@/ui/components/common/MediaHeader'
import MainLayout from '@/ui/components/MainLayout'
import CreateVideoButton from '@/ui/components/modules/Videos/CreateVideoButton'
import VideoCard from '@/ui/components/modules/Videos/VideoCard'

export default function VideosPage() {
  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <MediaHeader searchPlaceholder="Search a video" actionSlot={<CreateVideoButton />} />

        <div className="flex justify-between">
          <div className="hidden lg:block">
            <div className="flex gap-3 items-center">
              <span className="w-10 h-0.5 bg-emerald-500 rounded-full" />
              <p className="text-emerald-500 font-black text-[10px] leading-3.5 uppercase">live now</p>
            </div>
            <h2 className="font-black text-[48px] leading-14 text-emerald-500 uppercase">videos</h2>
          </div>
        </div>

        {/* Room list goes here */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3">
          {Array.from({ length: 10 }).map((_, index) => (
            <VideoCard key={index} />
          ))}
        </div>
      </div>
    </MainLayout>
  )
}
