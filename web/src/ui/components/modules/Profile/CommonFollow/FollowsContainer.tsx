import { getUserFollows } from '@/app/api/user'

import FollowItem from './FollowItem'

type Props = {
  id: string
}

export default async function FollowsContainer({ id }: Props) {
  const follows = await getUserFollows(id)

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
      {follows && follows.length > 0 && follows.map(follow => <FollowItem key={follow.userId} {...follow} />)}
    </div>
  )
}
