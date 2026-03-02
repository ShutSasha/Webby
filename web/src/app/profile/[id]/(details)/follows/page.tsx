import { Suspense } from 'react'

import FollowsContainer from '@/ui/components/modules/Profile/CommonFollow/FollowsContainer'
import { FollowsContainerSkeleton } from '@/ui/components/modules/Profile/CommonFollow/FollowsContainerSkeleton'
import { FollowsToggle } from '@/ui/components/modules/Profile/CommonFollow/FollowToggle'
import HeaderNavigation from '@/ui/components/modules/Profile/CommonFollow/HeaderNavigation'

type Props = {
  params: Promise<{ id: string }>
}

export default async function Follows({ params }: Props) {
  const { id } = await params

  return (
    <div className="flex flex-col gap-2">
      {/* Header Navigation */}
      <HeaderNavigation id={id} image="https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg">
        <FollowsToggle id={id} />
      </HeaderNavigation>

      <hr className="border-emerald-400/50 mb-1" />

      {/* Grid List */}

      <Suspense fallback={<FollowsContainerSkeleton />}>
        <FollowsContainer />
      </Suspense>
    </div>
  )
}
