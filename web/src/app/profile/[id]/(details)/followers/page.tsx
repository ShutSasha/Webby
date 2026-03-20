import { Suspense } from 'react'

import { notFound } from 'next/navigation'

import { getUser } from '@/app/api/user'
import FollowersContainer from '@/ui/components/modules/Profile/CommonFollow/FollowersContainer'
import { FollowsContainerSkeleton } from '@/ui/components/modules/Profile/CommonFollow/FollowsContainerSkeleton'
import { FollowersToggle } from '@/ui/components/modules/Profile/CommonFollow/FollowToggle'
import HeaderNavigation from '@/ui/components/modules/Profile/CommonFollow/HeaderNavigation'

type Props = {
  params: Promise<{ id: string }>
}

export default async function Followers({ params }: Props) {
  const { id } = await params
  const user = await getUser(id)

  if (!user) {
    notFound()
  }

  return (
    <div className="flex flex-col gap-2">
      {/* Header Navigation */}

      <HeaderNavigation id={id} image={user.user.avatarUrl} username={user.user.username}>
        <FollowersToggle id={id} />
      </HeaderNavigation>

      <hr className="border-emerald-400/50 mb-1" />

      {/* Grid List */}
      <Suspense fallback={<FollowsContainerSkeleton />}>
        <FollowersContainer id={id} />
      </Suspense>
    </div>
  )
}
