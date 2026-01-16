import Image from 'next/image'
import Link from 'next/link'

import UsersIcon from '@/assets/icons/ic_users.svg'

type Props = {
  id: string
}

export default function RoomCard({ id }: Props) {
  return (
    <Link
      href={`/rooms/${id}`}
      className="relative overflow-hidden flex rounded-xl h-[180px] py-3 px-2.5 cursor-pointer group"
    >
      {/* Background Image */}
      <Image
        src="https://cdn.magicdecor.in/com/2023/10/20174720/Anime-Scenery-Wallpaper-for-Walls-710x488.jpg"
        alt="Background"
        fill
        priority
        className="object-cover z-0 transition-transform duration-600 group-hover:scale-115"
      />
      {/* Overlay */}
      <div className="absolute inset-0 z-1 bg-black/40 transition-all duration-500 group-hover:bg-black/30" />

      <div className="z-10 flex-1 flex flex-col justify-between">
        {/* Card Header */}
        <div className="flex items-center justify-between">
          {/* Live */}
          <div className="bg-red-600 flex items-center gap-1 px-2 py-1 rounded-lg">
            <span className="h-1.5 w-1.5 rounded-full bg-neutral-100 block" />
            <p className="uppercase font-black text-[9px] leading-3 text-neutral-100">live</p>
          </div>

          {/* Viewers */}
          <div className="flex items-center gap-1 py-1 px-2 rounded-lg bg-neutral-800 overflow-hidden">
            <UsersIcon className="text-emerald-500" />
            <span className="text-neutral-100 font-black text-[9px] leading-3">1.2K</span>
          </div>
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
