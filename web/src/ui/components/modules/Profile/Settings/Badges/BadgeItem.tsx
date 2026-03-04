import Image from 'next/image'

import { BLUR_DATA_URLS } from '@/ui/images'

type BadgeItemProps = {
  title: string
  description: string
  image: string
  className?: string
}

export default function BadgeItem({ title, description, image, className }: BadgeItemProps) {
  return (
    <div
      className={`rounded-[20px] border border-border px-5 py-3 w-fit flex flex-col gap-1 items-center max-w-[200px]
        max-h-[200px] ${className}`}
    >
      <Image
        src={image}
        alt=""
        width={100}
        height={100}
        className="cursor-pointer size-[100px] object-cover rounded-full"
        loading="lazy"
        placeholder="blur"
        blurDataURL={BLUR_DATA_URLS['neutral800']}
      />
      <p className="font-bold">{title}</p>
      <p className="text-neutral-500 text-sm text-center line-clamp-2" title={description}>
        {description}
      </p>
    </div>
  )
}
