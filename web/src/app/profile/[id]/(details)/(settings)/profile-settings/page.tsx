import Image from 'next/image'
import { redirect } from 'next/navigation'

import { BLUR_DATA_URLS } from '@/ui/images'
import { auth } from '@/workspace/auth'

export default async function Settings() {
  const session = await auth()

  if (!session?.user) {
    redirect('/login')
  }

  return (
    <div className="flex flex-col">
      <div className="flex flex-col justify-center items-center gap-2">
        <Image
          src={session.user.image}
          alt=""
          width={125}
          height={125}
          className="size-[125px] rounded-full"
          loading="lazy"
          placeholder="blur"
          blurDataURL={BLUR_DATA_URLS['neutral800']}
        />
        <p className="text-[24px] font-medium leading-[30px]">{session.user.username}</p>
      </div>
      <p>settings content here</p>
    </div>
  )
}
