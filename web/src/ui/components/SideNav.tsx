import LogoIconSource from '@/assets/icons/ic_logo.svg'
import SearchIconSource from '@/assets/icons/ic_search.svg'
import Image from 'next/image'
import { FC, SVGProps } from 'react'

export default function SideNav() {
  // TODO сделать реюсабельным этот костыль
  const LogoIcon = LogoIconSource as unknown as FC<SVGProps<SVGSVGElement>>
  const SearchIcon = SearchIconSource as unknown as FC<SVGProps<SVGSVGElement>>

  return (
    <div className="flex flex-col items-center bg-neutral-900 p-2">
      <LogoIcon className="h-8 w-8" />
      <SearchIcon className="h-6 w-6 text-emerald-500" />
      <Image src="/icons/ic_search.svg" alt="ww" width={50} height={50} className="text-amber-400" />
      <div>right</div>
      <div>2</div>
    </div>
  )
}
