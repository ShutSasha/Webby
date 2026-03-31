import { getUserFollowers } from '@/lib/actions/user.actions'
import EmptyState from '@/ui/components/shared/EmptyState'

import FollowItem from './FollowItem'

type Props = {
  id: string
}

export default async function FollowersContainer({ id }: Props) {
  const followers = await getUserFollowers(id)

  if (!followers.data || !followers.success) {
    return (
      <div className="flex flex-col h-full w-full">
        <EmptyState
          title="Couldn't load the list"
          description="We ran into a problem while loading these users. Please try refreshing the page."
        />
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {followers.data.length > 0 && followers.data.map(follower => <FollowItem key={follower.userId} {...follower} />)}
    </div>
  )
}
