'use client'

import Image from 'next/image'

import { cn } from '@/lib/utils/utils'
import { useProfileStore } from '@/stores/profile.store'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  image: string
  username: string
}

export default function UserHeader({ image, username }: Props) {
  const isLoading = useProfileStore(state => state.isLoading)

  return (
    <div className="flex flex-col justify-center items-center gap-2">
      <div className="relative size-[125px] rounded-full overflow-hidden">
        <Image
          src={image}
          alt={username}
          width={400}
          height={400}
          className={cn(
            'size-full object-cover transition-all duration-500',

            isLoading ? 'scale-110 blur-[2px]' : 'scale-100 blur-0',
          )}
          loading="lazy"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
        />

        <div
          className={cn(
            'absolute inset-0 bg-black/40 flex items-center justify-center z-10 transition-opacity duration-300',
            isLoading ? 'opacity-100 visible' : 'opacity-0 invisible',
          )}
        >
          <div className="size-8 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      </div>

      <p
        className={cn(
          'text-[24px] font-medium leading-[30px] transition-all duration-300',
          isLoading ? 'opacity-50 animate-pulse' : 'opacity-100',
        )}
      >
        {username}
      </p>
    </div>
  )
}
