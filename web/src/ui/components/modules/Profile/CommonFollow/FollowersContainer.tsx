'use client'

import { useUserFollowersQuery } from '@/lib/hooks/api/user/useUserFollowers'
import EmptyState from '@/ui/components/shared/EmptyState'

import FollowItem from './FollowItem'
import { FollowsContainerSkeleton } from './FollowsContainerSkeleton'

type Props = {
  id: string
}

export default function FollowersContainer({ id }: Props) {
  const { data: followers, isError, isLoading } = useUserFollowersQuery(id)

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

  if (!followers || followers.length === 0) {
    return (
      <div className="flex flex-col h-full w-full">
        <EmptyState title="No followers yet" description="This user doesn’t have any followers yet." />
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {followers.map(follower => (
        <FollowItem key={follower.userId} {...follower} isFollow={false} userProfileId={id} />
      ))}
    </div>
  )
}
