import Image from 'next/image'

export default async function UserVideos() {
  await new Promise(resolve => {
    setTimeout(() => {
      resolve('')
    }, 1300)
  })

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {[...new Array(10)].map((item, index) => (
        <div key={index} className="group cursor-pointer">
          <Image
            src={'https://i.ibb.co/60Ns8j8r/cd4af4dd04fcfba0a358cfdee5c039f7.jpg'}
            alt=""
            height={1920}
            width={1080}
            className="w-full h-[260px] md:h-[200px] xl:h-[180px] rounded-2xl mb-1 group-hover:scale-95 transition-all
              duration-300"
            loading="lazy"
          />
          <p className="text-sm font-medium">さめ甘</p>
          <div className="flex gap-2 text-[12px] text-neutral-500">
            <p>1.3K views</p>
            <p>03/11/2025</p>
          </div>
        </div>
      ))}
    </div>
  )
}

export function UserVideosSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {Array.from({ length: 20 }).map((_, index) => (
        <div key={index}>
          {/* Thumbnail */}
          <div
            className="w-full h-[260px] md:h-[200px] xl:h-[180px] rounded-2xl mb-2 bg-neutral-800
              animate-[shimmer_1.5s_infinite]"
          />

          {/* Title */}
          <div className="h-4 w-2/3 bg-neutral-800 rounded-md mb-2 animate-pulse" />

          {/* Meta info */}
          <div className="flex gap-2">
            <div className="h-3 w-16 bg-neutral-800 rounded-md animate-pulse" />
            <div className="h-3 w-20 bg-neutral-800 rounded-md animate-pulse" />
          </div>
        </div>
      ))}
    </div>
  )
}
