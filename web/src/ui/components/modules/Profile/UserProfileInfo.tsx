import Image from 'next/image'

import MailIcon from '@/assets/icons/ic_mail.svg'
import { checkFollowing, getUser } from '@/lib/actions/user.actions'
import { cn } from '@/lib/utils/general.utils'
import EditProfileBtn from '@/ui/components/modules/Profile/EditProfileBtn'
import { UserAchievements } from '@/ui/components/modules/Profile/UserAchivments'
import UserBioSection from '@/ui/components/modules/Profile/UserBioSection'
import { BLUR_DATA_URLS } from '@/ui/images'
import { auth } from '@/workspace/auth'

import ComplaintButton from './ComplaintButton'
import FollowButton from './FollowButton'
import ActionButton from '../../shared/ActionButton'
import EmptyState from '../../shared/EmptyState'

export default async function UserProfileInfo({ id }: { id: string }) {
  const userData = await getUser(id)
  const session = await auth()
  const isFollowingResponse = await checkFollowing(id)
  const isFollowing = isFollowingResponse.data?.isFollowing ?? false
  const isOwner = session?.user.id === id

  if (!userData) {
    return (
      <div className="bg-neutral-900 rounded-[20px] p-5 flex items-center justify-center min-h-[250px]">
        <EmptyState
          title="User not found"
          description="This profile doesn't exist, has been deleted, or is temporarily unavailable."
        />
      </div>
    )
  }

  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col md:flex-row justify-between gap-4">
      {/* Left Part of user profle*/}
      <div className="flex flex-col gap-4">
        <div className="flex gap-4">
          <Image
            src={userData.user.avatarUrl}
            alt=""
            width={300}
            height={300}
            placeholder="blur"
            blurDataURL={BLUR_DATA_URLS['neutral800']}
            className={cn(
              'aspect-square rounded-full object-cover shrink-0 transition-all duration-500',
              'size-20',
              isOwner ? 'md:size-[175px]' : 'md:size-[125px]',
            )}
            preload
            loading="eager"
          />

          <UserBioSection userId={id} initialUserData={userData} />
        </div>

        {!isOwner && (
          <div className="flex items-center gap-2">
            <ActionButton label="Chat">
              <MailIcon className="w-4 h-4" />
            </ActionButton>

            <FollowButton targetUserId={id} initialIsFollowing={isFollowing} currentUserId={session?.user.id} />
          </div>
        )}
      </div>

      {/* Right Part of user profile*/}
      <div className="flex flex-col gap-2">
        {session?.user.id === id && <EditProfileBtn userId={id} />}
        {session && session?.user.id !== id && (
          <ComplaintButton authorId={session.user.id} targetId={id} targetType="User" />
        )}

        {userData.pinnedUserAchievements.length > 0 && (
          <UserAchievements achivements={userData.pinnedUserAchievements} />
        )}
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
          {/* Label "Achievement" */}
          <div className="h-4 w-16 bg-neutral-800 rounded-md animate-pulse" />

          {/* HR line */}
          <div className="h-px w-full bg-neutral-800" />

          {/* Achievements icons */}
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
