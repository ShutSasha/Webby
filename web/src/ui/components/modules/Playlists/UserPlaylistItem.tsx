'use client'

import Link from 'next/link'

import ImageBackground from '../Profile/ImageBackground'

type Props = {
  id: string
  src: string
  name: string
  isPrivate: boolean
  videoCount: number
}

export default function UserPlaylistItem(props: Props) {
  return (
    <Link href={`/playlists/${props.id}`} className="relative group cursor-pointer">
      <ImageBackground src={props.src} />
      <p className="text-sm font-medium text-neutral-100 line-clamp-1">{props.name}</p>
      <p className="text-[12px] text-neutral-500">{props.isPrivate ? 'Private' : 'Public'} &bull; Playlist</p>
      <div
        className="absolute bg-black/80 rounded-lg px-2 py-1 top-1/35 right-1/40 group-hover:top-1/20
          group-hover:right-1/25 transition-all duration-300 text-neutral-200 text-[12px]"
      >
        {props.videoCount} videos
      </div>
    </Link>
  )
}
