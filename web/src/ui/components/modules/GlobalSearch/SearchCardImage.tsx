'use client'

import { useState } from 'react'

import Image from 'next/image'

import VideosIcon from '@/assets/icons/Nav/tv-minimal-play.svg'
import UsersIcon from '@/assets/icons/Notifications/users.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  src: string
  alt: string
  isUser: boolean
}

export default function SearchCardImage({ src, alt, isUser }: Props) {
  const [imageError, setImageError] = useState(false)

  const FallbackIcon = isUser ? UsersIcon : VideosIcon

  return (
    <div
      className={`relative shrink-0 overflow-hidden bg-neutral-800 flex items-center justify-center
        ${isUser ? 'w-12 h-12 rounded-full' : 'w-24 h-14 rounded-lg'}`}
    >
      {!imageError && src ? (
        <Image
          src={src}
          alt={alt}
          width={160}
          height={90}
          className="object-cover aspect-video h-full w-full"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
          onError={() => setImageError(true)}
        />
      ) : (
        <FallbackIcon className="size-6 text-neutral-600 stroke-[1.5px]" />
      )}
    </div>
  )
}
