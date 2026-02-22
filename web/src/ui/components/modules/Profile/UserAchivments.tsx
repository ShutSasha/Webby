import Image from 'next/image'

export async function UserAchievements() {
  return (
    <div className="flex flex-col gap-1.5">
      <p>Badges</p>
      <hr className="text-emerald-400 w-full" />
      <div className="flex items-center gap-4">
        {[...new Array(3)].map((_, index) => (
          <Image
            key={index}
            src="https://i.ibb.co/ch9ZCDTr/badge1.png"
            alt=""
            width={200}
            height={200}
            loading="lazy"
            className="cursor-pointer h-[60px] w-[60px] lg:h-[90px] lg:w-[90px]"
          />
        ))}
      </div>
    </div>
  )
}
