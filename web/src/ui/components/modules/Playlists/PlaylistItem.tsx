'use client'

import Link from 'next/link'

import { PlaylistData } from '@/lib/placeholder-data/profile'

import ImageBackground from '../Profile/ImageBackground'

export default function PlaylistItem(props: PlaylistData) {
  return (
    <Link href={`/playlists/${props.id}`} className="relative group cursor-pointer">
      <ImageBackground src={props.src} />
      <p className="text-sm font-medium">{props.name}</p>
      <p className="text-[12px] text-neutral-500">{props.creator}</p>
      <div
        className="absolute bg-neutral-800 rounded-lg px-2 py-1 top-1/35 right-1/40 text-[12px] group-hover:top-1/20
          group-hover:right-1/25 transition-all duration-300 border border-border"
      >
        {props.videoCount} videos
      </div>
    </Link>
  )
}
