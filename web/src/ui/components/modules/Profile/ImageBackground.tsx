import Image from 'next/image'

import { BLUR_DATA_URLS } from '@/ui/images'

export default function ImageBackground({ src }: { src: string }) {
  return (
    <Image
      src={src}
      alt=""
      width={700}
      height={394}
      className="w-full rounded-2xl mb-1 group-hover:scale-95 transition-all duration-300 aspect-video object-cover"
      loading="lazy"
      placeholder="blur"
      blurDataURL={BLUR_DATA_URLS['neutral800']}
    />
  )
}
