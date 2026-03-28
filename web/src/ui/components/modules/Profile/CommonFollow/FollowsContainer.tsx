import { getUserFollows } from '@/app/api/user'
import EmptyState from '@/ui/components/shared/EmptyState'

import FollowItem from './FollowItem'

type Props = {
  id: string
}

export default async function FollowsContainer({ id }: Props) {
  const follows = await getUserFollows(id)

  if (!follows.data || !follows.success) {
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
      {follows.data.length > 0 && follows.data.map(follow => <FollowItem key={follow.userId} {...follow} />)}
    </div>
  )
}
