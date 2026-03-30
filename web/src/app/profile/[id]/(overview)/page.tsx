import { Suspense } from 'react'

import UserPlaylists, { UserPlaylistsSkeleton } from '@/ui/components/modules/Profile/UserPlaylists'
import UserVideos, { UserVideosSkeleton } from '@/ui/components/modules/Profile/UserVideos'
import AnimatedTabs from '@/ui/components/shared/AnimatedTabs'

type ProfileProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

type Tab = 'Video' | 'Playlist'
const allowedTabs = ['Video', 'Playlist'] as const

export default async function Profile({ searchParams }: ProfileProps) {
  const { tab } = await searchParams
  const currentTab = allowedTabs.includes(tab as Tab) ? tab : 'Video'

  // TODO: extract logic to constants
  const profileTabs = [
    {
      label: 'Video',
      href: '?tab=Video',
      isActive: currentTab === 'Video',
    },
    {
      label: 'Playlist',
      href: '?tab=Playlist',
      isActive: currentTab === 'Playlist',
    },
  ]

  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col gap-5">
      <AnimatedTabs tabs={profileTabs} layoutId="profile-tabs" linkClassName="py-1" />

      {currentTab === 'Video' && (
        <Suspense fallback={<UserVideosSkeleton />}>
          <UserVideos />
        </Suspense>
      )}
      {currentTab === 'Playlist' && (
        <Suspense fallback={<UserPlaylistsSkeleton />}>
          <UserPlaylists />
        </Suspense>
      )}
    </div>
  )
}
