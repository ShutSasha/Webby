import MainLayout from '@/ui/components/layouts/MainLayout'
import VideoPageHeader from '@/ui/components/modules/Videos/VideoPageHeader'
import VideosPublicSearchContainer from '@/ui/components/modules/Videos/VideosPublicSearchContainer'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function VideosPage({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''

  return (
    <MainLayout>
      <div className="flex flex-col w-full bg-neutral-900 rounded-[20px] p-5 gap-4 box-border">
        <VideoPageHeader />

        <VideosPublicSearchContainer query={safeQuery} />
      </div>
    </MainLayout>
  )
}
