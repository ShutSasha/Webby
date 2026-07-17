'use client'

import { useSession } from 'next-auth/react'

import { useUserFollowsQuery } from '@/lib/hooks/api/user/useUserFollows'
import EmptyState from '@/ui/components/shared/EmptyState'

import FollowItem from './FollowItem'
import { FollowsContainerSkeleton } from './FollowsContainerSkeleton'

type Props = {
  id: string
}

export default function FollowsContainer({ id }: Props) {
  const { data: session } = useSession()
  const { data: follows, isError, isLoading } = useUserFollowsQuery(id)

  if (isLoading) {
    return <FollowsContainerSkeleton />
  }

  if (isError) {
    return (
      <div className="flex flex-col h-full w-full">
        <EmptyState
          title="Couldn't load the list"
          description="We ran into a problem while loading these users. Please try refreshing the page."
        />
      </div>
    )
  }

  if (!follows || follows.length === 0) {
    return (
      <div className="flex flex-col h-full w-full">
        <EmptyState title="No follows yet" description="This user doesn’t have any follows yet." />
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {follows.length > 0 &&
        follows.map(follow => (
          <FollowItem
            key={follow.userId}
            {...follow}
            isFollow={true}
            currentUserId={session?.user.id}
            userProfileId={id}
          />
        ))}
    </div>
  )
}
