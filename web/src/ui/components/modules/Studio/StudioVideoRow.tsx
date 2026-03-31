'use client'

import Image from 'next/image'
import Link from 'next/link'

import MoreVertical from '@/assets/icons/shared/more-vertical.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  id: string
}

export function StudioVideoRow({ id }: Props) {
  const isPrivate = Number(id) % 2 === 0

  return (
    <div
      className="flex items-center border-b border-neutral-800 py-3 px-4 hover:bg-neutral-800/40 transition-colors
        group"
    >
      <div className="flex-1 flex gap-4 min-w-[300px]">
        <div className="relative w-32 aspect-video bg-neutral-800 rounded-lg overflow-hidden shrink-0">
          <Image
            src="https://cdn.magicdecor.in/com/2023/10/20174720/Anime-Scenery-Wallpaper-for-Walls-710x488.jpg"
            alt="Video thumbnail"
            fill
            className="object-cover"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral900']}
          />
          <div className="absolute bottom-1 right-1 bg-black/80 px-1 py-0.5 rounded text-[10px] font-medium text-white">
            14:20
          </div>
        </div>

        <div className="flex flex-col justify-center overflow-hidden">
          <Link
            href={`/videos/${id}`}
            className="text-sm font-semibold text-neutral-100 line-clamp-2 hover:text-emerald-400 transition-colors"
          >
            Amazing Anime Moon Landscape Speedart - Full Process {id}
          </Link>
          <p className="text-xs text-neutral-500 mt-1 line-clamp-1">Add description</p>
        </div>
      </div>

      <div className="w-28 flex justify-center shrink-0">
        <div className="flex items-center gap-1.5 text-xs text-neutral-300">
          {isPrivate ? (
            <>
              <span>Private</span>
            </>
          ) : (
            <>
              <span className="text-emerald-500">Public</span>
            </>
          )}
        </div>
      </div>

      <div className="w-32 flex flex-col items-center justify-center shrink-0">
        <p className="text-xs text-neutral-200">23 Mar 2026</p>
        <p className="text-[10px] text-neutral-500 mt-0.5">Uploaded</p>
      </div>

      <div className="w-24 text-center text-xs text-neutral-300 shrink-0">1.2K</div>

      <div className="w-12 flex justify-end shrink-0">
        <button
          className="p-1.5 rounded-full hover:bg-neutral-700/50 text-neutral-400 hover:text-neutral-100 opacity-0
            group-hover:opacity-100 transition-all cursor-pointer"
        >
          <MoreVertical className="size-5" />
        </button>
      </div>
    </div>
  )
}
