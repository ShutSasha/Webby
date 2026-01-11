import MediaHeader from '@/ui/components/common/MediaHeader'
import MainLayout from '@/ui/components/MainLayout'

export default function VideosPage() {
  return (
    <MainLayout>
      <div className="w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <MediaHeader searchPlaceholder="Search videos..." actionLabel="Create a video" />
      </div>
    </MainLayout>
  )
}
