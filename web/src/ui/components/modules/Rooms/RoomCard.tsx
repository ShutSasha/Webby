import Image from 'next/image'
import Link from 'next/link'

import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  id: string
}

export default function RoomCard({ id }: Props) {
  return (
    <Link href={`/rooms/${id}`} className="flex flex-col gap-3 group cursor-pointer">
      <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-neutral-800">
        <Image
          src="https://cdn.magicdecor.in/com/2023/10/20174720/Anime-Scenery-Wallpaper-for-Walls-710x488.jpg"
          alt="Room preview"
          fill
          loading="lazy"
          className="object-cover z-0 transition-transform duration-500 ease-out group-hover:scale-105"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral900']}
        />

        <div className="absolute top-2 left-2 z-10 bg-black/60 backdrop-blur-md px-2 py-1 rounded-md">
          <p className="uppercase font-bold text-[10px] tracking-wider text-neutral-200">Anime</p>
        </div>
      </div>

      <div className="flex gap-3 items-start px-1">
        <Image
          src="https://static.wikia.nocookie.net/madagascar/images/3/30/37455825.jpg/revision/latest?cb=20150512133950&path-prefix=ru"
          alt="User avatar"
          width={36}
          height={36}
          className="size-9 rounded-full object-cover shrink-0 mt-0.5"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral700']}
        />

        <div className="flex flex-col overflow-hidden">
          <h3
            className="text-neutral-100 text-sm font-semibold leading-snug line-clamp-2 transition-colors duration-200
              group-hover:text-emerald-400"
          >
            Roadside. Jakob & Ryan watching some cool anime episodes together
          </h3>

          <p className="text-neutral-400 text-xs mt-1 truncate hover:text-neutral-300 transition-colors">
            guuuntersteam
          </p>
        </div>
      </div>
    </Link>
  )
}
