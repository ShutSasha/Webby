'use client'

import Image from 'next/image'

import { cn } from '@/lib/utils/general.utils'
import { useProfileStore } from '@/stores/profile.store'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  image: string
  username: string
}

export default function UserHeader({ image, username }: Props) {
  const isLoading = useProfileStore(state => state.isLoading)

  return (
    <div className="flex flex-col justify-center items-center gap-4 py-4">
      <div
        className="relative size-[120px] rounded-full overflow-hidden ring-4 ring-neutral-800/50 shadow-xl
          bg-neutral-900"
      >
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
            `absolute inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-10 transition-all
            duration-300`,
            isLoading ? 'opacity-100 visible' : 'opacity-0 invisible',
          )}
        >
          <div className="size-8 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      </div>

      <h1
        className={cn(
          'text-2xl font-bold tracking-tight text-foreground-secondary transition-all duration-300',
          isLoading ? 'opacity-50 animate-pulse' : 'opacity-100',
        )}
      >
        {username}
      </h1>
    </div>
  )
}
