'use client'

import { GetUserResponse as UserBio } from '@/lib/actions/user.actions'
import { useUserProfileQuery } from '@/lib/hooks/api/user/useUserProfile'

import FollowConnections from './FollowConnections'

type UserBioSectionProps = {
  userId: string
  initialUserData: UserBio
}

export default function UserBioSection({ userId, initialUserData }: UserBioSectionProps) {
  const { data: userData } = useUserProfileQuery(userId, initialUserData)

  const currentData = userData ?? initialUserData

  const { username, about } = currentData.user
  const { followers, following } = currentData.userFollowStats

  return (
    <div className="flex flex-col justify-between">
      <div className="flex flex-col">
        <div className="flex items-center gap-3">
          <h2 className="font-medium text-xl md:text-[24px] text-white leading-tight">{username}</h2>
          {userData?.isPremiumUser && (
            <span className="py-1 px-2 bg-emerald-500 rounded-lg text-sm font-medium text-neutral-100">Premium</span>
          )}
        </div>
        {about && <p className="text-sm max-w-[400px] text-neutral-400">{about}</p>}
      </div>

      <FollowConnections followersCount={followers} followsCount={following} userId={userId} />
    </div>
  )
}
