'use client'

import { useEffect, useRef, useState } from 'react'

import { useRouter } from 'next/navigation'

import { pinAchievement, unpinAchievement } from '@/app/api/achievements'
import MoreOptions from '@/assets/icons/shared/more-vertical.svg'
import { cn } from '@/lib/utils/utils'
import { useToastStore } from '@/stores/toast-store'
import SafeImage from '@/ui/components/shared/SafeImage'
import { BLUR_DATA_URLS } from '@/ui/images'

type AchievementItemProps = {
  achievementId: string
  title: string
  description: string
  image: string
  isPinned: boolean
  isUnlocked: boolean
  className?: string
}

export default function AchievementItem({
  achievementId,
  title,
  description,
  image,
  isPinned,
  isUnlocked,
  className,
}: AchievementItemProps) {
  const [isOpen, setOpen] = useState<boolean>(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const optionsRef = useRef<SVGSVGElement>(null)
  const addToast = useToastStore(state => state.addToast)
  const router = useRouter()

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (!menuRef?.current?.contains(event.target as Node) && !optionsRef?.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isOpen])

  const handlePin = async () => {
    setOpen(false)
    const response = await pinAchievement(achievementId)

    if (!response.success) {
      const msg = response.errors?.message
      addToast(msg ? msg : 'Unexpected error has occurred.', 'error')
      return
    }

    router.refresh()
  }

  const handleUnpin = async () => {
    setOpen(false)
    const response = await unpinAchievement(achievementId)

    if (!response.success) {
      const msg = response.errors?.message
      addToast(msg ? msg : 'Unexpected error has occurred.', 'error')
      return
    }

    router.refresh()
  }

  return (
    <div
      className={`relative rounded-[20px] border border-border px-5 py-3 w-fit flex flex-col gap-1 items-center
        max-w-[200px] max-h-[200px] ${className} select-none`}
    >
      <SafeImage
        src={image}
        alt=""
        width={150}
        height={150}
        className={cn('size-[100px] object-cover rounded-full', { grayscale: !isUnlocked })}
        loading="lazy"
        placeholder="blur"
        blurDataURL={BLUR_DATA_URLS['neutral800']}
      />
      <p className="font-bold select-text">{title}</p>
      <p className="text-neutral-500 text-sm text-center line-clamp-2" title={description}>
        {description}
      </p>

      <MoreOptions
        ref={optionsRef}
        className="absolute size-5 right-2 top-2.5 text-neutral-300 cursor-pointer"
        onClick={(e: React.MouseEvent) => {
          e.stopPropagation()
          setOpen(prev => !prev)
        }}
      />

      {isOpen && (
        <div
          ref={menuRef}
          className="absolute top-10 right-2 border border-border bg-neutral-800 rounded-lg p-2 z-50 shadow-xl"
        >
          {!isPinned && (
            <p
              className="hover:bg-neutral-700 px-2 py-1 rounded cursor-pointer transition-colors text-nowrap"
              onClick={handlePin}
            >
              Pin the achievement
            </p>
          )}
          {isPinned && (
            <p
              className="hover:bg-neutral-700 px-2 py-1 rounded cursor-pointer transition-colors text-nowrap"
              onClick={handleUnpin}
            >
              Unpin the achievement
            </p>
          )}
        </div>
      )}
    </div>
  )
}
