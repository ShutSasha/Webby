import { Suspense } from 'react'

import Image from 'next/image'

import MailIcon from '@/assets/icons/ic_mail.svg'
import ComplaintIcon from '@/assets/icons/Profile/ic_complaint.svg'
import UserPlusIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import EditProfileBtn from '@/ui/components/modules/Profile/EditProfileBtn'
import ProfileActionButton from '@/ui/components/modules/Profile/ProfileActionButton'
import { UserAchievements } from '@/ui/components/modules/Profile/UserAchivments'
import UserBioSection, { UserBioSectionSkeleton } from '@/ui/components/modules/Profile/UserBioSection'
import { auth } from '@/workspace/auth'

export default async function UserProfileInfo({ id }: { id: string }) {
  // const user = await getUser(id)
  const session = await auth()
  await new Promise(resolve => {
    setTimeout(() => {
      resolve('')
    }, 2000)
  })

  return (
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

          <Suspense fallback={<UserBioSectionSkeleton />}>
            <UserBioSection username="username1" userId={id} bio="bio" />
          </Suspense>
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
  )
}

export function UserProfileInfoSkeleton() {
  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col md:flex-row justify-between gap-4">
      <div className="flex flex-col gap-4">
        <div className="flex gap-4">
          {/* Avatar */}
          <div className="h-20 w-20 md:h-[125px] md:w-[125px] rounded-full bg-neutral-800 animate-pulse" />

          {/* Bio block */}
          <div className="flex flex-col gap-2">
            <div className="h-5 w-32 bg-neutral-800 rounded-md animate-pulse" />
            <div className="h-4 w-48 bg-neutral-800 rounded-md animate-pulse" />
          </div>
        </div>

        <div className="flex gap-2">
          <div className="h-9 w-20 bg-neutral-800 rounded-lg animate-pulse" />
          <div className="h-9 w-24 bg-neutral-800 rounded-lg animate-pulse" />
        </div>
      </div>

      <div className="flex flex-col gap-2">
        <div className="h-9 w-32 bg-neutral-800 rounded-lg self-end animate-pulse" />
        <div className="h-20 w-40 bg-neutral-800 rounded-lg animate-pulse" />
      </div>
    </div>
  )
}
