import Image from 'next/image'
import Link from 'next/link'

import MailIcon from '@/assets/icons/ic_mail_with_background.svg'
import TrashIcon from '@/assets/icons/ic_trash.svg'
import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  user: {
    image: string
  }
}

export default function FollowItem({ user }: Props) {
  return (
    <div
      className="bg-neutral-800 hover:bg-neutral-300/10 transition-all duration-200 ease-in rounded-xl p-2 flex
        items-center justify-between gap-2 border border-transparent hover:border-emerald-400/25"
    >
      <Link href={'/'} className="flex items-center gap-2">
        {/* User info */}
        <Image
          src={user.image}
          alt=""
          width={60}
          height={60}
          className="w-13 h-13 rounded-full"
          loading="lazy"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral900']}
        />
        <div className="flex flex-col">
          <p className="text-[16px] font-medium">Username</p>
          <p className="text-neutral-500 text-[12px]">152 followers</p>
        </div>
      </Link>

      {/* Actions */}
      <div className="flex items-center gap-5 pr-2">
        <div className="relative group/mail cursor-pointer">
          <MailIcon
            className="w-4.5 h-4.5 text-neutral-500/90 group-hover/mail:text-emerald-500/90 transition-all duration-300"
          />
          <div
            className="absolute w-8 h-8 left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2
              group-hover/mail:bg-emerald-500/20 transition-all duration-300 rounded-full"
          />
        </div>
        <div className="relative group/trash cursor-pointer">
          <TrashIcon
            className="w-4.5 h-4.5 text-red-500/80 group-hover/trash:text-red-500/90 transition-all duration-300"
          />
          <div
            className="absolute w-8 h-8 left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2
              group-hover/trash:bg-red-500/20 transition-all duration-300 rounded-full"
          />
        </div>
      </div>
    </div>
  )
}
