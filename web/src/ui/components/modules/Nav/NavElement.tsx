import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

type Props = {
  href: ComponentProps<typeof Link>['href']
  isExpanded: boolean
  Icon: FC<SVGProps<SVGSVGElement>>
  text: string
  className?: string
  iconSize: string
  isLogo?: boolean
}

export default function NavElement({ isExpanded, href, Icon, text, className, iconSize, isLogo }: Props) {
  return (
    <Link
      href={href}
      className={`flex h-fit cursor-pointer flex-col items-center rounded-sm ${className} ${
        isExpanded
          ? 'flex-row transition-colors duration-300 ease-out hover:bg-emerald-500 hover:text-neutral-800'
          : 'flex-col hover:text-emerald-500'
        } ${isLogo ? 'py-3' : ''}`}
    >
      <div
        className={`w-full flex-1 flex flex-col md:flex md:flex-row items-center
          ${isExpanded ? 'gap-2 p-1 ' : 'gap-0 p-0 '}`}
      >
        <Icon className={`${iconSize} shrink-0`} />
        {/* text for desktop */}
        <p className={`overflow-hidden transition-opacity ${isExpanded ? 'opacity-100 ' : 'w-0 h-0 opacity-0'}`}>
          {text}
        </p>
        {/* text for mobile */}
        <p className="w-full overflow-hidden block md:hidden truncate max-[380px]:text-[12px] text-center">{text}</p>
      </div>
    </Link>
  )
}
