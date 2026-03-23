import { getVideoInfo } from '@/app/api/videos'

import CustomPlayer from '../Player/CustomPlayer'

type Props = {
  videoId: string
}

export default async function PlayerContainer({ videoId }: Props) {
  const video = await getVideoInfo(videoId)

  if (!video.success || !video.data) {
    return <div>video not found</div>
  }

  return <CustomPlayer videoUrl={video.data.videoUrl} />
}
