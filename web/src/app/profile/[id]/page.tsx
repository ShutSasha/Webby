import { Suspense } from 'react'

import ProfileTab from '@/ui/components/modules/Profile/ProfileTab'
import UserPlaylists, { UserPlaylistsSkeleton } from '@/ui/components/modules/Profile/UserPlaylists'
import UserVideos, { UserVideosSkeleton } from '@/ui/components/modules/Profile/UserVideos'

type ProfileProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const allowedTabs = ['Video', 'Playlist'] as const

export default async function Profile({ searchParams }: ProfileProps) {
  const { tab } = await searchParams
  const currentTab = allowedTabs.includes(tab as any) ? tab : 'Video'

  return (
    <div className="bg-neutral-900 rounded-[20px] p-5 flex flex-col gap-5">
      <div className="flex gap-2 items-center">
        <ProfileTab label="Video" />
        <ProfileTab label="Playlist" />
      </div>

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
