import Image from 'next/image'

import MailIcon from '@/assets/icons/ic_mail.svg'
import ComplaintIcon from '@/assets/icons/Profile/ic_complaint.svg'
import UserPlusIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import { clog } from '@/lib/utils/utils'
import MainLayout from '@/ui/components/MainLayout'
import EditProfileBtn from '@/ui/components/modules/Profile/EditProfileBtn'
import ProfileActionButton from '@/ui/components/modules/Profile/ProfileActionButton'
import ProfileTab from '@/ui/components/modules/Profile/ProfileTab'
import { UserAchievements } from '@/ui/components/modules/Profile/UserAchivments'
import UserBioSection from '@/ui/components/modules/Profile/UserBioSection'
import UserPlaylists from '@/ui/components/modules/Profile/UserPlaylists'
import UserVideos from '@/ui/components/modules/Profile/UserVideos'
import { auth } from '@/workspace/auth'

type ProfileProps = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const allowedTabs = ['Video', 'Playlist'] as const

export default async function Profile({ params, searchParams }: ProfileProps) {
  const { id } = await params
  const session = await auth()
  const { tab } = await searchParams
  const currentTab = allowedTabs.includes(tab as any) ? tab : 'Video'

  // await, sync user data
  // await, sync user folowers and follows
  // await, sync user pinned badges

  // await, sync videos - Lates/Popular
  // await, sync Public playlists

  return (
    <MainLayout>
      <div className="flex flex-col gap-4 w-full max-w-5xl 2xl:max-w-7xl mx-auto">
        <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col md:flex-row justify-between gap-4">
          {/* Left Part of user profle*/}
          <div className="flex flex-col gap-4">
            <div className="flex gap-4">
              <Image
                src={'https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg'}
                alt=""
                width={500}
                height={500}
                className="h-20 w-20 md:h-[125px] md:w-[125px] rounded-full"
              />
              <UserBioSection username="username1" userId={id} bio="bio" />
            </div>
            <div className="flex items-center gap-2">
              <ProfileActionButton Icon={MailIcon} label="Chat" iconClassName="w-4 h-4" />

              <ProfileActionButton Icon={UserPlusIcon} label="Follow" iconClassName="w-4 h-4" />
            </div>
          </div>

          {/* Right Part of user profile*/}
          <div className="flex flex-col gap-2">
            {session?.user.id === id ? (
              <EditProfileBtn userId={id} />
            ) : (
              <ProfileActionButton
                Icon={ComplaintIcon}
                label="Leave complaint"
                iconClassName="w-4 h-4"
                btnClassName="self-end"
              />
            )}
            <UserAchievements />
          </div>
        </div>

        {/* Block - Dynamic content: Videos, Playlists */}
        <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col gap-5">
          {/* Videos, playlists selector */}

          <div className="flex gap-2 items-center">
            <ProfileTab label="Video" />
            <ProfileTab label="Playlist" />
          </div>

          {currentTab === 'Video' && <UserVideos />}
          {currentTab === 'Playlist' && <UserPlaylists />}
        </div>
      </div>
    </MainLayout>
  )
}
