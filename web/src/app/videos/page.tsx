import VideoPageHeader from '@/ui/components/modules/Videos/VideoPageHeader'
import VideosPublicSearchContainer from '@/ui/components/modules/Videos/VideosPublicSearchContainer'

type Props = {
  searchParams: Promise<{ query?: string }>
}

export default async function VideosPage({ searchParams }: Props) {
  const { query } = await searchParams
  const safeQuery = query || ''

  return (
    <>
      <VideoPageHeader />
      <VideosPublicSearchContainer query={safeQuery} />
    </>
  )
}
