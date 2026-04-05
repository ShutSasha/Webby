import { Suspense } from 'react'

import { ALLOWED_PROFILE_TABS, PROFILE_TABS_CONFIG, ProfileTabValue } from '@/lib/constants/profile.constants'
import UserPlaylists, { UserPlaylistsSkeleton } from '@/ui/components/modules/Profile/UserPlaylists'
import UserVideos, { UserVideosSkeleton } from '@/ui/components/modules/Profile/UserVideos'
import AnimatedTabs from '@/ui/components/shared/AnimatedTabs'

type ProfileProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function Profile({ searchParams }: ProfileProps) {
  const { tab } = await searchParams
  const currentTab = ALLOWED_PROFILE_TABS.includes(tab as ProfileTabValue) ? (tab as ProfileTabValue) : 'video'

  const profileTabs = PROFILE_TABS_CONFIG.map(config => ({
    label: config.label,
    href: `?tab=${config.value}`,
    isActive: currentTab === config.value,
  }))

  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col gap-5">
      <AnimatedTabs tabs={profileTabs} layoutId="profile-tabs" linkClassName="py-1" />

      {currentTab === 'video' && (
        <Suspense fallback={<UserVideosSkeleton />}>
          <UserVideos />
        </Suspense>
      )}
      {currentTab === 'playlist' && (
        <Suspense fallback={<UserPlaylistsSkeleton />}>
          <UserPlaylists />
        </Suspense>
      )}
    </div>
  )
}
