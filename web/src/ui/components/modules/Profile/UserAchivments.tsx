import Image from 'next/image'

import { cn } from '@/lib/utils/general.utils'
import { Achievement } from '@/types/achivement.types'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  achivements: Achievement[]
}

export async function UserAchievements({ achivements }: Props) {
  return (
    <div className="flex flex-col gap-2 self-start md:self-end mt-4 md:mt-auto md:items-end">
      <span className="text-xs font-semibold text-neutral-400 uppercase tracking-wider">Achievements</span>

      <div className="flex items-center gap-3 md:justify-end">
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
            className={cn(
              'cursor-pointer size-[50px] md:size-[60px] object-cover rounded-full',
              'border border-neutral-800 hover:border-neutral-600 transition-colors ',
            )}
          />
        ))}
      </div>
    </div>
  )
}
