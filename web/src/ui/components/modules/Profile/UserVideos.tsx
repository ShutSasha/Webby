import { videos } from '@/lib/placeholder-data/profile'

import ImageBackground from './ImageBackground'

export default async function UserVideos() {
  // await, sync videos - Lates/Popular
  await new Promise(r => setTimeout(r, 300))

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {videos.map(video => (
        <div key={video.id} className="group cursor-pointer">
          <ImageBackground src={video.src} />
          <p className="text-sm font-medium">{video.title}</p>
          <div className="flex gap-2 text-[12px] text-neutral-500">
            <p>{video.views}</p>
            <p>{video.date}</p>
          </div>
        </div>
      ))}
    </div>
  )
}

export function UserVideosSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {Array.from({ length: 12 }).map((_, index) => (
        <div key={index}>
          <div className="w-full aspect-video rounded-2xl mb-1 bg-neutral-800 animate-pulse" />

          {/* Title */}
          <div className="h-4 w-3/4 bg-neutral-800 rounded-md mb-1.5 animate-pulse mt-1" />

          {/* Meta info (Views and Date) */}
          <div className="flex gap-2 items-center">
            <div className="h-3 w-16 bg-neutral-800 rounded-md animate-pulse" />
            <div className="h-3 w-20 bg-neutral-800 rounded-md animate-pulse" />
          </div>
        </div>
      ))}
    </div>
  )
}
