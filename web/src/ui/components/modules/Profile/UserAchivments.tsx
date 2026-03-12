import Image from 'next/image'

import { Achievement } from '@/types/achivement'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  achivements: Achievement[]
}

export async function UserAchievements({ achivements }: Props) {
  return (
    <div className="flex flex-col gap-1.5 self-end">
      <p>Badges</p>
      <hr className="text-emerald-400 w-full" />
      <div className="flex items-center gap-4">
        {achivements.map(achivement => (
          <Image
            key={achivement.achievementId}
            src={achivement.iconUrl}
            alt=""
            width={180}
            height={180}
            preload
            loading="eager"
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral800']}
            className="cursor-pointer size-[60px] md:size-[90px] object-cover rounded-full"
          />
        ))}
      </div>
    </div>
  )
}
