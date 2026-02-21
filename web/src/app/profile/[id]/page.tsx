import Image from 'next/image'
import Link from 'next/link'

import MainLayout from '@/ui/components/MainLayout'
import FollowConnections from '@/ui/components/modules/Profile/FollowConnections'

type ProfileProps = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function Profile({ params, searchParams }: ProfileProps) {
  const { id } = await params // наприклад, '123'
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
          <div className="flex gap-4">
            <Image
              src={'https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg'}
              alt=""
              width={125}
              height={125}
              className="h-[125px] w-[125px] rounded-full"
            />
            {/* User Info */}
            <div className="flex flex-col justify-between">
              <div className="flex flex-col">
                <p className="font-medium text-[24px]">username</p>
                <p className="text-sm max-w-[400px]">Lorem ipsum dolor sit amet consectetur adipisicing</p>
              </div>
              {/* followers/follows block */}
              <FollowConnections userId={id} />
            </div>
          </div>

          {/* Right Part of user profile*/}
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
