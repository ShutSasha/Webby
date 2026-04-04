import Image from 'next/image'
import Link from 'next/link'

import { Video } from '@/lib/actions/video.actions'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  video: Pick<Video, 'videoId' | 'previewUrl'>
}

export default function AsideVideoCard({ video }: Props) {
  return (
    <Link href={`/videos/${video.videoId}`} className="group flex gap-3">
      <div className="overflow-hidden rounded-lg relative w-[168px] h-fit shrink-0">
        <Image
          src={video.previewUrl}
          width={400}
          height={400}
          alt=""
          className="aspect-video rounded-lg group-hover:scale-115 transition-transform duration-600"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral900']}
        />
        <div className="absolute inset-0 z-1 bg-black/15 transition-all duration-500 group-hover:bg-black/5" />
      </div>

      {/* description block */}
      <div className="flex flex-col gap-1">
        <p className="text-sm line-clamp-2" title="Survive 30 Days Stranded With Your Ex, Win $250,000">
          Survive 30 Days Stranded With Your Ex, Win $250,000
        </p>
        <p className="text-sm text-neutral-400" title="guuuntersteam">
          guuuntersteam
        </p>
        <p className="text-[12px] text-neutral-400 line-clamp-1" title="100 views">
          100 views | 20/12/2025
        </p>
      </div>
    </Link>
  )
}
