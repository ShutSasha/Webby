import { ReactNode } from 'react'

import Image from 'next/image'
import Link from 'next/link'

type Props = {
  id: string
  image: string
  children: ReactNode
}

export default function HeaderNavigation({ id, image, children }: Props) {
  return (
    <div className="relative flex items-center">
      <Link href={`/profile/${id}`} className="flex items-center gap-2">
        <Image src={image} alt="" width={400} height={400} className="w-8 h-8 rounded-full" />
        <p className="text-[16px] font-semibold underline">Username</p>
      </Link>

      <div className="absolute left-1/2 -translate-x-1/2">{children}</div>
    </div>
  )
}
