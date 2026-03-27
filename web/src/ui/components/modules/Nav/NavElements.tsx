'use client'

import { ComponentProps, FC, SVGProps } from 'react'

import { motion } from 'framer-motion'
import Image from 'next/image'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils/utils'

type DesktopNavProps = {
  href: ComponentProps<typeof Link>['href']
  isExpanded: boolean
  Icon: FC<SVGProps<SVGSVGElement>> | string
  text: string
  className?: string
  iconSize: string
  isUserProfile?: boolean
}

type MobileNavProps = {
  href: ComponentProps<typeof Link>['href']
  Icon: FC<SVGProps<SVGSVGElement>>
  text: string
  className?: string
  iconSize: string
  toggleMobileNav: () => void
}

export function DesktopNavElement({
  isExpanded,
  href,
  Icon,
  text,
  className,
  iconSize,
  isUserProfile,
}: DesktopNavProps) {
  const pathname = usePathname()

  const isActive = href === '/' ? pathname === href : pathname.startsWith(href as string)

  return (
    <Link
      href={href}
      className={cn(
        'relative group flex items-center rounded-xl transition-colors duration-300 ease-out cursor-pointer',
        isExpanded ? 'px-3 py-2.5' : 'p-2.5 justify-center',
        isActive ? 'text-emerald-500 font-semibold' : 'text-neutral-400 hover:bg-neutral-800/60 hover:text-neutral-100',
        className,
      )}
    >
      {isActive && (
        <motion.div
          layoutId="sidebar-active-indicator"
          className="absolute inset-0 bg-emerald-500/10 rounded-xl border border-emerald-500/20"
          initial={false}
          transition={{
            type: 'spring',
            stiffness: 350,
            damping: 30,
          }}
        />
      )}

      <div className={cn('relative z-10 flex items-center w-full', !isExpanded && 'justify-center')}>
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

        <span
          className={cn(
            'overflow-hidden whitespace-nowrap transition-all duration-300 ease-in-out',
            isExpanded ? 'w-auto opacity-100 ml-3' : 'w-0 opacity-0 ml-0',
          )}
        >
          {text}
        </span>
      </div>
    </Link>
  )
}

export function MobileNavElement({ href, Icon, text, className, iconSize, toggleMobileNav }: MobileNavProps) {
  const elementWrapperClasses = `flex h-fit cursor-pointer flex-col items-center rounded-sm hover:text-emerald-500 ${className}`

  if (text === 'Menu') {
    return (
      <div className={elementWrapperClasses} onClick={toggleMobileNav}>
        <div className={'w-full flex-1 flex flex-col items-center'}>
          <Icon className={`${iconSize} shrink-0`} />
          <p className="w-full overflow-hidden truncate max-[380px]:text-[12px] text-center">{text}</p>
        </div>
      </div>
    )
  }

  return (
    <Link href={href} className={elementWrapperClasses}>
      <div className={'w-full flex-1 flex flex-col items-center'}>
        <Icon className={`${iconSize} shrink-0`} />
        <p className="w-full overflow-hidden truncate max-[380px]:text-[12px] text-center">{text}</p>
      </div>
    </Link>
  )
}
