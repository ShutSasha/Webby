import PlaylistMediaVideoContainer from '@/ui/components/modules/Playlists/PlaylistMediaVideoContainer'
import VideoRecommendationsList from '@/ui/components/modules/Videos/VideoRecommendationsList'

type Props = { params: Promise<{ id: string }> }

export default async function VideoPage({ params }: Props) {
  const { id } = await params

  return (
    <div className="flex flex-col xl:flex-row gap-5">
      <PlaylistMediaVideoContainer videoId={id} />
      <div
        className="w-full xl:w-[360px] 2xl:w-[400px] grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-1 gap-3
          shrink-0 h-fit content-start"
      >
        <VideoRecommendationsList videoId={id} />
      </div>
    </div>
  )
}
