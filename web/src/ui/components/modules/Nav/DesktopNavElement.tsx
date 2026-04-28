'use client'

import { ComponentProps, FC, SVGProps } from 'react'

import Image from 'next/image'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils/general.utils'

type DesktopNavProps = {
  isExpanded: boolean
  Icon: FC<SVGProps<SVGSVGElement>> | string
  text: string
  iconSize: string
  href?: ComponentProps<typeof Link>['href']
  className?: string
  isUserProfile?: boolean
  onClick?: () => void
  badgeCount?: number
}

export function DesktopNavElement({
  isExpanded,
  href,
  onClick,
  Icon,
  text,
  className,
  iconSize,
  isUserProfile,
  badgeCount,
}: DesktopNavProps) {
  const pathname = usePathname()

  const isActive = href ? (href === '/' ? pathname === href : pathname.startsWith(href as string)) : false

  const elementClasses = cn(
    'relative group flex items-center rounded-xl transition-colors duration-300 ease-out cursor-pointer',
    isExpanded ? 'px-3 py-2.5' : '2xl:p-2.5 xl:p-2 p-1.5 py-2 justify-center',
    isActive ? 'text-emerald-500 font-semibold' : 'text-neutral-400 hover:bg-neutral-800/60 hover:text-neutral-100',
    className,
  )

  const content = (
    <div className={cn('relative z-10 flex items-center w-full', !isExpanded && 'justify-center')}>
      <div className="relative shrink-0 flex items-center justify-center">
        {typeof Icon === 'string' ? (
          <Image
            src={Icon}
            alt={text}
            width={24}
            height={24}
            className={cn(iconSize, 'object-cover rounded-full transition-transform group-hover:scale-105')}
          />
        ) : (
          <Icon
            className={cn(
              iconSize,
              `shrink-0 transition-all duration-300 ${isUserProfile ? 'stroke-1' : 'stroke-[1.5px]'}`,
              isActive ? 'text-emerald-500' : 'group-hover:text-emerald-400',
            )}
          />
        )}

        {!isExpanded && badgeCount !== undefined && badgeCount > 0 && (
          <div
            className="absolute -top-1 -right-1 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-emerald-500
              px-1 text-[10px] font-bold text-neutral-900 ring-1 ring-[#0A0A0A]"
          >
            {badgeCount > 9 ? '9+' : badgeCount}
          </div>
        )}

        <span
          className={cn(
            'overflow-hidden whitespace-nowrap transition-all duration-300 ease-in-out',
            isExpanded ? 'w-auto opacity-100 ml-3' : 'w-0 opacity-0 ml-0',
          )}
        >
          {text}
        </span>

        {isExpanded && badgeCount !== undefined && badgeCount > 0 && (
          <div
            className="ml-4 flex h-[22px] min-w-[22px] items-center justify-center rounded-full bg-emerald-500 px-1.5
              text-[12px] font-bold text-neutral-900 shrink-0 shadow-sm shadow-emerald-500/20"
          >
            {badgeCount > 99 ? '99+' : badgeCount}
          </div>
        )}
      </div>
    </div>
  )

  if (href) {
    return (
      <Link href={href} className={elementClasses}>
        {isActive && <div className="absolute inset-0 bg-emerald-500/10 rounded-xl border border-emerald-500/20" />}
        {content}
      </Link>
    )
  }

  return (
    <button onClick={onClick} className={cn(elementClasses, 'w-full text-left bg-transparent border-none')}>
      {content}
    </button>
  )
}
