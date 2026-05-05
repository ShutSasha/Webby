import { Suspense } from 'react'

import { VideoSource } from '@/types/video.types'
import VideoDetails from '@/ui/components/modules/Playlists/VideoDetails'
import VideoDetailsSkeleton from '@/ui/components/modules/Playlists/VideoDetailsSkeleton'
import AsideVideoCard from '@/ui/components/modules/Videos/AsideVideoCard'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ source?: VideoSource }>
}

export default async function VideoPage({ params, searchParams }: Props) {
  const { id } = await params
  const { source } = await searchParams
  const session = await auth()

  return (
    <div className="flex flex-col xl:flex-row gap-5">
      <Suspense key={id} fallback={<VideoDetailsSkeleton />}>
        <VideoDetails v={id} currentUserId={session?.user.id} source={source} />
      </Suspense>

      <div
        className="w-full xl:w-[360px] 2xl:w-[400px] grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-1 gap-3
          shrink-0"
      >
        {[...new Array(20)].map((_, index) => (
          <AsideVideoCard
            key={index}
            video={{
              videoId: 'e68b32e3-473f-40d8-ad40-b5a64a9436d8',
              previewUrl:
                'https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/previews/e68b32e3-473f-40d8-ad40-b5a64a9436d8/14:26:38videoframe_3294840.png',
            }}
          />
        ))}
      </div>
    </div>
  )
}
