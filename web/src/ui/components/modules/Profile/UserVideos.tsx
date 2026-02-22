import Image from 'next/image'

export default async function UserVideos() {
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
