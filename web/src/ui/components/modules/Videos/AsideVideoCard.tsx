import Link from 'next/link'

import { formatTimeAgo } from '@/lib/utils/date.utils'
import { formatViews } from '@/lib/utils/video.utils'
import { RecommendedVideo } from '@/types/video.types'
import SafeImage from '@/ui/components/shared/SafeImage'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  video: RecommendedVideo
}

export default function AsideVideoCard({ video }: Props) {
  return (
    <Link href={`/videos/${video.videoId}`} className="group flex flex-col xl:flex-row gap-3">
      <div className="overflow-hidden rounded-lg relative xl:w-[168px] h-fit shrink-0 aspect-video">
        <SafeImage
          src={video.previewUrl}
          fallbackType="video"
          width={400}
          height={400}
          alt={video.name}
          className="aspect-video rounded-lg group-hover:scale-110 transition-transform duration-600 w-full
            object-cover"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral900']}
        />
        <div className="absolute inset-0 z-1 bg-black/15 transition-all duration-500 group-hover:bg-black/5" />
      </div>

      <div className="flex flex-col gap-1 flex-1 overflow-hidden">
        <p className="text-sm line-clamp-2 text-neutral-200 font-medium" title={video.name}>
          {video.name}
        </p>
        <p className="text-sm text-neutral-400 truncate" title={video.user.username}>
          {video.user.username}
        </p>
        <p className="text-[12px] text-neutral-400 line-clamp-1" title={`${video.views} views`}>
          {formatViews(video.views)} views • {formatTimeAgo(video.createdAt)}
        </p>
      </div>
    </Link>
  )
}

export function AsideVideoCardSkeleton() {
  return (
    <div className="flex flex-col xl:flex-row gap-3 animate-pulse">
      <div className="xl:w-[168px] aspect-video bg-neutral-800/80 rounded-lg shrink-0" />
      <div className="flex flex-col gap-2 flex-1 pt-1">
        <div className="h-3.5 bg-neutral-800/80 rounded w-[90%]" />
        <div className="h-3.5 bg-neutral-800/80 rounded w-[60%]" />
        <div className="h-2.5 bg-neutral-800/60 rounded w-[40%] mt-1" />
      </div>
    </div>
  )
}
