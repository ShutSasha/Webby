import Image from 'next/image'

import MainLayout from '@/ui/components/MainLayout'
import EditProfileBtn from '@/ui/components/modules/Profile/EditProfileBtn'
import UserBioSection from '@/ui/components/modules/Profile/UserBioSection'
import { auth } from '@/workspace/auth'

type ProfileProps = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function Profile({ params, searchParams }: ProfileProps) {
  const { id } = await params
  const session = await auth()
  // const { tab } = await searchParams // наприклад, '?tab=videos'

  // await, sync user data
  // await, sync user folowers and follows
  // await, sync user pinned badges

  // await, sync videos - Lates/Popular
  // await, sync Public playlists

  return (
    <MainLayout>
      <div className="flex flex-col gap-4 w-full max-w-5xl 2xl:max-w-7xl mx-auto">
        <div className="bg-neutral-900 rounded-[20px] p-5 flex justify-between gap-4">
          {/* Left Part of user profle*/}
          <div className="flex gap-4 max-h-[125px]">
            <Image
              src={'https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg'}
              alt=""
              width={125}
              height={125}
              className="h-[125px] w-[125px] rounded-full"
            />
            <UserBioSection username="username1" userId={id} bio="bio" />
          </div>

          {/* Right Part of user profile*/}
          <div className="flex flex-col gap-2">
            {session?.user.id === id && <EditProfileBtn userId={id} />}

            <div className="flex flex-col gap-1.5">
              <p>Badges</p>
              <hr className="text-emerald-400 w-full" />
              <div className="flex items-center gap-4">
                {[...new Array(3)].map((_, index) => (
                  <Image
                    key={index}
                    src="https://i.ibb.co/ch9ZCDTr/badge1.png"
                    alt=""
                    width={80}
                    height={80}
                    loading="lazy"
                    className="cursor-pointer"
                  />
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Videos, playlists selector */}
        <div>selector for videos/playlists</div>

        {/* Block - Dynamic content: Videos, Playlists */}
        <div className="bg-neutral-900 rounded-[20px] p-5">
          <p>Videos, Playlists CONTENT</p>
        </div>
      </div>
    </MainLayout>
  )
}
