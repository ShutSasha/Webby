import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

type DesktopNavProps = {
  href: ComponentProps<typeof Link>['href']
  isExpanded: boolean
  Icon: FC<SVGProps<SVGSVGElement>>
  text: string
  className?: string
  iconSize: string
}

type MobileNavProps = {
  href: ComponentProps<typeof Link>['href']
  Icon: FC<SVGProps<SVGSVGElement>>
  text: string
  className?: string
  iconSize: string
  toggleMobileNav: () => void
}

export function DesktopNavElement({ isExpanded, href, Icon, text, className, iconSize }: DesktopNavProps) {
  return (
    <Link
      href={href}
      className={`flex h-fit cursor-pointer flex-col items-center rounded-sm ${className} ${
        isExpanded
          ? 'flex-row transition-colors duration-300 ease-out hover:bg-emerald-500 hover:text-neutral-800'
          : 'flex-col hover:text-emerald-500'
        }`}
    >
      <div className={`w-full flex-1 flex flex-row items-center ${isExpanded ? 'gap-2 p-1 ' : 'gap-0 p-0 '}`}>
        <Icon className={`${iconSize} shrink-0`} />
        <p
          className={`overflow-hidden whitespace-nowrap transition-opacity font-medium
            ${isExpanded ? 'opacity-100 ' : 'w-0 h-0 opacity-0'}`}
        >
          {text}
        </p>
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
