'use client'

import Link from 'next/link'

import ImageBackground from '../Profile/ImageBackground'

type Props = {
  id: string
  src: string
  name: string
  creator: string
  videoCount: number
}

export default function PlaylistItem(props: Props) {
  return (
    <Link href={`/playlists/${props.id}`} className="relative group cursor-pointer">
      <ImageBackground src={props.src} />
      <p className="text-sm font-medium text-foreground-secondary line-clamp-1">{props.name}</p>
      <p className="text-[12px] text-foreground-muted">{props.creator}</p>
      <div
        className="absolute bg-black/60 backdrop-blur-sm rounded-lg px-2 py-1 top-1/35 right-1/40 group-hover:top-1/20
          group-hover:right-1/25 transition-all duration-300 text-neutral-100 text-[12px]"
      >
        {props.videoCount} videos
      </div>
    </Link>
  )
}

export function PlaylistItemSkeleton() {
  return (
    <div className="relative animate-pulse w-full">
      <div className="w-full rounded-2xl mb-1 aspect-video bg-background/50" />

      <div className="h-4 bg-background/50 rounded-md w-3/4 mt-1.5 mb-1.5" />

      <div className="h-3 bg-background/50 rounded-md w-1/2" />

      <div className="absolute bg-surface/80 rounded-lg w-14 h-6 top-1/35 right-1/40" />
    </div>
  )
}
