'use client'

import { useEffect, useRef, useState, useTransition } from 'react'

import { useRouter } from 'next/navigation'

import MoreOptions from '@/assets/icons/shared/more-vertical.svg'
import { pinAchievement, unpinAchievement } from '@/lib/actions/achievement.actions'
import { formatDate } from '@/lib/utils/date.utils'
import { cn } from '@/lib/utils/general.utils'
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
  achievementProgressValue: number
  targetValue: number
  unlockedAt: string
  className?: string
}

export default function AchievementItem({
  achievementId,
  title,
  description,
  image,
  isPinned,
  isUnlocked,
  achievementProgressValue,
  targetValue,
  unlockedAt,
  className,
}: AchievementItemProps) {
  const [isOpen, setOpen] = useState<boolean>(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const optionsRef = useRef<SVGSVGElement>(null)
  const addToast = useToastStore(state => state.addToast)
  const router = useRouter()

  const [isFetching, setIsFetching] = useState(false)
  const [isPending, startTransition] = useTransition()

  const isLoading = isFetching || isPending

  const progressPercentage =
    targetValue > 0 ? Math.min(100, Math.round((achievementProgressValue / targetValue) * 100)) : 0

  const formattedDate = unlockedAt ? formatDate(unlockedAt) : ''

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
    setIsFetching(true)

    const response = await pinAchievement(achievementId)

    if (!response.success) {
      const msg = response.errors?.message
      addToast(msg ? msg : 'Unexpected error has occurred.', 'error')
      setIsFetching(false)
      return
    }

    setIsFetching(false)

    startTransition(() => {
      router.refresh()
    })
  }

  const handleUnpin = async () => {
    setOpen(false)
    setIsFetching(true)

    const response = await unpinAchievement(achievementId)

    if (!response.success) {
      const msg = response.errors?.message
      addToast(msg ? msg : 'Unexpected error has occurred.', 'error')
      setIsFetching(false)
      return
    }

    setIsFetching(false)

    startTransition(() => {
      router.refresh()
    })
  }

  return (
    <div
      className={cn(
        `relative rounded-[20px] border border-border px-4 py-4 flex flex-col gap-1 items-center w-[200px] min-h-[230px]
        select-none transition-all shrink-0`,
        className,
        isLoading && 'pointer-events-none',
      )}
    >
      {isLoading && (
        <div
          className="absolute inset-0 bg-black/20 z-20 flex items-center justify-center backdrop-blur-[1px]
            transition-all rounded-[20px]"
        >
          <div className="size-8 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
        </div>
      )}

      <SafeImage
        src={image}
        alt=""
        width={150}
        height={150}
        className={cn('size-[90px] object-cover rounded-full mt-1', { grayscale: !isUnlocked })}
        loading="lazy"
        placeholder="blur"
        blurDataURL={BLUR_DATA_URLS['neutral800']}
      />
      <p className="font-bold select-text mt-1">{title}</p>
      <p className="text-foreground0 text-sm text-center line-clamp-2" title={description}>
        {description}
      </p>

      {!isUnlocked ? (
        <div className="w-full mt-auto pt-3">
          <div className="flex justify-between items-center text-[11px] text-foreground-muted mb-1.5 px-1 font-medium">
            <span>
              {achievementProgressValue} / {targetValue}
            </span>
            <span>{progressPercentage}%</span>
          </div>
          <div className="w-full h-1.5 bg-neutral-800 rounded-full overflow-hidden">
            <div
              className={cn(
                'h-full rounded-full transition-all duration-500',
                isUnlocked ? 'bg-emerald-500' : 'bg-emerald-500/60',
              )}
              style={{ width: `${progressPercentage}%` }}
            />
          </div>
        </div>
      ) : (
        <div className="w-full mt-auto pt-3 flex justify-center">
          {formattedDate !== 'Unknown time' && formattedDate !== '' && (
            <div
              className="inline-flex items-center gap-1.5 bg-emerald-500/10 border border-emerald-500/20
                text-emerald-400 text-[11px] font-medium px-2.5 py-1 rounded-full"
            >
              <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
              </svg>
              Unlocked {formattedDate}
            </div>
          )}
        </div>
      )}

      {isUnlocked && (
        <MoreOptions
          ref={optionsRef}
          className="absolute size-5 right-3 top-3 text-foreground-subtle cursor-pointer hover:text-foreground-secondary
            transition-colors z-10"
          onClick={(e: React.MouseEvent) => {
            e.stopPropagation()
            setOpen(prev => !prev)
          }}
        />
      )}

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
