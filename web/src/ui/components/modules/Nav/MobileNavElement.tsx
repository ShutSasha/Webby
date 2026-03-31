import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

type MobileNavProps = {
  href: ComponentProps<typeof Link>['href']
  Icon: FC<SVGProps<SVGSVGElement>>
  text: string
  className?: string
  iconSize: string
  toggleMobileNav: () => void
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
