import Image from 'next/image'

export default async function PlaylistItem() {
  return (
    <div className="relative group cursor-pointer">
      <Image
        src={'https://i.ibb.co/FShJxq8/anime-style-clouds.jpg'}
        alt=""
        height={1920}
        width={1080}
        className="w-full rounded-2xl group-hover:scale-95 transition-all duration-300 mb-1"
        loading="lazy"
      />
      <p className="text-sm font-medium">Playlist Name</p>
      <p className="text-[12px] text-neutral-500">Creator</p>
      <div
        className="absolute bg-neutral-800 rounded-lg px-2 py-1 top-2 right-2 text-[12px] group-hover:top-3.5
          group-hover:right-3.5 transition-all duration-300"
      >
        3 videos
      </div>
    </div>
  )
}
