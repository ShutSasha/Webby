import Image from 'next/image'
import Link from 'next/link'

import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  id: string
  name: string
  thumbnail: string
  category: string
  hostUsername: string
  hostAvatarUrl: string
}

export default function RoomCard({ id, name, thumbnail, category, hostUsername, hostAvatarUrl }: Props) {
  return (
    <Link href={`/rooms/${id}`} className="flex flex-col gap-3 group cursor-pointer">
      <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-card">
        <Image
          src={thumbnail}
          alt="Room preview"
          width={640}
          height={480}
          loading="lazy"
          className="object-cover z-0 transition-transform duration-500 ease-out group-hover:scale-105"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
        />

        <div className="absolute top-2 left-2 z-10 bg-black/60 backdrop-blur-md px-2 py-1 rounded-md">
          <p className="uppercase font-bold text-[10px] tracking-wider text-neutral-200 line-clamp-1 max-w-[100px]">
            {category}
          </p>
        </div>
      </div>

      <div className="flex gap-3 items-start px-1">
        <Image
          src={hostAvatarUrl}
          alt="User avatar"
          width={80}
          height={80}
          className="size-10 rounded-full object-cover shrink-0 mt-0.5"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral700']}
        />

        <div className="flex flex-col overflow-hidden">
          <h3
            className="text-foreground-strong text-sm font-semibold leading-snug line-clamp-2 transition-colors
              duration-200 group-hover:text-emerald-400"
          >
            {name}
          </h3>

          <p className="text-muted-active text-xs mt-1 truncate hover:text-foreground transition-colors">
            {hostUsername}
          </p>
        </div>
      </div>
    </Link>
  )
}

export function RoomCardSkeleton() {
  return (
    <div className="flex flex-col gap-3 w-full">
      {/* Thumbnail Skeleton */}
      <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-card/80 animate-pulse">
        {/* Category Badge Skeleton */}
        <div className="absolute top-2 left-2 z-10 bg-neutral-700/50 h-5 w-12 rounded-md" />
      </div>

      <div className="flex gap-3 items-start px-1">
        {/* Avatar Skeleton */}
        <div className="size-9 rounded-full bg-card/80 shrink-0 mt-0.5 animate-pulse" />

        <div className="flex flex-1 flex-col mt-0.5">
          <div className="flex flex-col gap-1.5">
            <div className="h-3.5 bg-card/80 rounded w-[90%] animate-pulse" />
            <div className="h-3.5 bg-card/80 rounded w-[70%] animate-pulse" />
          </div>

          {/* Username Skeleton */}
          <div className="h-2.5 bg-card/60 rounded w-[40%] animate-pulse mt-2" />
        </div>
      </div>
    </div>
  )
}
