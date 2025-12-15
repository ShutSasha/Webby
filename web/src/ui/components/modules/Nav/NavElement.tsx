import { ComponentProps, FC, SVGProps } from 'react'

import Link from 'next/link'

type Props = {
  href: ComponentProps<typeof Link>['href']
  isExpanded: boolean
  Icon: FC<SVGProps<SVGSVGElement>>
  text: string
  className?: string
}

export default function NavElement({ isExpanded, href, Icon, text, className }: Props) {
  return (
    <Link
      href={href}
      className={`flex w-full h-fit cursor-pointer flex-col items-center rounded-sm ${className} ${
        isExpanded
          ? 'flex-row gap-2 p-1 transition-colors duration-300 ease-out hover:bg-emerald-500 hover:text-neutral-800'
          : 'flex-col gap-0 p-0 hover:text-emerald-500'
        }`}
    >
      <Icon className={'h-6 w-6'} />
      <p className={`overflow-hidden transition-opacity ${isExpanded ? 'opacity-100 ' : 'w-0 h-0 opacity-0'}`}>
        {text}
      </p>
    </Link>
  )
}
