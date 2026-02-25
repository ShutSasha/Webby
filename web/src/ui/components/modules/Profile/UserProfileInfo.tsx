import Image from 'next/image'

import MailIcon from '@/assets/icons/ic_mail.svg'
import ComplaintIcon from '@/assets/icons/Profile/ic_complaint.svg'
import UserPlusIcon from '@/assets/icons/Profile/ic_user_plus.svg'
import EditProfileBtn from '@/ui/components/modules/Profile/EditProfileBtn'
import ProfileActionButton from '@/ui/components/modules/Profile/ProfileActionButton'
import { UserAchievements } from '@/ui/components/modules/Profile/UserAchivments'
import UserBioSection from '@/ui/components/modules/Profile/UserBioSection'
import { auth } from '@/workspace/auth'

export default async function UserProfileInfo({ id }: { id: string }) {
  // const user = await getUser(id)
  // await, sync user data
  // await, sync user folowers and follows
  // await, sync user pinned badges
  const session = await auth()
  await new Promise(r => setTimeout(r, 800))

  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col md:flex-row justify-between gap-4">
      {/* Left Part of user profle*/}
      <div className="flex flex-col gap-4">
        <div className="flex gap-4">
          <Image
            src={'https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg'}
            alt=""
            width={250}
            height={250}
            placeholder="blur"
            blurDataURL="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mMUrwcAALMAmGjO2MQAAAAASUVORK5CYII="
            className="h-20 w-20 md:h-[125px] md:w-[125px] rounded-full object-cover"
            preload
            loading="eager"
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
  )
}

export function UserProfileInfoSkeleton() {
  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col md:flex-row justify-between gap-4">
      {/* --- Left Part (Avatar + Bio + Actions) --- */}
      <div className="flex flex-col gap-4">
        <div className="flex gap-4">
          {/* Avatar Skeleton */}
          <div className="shrink-0 h-20 w-20 md:h-[125px] md:w-[125px] rounded-full bg-neutral-800 animate-pulse" />

          {/* Bio Section Skeleton */}
          <div className="flex flex-col justify-between py-1">
            <div className="flex flex-col gap-2">
              {/* Username */}
              <div className="h-7 w-32 md:w-48 bg-neutral-800 rounded-md animate-pulse" />
              {/* Bio line */}
              <div className="h-4 w-40 md:w-64 bg-neutral-800 rounded-md animate-pulse" />
            </div>

            {/* Followers/Follows line */}
            <div className="flex gap-2 items-center mt-2">
              <div className="h-4 w-20 bg-neutral-800 rounded-md animate-pulse" />
              <div className="h-1 w-1 bg-neutral-800 rounded-full" />
              <div className="h-4 w-20 bg-neutral-800 rounded-md animate-pulse" />
            </div>
          </div>
        </div>

        {/* Action Buttons (Chat / Follow) */}
        <div className="flex items-center gap-2">
          <div className="h-8.5 w-24 bg-neutral-800 rounded-full animate-pulse" />
          <div className="h-8.5 w-24 bg-neutral-800 rounded-full animate-pulse" />
        </div>
      </div>

      {/* --- Right Part (Edit Btn + Achievements) --- */}
      <div className="flex flex-col gap-4 md:gap-2 min-w-[200px]">
        {/* Edit/Complaint Button (Aligned end on desktop) */}
        <div className="h-9 w-32 bg-neutral-800 rounded-full animate-pulse self-start md:self-end" />

        {/* Achievements Skeleton */}
        <div className="flex flex-col gap-1.5 mt-auto">
          {/* Label "Badges" */}
          <div className="h-4 w-16 bg-neutral-800 rounded-md animate-pulse" />

          {/* HR line */}
          <div className="h-px w-full bg-neutral-800" />

          {/* Badges icons */}
          <div className="flex items-center gap-4">
            {[1, 2, 3].map(i => (
              <div
                key={i}
                className="h-[60px] w-[60px] lg:h-[90px] lg:w-[90px] bg-neutral-800 rounded-full animate-pulse"
              />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
