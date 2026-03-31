import { Suspense } from 'react'

import PlayerSkeleton from '@/ui/components/modules/Playlists/PlayerSkeleton'
import VideoDetails from '@/ui/components/modules/Playlists/VideoDetails'
import AsideVideoCard from '@/ui/components/modules/Videos/AsideVideoCard'
import { auth } from '@/workspace/auth'

type Props = {
  params: Promise<{ id: string }>
}

export default async function VideoPage({ params }: Props) {
  const { id } = await params
  const session = await auth()

  return (
    <div className="flex gap-5">
      <Suspense key={id} fallback={<PlayerSkeleton />}>
        <VideoDetails v={id} currentUserId={session?.user.id} />
      </Suspense>

      <div className="w-[418px] flex flex-col gap-3">
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
