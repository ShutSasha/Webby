import Image from 'next/image'
import Link from 'next/link'

import UsersIcon from '@/assets/icons/ic_users.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  id: string
}

export default function RoomCard({ id }: Props) {
  return (
    <Link
      href={`/rooms/${id}`}
      className="relative aspect-video overflow-hidden flex rounded-xl py-3 px-2.5 cursor-pointer group"
    >
      {/* Background Image */}
      <Image
        src="https://cdn.magicdecor.in/com/2023/10/20174720/Anime-Scenery-Wallpaper-for-Walls-710x488.jpg"
        alt="Background"
        fill
        loading="lazy"
        className="object-cover z-0 transition-transform duration-600 group-hover:scale-115"
        placeholder="blur"
        blurDataURL={BLUR_DATA_URLS['neutral900']}
      />
      {/* Overlay */}
      <div className="absolute inset-0 z-1 bg-black/40 transition-all duration-500 group-hover:bg-black/30" />

      <div className="z-10 flex-1 flex flex-col justify-between">
        {/* Tag */}
        <div className="bg-neutral-900 flex items-center gap-1 px-2 py-2 rounded-lg w-fit">
          <p className="uppercase font-bold text-[12px] leading-3 text-neutral-300 tracking-wide">anime</p>
        </div>

        <div className="flex flex-col gap-1">
          {/* User */}
          <div className="flex gap-1.5 items-center">
            <Image
              src="https://static.wikia.nocookie.net/madagascar/images/3/30/37455825.jpg/revision/latest?cb=20150512133950&path-prefix=ru"
              alt="User avatar"
              width={24}
              height={24}
              className="rounded-full object-cover h-7 w-7 border border-emerald-500"
              placeholder="blur"
              blurDataURL={BLUR_DATA_URLS['neutral700']}
            />
            <p className="text-neutral-100/85 font-black text-[12px] leading-4">Username</p>
          </div>

          {/* Title */}
          <p className="text-neutral-100 text-sm leading-4.5 font-black">Roadside. Jakob & Ryan</p>
        </div>
      </div>
    </Link>
  )
}
