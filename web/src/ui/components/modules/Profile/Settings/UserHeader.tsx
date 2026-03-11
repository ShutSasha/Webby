import Image from 'next/image'

import { BLUR_DATA_URLS } from '@/ui/images'

type Props = {
  image: string
  username: string
}

export default function UserHeader({ image, username }: Props) {
  return (
    <div className="flex flex-col justify-center items-center gap-2">
      <Image
        src={image}
        alt=""
        width={400}
        height={400}
        className="size-[125px] rounded-full object-cover"
        loading="lazy"
        placeholder="blur"
        blurDataURL={BLUR_DATA_URLS['neutral800']}
      />
      <p className="text-[24px] font-medium leading-[30px]">{username}</p>
    </div>
  )
}
