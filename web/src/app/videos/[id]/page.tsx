import StandaloneVideoContainer from '@/ui/components/modules/Videos/StandaloneVideoContainer'
import VideoRecommendationsList from '@/ui/components/modules/Videos/VideoRecommendationsList'
import { auth } from '@/workspace/auth'

type Props = { params: Promise<{ id: string }> }

export default async function VideoPage({ params }: Props) {
  const { id } = await params
  const session = await auth()

  return (
    <div className="flex flex-col xl:flex-row gap-5">
      <StandaloneVideoContainer videoId={id} currentUserId={session?.user?.id} />
      <div
        className="w-full xl:w-[360px] 2xl:w-[400px] grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-1 gap-3
          shrink-0 h-fit content-start"
      >
        <VideoRecommendationsList videoId={id} />
      </div>
    </div>
  )
}
