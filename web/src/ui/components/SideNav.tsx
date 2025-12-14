import Image from 'next/image'
import Link from 'next/link'

import LogoIcon from '@/assets/icons/ic_logo.svg'
import SearchIcon from '@/assets/icons/ic_search.svg'

export default function SideNav() {
  return (
    <div className="flex flex-col items-center bg-neutral-900 p-2">
      <Link href={'/'}>
        <LogoIcon className="h-8 w-8" />
      </Link>
      <SearchIcon className="h-6 w-6 text-emerald-500" />
      <Image src="/icons/ic_search.svg" alt="ww" width={50} height={50} className="text-amber-400" />
      <div>right</div>
      <div>2</div>
    </div>
  )
}
